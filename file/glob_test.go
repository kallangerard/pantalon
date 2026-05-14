package file

import (
	"testing"

	"github.com/kallangerard/pantalon/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var globItems = []api.ConfigurationItem{
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

func TestGlobFilter_DoublestarRecursive(t *testing.T) {
	result, err := GlobFilter(globItems, []string{"terraform/compute/**"})
	require.NoError(t, err)
	assert.Equal(t, []api.ConfigurationItem{globItems[0], globItems[1]}, result)
}

func TestGlobFilter_SingleSegmentWildcard(t *testing.T) {
	result, err := GlobFilter(globItems, []string{"terraform/*/environments/dev"})
	require.NoError(t, err)
	assert.Equal(t, []api.ConfigurationItem{globItems[0], globItems[2]}, result)
}

func TestGlobFilter_MultiplePatterns_OR(t *testing.T) {
	result, err := GlobFilter(globItems, []string{"terraform/compute/**", "terraform/network/environments/prod"})
	require.NoError(t, err)
	assert.Equal(t, []api.ConfigurationItem{globItems[0], globItems[1], globItems[3]}, result)
}

func TestGlobFilter_NoPatterns_ReturnsAll(t *testing.T) {
	result, err := GlobFilter(globItems, []string{})
	require.NoError(t, err)
	assert.Equal(t, globItems, result)
}

func TestGlobFilter_NoMatch(t *testing.T) {
	result, err := GlobFilter(globItems, []string{"terraform/storage/**"})
	require.NoError(t, err)
	assert.Equal(t, []api.ConfigurationItem{}, result)
}

func TestGlobFilter_ExactMatch(t *testing.T) {
	result, err := GlobFilter(globItems, []string{"terraform/compute/environments/dev"})
	require.NoError(t, err)
	assert.Equal(t, []api.ConfigurationItem{globItems[0]}, result)
}

func TestGlobFilter_InvalidPattern(t *testing.T) {
	_, err := GlobFilter(globItems, []string{"["})
	assert.Error(t, err)
}
