package file

import (
	"fmt"
	"strings"
	"testing"

	"pgregory.net/rapid"

	"github.com/kallangerard/pantalon/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// genSegment produces a path segment: starts with a-z, then any mix of
// lowercase letters, digits, hyphens, underscores, and dots (0-19 extra chars).
func genSegment() *rapid.Generator[string] {
	return rapid.StringMatching(`[a-z][a-z0-9\-_.]{0,19}`)
}

// genDirPath produces a slash-joined directory path with 1-8 segments.
func genDirPath() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		depth := rapid.IntRange(1, 8).Draw(t, "depth")
		segments := make([]string, depth)
		for i := range segments {
			segments[i] = genSegment().Draw(t, fmt.Sprintf("seg%d", i))
		}
		return strings.Join(segments, "/")
	})
}

// genItem produces a ConfigurationItem with a generated directory path.
func genItem() *rapid.Generator[api.ConfigurationItem] {
	return rapid.Custom(func(t *rapid.T) api.ConfigurationItem {
		dir := genDirPath().Draw(t, "dir")
		return api.ConfigurationItem{
			Name: genSegment().Draw(t, "name"),
			Dir:  dir,
			Path: dir + "/pantalon.yaml",
		}
	})
}

// genItems produces a slice of 0-10 ConfigurationItems.
func genItems() *rapid.Generator[[]api.ConfigurationItem] {
	return rapid.SliceOfN(genItem(), 0, 10)
}

// Property: no patterns returns all items unchanged.
func TestGlobFilter_Property_NoPatternsReturnsAll(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		items := genItems().Draw(t, "items")
		result, err := GlobFilter(items, []string{})
		require.NoError(t, err)
		assert.Equal(t, items, result)
	})
}

// Property: "**" matches every directory path regardless of depth or characters.
func TestGlobFilter_Property_DoubleStarMatchesAll(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		items := genItems().Draw(t, "items")
		result, err := GlobFilter(items, []string{"**"})
		require.NoError(t, err)
		assert.Equal(t, items, result)
	})
}

// Property: an exact-dir pattern always matches that item.
func TestGlobFilter_Property_ExactPatternMatchesItem(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		item := genItem().Draw(t, "item")
		result, err := GlobFilter([]api.ConfigurationItem{item}, []string{item.Dir})
		require.NoError(t, err)
		assert.Contains(t, result, item)
	})
}

// Property: every item in the result was present in the input.
func TestGlobFilter_Property_ResultIsSubsetOfInput(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		items := genItems().Draw(t, "items")
		pattern := genDirPath().Draw(t, "pattern")
		result, err := GlobFilter(items, []string{pattern})
		require.NoError(t, err)
		for _, r := range result {
			assert.Contains(t, items, r)
		}
	})
}

// Property: an item matching multiple patterns appears exactly once (no duplicates).
func TestGlobFilter_Property_NoDuplicatesOnMultiplePatternMatch(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		item := genItem().Draw(t, "item")
		// item.Dir and "**" both match the item.
		result, err := GlobFilter([]api.ConfigurationItem{item}, []string{item.Dir, "**"})
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

// Property: adding more patterns never reduces the result size (OR semantics).
func TestGlobFilter_Property_ORMonotonicity(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		items := genItems().Draw(t, "items")
		p1 := genDirPath().Draw(t, "p1")
		p2 := genDirPath().Draw(t, "p2")

		r1, err := GlobFilter(items, []string{p1})
		require.NoError(t, err)
		r12, err := GlobFilter(items, []string{p1, p2})
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(r12), len(r1))
	})
}

// Property: "*" does not match a multi-segment path.
func TestGlobFilter_Property_SingleStarDoesNotCrossSegments(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a dir with at least 2 segments.
		seg1 := genSegment().Draw(t, "seg1")
		seg2 := genSegment().Draw(t, "seg2")
		dir := seg1 + "/" + seg2
		item := api.ConfigurationItem{Name: "x", Dir: dir, Path: dir + "/pantalon.yaml"}

		result, err := GlobFilter([]api.ConfigurationItem{item}, []string{"*"})
		require.NoError(t, err)
		assert.NotContains(t, result, item)
	})
}

// Property: "prefix/**" matches any path whose first segment equals prefix.
func TestGlobFilter_Property_PrefixDoubleStarMatchesChildren(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		prefix := genSegment().Draw(t, "prefix")
		// Generate the remainder (at least one more segment).
		rest := genDirPath().Draw(t, "rest")
		dir := prefix + "/" + rest
		item := api.ConfigurationItem{Name: "x", Dir: dir, Path: dir + "/pantalon.yaml"}

		result, err := GlobFilter([]api.ConfigurationItem{item}, []string{prefix + "/**"})
		require.NoError(t, err)
		assert.Contains(t, result, item)
	})
}
