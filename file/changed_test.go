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
			Name:         "item1",
			Path:         "a/pantalon.yaml",
			Dir:          "a",
			Context:      nil,
			Dependencies: nil,
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
			Name:         "item1",
			Path:         "a/pantalon.yaml",
			Dir:          "a",
			Context:      nil,
			Dependencies: nil,
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
			Name:         "item1",
			Path:         "a/b/c/pantalon.yaml",
			Dir:          "a/b/c",
			Context:      nil,
			Dependencies: nil,
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
			Name:         "item1",
			Path:         "a/b/1/pantalon.yaml",
			Dir:          "a/b/1",
			Context:      nil,
			Dependencies: nil,
		},
		{
			Name:         "item2",
			Path:         "a/b/2/pantalon.yaml",
			Dir:          "a/b/2",
			Context:      nil,
			Dependencies: nil,
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
// Test that configurations are included when a dependency path changes
func TestChangedDirs_DependencyChanged(t *testing.T) {
	cfgs := []api.TerraformConfiguration{
		{
			Metadata:     api.Metadata{Name: "app1"},
			Path:         "environments/dev/pantalon.yaml",
			Dependencies: []string{"shared/modules/vpc", "shared/modules/database"},
		},
		{
			Metadata: api.Metadata{Name: "app2"},
			Path:     "environments/prod/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []api.ConfigurationItem{
		{
			Name:         "app1",
			Path:         "environments/dev/pantalon.yaml",
			Dir:          "environments/dev",
			Context:      nil,
			Dependencies: []string{"shared/modules/vpc", "shared/modules/database"},
		},
	}

	changedFilesJson := []byte(`["shared/modules/vpc"]`)

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

// Test that configurations are included when a nested dependency path changes
func TestChangedDirs_NestedDependencyChanged(t *testing.T) {
	cfgs := []api.TerraformConfiguration{
		{
			Metadata:     api.Metadata{Name: "web-app"},
			Path:         "apps/web/pantalon.yaml",
			Dependencies: []string{"modules/shared"},
		},
	}

	expectedFilteredCfgs := []api.ConfigurationItem{
		{
			Name:         "web-app",
			Path:         "apps/web/pantalon.yaml",
			Dir:          "apps/web",
			Context:      nil,
			Dependencies: []string{"modules/shared"},
		},
	}

	// A change inside the dependency should trigger the configuration
	changedFilesJson := []byte(`["modules/shared/vpc"]`)

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

// Test that configurations are included when change is nested inside dependency path
func TestChangedDirs_DependencyContainsChange(t *testing.T) {
	cfgs := []api.TerraformConfiguration{
		{
			Metadata:     api.Metadata{Name: "backend"},
			Path:         "services/backend/pantalon.yaml",
			Dependencies: []string{"common/modules/database/postgres"},
		},
	}

	expectedFilteredCfgs := []api.ConfigurationItem{
		{
			Name:         "backend",
			Path:         "services/backend/pantalon.yaml",
			Dir:          "services/backend",
			Context:      nil,
			Dependencies: []string{"common/modules/database/postgres"},
		},
	}

	// A change at the parent level should trigger configurations that depend on nested paths
	changedFilesJson := []byte(`["common/modules/database"]`)

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

// Test that multiple configurations with overlapping dependencies work correctly
func TestChangedDirs_MultipleDependencies(t *testing.T) {
	cfgs := []api.TerraformConfiguration{
		{
			Metadata:     api.Metadata{Name: "frontend"},
			Path:         "apps/frontend/pantalon.yaml",
			Dependencies: []string{"shared/modules/cdn", "shared/modules/dns"},
		},
		{
			Metadata:     api.Metadata{Name: "backend"},
			Path:         "apps/backend/pantalon.yaml",
			Dependencies: []string{"shared/modules/database", "shared/modules/dns"},
		},
		{
			Metadata: api.Metadata{Name: "monitoring"},
			Path:     "infra/monitoring/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []api.ConfigurationItem{
		{
			Name:         "frontend",
			Path:         "apps/frontend/pantalon.yaml",
			Dir:          "apps/frontend",
			Context:      nil,
			Dependencies: []string{"shared/modules/cdn", "shared/modules/dns"},
		},
		{
			Name:         "backend",
			Path:         "apps/backend/pantalon.yaml",
			Dir:          "apps/backend",
			Context:      nil,
			Dependencies: []string{"shared/modules/database", "shared/modules/dns"},
		},
	}

	// A change to shared DNS module should affect both frontend and backend
	changedFilesJson := []byte(`["shared/modules/dns"]`)

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

// Test that normal directory changes still work with dependencies
func TestChangedDirs_MixedDirectoryAndDependencyChanges(t *testing.T) {
	cfgs := []api.TerraformConfiguration{
		{
			Metadata:     api.Metadata{Name: "app1"},
			Path:         "apps/app1/pantalon.yaml",
			Dependencies: []string{"shared/vpc"},
		},
		{
			Metadata:     api.Metadata{Name: "app2"},
			Path:         "apps/app2/pantalon.yaml",
			Dependencies: []string{"shared/database"},
		},
	}

	expectedFilteredCfgs := []api.ConfigurationItem{
		{
			Name:         "app1",
			Path:         "apps/app1/pantalon.yaml",
			Dir:          "apps/app1",
			Context:      nil,
			Dependencies: []string{"shared/vpc"},
		},
		{
			Name:         "app2",
			Path:         "apps/app2/pantalon.yaml",
			Dir:          "apps/app2",
			Context:      nil,
			Dependencies: []string{"shared/database"},
		},
	}

	// Changes to app1 directory and shared/database should affect both apps
	changedFilesJson := []byte(`["apps/app1", "shared/database"]`)

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

func TestChangedDirs_ChangedDirInsidePantalonCfg(t *testing.T) {

	cfgs := []api.TerraformConfiguration{
		{
			Metadata: api.Metadata{Name: "item1"},
			Path:     "a/b/1/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []api.ConfigurationItem{
		{
			Name:         "item1",
			Path:         "a/b/1/pantalon.yaml",
			Dir:          "a/b/1",
			Context:      nil,
			Dependencies: nil,
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
