package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectorFlag_RejectsBlank(t *testing.T) {
	var f selectorFlag
	assert.Error(t, f.Set("  "))
}

func TestSelectorFlag_RejectsRepeat(t *testing.T) {
	var f selectorFlag
	assert.NoError(t, f.Set("environment=prod"))
	assert.Equal(t, "environment=prod", f.value)
	assert.Error(t, f.Set("team=a"))
}
