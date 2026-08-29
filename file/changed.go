package file

import (
	"strings"

	"github.com/kallangerard/pantalon/api"
)

// ChangedFiles filters the list of unfiltered configuration files based on the list of changed directories.
func ChangedFiles(allItems []api.ConfigurationItem, changedDirs []string) ([]api.ConfigurationItem, error) {
	if len(changedDirs) == 0 {
		return make([]api.ConfigurationItem, 0), nil
	}

	// A change at the repository root matches every configuration, so there is
	// nothing to compare per item.
	for _, dir := range changedDirs {
		if dir == "." {
			filtered := make([]api.ConfigurationItem, len(allItems))
			copy(filtered, allItems)
			return filtered, nil
		}
	}

	// Deliberately not preallocated to len(allItems): a changed-dirs filter
	// usually keeps a small fraction of the configurations, and reserving the
	// full slice costs more than append's growth.
	filteredCfgs := make([]api.ConfigurationItem, 0)

	for _, cfg := range allItems {
		for _, dir := range changedDirs {
			if strings.HasPrefix(dir, cfg.Dir) || strings.HasPrefix(cfg.Dir, dir) {
				filteredCfgs = append(filteredCfgs, cfg)
				break
			}
		}
	}

	return filteredCfgs, nil
}
