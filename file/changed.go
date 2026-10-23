package file

import (
	"path"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/kallangerard/pantalon/api"
)

// ChangedFiles filters the list of unfiltered configuration files based on the list of changed directories.
//
// A configuration is selected when a changed directory is within its own
// directory (or contains it), or matches its dependencies.paths.
func ChangedFiles(allItems []api.ConfigurationItem, changedDirs []string) ([]api.ConfigurationItem, error) {
	filteredCfgs := make([]api.ConfigurationItem, 0)

	for _, cfg := range allItems {
		for _, dir := range changedDirs {
			matched := isWithin(dir, cfg.Dir) || isWithin(cfg.Dir, dir)
			if !matched && cfg.Dependencies != nil {
				var err error
				matched, err = matchesDependencyPaths(dir, cfg.Dependencies.Paths)
				if err != nil {
					return nil, err
				}
			}
			if matched {
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

// matchesDependencyPaths reports whether dir matches the patterns, following
// GitHub Actions paths filters: patterns are evaluated in order, a "!" prefix
// excludes, and the last matching pattern wins. Each pattern matches a
// directory and everything below it.
func matchesDependencyPaths(dir string, patterns []string) (bool, error) {
	matched := false
	for _, p := range patterns {
		pattern, negated := strings.CutPrefix(p, "!")
		pattern = path.Clean(pattern)

		ok, err := matchesDir(pattern, dir)
		if err != nil {
			return false, err
		}
		if ok {
			matched = !negated
		}
	}
	return matched, nil
}

// matchesDir reports whether dir matches pattern or is below a directory that does.
func matchesDir(pattern, dir string) (bool, error) {
	if pattern == "." {
		return true, nil
	}
	if ok, err := doublestar.Match(pattern, dir); ok || err != nil {
		return ok, err
	}
	return doublestar.Match(pattern+"/**", dir)
}
