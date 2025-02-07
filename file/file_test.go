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
