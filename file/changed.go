package file

import (
	"strings"

	"github.com/kallangerard/pantalon/api"
)

// ChangedFiles filters the list of unfiltered configuration files based on the list of changed directories.
func ChangedFiles(allItems []api.ConfigurationItem, changedDirs []string) ([]api.ConfigurationItem, error) {
	filteredCfgs := make([]api.ConfigurationItem, 0)

	for _, cfg := range allItems {
		for _, dir := range changedDirs {
			if isWithin(dir, cfg.Dir) || isWithin(cfg.Dir, dir) {
				filteredCfgs = append(filteredCfgs, cfg)
				break
			}
		}
	}

	return filteredCfgs, nil
}

// isWithin reports whether dir equals parent or is a descendant of it,
// comparing whole path segments so "a/b-2" is not treated as being within
// "a/b". Both paths must be clean and slash-separated.
func isWithin(dir, parent string) bool {
	return parent == "." || dir == parent || strings.HasPrefix(dir, parent+"/")
}
