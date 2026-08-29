package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/kallangerard/pantalon/api"
)

const benchMarker = `apiVersion: pantalon.kallan.dev/v1alpha1
kind: TerraformConfiguration
metadata:
  name: %s
context:
  service_account: ci@example.iam.gserviceaccount.com
`

// buildRepo lays out a synthetic repository under root: `stacks` stacks, each
// with `envs` root modules, plus a deep `.git`-shaped tree of noise directories
// that contain no marker files.
func buildRepo(tb testing.TB, root string, stacks, envs, noiseDirs int) {
	tb.Helper()

	for s := range stacks {
		for e := range envs {
			dir := filepath.Join(root, "terraform", "stack"+strconv.Itoa(s), "environments", "env"+strconv.Itoa(e))
			if err := os.MkdirAll(dir, 0o755); err != nil {
				tb.Fatal(err)
			}
			name := fmt.Sprintf("stack%d-env%d", s, e)
			if err := os.WriteFile(filepath.Join(dir, "pantalon.yaml"), []byte(fmt.Sprintf(benchMarker, name)), 0o644); err != nil {
				tb.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "main.tf"), []byte("# main\n"), 0o644); err != nil {
				tb.Fatal(err)
			}
			// Terraform's local cache: many directories, never root modules.
			cache := filepath.Join(dir, ".terraform", "modules", "vendored")
			if err := os.MkdirAll(cache, 0o755); err != nil {
				tb.Fatal(err)
			}
		}
	}

	// Noise that a real checkout carries: .git object shards and docs.
	for i := range noiseDirs {
		dir := filepath.Join(root, ".git", "objects", fmt.Sprintf("%02x", i%256), strconv.Itoa(i))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			tb.Fatal(err)
		}
	}
}

func benchInRepo(tb testing.TB, stacks, envs, noiseDirs int) {
	tb.Helper()
	root := tb.TempDir()
	buildRepo(tb, root, stacks, envs, noiseDirs)

	cwd, err := os.Getwd()
	if err != nil {
		tb.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(func() { _ = os.Chdir(cwd) })
}

func BenchmarkFindFiles(b *testing.B) {
	benchInRepo(b, 40, 5, 400)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		paths, err := findFiles()
		if err != nil {
			b.Fatal(err)
		}
		if len(paths) == 0 {
			b.Fatal("no marker files found")
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	benchInRepo(b, 40, 5, 400)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfgs, err := Search()
		if err != nil {
			b.Fatal(err)
		}
		if len(cfgs) == 0 {
			b.Fatal("no configurations found")
		}
	}
}

func benchItems(n int) []api.ConfigurationItem {
	items := make([]api.ConfigurationItem, n)
	for i := range items {
		dir := fmt.Sprintf("terraform/stack%d/environments/env%d", i/5, i%5)
		items[i] = api.ConfigurationItem{
			Name: fmt.Sprintf("stack%d-env%d", i/5, i%5),
			Path: dir + "/pantalon.yaml",
			Dir:  dir,
		}
	}
	return items
}

func BenchmarkChangedFiles(b *testing.B) {
	items := benchItems(200)
	changed := []string{
		"terraform/stack3/environments/env1",
		"terraform/stack17/environments/env4",
		"docs",
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ChangedFiles(items, changed); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkChangedFilesBroad(b *testing.B) {
	items := benchItems(200)
	changed := []string{"terraform"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ChangedFiles(items, changed); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGlobFilterLiteral(b *testing.B) {
	items := benchItems(200)
	patterns := []string{
		"terraform/stack1/environments/env1",
		"terraform/stack2/environments/env2",
		"terraform/stack3/environments/env3",
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := GlobFilter(items, patterns); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGlobFilter(b *testing.B) {
	items := benchItems(200)
	patterns := []string{"terraform/stack1/**", "terraform/stack2/**", "terraform/stack3/**"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := GlobFilter(items, patterns); err != nil {
			b.Fatal(err)
		}
	}
}
