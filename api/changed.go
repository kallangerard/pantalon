package api

import (
	"encoding/json"
	"errors"
	"path/filepath"
)

var ErrChangedFileIsNotDirectory = errors.New("changed file is not a directory")

func UnmarshalChangedFileJson(j []byte) ([]string, error) {
	var actualPaths []string
	err := json.Unmarshal(j, &actualPaths)
	if err != nil {
		return nil, err
	}

	var errs []error

	for _, path := range actualPaths {
		// If path is not a dir, return an error
		if !isDir(path) {
			errs = append(errs, errors.New(path+": "+ErrChangedFileIsNotDirectory.Error()))
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return actualPaths, nil
}

func isDir(path string) bool {
	return path == "." || filepath.Ext(path) == ""
}
