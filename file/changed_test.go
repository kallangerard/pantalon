package file

import (
	"testing"

	"github.com/kallangerard/pantalon/api"
	"github.com/stretchr/testify/assert"
)

func TestChangedDirs_SimpleDir(t *testing.T) {
	cfgs := []api.TerraformConfiguration{
		{
			Metadata: api.Metadata{Name: "item1"},
			Path:     "a/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []api.ConfigurationItem{
		{
			Name: "item1",
			Path: "a/pantalon.yaml",
			Dir:  "a",
		},
	}

	changedFilesJson := []byte(`["a"]`)

	items, err := api.MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := api.UnmarshalChangedFileJson(changedFilesJson)
	if err != nil {
		t.Fatal(err)
	}

	filteredCfgs, err := ChangedFiles(items, changedDirs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedFilteredCfgs, filteredCfgs)
}

func TestChangedDirs_NoChanges(t *testing.T) {

	cfgs := []api.TerraformConfiguration{
		{
			Metadata: api.Metadata{Name: "item1"},
			Path:     "a/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []api.ConfigurationItem{}

	changedFilesJson := []byte(`[]`)

	items, err := api.MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := api.UnmarshalChangedFileJson(changedFilesJson)
	if err != nil {
		t.Fatal(err)
	}

	filteredCfgs, err := ChangedFiles(items, changedDirs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedFilteredCfgs, filteredCfgs)
}

func TestChangedDirs_UnrelatedChange(t *testing.T) {

	cfgs := []api.TerraformConfiguration{
		{
			Metadata: api.Metadata{Name: "item1"},
			Path:     "path/a/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []api.ConfigurationItem{}

	changedFilesJson := []byte(`["path/b"]`)

	items, err := api.MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := api.UnmarshalChangedFileJson(changedFilesJson)
	if err != nil {
		t.Fatal(err)
	}

	filteredCfgs, err := ChangedFiles(items, changedDirs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedFilteredCfgs, filteredCfgs)
}

func TestChangedDirs_HandleRootDir(t *testing.T) {

	cfgs := []api.TerraformConfiguration{
		{
			Metadata: api.Metadata{Name: "item1"},
			Path:     "a/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []api.ConfigurationItem{
		{
			Name: "item1",
			Path: "a/pantalon.yaml",
			Dir:  "a",
		},
	}

	changedFilesJson := []byte(`["."]`)

	items, err := api.MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := api.UnmarshalChangedFileJson(changedFilesJson)
	if err != nil {
		t.Fatal(err)
	}

	filteredCfgs, err := ChangedFiles(items, changedDirs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedFilteredCfgs, filteredCfgs)
}

func TestChangedDirs_NestedPantalonFile(t *testing.T) {

	cfgs := []api.TerraformConfiguration{
		{
			Metadata: api.Metadata{Name: "item1"},
			Path:     "a/b/c/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []api.ConfigurationItem{
		{
			Name: "item1",
			Path: "a/b/c/pantalon.yaml",
			Dir:  "a/b/c",
		},
	}

	changedFilesJson := []byte(`["a/b"]`)

	items, err := api.MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := api.UnmarshalChangedFileJson(changedFilesJson)
	if err != nil {
		t.Fatal(err)
	}

	filteredCfgs, err := ChangedFiles(items, changedDirs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedFilteredCfgs, filteredCfgs)
}

func TestChangedDirs_Mixed(t *testing.T) {

	cfgs := []api.TerraformConfiguration{
		{
			Metadata: api.Metadata{Name: "item1"},
			Path:     "a/b/1/pantalon.yaml",
		},
		{
			Metadata: api.Metadata{Name: "item2"},
			Path:     "a/b/2/pantalon.yaml",
		},
		{
			Metadata: api.Metadata{Name: "item3"},
			Path:     "a/b/3/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []api.ConfigurationItem{
		{
			Name: "item1",
			Path: "a/b/1/pantalon.yaml",
			Dir:  "a/b/1",
		},
		{
			Name: "item2",
			Path: "a/b/2/pantalon.yaml",
			Dir:  "a/b/2",
		},
	}

	changedFilesJson := []byte(`["a/b/1","a/b/2"]`)

	items, err := api.MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := api.UnmarshalChangedFileJson(changedFilesJson)
	if err != nil {
		t.Fatal(err)
	}

	filteredCfgs, err := ChangedFiles(items, changedDirs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedFilteredCfgs, filteredCfgs)
}

// If a directory changed is nested inside a Pantalon directory, the pantalon directory should be returned
func TestChangedDirs_ChangedDirInsidePantalonCfg(t *testing.T) {

	cfgs := []api.TerraformConfiguration{
		{
			Metadata: api.Metadata{Name: "item1"},
			Path:     "a/b/1/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []api.ConfigurationItem{
		{
			Name: "item1",
			Path: "a/b/1/pantalon.yaml",
			Dir:  "a/b/1",
		},
	}

	changedFilesJson := []byte(`["a/b/1/modules/foo"]`)

	items, err := api.MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := api.UnmarshalChangedFileJson(changedFilesJson)
	if err != nil {
		t.Fatal(err)
	}

	filteredCfgs, err := ChangedFiles(items, changedDirs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedFilteredCfgs, filteredCfgs)
}
