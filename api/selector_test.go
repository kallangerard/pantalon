package api

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelector_Matches(t *testing.T) {
	labels := map[string]string{"env": "prod", "team": "data", "example.com/tier": "1"}

	tests := []struct {
		selector string
		want     bool
	}{
		{"", true},
		{"env=prod", true},
		{"env==prod", true},
		{"env=dev", false},
		{"missing=prod", false},
		{"env!=dev", true},
		{"env!=prod", false},
		{"missing!=prod", true},
		{"env in (prod,qa)", true},
		{"env in (dev,qa)", false},
		{"missing in (prod)", false},
		{"env notin (dev,qa)", true},
		{"env notin (prod)", false},
		{"missing notin (prod)", true},
		{"env", true},
		{"missing", false},
		{"!missing", true},
		{"!env", false},
		{"example.com/tier=1", true},
		{"env=prod,team=data", true},
		{"env=prod,team=compute", false},
		{" env = prod , team in ( data , compute ) ", true},
		{"team in (data,compute),!legacy,env!=dev", true},
	}

	for _, tt := range tests {
		t.Run(tt.selector, func(t *testing.T) {
			sel, err := ParseSelector(tt.selector)
			require.NoError(t, err)
			assert.Equal(t, tt.want, sel.Matches(labels))
		})
	}
}

func TestSelector_MatchesNilLabels(t *testing.T) {
	sel, err := ParseSelector("env!=prod,!team")
	require.NoError(t, err)
	assert.True(t, sel.Matches(nil))
}

func TestParseSelector_Invalid(t *testing.T) {
	tests := []string{
		"env=prod,",
		",env=prod",
		"=prod",
		"!",
		"env in (a",
		"env in a)",
		"env in a,b",
		"bad key=1",
		"env=bad value",
		"env=-bad",
		"/name=1",
		"Example.com/name=1",
		"env=" + strings.Repeat("a", 64),
		strings.Repeat("a", 64) + "=1",
	}

	for _, s := range tests {
		t.Run(s, func(t *testing.T) {
			_, err := ParseSelector(s)
			assert.Error(t, err)
		})
	}
}
