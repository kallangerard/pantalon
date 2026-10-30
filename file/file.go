package file

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kallangerard/pantalon/api"
)

func Search() ([]api.TerraformConfiguration, error) {
	paths, err := findFiles()
	if err != nil {
		return nil, err
	}

	var result []api.TerraformConfiguration
	seen := make(map[string]string)
	for _, path := range paths {
		tfCfg, err := readFile(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		// Names identify CI jobs, so they must be unique across the repository
		// regardless of any later filtering. Compare case-insensitively so
		// that relaxing name validation can't let "Dev" and "dev" coexist.
		key := strings.ToLower(tfCfg.Metadata.Name)
		if other, ok := seen[key]; ok {
			return nil, fmt.Errorf("duplicate metadata.name %q in %s and %s", tfCfg.Metadata.Name, other, path)
		}
		seen[key] = path
		tfCfg.Path = path
		result = append(result, tfCfg)
	}
	return result, nil
}

func findFiles() ([]string, error) {
	var result []string
	err := filepath.WalkDir(".", func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}

		if !entry.IsDir() {
			return nil
		}

		pantalonPath := filepath.Join(path, "pantalon.yaml")

		_, err = os.Stat(pantalonPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// Skip this directory if it doesn't contain a pantalon.yaml file
				return nil
			}
			return err
		}

		result = append(result, pantalonPath)
		return filepath.SkipDir
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func readFile(path string) (api.TerraformConfiguration, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return api.TerraformConfiguration{}, err
	}

	warnings := api.Warnings(file)

	cfg := api.New()
	tfCfg, err := cfg.Unmarshal(file)
	if err != nil {
		// A misspelt required field surfaces as a validation error, so name
		// the likely cause.
		if len(warnings) > 0 {
			return api.TerraformConfiguration{}, fmt.Errorf("%w (%s)", err, joinWarnings(warnings))
		}
		return api.TerraformConfiguration{}, err
	}

	for _, w := range warnings {
		if os.Getenv("GITHUB_ACTIONS") == "true" {
			// Workflow commands must start the line, so bypass the log prefix.
			fmt.Fprintf(log.Writer(), "::warning file=%s,line=%d::%s\n", path, w.Line, w.Message)
		} else {
			log.Printf("warning: %s:%d: %s", path, w.Line, w.Message)
		}
	}
	return tfCfg, nil
}

func joinWarnings(warnings []api.Warning) string {
	msgs := make([]string, len(warnings))
	for i, w := range warnings {
		msgs[i] = fmt.Sprintf("line %d: %s", w.Line, w.Message)
	}
	return strings.Join(msgs, "; ")
}
