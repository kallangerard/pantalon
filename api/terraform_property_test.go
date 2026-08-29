package api

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

// referenceSubdomainLabel is the regular expression isValidSubdomainLabel
// replaced. It is kept here as the oracle for the hand rolled implementation.
var referenceSubdomainLabel = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

func referenceIsValidSubdomainLabel(s string) bool {
	return referenceSubdomainLabel.MatchString(s) && len(s) <= 253
}

// Property: the hand rolled validator agrees with the regular expression for
// any string, including ones drawn from the alphabet it accepts.
func TestIsValidSubdomainLabel_Property_MatchesReference(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		s := rapid.OneOf(
			rapid.String(),
			rapid.StringMatching(`[a-z0-9\-]{0,30}`),
			rapid.StringMatching(`[a-zA-Z0-9\-_.]{0,10}`),
		).Draw(t, "s")

		assert.Equal(t, referenceIsValidSubdomainLabel(s), isValidSubdomainLabel(s), "input: %q", s)
	})
}

func TestIsValidSubdomainLabel_Edges(t *testing.T) {
	cases := map[string]bool{
		"":                      false,
		"a":                     true,
		"0":                     true,
		"-":                     false,
		"--":                    false,
		"a-":                    false,
		"-a":                    false,
		"a-b":                   true,
		"a--b":                  true,
		"hello-world":           true,
		"Hello":                 false,
		"hello_world":           false,
		"hello.world":           false,
		"hello world":           false,
		"héllo":                 false,
		string(make([]byte, 1)): false,
	}

	for input, want := range cases {
		assert.Equal(t, want, isValidSubdomainLabel(input), "input: %q", input)
	}

	long := make([]byte, 253)
	for i := range long {
		long[i] = 'a'
	}
	assert.True(t, isValidSubdomainLabel(string(long)))
	assert.False(t, isValidSubdomainLabel(string(long)+"a"))
}
