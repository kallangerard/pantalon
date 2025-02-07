package file

import (
	"log"
	"os"
	"path/filepath"
)

func FindPantalonFiles(root string) ([]string, error) {
	var result []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
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
			// Skip this directory if it doesn't contain a pantalon.yaml file
			return nil
		}

		result = append(result, pantalonPath)
		return filepath.SkipDir
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
