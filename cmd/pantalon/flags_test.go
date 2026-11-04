package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectorFlag(t *testing.T) {
	var f selectorFlag
	assert.Error(t, f.Set("  "))
	assert.NoError(t, f.Set("environment=prod"))
	assert.Equal(t, "environment=prod", f.value)
	assert.Error(t, f.Set("team=a"))
}
