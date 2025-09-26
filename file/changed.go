package file

import (
	"path/filepath"
	"strings"

	"github.com/kallangerard/pantalon/api"
)

// ChangedFiles filters the list of unfiltered configuration files based on the list of changed directories.
// It includes configurations if:
// 1. The changed directory overlaps with the configuration directory
// 2. The changed directory overlaps with any of the configuration's dependencies
func ChangedFiles(allItems []api.ConfigurationItem, changedDirs []string) ([]api.ConfigurationItem, error) {
	filteredCfgs := make([]api.ConfigurationItem, 0)

	for _, cfg := range allItems {
		included := false
		
		for _, dir := range changedDirs {
			// Check if changed dir overlaps with config dir
			if dir == "." || strings.HasPrefix(dir, cfg.Dir) || strings.HasPrefix(cfg.Dir, dir) {
				filteredCfgs = append(filteredCfgs, cfg)
				included = true
				break
			}
			
			// Check if changed dir overlaps with any dependency
			for _, dependency := range cfg.Dependencies {
				var resolvedDependency string
				
				// If dependency is relative (starts with ./ or ../), resolve relative to config directory
				// Otherwise, treat as absolute path from repository root
				if strings.HasPrefix(dependency, "./") || strings.HasPrefix(dependency, "../") {
					resolvedDependency = filepath.Join(cfg.Dir, dependency)
					resolvedDependency = filepath.Clean(resolvedDependency)
				} else {
					resolvedDependency = dependency
				}
				
				if strings.HasPrefix(dir, resolvedDependency) || strings.HasPrefix(resolvedDependency, dir) {
					filteredCfgs = append(filteredCfgs, cfg)
					included = true
					break
				}
			}
			
			if included {
				break
			}
		}
	}

	return filteredCfgs, nil
}
