package pantalon

import (
	"github.com/bmatcuk/doublestar/v4"
)

// GlobFilter returns items whose Dir matches any of the provided doublestar glob patterns.
func GlobFilter(items []ConfigurationItem, patterns []string) ([]ConfigurationItem, error) {
	if len(patterns) == 0 {
		return items, nil
	}

	filtered := make([]ConfigurationItem, 0)
	for _, item := range items {
		for _, pattern := range patterns {
			matched, err := doublestar.Match(pattern, item.Dir)
			if err != nil {
				return nil, err
			}
			if matched {
				filtered = append(filtered, item)
				break
			}
		}
	}
	return filtered, nil
}
