package pantalon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseChangedDirs_ReturnsArrayOfDirectoryPaths(t *testing.T) {
	j := `[".", "foo", "foo/bar"]`

	expectedPaths := []string{".", "foo", "foo/bar"}

	actualPaths, err := ParseChangedDirs([]byte(j))
	assert.NoError(t, err)

	assert.Equal(t, expectedPaths, actualPaths)
}

// Equivalent spellings of a path must compare equal when matching, otherwise
// configurations are silently skipped.
func TestParseChangedDirs_NormalisesPaths(t *testing.T) {
	j := `["a/dev/", "./a/dev", "a//dev", "", "./"]`

	expectedPaths := []string{"a/dev", "a/dev", "a/dev", ".", "."}

	actualPaths, err := ParseChangedDirs([]byte(j))
	assert.NoError(t, err)

	assert.Equal(t, expectedPaths, actualPaths)
}

// Directories whose names contain dots, such as .github, are common in
// changed-files output and must be accepted.
func TestParseChangedDirs_AcceptsDotDirectories(t *testing.T) {
	j := `[".github", ".github/workflows", ".tflint.d", "env.v2"]`

	expectedPaths := []string{".github", ".github/workflows", ".tflint.d", "env.v2"}

	actualPaths, err := ParseChangedDirs([]byte(j))
	assert.NoError(t, err)

	assert.Equal(t, expectedPaths, actualPaths)
}

// Paths outside the repository root can never match a configuration, so they
// are rejected rather than silently ignored.
func TestParseChangedDirs_RejectsPathsOutsideRepo(t *testing.T) {
	j := `["/home/runner/work/a/dev", "../a/dev", "..", "a/../../b"]`

	actualPaths, err := ParseChangedDirs([]byte(j))
	assert.ErrorContains(t, err, "/home/runner/work/a/dev: "+ErrChangedDirOutsideRepo.Error())
	assert.ErrorContains(t, err, "../a/dev: "+ErrChangedDirOutsideRepo.Error())
	assert.ErrorContains(t, err, "..: "+ErrChangedDirOutsideRepo.Error())
	assert.ErrorContains(t, err, "a/../../b: "+ErrChangedDirOutsideRepo.Error())
	assert.Nil(t, actualPaths, "actualPaths should be nil")
}

