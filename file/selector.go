package file

import "github.com/kallangerard/pantalon/api"

// SelectorFilter returns items whose labels satisfy the selector.
func SelectorFilter(items []api.ConfigurationItem, sel api.Selector) []api.ConfigurationItem {
	filtered := make([]api.ConfigurationItem, 0)
	for _, item := range items {
		if sel.Matches(item.Labels) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
