package file

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalkDir(t *testing.T) {
	root := path.Join("..", "testdata", "terraform", "single-dir")
	expectedPaths := []string{path.Join(root, "pantalon.yaml")}

	paths, err := FindPantalonFiles(root)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPaths, paths)
}

func TestWalkDir_NestedFileShouldBeFound(t *testing.T) {
	root := path.Join("..", "testdata", "terraform", "nested-dir")
	expectedPath := path.Join(root, "parent", "pantalon.yaml")

	paths, err := FindPantalonFiles(root)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPath, paths[0])
}

func TestWalkDir_ChildDirectoriesShouldNotBeSearched(t *testing.T) {
	root := path.Join("..", "testdata", "terraform", "nested-dir")
	expectedPaths := []string{path.Join(root, "parent", "pantalon.yaml")}

	paths, err := FindPantalonFiles(root)
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

	paths, err := FindPantalonFiles(root)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPaths, paths)
}