func TestChangedDirs_SimpleDir(t *testing.T) {
	cfgs := []TerraformConfiguration{
		{
			Metadata: Metadata{Name: "item1"},
			Path:     "a/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []ConfigurationItem{
		{
			Name: "item1",
			Path: "a/pantalon.yaml",
			Dir:  "a",
		},
	}

	changedFilesJson := []byte(`["a"]`)

	items, err := MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := ParseChangedDirs(changedFilesJson)
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

	cfgs := []TerraformConfiguration{
		{
			Metadata: Metadata{Name: "item1"},
			Path:     "a/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []ConfigurationItem{}

	changedFilesJson := []byte(`[]`)

	items, err := MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := ParseChangedDirs(changedFilesJson)
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

	cfgs := []TerraformConfiguration{
		{
			Metadata: Metadata{Name: "item1"},
			Path:     "path/a/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []ConfigurationItem{}

	changedFilesJson := []byte(`["path/b"]`)

	items, err := MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := ParseChangedDirs(changedFilesJson)
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

	cfgs := []TerraformConfiguration{
		{
			Metadata: Metadata{Name: "item1"},
			Path:     "a/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []ConfigurationItem{
		{
			Name: "item1",
			Path: "a/pantalon.yaml",
			Dir:  "a",
		},
	}

	changedFilesJson := []byte(`["."]`)

	items, err := MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := ParseChangedDirs(changedFilesJson)
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

	cfgs := []TerraformConfiguration{
		{
			Metadata: Metadata{Name: "item1"},
			Path:     "a/b/c/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []ConfigurationItem{
		{
			Name: "item1",
			Path: "a/b/c/pantalon.yaml",
			Dir:  "a/b/c",
		},
	}

	changedFilesJson := []byte(`["a/b"]`)

	items, err := MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := ParseChangedDirs(changedFilesJson)
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

	cfgs := []TerraformConfiguration{
		{
			Metadata: Metadata{Name: "item1"},
			Path:     "a/b/1/pantalon.yaml",
		},
		{
			Metadata: Metadata{Name: "item2"},
			Path:     "a/b/2/pantalon.yaml",
		},
		{
			Metadata: Metadata{Name: "item3"},
			Path:     "a/b/3/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []ConfigurationItem{
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

	items, err := MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := ParseChangedDirs(changedFilesJson)
	if err != nil {
		t.Fatal(err)
	}

	filteredCfgs, err := ChangedFiles(items, changedDirs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedFilteredCfgs, filteredCfgs)
}

// If multiple changed directories fall within the same Pantalon configuration,
// that configuration must only be returned once. Otherwise the same
// pantalon.yaml is emitted multiple times, producing duplicate matrix jobs.
func TestChangedDirs_MultipleChangedDirsWithinSameConfig(t *testing.T) {

	cfgs := []TerraformConfiguration{
		{
			Metadata: Metadata{Name: "item1"},
			Path:     "a/b/1/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []ConfigurationItem{
		{
			Name: "item1",
			Path: "a/b/1/pantalon.yaml",
			Dir:  "a/b/1",
		},
	}

	// Two files changed in different subdirectories of the same Pantalon
	// configuration (e.g. two separate modules under a/b/1) result in two
	// entries in changed-dirs that both match a/b/1.
	changedFilesJson := []byte(`["a/b/1/modules/foo", "a/b/1/modules/bar"]`)

	items, err := MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := ParseChangedDirs(changedFilesJson)
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

	cfgs := []TerraformConfiguration{
		{
			Metadata: Metadata{Name: "item1"},
			Path:     "a/b/1/pantalon.yaml",
		},
	}

	expectedFilteredCfgs := []ConfigurationItem{
		{
			Name: "item1",
			Path: "a/b/1/pantalon.yaml",
			Dir:  "a/b/1",
		},
	}

	changedFilesJson := []byte(`["a/b/1/modules/foo"]`)

	items, err := MarshalItems(cfgs)
	if err != nil {
		t.Fatal(err)
	}

	changedDirs, err := ParseChangedDirs(changedFilesJson)
	if err != nil {
		t.Fatal(err)
	}

	filteredCfgs, err := ChangedFiles(items, changedDirs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedFilteredCfgs, filteredCfgs)
}

// Directories that merely share a string prefix are siblings, not nested, and
// must not match each other.
func TestChangedDirs_SegmentBoundaries(t *testing.T) {
	tests := []struct {
		name       string
		cfgPath    string
		changedDir string
		want       bool
	}{
		{"sibling with shared prefix", "a/dev/pantalon.yaml", "a/dev-2", false},
		{"sibling config with shared prefix", "a/dev-2/pantalon.yaml", "a/dev", false},
		{"sibling with shared prefix nested", "a/dev/pantalon.yaml", "a/dev-2/modules", false},
		{"exact match", "a/dev/pantalon.yaml", "a/dev", true},
		{"changed dir inside config", "a/dev/pantalon.yaml", "a/dev/modules", true},
		{"config inside changed dir", "a/dev/pantalon.yaml", "a", true},
		{"root config", "pantalon.yaml", "a/dev", true},
		{"root changed dir", "a/dev/pantalon.yaml", ".", true},
		{"trailing slash on parent", "a/dev/pantalon.yaml", "a/", true},
		{"trailing slash on exact", "a/dev/pantalon.yaml", "a/dev/", true},
		{"dot slash prefix", "a/dev/pantalon.yaml", "./a/dev", true},
		{"empty changed dir", "a/dev/pantalon.yaml", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := MarshalItems([]TerraformConfiguration{
				{Metadata: Metadata{Name: "item1"}, Path: tt.cfgPath},
			})
			if err != nil {
				t.Fatal(err)
			}

			changedDirs, err := ParseChangedDirs([]byte(`["` + tt.changedDir + `"]`))
			if err != nil {
				t.Fatal(err)
			}

			filteredCfgs, err := ChangedFiles(items, changedDirs)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.want, len(filteredCfgs) == 1)
		})
	}
}

func TestChangedDirs_Dependencies(t *testing.T) {
	tests := []struct {
		name       string
		deps       []string
		changedDir string
		want       bool
	}{
		{"plain path exact", []string{"modules/net"}, "modules/net", true},
		{"plain path descendant", []string{"modules/net"}, "modules/net/v2", true},
		{"plain path sibling with shared prefix", []string{"modules/net"}, "modules/net-2", false},
		{"plain path parent changed", []string{"modules/net"}, "modules", false},
		{"double star", []string{"modules/**"}, "modules/net/v2", true},
		{"double star matches base dir", []string{"modules/**"}, "modules", true},
		{"single star segment", []string{"modules/*/v2"}, "modules/net/v2", true},
		{"single star segment descendant", []string{"modules/*/v2"}, "modules/net/v2/sub", true},
		{"single star segment other version", []string{"modules/*/v2"}, "modules/net/v1", false},
		{"trailing slash pattern", []string{"modules/net/"}, "modules/net", true},
		{"negation excludes", []string{"modules/**", "!modules/*/v1"}, "modules/net/v1", false},
		{"negation excludes descendant", []string{"modules/**", "!modules/*/v1"}, "modules/net/v1/sub", false},
		{"negation leaves others", []string{"modules/**", "!modules/*/v1"}, "modules/net/v2", true},
		{"last match wins re-include", []string{"modules/**", "!modules/*/v1", "modules/net/v1"}, "modules/net/v1", true},
		{"unrelated", []string{"modules/**"}, "other", false},
		{"no dependencies", nil, "modules/net", false},
		// Negations only apply to dependencies, never to the configuration's own directory.
		{"negation does not exclude own dir", []string{"modules/**", "!env/**"}, "env/dev", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := MarshalItems([]TerraformConfiguration{
				{
					Metadata:     Metadata{Name: "item1"},
					Dependencies: Dependencies{Paths: tt.deps},
					Path:         "env/dev/pantalon.yaml",
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			changedDirs, err := ParseChangedDirs([]byte(`["` + tt.changedDir + `"]`))
			if err != nil {
				t.Fatal(err)
			}

			filteredCfgs, err := ChangedFiles(items, changedDirs)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.want, len(filteredCfgs) == 1)
		})
	}
}
