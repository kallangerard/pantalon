package pantalon

import "fmt"

// FilterOptions narrows items; zero values select everything.
type FilterOptions struct {
	// ChangedDirs is a JSON array of changed directories, as given to --changed-dirs.
	ChangedDirs string
	PathGlobs   []string
	// Selector is a label selector such as "env=prod,team in (a,b)".
	Selector string
}

// Filter applies changed-dirs, path globs and selector in that order.
func Filter(items []ConfigurationItem, opts FilterOptions) ([]ConfigurationItem, error) {
	if opts.ChangedDirs != "" {
		changedDirs, err := ParseChangedDirs([]byte(opts.ChangedDirs))
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling changed dirs: %w", err)
		}
		items, err = ChangedFiles(items, changedDirs)
		if err != nil {
			return nil, fmt.Errorf("error filtering changed files: %w", err)
		}
	}

	if len(opts.PathGlobs) > 0 {
		var err error
		items, err = GlobFilter(items, opts.PathGlobs)
		if err != nil {
			return nil, fmt.Errorf("error filtering by path glob: %w", err)
		}
	}

	if opts.Selector != "" {
		sel, err := ParseSelector(opts.Selector)
		if err != nil {
			return nil, err
		}
		items = SelectorFilter(items, sel)
	}

	return items, nil
}
