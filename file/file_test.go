package file

import (
	"path"
	"testing"

	"github.com/kallangerard/pantalon/api"
	"github.com/stretchr/testify/assert"
)

func TestWalkDir(t *testing.T) {
	root := path.Join("..", "testdata", "terraform", "single-dir")
	expectedPaths := []string{path.Join(root, "pantalon.yaml")}

	paths, err := findFiles(root)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPaths, paths)
}

func TestWalkDir_NestedFileShouldBeFound(t *testing.T) {
	root := path.Join("..", "testdata", "terraform", "nested-dir")
	expectedPath := path.Join(root, "parent", "pantalon.yaml")

	paths, err := findFiles(root)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPath, paths[0])
}

func TestWalkDir_ChildDirectoriesShouldNotBeSearched(t *testing.T) {
	root := path.Join("..", "testdata", "terraform", "nested-dir")
	expectedPaths := []string{path.Join(root, "parent", "pantalon.yaml")}

	paths, err := findFiles(root)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPaths, paths)
}

func TestWalkDir_SiblingDirectories(t *testing.T) {
	root := path.Join("..", "testdata", "terraform", "sibling-dir")
	expectedPaths := []string{
		path.Join(root, "a", "pantalon.yaml"),
		path.Join(root, "b", "pantalon.yaml"),
		path.Join(root, "c", "pantalon.yaml"),
	}

	paths, err := findFiles(root)
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
	root := path.Join("..", "testdata", "terraform", "single-dir")
	expected := []api.TerraformConfiguration{
		{
			ApiVersion: "pantalon.kallan.dev/v1alpha1",
			Kind:       "TerraformConfiguration",
			Metadata: api.Metadata{
				Name: "single-dir",
			},
			Path: path.Join(root, "pantalon.yaml"),
		},
	}

	result, err := Search(root)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

// If multiple valid files exist in the path the Search function should return all of them as api.TerraformConfiguration.
func TestSearch_MultipleFiles(t *testing.T) {
	root := path.Join("..", "testdata", "terraform", "sibling-dir")
	expected := []api.TerraformConfiguration{
		{
			ApiVersion: "pantalon.kallan.dev/v1alpha1",
			Kind:       "TerraformConfiguration",
			Metadata: api.Metadata{
				Name: "sibling-dir-a",
			},
			Path: path.Join(root, "a", "pantalon.yaml"),
		},
		{
			ApiVersion: "pantalon.kallan.dev/v1alpha1",
			Kind:       "TerraformConfiguration",
			Metadata: api.Metadata{
				Name: "sibling-dir-b",
			},
			Path: path.Join(root, "b", "pantalon.yaml"),
		},
		{
			ApiVersion: "pantalon.kallan.dev/v1alpha1",
			Kind:       "TerraformConfiguration",
			Metadata: api.Metadata{
				Name: "sibling-dir-c",
			},
			Path: path.Join(root, "c", "pantalon.yaml"),
		},
	}

	result, err := Search(root)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

// If no files are found the Search function should return an empty slice.
func TestSearch_NoFilesFound(t *testing.T) {
	root := path.Join("..", "testdata", "terraform", "empty-dir")

	result, err := Search(root)
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, result)
}

// If the root directory does not exist the Search function should return an error.
func TestSearch_InvalidDirectory(t *testing.T) {
	root := path.Join("..", "testdata", "terraform", "nonexistent-dir")

	_, err := Search(root)
	assert.Error(t, err)
}
