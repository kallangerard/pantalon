package file

import (
	"testing"

	"github.com/kallangerard/pantalon/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectorFilter(t *testing.T) {
	items := []api.ConfigurationItem{
		{Name: "a", Labels: map[string]string{"env": "dev"}},
		{Name: "b", Labels: map[string]string{"env": "prod"}},
		{Name: "c"},
	}

	sel, err := api.ParseSelector("env=prod")
	require.NoError(t, err)
	assert.Equal(t, []api.ConfigurationItem{items[1]}, SelectorFilter(items, sel))

	sel, err = api.ParseSelector("env!=prod")
	require.NoError(t, err)
	assert.Equal(t, []api.ConfigurationItem{items[0], items[2]}, SelectorFilter(items, sel))
}
