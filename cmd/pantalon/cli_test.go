package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests run the built binary and compare exit code, stdout and stderr
// with testdata/golden/*.golden. They pin what users see, so refactors must
// not change them. Regenerate deliberately with: go test ./cmd/pantalon -update
var update = flag.Bool("update", false, "rewrite golden files")

var binary string

func TestMain(m *testing.M) {
	flag.Parse()
	dir, err := os.MkdirTemp("", "pantalon-bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	binary = filepath.Join(dir, "pantalon")
	if out, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n%s", err, out)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

const header = "apiVersion: pantalon.kallan.dev/v1alpha1\nkind: TerraformConfiguration\n"

type cliCase struct {
	name  string
	args  []string
	env   []string
	files map[string]string // fixture written to a temp dir; if nil, run in examples/
}

var cliCases = []cliCase{
	{name: "default-json"},
	{name: "yaml", args: []string{"--output-format=yaml"}},
	{name: "changed-dirs", args: []string{`--changed-dirs=["terraform/compute/environments/dev"]`}},
	{name: "changed-dir-parent", args: []string{`--changed-dirs=["terraform/compute"]`}},
	{name: "changed-dir-unrelated", args: []string{`--changed-dirs=["docs"]`}},
	{name: "changed-dir-outside-repo", args: []string{`--changed-dirs=["../x"]`}},
	{name: "changed-dirs-bad-json", args: []string{`--changed-dirs=nope`}},
	{name: "path-glob", args: []string{"--path-glob=terraform/data/**"}},
	{name: "path-glob-repeated", args: []string{"--path-glob=terraform/data/**/dev", "--path-glob=terraform/compute/**/prod"}},
	{name: "path-glob-no-match", args: []string{"--path-glob=nothing/**"}},
	{name: "selector-short", args: []string{"-l", "environment=prod"}},
	{name: "selector-set", args: []string{"--selector=environment in (dev,qa),component!=data"}},
	{name: "selector-bad", args: []string{"--selector=environment in"}},
	{name: "selector-empty", args: []string{"--selector="}},
	{name: "selector-twice", args: []string{"-l", "a=b", "--selector=c=d"}},
	{name: "filters-combined", args: []string{
		`--changed-dirs=["terraform/data","terraform/compute"]`,
		"--path-glob=terraform/**/dev", "-l", "component=data",
	}},
	{name: "empty-result-json", args: []string{"-l", "nope=nope"}},
	{name: "empty-result-yaml", args: []string{"-l", "nope=nope", "--output-format=yaml"}},
	{name: "bad-output-format", args: []string{"--output-format=xml"}},
	{name: "unknown-flag", args: []string{"--nope"}},
	{name: "help", args: []string{"--help"}},
	{
		name:  "no-configs",
		files: map[string]string{"README.md": "hi\n"},
	},
	{
		name: "full-config",
		files: map[string]string{"a/pantalon.yaml": header + `
metadata:
  name: a
  labels: {env: prod}
dependencies:
  paths: ["modules/**", "!modules/old"]
context:
  sa: x@y.z
`},
	},
	{
		name: "dependency-path-selects-config",
		args: []string{`--changed-dirs=["modules/shared"]`},
		files: map[string]string{
			"a/pantalon.yaml": header + "metadata: {name: a}\ndependencies: {paths: [\"modules/**\"]}\n",
			"b/pantalon.yaml": header + "metadata: {name: b}\n",
		},
	},
	{
		name: "duplicate-name",
		files: map[string]string{
			"a/pantalon.yaml": header + "metadata: {name: dup}\n",
			"b/pantalon.yaml": header + "metadata: {name: dup}\n",
		},
	},
	{
		name:  "unknown-field-warning",
		files: map[string]string{"a/pantalon.yaml": header + "metadata: {name: a}\nfromTheFuture: true\n"},
	},
	{
		name:  "unknown-field-warning-github-actions",
		env:   []string{"GITHUB_ACTIONS=true"},
		files: map[string]string{"a/pantalon.yaml": header + "metadata: {name: a}\nfromTheFuture: true\n"},
	},
	{
		name:  "misspelt-required-field",
		files: map[string]string{"a/pantalon.yaml": header + "metadta: {name: a}\n"},
	},
	{
		name:  "invalid-name",
		files: map[string]string{"a/pantalon.yaml": header + "metadata: {name: Not_Valid}\n"},
	},
	{
		name:  "wrong-kind",
		files: map[string]string{"a/pantalon.yaml": "apiVersion: pantalon.kallan.dev/v1alpha1\nkind: Other\nmetadata: {name: a}\n"},
	},
	{
		name:  "invalid-yaml",
		files: map[string]string{"a/pantalon.yaml": "metadata: [unclosed\n"},
	},
	{
		name:  "yaml-anchor",
		files: map[string]string{"a/pantalon.yaml": header + "metadata: &m {name: a}\n"},
	},
	{
		name:  "bad-dependency-path",
		files: map[string]string{"a/pantalon.yaml": header + "metadata: {name: a}\ndependencies: {paths: [\"../x\"]}\n"},
	},
	{
		name: "stops-descent-at-first-match",
		files: map[string]string{
			"a/pantalon.yaml":     header + "metadata: {name: a}\n",
			"a/sub/pantalon.yaml": header + "metadata: {name: nested}\n",
		},
	},
}

var logTimestamp = regexp.MustCompile(`(?m)^\d{4}/\d\d/\d\d \d\d:\d\d:\d\d `)

func runCLI(t *testing.T, c cliCase) string {
	t.Helper()
	dir := t.TempDir()
	if c.files == nil {
		abs, err := filepath.Abs("../../examples")
		require.NoError(t, err)
		dir = abs
	}
	for name, content := range c.files {
		p := filepath.Join(dir, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o755))
		require.NoError(t, os.WriteFile(p, []byte(content), 0o644))
	}

	cmd := exec.Command(binary, c.args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GITHUB_ACTIONS=")
	cmd.Env = append(cmd.Env, c.env...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	code := 0
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		require.True(t, errors.As(err, &ee), "run: %v", err)
		code = ee.ExitCode()
	}

	// The log package timestamps stderr lines; that is noise, not contract.
	stderrText := logTimestamp.ReplaceAllString(stderr.String(), "")
	return fmt.Sprintf("exit: %d\n--- stdout ---\n%s--- stderr ---\n%s", code, stdout.String(), stderrText)
}

func TestCLIGolden(t *testing.T) {
	for _, c := range cliCases {
		t.Run(c.name, func(t *testing.T) {
			got := runCLI(t, c)
			path := filepath.Join("..", "..", "testdata", "golden", c.name+".golden")
			if *update {
				require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
				require.NoError(t, os.WriteFile(path, []byte(got), 0o644))
				return
			}
			want, err := os.ReadFile(path)
			require.NoError(t, err, "missing golden; run with -update")
			require.Equal(t, string(want), got)
		})
	}
}

// TestJSONSchema pins the JSON keys independently of the goldens, so a
// carelessly regenerated golden cannot approve a renamed field.
func TestJSONSchema(t *testing.T) {
	out := runCLI(t, cliCase{name: "schema", files: map[string]string{
		"a/pantalon.yaml": header + "metadata: {name: a}\ndependencies: {paths: [x]}\n",
	}})
	body := strings.TrimSuffix(strings.SplitN(strings.SplitN(out, "--- stdout ---\n", 2)[1], "--- stderr ---", 2)[0], "\n")
	var items []map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &items))
	require.Len(t, items, 1)
	keys := make([]string, 0)
	for k := range items[0] {
		keys = append(keys, k)
	}
	require.ElementsMatch(t, []string{"name", "path", "dir", "labels", "dependencies", "context"}, keys)
}
