package file

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/kallangerard/pantalon/api"
)

func Search() ([]api.TerraformConfiguration, error) {
	paths, err := findFiles()
	if err != nil {
		return nil, err
	}

	var result []api.TerraformConfiguration
	for _, path := range paths {
		tfCfg, err := readFile(path)
		tfCfg.Path = path
		if err != nil {
			return nil, err
		}
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

	cfg := api.New()
	tfCfg, err := cfg.Unmarshal(file)
	if err != nil {
		return api.TerraformConfiguration{}, err
	}
	return tfCfg, nil
}
