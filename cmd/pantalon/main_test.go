package main

import (
	"testing"

	"github.com/kallangerard/pantalon/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var filterTestItems = []api.ConfigurationItem{
	{
		Name:   "compute-dev",
		Path:   "terraform/compute/environments/dev/pantalon.yaml",
		Dir:    "terraform/compute/environments/dev",
		Labels: map[string]string{"environment": "dev", "team": "compute"},
	},
	{
		Name:   "compute-prod",
		Path:   "terraform/compute/environments/prod/pantalon.yaml",
		Dir:    "terraform/compute/environments/prod",
		Labels: map[string]string{"environment": "prod", "team": "compute"},
	},
	{
		Name:   "network-dev",
		Path:   "terraform/network/environments/dev/pantalon.yaml",
		Dir:    "terraform/network/environments/dev",
		Labels: map[string]string{"environment": "dev", "team": "network"},
	},
	{
		Name:   "network-prod",
		Path:   "terraform/network/environments/prod/pantalon.yaml",
		Dir:    "terraform/network/environments/prod",
		Labels: map[string]string{"environment": "prod", "team": "network"},
	},
}

func TestFilterItems_PathGlobDefined_NoChangedDirs(t *testing.T) {
	result, err := filterItems(filterTestItems, "", []string{"terraform/compute/**"}, "")
	require.NoError(t, err)
	assert.Equal(t, []api.ConfigurationItem{filterTestItems[0], filterTestItems[1]}, result)
}

func TestFilterItems_NoChangedDirs_ReturnsAllItems(t *testing.T) {
	result, err := filterItems(filterTestItems, "", nil, "")
	require.NoError(t, err)
	assert.Equal(t, filterTestItems, result)
}

func TestFilterItems_Selector(t *testing.T) {
	result, err := filterItems(filterTestItems, "", nil, "environment=prod")
	require.NoError(t, err)
	assert.Equal(t, []api.ConfigurationItem{filterTestItems[1], filterTestItems[3]}, result)
}

func TestFilterItems_SelectorCombinedWithGlobAndChangedDirs(t *testing.T) {
	changed := `["terraform/compute/environments/prod","terraform/network/environments/prod"]`
	result, err := filterItems(filterTestItems, changed, []string{"terraform/network/**"}, "team in (network,data),environment=prod")
	require.NoError(t, err)
	assert.Equal(t, []api.ConfigurationItem{filterTestItems[3]}, result)
}

func TestFilterItems_InvalidSelector(t *testing.T) {
	_, err := filterItems(filterTestItems, "", nil, "environment in (prod")
	assert.Error(t, err)
}
