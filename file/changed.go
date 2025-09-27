package file

import (
	"path/filepath"
	"strings"

	"github.com/kallangerard/pantalon/api"
)

// ChangedFiles filters the list of unfiltered configuration files based on the list of changed directories.
// It includes configurations if:
// 1. The changed directory overlaps with the configuration directory
// 2. The changed directory overlaps with any of the configuration's dependencies (absolute paths from repo root)
// Dependencies support glob patterns (e.g., /terraform/modules/**)
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
				// Dependencies must be absolute paths from repository root (starting with /)
				// Remove leading slash for comparison with changed directories
				dependencyPath := strings.TrimPrefix(dependency, "/")
				
				// Check for glob pattern support
				if strings.Contains(dependency, "*") {
					// Use filepath.Match for glob patterns
					if matched, err := filepath.Match(dependencyPath, dir); err == nil && matched {
						filteredCfgs = append(filteredCfgs, cfg)
						included = true
						break
					}
					// Also check if the changed dir is a parent of the glob pattern
					// For patterns like /terraform/modules/**, a change in /terraform/modules should match
					globBase := strings.Split(dependencyPath, "*")[0]
					globBase = strings.TrimSuffix(globBase, "/")
					if strings.HasPrefix(dir, globBase) || strings.HasPrefix(globBase, dir) {
						filteredCfgs = append(filteredCfgs, cfg)
						included = true
						break
					}
				} else {
					// Exact path matching for non-glob dependencies
					if strings.HasPrefix(dir, dependencyPath) || strings.HasPrefix(dependencyPath, dir) {
						filteredCfgs = append(filteredCfgs, cfg)
						included = true
						break
					}
				}
			}
			
			if included {
				break
			}
		}
	}

	return filteredCfgs, nil
}
