package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalChangedFileJson_ReturnsArrayOfDirectoryPaths(t *testing.T) {
	j := `[".", "foo", "foo/bar"]`

	expectedPaths := []string{".", "foo", "foo/bar"}

	actualPaths, err := UnmarshalChangedFileJson([]byte(j))
	assert.NoError(t, err)

	assert.Equal(t, expectedPaths, actualPaths)
}

// Equivalent spellings of a path must compare equal when matching, otherwise
// configurations are silently skipped.
func TestUnmarshalChangedFileJson_NormalisesPaths(t *testing.T) {
	j := `["a/dev/", "./a/dev", "a//dev", "", "./"]`

	expectedPaths := []string{"a/dev", "a/dev", "a/dev", ".", "."}

	actualPaths, err := UnmarshalChangedFileJson([]byte(j))
	assert.NoError(t, err)

	assert.Equal(t, expectedPaths, actualPaths)
}

// Directories whose names contain dots, such as .github, are common in
// changed-files output and must be accepted.
func TestUnmarshalChangedFileJson_AcceptsDotDirectories(t *testing.T) {
	j := `[".github", ".github/workflows", ".tflint.d", "env.v2"]`

	expectedPaths := []string{".github", ".github/workflows", ".tflint.d", "env.v2"}

	actualPaths, err := UnmarshalChangedFileJson([]byte(j))
	assert.NoError(t, err)

	assert.Equal(t, expectedPaths, actualPaths)
}

// Paths outside the repository root can never match a configuration, so they
// are rejected rather than silently ignored.
func TestUnmarshalChangedFileJson_RejectsPathsOutsideRepo(t *testing.T) {
	j := `["/home/runner/work/a/dev", "../a/dev", "..", "a/../../b"]`

	actualPaths, err := UnmarshalChangedFileJson([]byte(j))
	assert.ErrorContains(t, err, "/home/runner/work/a/dev: "+ErrChangedDirOutsideRepo.Error())
	assert.ErrorContains(t, err, "../a/dev: "+ErrChangedDirOutsideRepo.Error())
	assert.ErrorContains(t, err, "..: "+ErrChangedDirOutsideRepo.Error())
	assert.ErrorContains(t, err, "a/../../b: "+ErrChangedDirOutsideRepo.Error())
	assert.Nil(t, actualPaths, "actualPaths should be nil")
}
