package file

import (
	"strings"

	"github.com/kallangerard/pantalon/api"
)

// changedFiles filters the list of unfiltered configuration files based on the list of changed directories.
func changedFiles(allItems []api.ConfigurationItem, changedDirs []string) ([]string, error) {
	filteredCfgs := make([]string, 0)

	for _, cfg := range allItems {
		for _, dir := range changedDirs {
			if dir == "." || strings.HasPrefix(dir, cfg.Dir) || strings.HasPrefix(cfg.Dir, dir) {
				filteredCfgs = append(filteredCfgs, cfg.Dir)
			}
		}
	}

	return filteredCfgs, nil
}
