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
			if dir == "." || strings.HasPrefix(dir, cfg.Dir) || strings.HasPrefix(cfg.Dir, dir) {
				filteredCfgs = append(filteredCfgs, cfg)
			}
		}
	}

	return filteredCfgs, nil
}
