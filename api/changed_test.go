package api

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalChangedFileJson_ReturnsArrayOfDirectoryPaths(t *testing.T) {
	j := `[".", "foo", "foo/bar"]`

	expectedPaths := []string{
		filepath.Join("."),
		filepath.Join("foo"),
		filepath.Join("foo", "bar"),
	}

	actualPaths, err := UnmarshalChangedFileJson([]byte(j))
	assert.NoError(t, err)

	assert.Equal(t, expectedPaths, actualPaths)
}

// ChangedFiles must only be directories, not files
func TestUnmarshalChangedFileJson_RejectChangedFilesIsFiles(t *testing.T) {
	j := `["README.md", "foo/main.tf"]`

	actualPaths, err := UnmarshalChangedFileJson([]byte(j))
	assert.ErrorContains(t, err, "README.md: "+ErrChangedFileIsNotDirectory.Error())
	assert.ErrorContains(t, err, "foo/main.tf: "+ErrChangedFileIsNotDirectory.Error())
	assert.Nil(t, actualPaths, "actualPaths should be nil")
}
