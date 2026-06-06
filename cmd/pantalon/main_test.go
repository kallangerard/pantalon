package main

import (
	"testing"

	"github.com/kallangerard/pantalon/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var filterTestItems = []api.ConfigurationItem{
	{
		Name: "compute-dev",
		Path: "terraform/compute/environments/dev/pantalon.yaml",
		Dir:  "terraform/compute/environments/dev",
	},
	{
		Name: "compute-prod",
		Path: "terraform/compute/environments/prod/pantalon.yaml",
		Dir:  "terraform/compute/environments/prod",
	},
	{
		Name: "network-dev",
		Path: "terraform/network/environments/dev/pantalon.yaml",
		Dir:  "terraform/network/environments/dev",
	},
	{
		Name: "network-prod",
		Path: "terraform/network/environments/prod/pantalon.yaml",
		Dir:  "terraform/network/environments/prod",
	},
}

func TestFilterItems_PathGlobDefined_NoChangedDirs(t *testing.T) {
	result, err := filterItems(filterTestItems, "", []string{"terraform/compute/**"})
	require.NoError(t, err)
	assert.Equal(t, []api.ConfigurationItem{filterTestItems[0], filterTestItems[1]}, result)
}

func TestFilterItems_NoChangedDirs_ReturnsAllItems(t *testing.T) {
	result, err := filterItems(filterTestItems, "", nil)
	require.NoError(t, err)
	assert.Equal(t, filterTestItems, result)
}
