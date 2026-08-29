package file

import (
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/kallangerard/pantalon/api"
)

// globMeta are the characters that make a doublestar pattern something other
// than a literal path.
const globMeta = `*?[{\`

// GlobFilter returns items whose Dir matches any of the provided doublestar glob patterns.
func GlobFilter(items []api.ConfigurationItem, patterns []string) ([]api.ConfigurationItem, error) {
	if len(patterns) == 0 {
		return items, nil
	}

	// doublestar has no compiled-pattern type and rescans the pattern on every
	// call, so classify each pattern once rather than once per configuration.
	// A pattern without metacharacters is just a string comparison.
	literal := make([]bool, len(patterns))
	for i, pattern := range patterns {
		literal[i] = !strings.ContainsAny(pattern, globMeta)
	}

	filtered := make([]api.ConfigurationItem, 0)
	for _, item := range items {
		for i, pattern := range patterns {
			var (
				matched bool
				err     error
			)
			if literal[i] {
				matched = pattern == item.Dir
			} else {
				matched, err = doublestar.Match(pattern, item.Dir)
				if err != nil {
					return nil, err
				}
			}
			if matched {
				filtered = append(filtered, item)
				break
			}
		}
	}
	return filtered, nil
}
