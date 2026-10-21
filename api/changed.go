package api

import (
	"encoding/json"
	"errors"
	"path"
	"path/filepath"
	"strings"
)

var ErrChangedDirOutsideRepo = errors.New("changed dir must be relative to the repository root")

// UnmarshalChangedFileJson parses a JSON array of changed directories and
// normalises each to a clean, slash-separated path relative to the repository
// root, so "a/dev/", "./a/dev" and "a/dev" are equivalent and "" becomes ".".
func UnmarshalChangedFileJson(j []byte) ([]string, error) {
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
