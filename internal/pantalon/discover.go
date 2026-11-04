package pantalon

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"
)

// Discover walks fsys for pantalon.yaml files, stopping descent at the first
// match in each branch. Warnings are returned even when err is non-nil.
func Discover(fsys fs.FS) ([]TerraformConfiguration, []Warning, error) {
	paths, err := findFiles(fsys)
	if err != nil {
		return nil, nil, err
	}

	var result []TerraformConfiguration
	var warnings []Warning
	seen := make(map[string]string)
	for _, p := range paths {
		tfCfg, ws, err := readFile(fsys, p)
		warnings = append(warnings, ws...)
		if err != nil {
			return nil, warnings, fmt.Errorf("%s: %w", p, err)
		}
		// Names identify CI jobs, so they must be unique across the repository
		// regardless of any later filtering. Compare case-insensitively so
		// that relaxing name validation can't let "Dev" and "dev" coexist.
		key := strings.ToLower(tfCfg.Metadata.Name)
		if other, ok := seen[key]; ok {
			return nil, warnings, fmt.Errorf("duplicate metadata.name %q in %s and %s", tfCfg.Metadata.Name, other, p)
		}
		seen[key] = p
		tfCfg.Path = p
		result = append(result, tfCfg)
	}
	return result, warnings, nil
}

func findFiles(fsys fs.FS) ([]string, error) {
	var result []string
	err := fs.WalkDir(fsys, ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			return nil
		}

		pantalonPath := path.Join(p, "pantalon.yaml")
		if _, err := fs.Stat(fsys, pantalonPath); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil // no marker here, keep descending
			}
			return err
		}

		result = append(result, pantalonPath)
		return fs.SkipDir
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func readFile(fsys fs.FS, p string) (TerraformConfiguration, []Warning, error) {
	data, err := fs.ReadFile(fsys, p)
	if err != nil {
		return TerraformConfiguration{}, nil, err
	}
	cfg, warnings, err := Parse(data)
	for i := range warnings {
		warnings[i].Path = p
	}
	return cfg, warnings, err
}
