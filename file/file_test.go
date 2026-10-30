package file

import (
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/kallangerard/pantalon/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalkDir(t *testing.T) {
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(originalCwd) })
	root := path.Join("..", "testdata", "terraform", "single-dir")
	os.Chdir(root)

	expectedPaths := []string{path.Join("pantalon.yaml")}

	paths, err := findFiles()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPaths, paths)
}

func TestWalkDir_NestedFileShouldBeFound(t *testing.T) {
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(originalCwd) })
	root := path.Join("..", "testdata", "terraform", "nested-dir")
	os.Chdir(root)

	expectedPath := path.Join("parent", "pantalon.yaml")

	paths, err := findFiles()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPath, paths[0])
}

func TestWalkDir_ChildDirectoriesShouldNotBeSearched(t *testing.T) {
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(originalCwd) })

	root := path.Join("..", "testdata", "terraform", "nested-dir")
	expectedPaths := []string{path.Join("parent", "pantalon.yaml")}
	os.Chdir(root)

	paths, err := findFiles()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPaths, paths)
}

func TestWalkDir_SiblingDirectories(t *testing.T) {
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(originalCwd) })
	root := path.Join("..", "testdata", "terraform", "sibling-dir")
	os.Chdir(root)

	expectedPaths := []string{
		path.Join("a", "pantalon.yaml"),
		path.Join("b", "pantalon.yaml"),
		path.Join("c", "pantalon.yaml"),
	}

	paths, err := findFiles()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPaths, paths)
}

func TestReadFile_Success(t *testing.T) {
	path := path.Join("..", "testdata", "terraform", "single-dir", "pantalon.yaml")

	cfg, err := readFile(path)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "pantalon.kallan.dev/v1alpha1", cfg.ApiVersion)
	assert.Equal(t, "TerraformConfiguration", cfg.Kind)
	assert.Equal(t, "single-dir", cfg.Metadata.Name)
}

// If a single valid file exists the readFile function should return a single api.TerraformConfiguration.
func TestSearch_Success(t *testing.T) {
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(originalCwd) })
	root := path.Join("..", "testdata", "terraform", "single-dir")
	os.Chdir(root)

	expected := []api.TerraformConfiguration{
		{
			ApiVersion: "pantalon.kallan.dev/v1alpha1",
			Kind:       "TerraformConfiguration",
			Metadata: api.Metadata{
				Name: "single-dir",
			},
			Path: path.Join("pantalon.yaml"),
		},
	}

	result, err := Search()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

// If multiple valid files exist in the path the Search function should return all of them as api.TerraformConfiguration.
func TestSearch_MultipleFiles(t *testing.T) {
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(originalCwd) })
	root := path.Join("..", "testdata", "terraform", "sibling-dir")
	os.Chdir(root)

	expected := []api.TerraformConfiguration{
		{
			ApiVersion: "pantalon.kallan.dev/v1alpha1",
			Kind:       "TerraformConfiguration",
			Metadata: api.Metadata{
				Name: "sibling-dir-a",
			},
			Path: path.Join("a", "pantalon.yaml"),
		},
		{
			ApiVersion: "pantalon.kallan.dev/v1alpha1",
			Kind:       "TerraformConfiguration",
			Metadata: api.Metadata{
				Name: "sibling-dir-b",
			},
			Path: path.Join("b", "pantalon.yaml"),
		},
		{
			ApiVersion: "pantalon.kallan.dev/v1alpha1",
			Kind:       "TerraformConfiguration",
			Metadata: api.Metadata{
				Name: "sibling-dir-c",
			},
			Path: path.Join("c", "pantalon.yaml"),
		},
	}

	result, err := Search()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

// If no files are found the Search function should return an empty slice.
func TestSearch_NoFilesFound(t *testing.T) {
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(originalCwd) })
	root := path.Join("..", "testdata", "terraform", "empty-dir")
	os.Chdir(root)

	result, err := Search()
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, result)
}

// writeConfigs writes each path's pantalon.yaml content under a temp dir and
// changes into it.
func writeConfigs(t *testing.T, files map[string]string) {
	dir := t.TempDir()
	for p, content := range files {
		full := filepath.Join(dir, p)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(originalCwd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
}

const configHeader = "apiVersion: pantalon.kallan.dev/v1alpha1\nkind: TerraformConfiguration\n"

func TestSearch_ErrorNamesFile(t *testing.T) {
	writeConfigs(t, map[string]string{"a/pantalon.yaml": configHeader + "metadata:\n  name: Not_Valid\n"})

	_, err := Search()
	assert.EqualError(t, err, `a/pantalon.yaml: invalid metadata.name "Not_Valid"`)
}

func TestSearch_ValidationErrorNamesUnknownField(t *testing.T) {
	writeConfigs(t, map[string]string{"a/pantalon.yaml": configHeader + "metadata:\n  nmae: a\n"})

	_, err := Search()
	assert.EqualError(t, err, `a/pantalon.yaml: invalid metadata.name "" (line 4: unknown field "metadata.nmae")`)
}

func TestSearch_AnchorErrorNamesFileAndLine(t *testing.T) {
	writeConfigs(t, map[string]string{"a/pantalon.yaml": configHeader + "metadata:\n  name: a\nbase: &b 1\n"})

	_, err := Search()
	assert.EqualError(t, err, "a/pantalon.yaml: line 5: YAML anchors, aliases and merge keys are not supported")
}

// captureLog returns the log output written during the test.
func captureLog(t *testing.T) *bytes.Buffer {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
	})
	return &buf
}

const typoConfig = configHeader + "metadata:\n  name: a\ndependecies:\n  paths: [b]\n"

func TestSearch_UnknownFieldWarns(t *testing.T) {
	t.Setenv("GITHUB_ACTIONS", "")
	writeConfigs(t, map[string]string{"a/pantalon.yaml": typoConfig})
	logs := captureLog(t)

	result, err := Search()
	require.NoError(t, err)

	assert.Equal(t, "a", result[0].Metadata.Name)
	assert.Equal(t, "warning: a/pantalon.yaml:5: unknown field \"dependecies\"\n", logs.String())
}

func TestSearch_UnknownFieldAnnotatesInGitHubActions(t *testing.T) {
	t.Setenv("GITHUB_ACTIONS", "true")
	writeConfigs(t, map[string]string{"a/pantalon.yaml": typoConfig})
	logs := captureLog(t)

	_, err := Search()
	require.NoError(t, err)

	assert.Equal(t, "::warning file=a/pantalon.yaml,line=5::unknown field \"dependecies\"\n", logs.String())
}

func TestSearch_DuplicateName(t *testing.T) {
	writeConfigs(t, map[string]string{
		"a/pantalon.yaml": configHeader + "metadata:\n  name: same\n",
		"b/pantalon.yaml": configHeader + "metadata:\n  name: same\n",
	})

	_, err := Search()
	assert.EqualError(t, err, `duplicate metadata.name "same" in a/pantalon.yaml and b/pantalon.yaml`)
}
