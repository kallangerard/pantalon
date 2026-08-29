package file

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/kallangerard/pantalon/api"
	"github.com/stretchr/testify/assert"
)

func TestWalkDir(t *testing.T) {
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalCwd); err != nil {
			t.Fatal(err)
		}
	})
	root := path.Join("..", "testdata", "terraform", "single-dir")
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

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
	t.Cleanup(func() {
		if err := os.Chdir(originalCwd); err != nil {
			t.Fatal(err)
		}
	})
	root := path.Join("..", "testdata", "terraform", "nested-dir")
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

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
	t.Cleanup(func() {
		if err := os.Chdir(originalCwd); err != nil {
			t.Fatal(err)
		}
	})

	root := path.Join("..", "testdata", "terraform", "nested-dir")
	expectedPaths := []string{path.Join("parent", "pantalon.yaml")}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

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
	t.Cleanup(func() {
		if err := os.Chdir(originalCwd); err != nil {
			t.Fatal(err)
		}
	})
	root := path.Join("..", "testdata", "terraform", "sibling-dir")
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

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
	t.Cleanup(func() {
		if err := os.Chdir(originalCwd); err != nil {
			t.Fatal(err)
		}
	})
	root := path.Join("..", "testdata", "terraform", "single-dir")
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

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
	t.Cleanup(func() {
		if err := os.Chdir(originalCwd); err != nil {
			t.Fatal(err)
		}
	})
	root := path.Join("..", "testdata", "terraform", "sibling-dir")
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

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
	t.Cleanup(func() {
		if err := os.Chdir(originalCwd); err != nil {
			t.Fatal(err)
		}
	})
	root := path.Join("..", "testdata", "terraform", "empty-dir")
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

	result, err := Search()
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, result)
}

// chdirTemp creates a temporary directory, changes into it, and restores the
// working directory when the test finishes.
func chdirTemp(t *testing.T) string {
	t.Helper()
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalCwd); err != nil {
			t.Fatal(err)
		}
	})
	return root
}

// writeMarker creates dir and places a minimal valid marker file in it.
func writeMarker(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	doc := "apiVersion: " + api.PantalonVersion + "\nkind: " + api.TerraformKind + "\nmetadata:\n  name: " + name + "\n"
	if err := os.WriteFile(filepath.Join(dir, MarkerFile), []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestWalkDir_SkipsGitAndTerraformDirectories(t *testing.T) {
	chdirTemp(t)

	writeMarker(t, filepath.Join("terraform", "compute"), "compute")
	// Marker files a real checkout can contain but that are not root modules
	// of this repository.
	writeMarker(t, filepath.Join(".git", "modules", "vendored"), "vendored")
	writeMarker(t, filepath.Join("stacks", ".terraform", "modules", "cached"), "cached")

	paths, err := findFiles()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []string{filepath.Join("terraform", "compute", MarkerFile)}, paths)
}

func TestWalkDir_LexicalOrderAcrossManyDirectories(t *testing.T) {
	chdirTemp(t)

	expected := make([]string, 0, 128)
	for i := range 128 {
		dir := filepath.Join("terraform", fmt.Sprintf("stack%03d", i), "dev")
		writeMarker(t, dir, fmt.Sprintf("stack%03d-dev", i))
		expected = append(expected, filepath.Join(dir, MarkerFile))
	}

	// The walk is concurrent, so run it repeatedly to catch ordering that
	// depends on which goroutine finishes first.
	for range 5 {
		paths, err := findFiles()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected, paths)
	}
}

func TestWalkDir_MissingDirectoryReturnsError(t *testing.T) {
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalCwd); err != nil {
			t.Fatal(err)
		}
	})

	w := walker{sem: make(chan struct{}, 1)}
	_, err = w.walk(filepath.Join(t.TempDir(), "does-not-exist"))
	assert.Error(t, err)
}
