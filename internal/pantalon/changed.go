package pantalon

import (
	"encoding/json"
	"errors"
	"path"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

var ErrChangedDirOutsideRepo = errors.New("changed dir must be relative to the repository root")

// ParseChangedDirs parses a JSON array of changed directories and
// normalises each to a clean, slash-separated path relative to the repository
// root, so "a/dev/", "./a/dev" and "a/dev" are equivalent and "" becomes ".".
func ParseChangedDirs(j []byte) ([]string, error) {
	var actualPaths []string
	err := json.Unmarshal(j, &actualPaths)
	if err != nil {
		return nil, err
	}

	var errs []error

	for i, p := range actualPaths {
		p = path.Clean(filepath.ToSlash(p))
		if path.IsAbs(p) || p == ".." || strings.HasPrefix(p, "../") {
			errs = append(errs, errors.New(actualPaths[i]+": "+ErrChangedDirOutsideRepo.Error()))
			continue
		}
		actualPaths[i] = p
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return actualPaths, nil
}

// ChangedFiles filters the list of unfiltered configuration files based on the list of changed directories.
//
// A configuration is selected when a changed directory is within its own
// directory (or contains it), or matches its dependencies.paths.
func ChangedFiles(allItems []ConfigurationItem, changedDirs []string) ([]ConfigurationItem, error) {
	filteredCfgs := make([]ConfigurationItem, 0)

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
