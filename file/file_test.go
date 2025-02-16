package file

import (
	"os"
	"path"
	"testing"

	"github.com/kallangerard/pantalon/api"
	"github.com/stretchr/testify/assert"
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
