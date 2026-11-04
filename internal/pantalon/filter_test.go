package pantalon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var filterTestItems = []ConfigurationItem{
	{
		Name:   "compute-dev",
		Path:   "terraform/compute/environments/dev/yaml",
		Dir:    "terraform/compute/environments/dev",
		Labels: map[string]string{"environment": "dev", "team": "compute"},
	},
	{
		Name:   "compute-prod",
		Path:   "terraform/compute/environments/prod/yaml",
		Dir:    "terraform/compute/environments/prod",
		Labels: map[string]string{"environment": "prod", "team": "compute"},
	},
	{
		Name:   "network-dev",
		Path:   "terraform/network/environments/dev/yaml",
		Dir:    "terraform/network/environments/dev",
		Labels: map[string]string{"environment": "dev", "team": "network"},
	},
	{
		Name:   "network-prod",
		Path:   "terraform/network/environments/prod/yaml",
		Dir:    "terraform/network/environments/prod",
		Labels: map[string]string{"environment": "prod", "team": "network"},
	},
}

func TestFilter_PathGlobDefined_NoChangedDirs(t *testing.T) {
	result, err := Filter(filterTestItems, FilterOptions{PathGlobs: []string{"terraform/compute/**"}})
	require.NoError(t, err)
	assert.Equal(t, []ConfigurationItem{filterTestItems[0], filterTestItems[1]}, result)
}

func TestFilter_NoChangedDirs_ReturnsAllItems(t *testing.T) {
	result, err := Filter(filterTestItems, FilterOptions{})
	require.NoError(t, err)
	assert.Equal(t, filterTestItems, result)
}

func TestFilter_Selector(t *testing.T) {
	result, err := Filter(filterTestItems, FilterOptions{Selector: "environment=prod"})
	require.NoError(t, err)
	assert.Equal(t, []ConfigurationItem{filterTestItems[1], filterTestItems[3]}, result)
}

func TestFilter_SelectorCombinedWithGlobAndChangedDirs(t *testing.T) {
	changed := `["terraform/compute/environments/prod","terraform/network/environments/prod"]`
	result, err := Filter(filterTestItems, FilterOptions{ChangedDirs: changed, PathGlobs: []string{"terraform/network/**"}, Selector: "team in (network,data),environment=prod"})
	require.NoError(t, err)
	assert.Equal(t, []ConfigurationItem{filterTestItems[3]}, result)
}

func TestFilter_InvalidSelector(t *testing.T) {
	_, err := Filter(filterTestItems, FilterOptions{Selector: "environment in (prod"})
	assert.Error(t, err)
}
