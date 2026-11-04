package pantalon

import (
	"io/fs"
	"os"
	"path"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalkDir(t *testing.T) {
	root := path.Join("..", "..", "testdata", "terraform", "single-dir")
	fsys := os.DirFS(root)

	expectedPaths := []string{path.Join("pantalon.yaml")}

	paths, err := findFiles(fsys)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPaths, paths)
}

func TestWalkDir_NestedFileShouldBeFound(t *testing.T) {
	root := path.Join("..", "..", "testdata", "terraform", "nested-dir")
	fsys := os.DirFS(root)

	expectedPath := path.Join("parent", "pantalon.yaml")

	paths, err := findFiles(fsys)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPath, paths[0])
}

func TestWalkDir_ChildDirectoriesShouldNotBeSearched(t *testing.T) {
	root := path.Join("..", "..", "testdata", "terraform", "nested-dir")
	expectedPaths := []string{path.Join("parent", "pantalon.yaml")}
	fsys := os.DirFS(root)

	paths, err := findFiles(fsys)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPaths, paths)
}

func TestWalkDir_SiblingDirectories(t *testing.T) {
	root := path.Join("..", "..", "testdata", "terraform", "sibling-dir")
	fsys := os.DirFS(root)

	expectedPaths := []string{
		path.Join("a", "pantalon.yaml"),
		path.Join("b", "pantalon.yaml"),
		path.Join("c", "pantalon.yaml"),
	}

	paths, err := findFiles(fsys)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedPaths, paths)
}

func TestReadFile_Success(t *testing.T) {
	fsys := os.DirFS(path.Join("..", "..", "testdata", "terraform", "single-dir"))

	cfg, _, err := readFile(fsys, "pantalon.yaml")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "pantalon.kallan.dev/v1alpha1", cfg.APIVersion)
	assert.Equal(t, "TerraformConfiguration", cfg.Kind)
	assert.Equal(t, "single-dir", cfg.Metadata.Name)
}

// If a single valid file exists the readFile function should return a single TerraformConfiguration.
func TestDiscover_Success(t *testing.T) {
	root := path.Join("..", "..", "testdata", "terraform", "single-dir")
	fsys := os.DirFS(root)

	expected := []TerraformConfiguration{
		{
			APIVersion: "pantalon.kallan.dev/v1alpha1",
			Kind:       "TerraformConfiguration",
			Metadata: Metadata{
				Name: "single-dir",
			},
			Path: path.Join("pantalon.yaml"),
		},
	}

	result, _, err := Discover(fsys)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

// If multiple valid files exist in the path the Discover function should return all of them as TerraformConfiguration.
func TestDiscover_MultipleFiles(t *testing.T) {
	root := path.Join("..", "..", "testdata", "terraform", "sibling-dir")
	fsys := os.DirFS(root)

	expected := []TerraformConfiguration{
		{
			APIVersion: "pantalon.kallan.dev/v1alpha1",
			Kind:       "TerraformConfiguration",
			Metadata: Metadata{
				Name: "sibling-dir-a",
			},
			Path: path.Join("a", "pantalon.yaml"),
		},
		{
			APIVersion: "pantalon.kallan.dev/v1alpha1",
			Kind:       "TerraformConfiguration",
			Metadata: Metadata{
				Name: "sibling-dir-b",
			},
			Path: path.Join("b", "pantalon.yaml"),
		},
		{
			APIVersion: "pantalon.kallan.dev/v1alpha1",
			Kind:       "TerraformConfiguration",
			Metadata: Metadata{
				Name: "sibling-dir-c",
			},
			Path: path.Join("c", "pantalon.yaml"),
		},
	}

	result, _, err := Discover(fsys)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

// If no files are found the Discover function should return an empty slice.
func TestDiscover_NoFilesFound(t *testing.T) {
	root := path.Join("..", "..", "testdata", "terraform", "empty-dir")
	fsys := os.DirFS(root)

	result, _, err := Discover(fsys)
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, result)
}

// writeConfigs returns an in-memory tree holding each path's content.
func writeConfigs(files map[string]string) fs.FS {
	m := fstest.MapFS{}
	for p, content := range files {
		m[p] = &fstest.MapFile{Data: []byte(content)}
	}
	return m
}

const configHeader = "apiVersion: pantalon.kallan.dev/v1alpha1\nkind: TerraformConfiguration\n"

func TestDiscover_ErrorNamesFile(t *testing.T) {
	fsys := writeConfigs(map[string]string{"a/pantalon.yaml": configHeader + "metadata:\n  name: Not_Valid\n"})

	_, _, err := Discover(fsys)
	assert.EqualError(t, err, `a/pantalon.yaml: invalid metadata.name "Not_Valid"`)
}

func TestDiscover_ValidationErrorNamesUnknownField(t *testing.T) {
	fsys := writeConfigs(map[string]string{"a/pantalon.yaml": configHeader + "metadata:\n  nmae: a\n"})

	_, _, err := Discover(fsys)
	assert.EqualError(t, err, `a/pantalon.yaml: invalid metadata.name "" (line 4: unknown field "metadata.nmae")`)
}

func TestDiscover_AnchorErrorNamesFileAndLine(t *testing.T) {
	fsys := writeConfigs(map[string]string{"a/pantalon.yaml": configHeader + "metadata:\n  name: a\nbase: &b 1\n"})

	_, _, err := Discover(fsys)
	assert.EqualError(t, err, "a/pantalon.yaml: line 5: YAML anchors, aliases and merge keys are not supported")
}

func TestDiscover_UnknownFieldWarns(t *testing.T) {
	fsys := writeConfigs(map[string]string{"a/pantalon.yaml": configHeader + "metadata:\n  name: a\ndependecies:\n  paths: [b]\n"})

	result, warnings, err := Discover(fsys)
	require.NoError(t, err)

	assert.Equal(t, "a", result[0].Metadata.Name)
	assert.Equal(t, []Warning{{Path: "a/pantalon.yaml", Line: 5, Message: `unknown field "dependecies"`}}, warnings)
}

func TestDiscover_DuplicateName(t *testing.T) {
	fsys := writeConfigs(map[string]string{
		"a/pantalon.yaml": configHeader + "metadata:\n  name: same\n",
		"b/pantalon.yaml": configHeader + "metadata:\n  name: same\n",
	})

	_, _, err := Discover(fsys)
	assert.EqualError(t, err, `duplicate metadata.name "same" in a/pantalon.yaml and b/pantalon.yaml`)
}
