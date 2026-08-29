package file

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/kallangerard/pantalon/api"
)

// MarkerFile is the file that marks a directory as a Terraform root module.
const MarkerFile = "pantalon.yaml"

// skippedDirs are never descended into. `.git` holds thousands of directories
// in a normal checkout and `.terraform` holds provider and module caches, which
// may contain copies of a marker file that are not root modules of this
// repository.
var skippedDirs = map[string]struct{}{
	".git":       {},
	".terraform": {},
}

// concurrency bounds the goroutines used to walk directories and to read the
// marker files that are found. Directory traversal is dominated by syscall
// latency rather than CPU, so oversubscribing the available cores pays off.
var concurrency = max(4, runtime.NumCPU()*4)

func Search() ([]api.TerraformConfiguration, error) {
	paths, err := findFiles()
	if err != nil {
		return nil, err
	}

	result := make([]api.TerraformConfiguration, len(paths))
	errs := make([]error, len(paths))

	forEach(len(paths), func(i int) {
		tfCfg, err := readFile(paths[i])
		if err != nil {
			errs[i] = err
			return
		}
		tfCfg.Path = paths[i]
		result[i] = tfCfg
	})

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

// forEach runs fn for each index in [0, n) across a bounded pool of workers.
func forEach(n int, fn func(i int)) {
	if n == 0 {
		return
	}
	if n == 1 {
		fn(0)
		return
	}

	workers := min(n, concurrency)
	var (
		wg   sync.WaitGroup
		next = make(chan int)
	)
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for i := range next {
				fn(i)
			}
		}()
	}
	for i := range n {
		next <- i
	}
	close(next)
	wg.Wait()
}

// findFiles returns the path of every marker file in the tree rooted at the
// working directory, in the same depth-first lexical order filepath.WalkDir
// produces. Descent stops at the first marker file found on a branch: the
// directory is a root module and its children belong to it.
func findFiles() ([]string, error) {
	w := walker{sem: make(chan struct{}, concurrency)}
	return w.walk(".")
}

type walker struct {
	// sem bounds the number of directories walked concurrently. A directory
	// that cannot claim a slot is walked by the calling goroutine, so the walk
	// never blocks waiting for one.
	sem chan struct{}
}

func (w walker) walk(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// os.ReadDir already reported whether each entry is a directory, so the
	// marker file is found without an extra stat per directory.
	subdirs := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() {
			if name == MarkerFile {
				return []string{join(dir, name)}, nil
			}
			continue
		}
		if _, skip := skippedDirs[name]; skip {
			continue
		}
		subdirs = append(subdirs, join(dir, name))
	}

	if len(subdirs) == 0 {
		return nil, nil
	}

	// Results are collected per subdirectory and concatenated in order, so a
	// concurrent walk still returns paths in lexical order.
	found := make([][]string, len(subdirs))
	errs := make([]error, len(subdirs))
	var wg sync.WaitGroup

	for i, subdir := range subdirs {
		select {
		case w.sem <- struct{}{}:
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() { <-w.sem }()
				found[i], errs[i] = w.walk(subdir)
			}()
		default:
			found[i], errs[i] = w.walk(subdir)
		}
	}
	wg.Wait()

	total := 0
	for i, paths := range found {
		if errs[i] != nil {
			return nil, errs[i]
		}
		total += len(paths)
	}
	if total == 0 {
		return nil, nil
	}

	result := make([]string, 0, total)
	for _, paths := range found {
		result = append(result, paths...)
	}
	return result, nil
}

// join appends name to an already-clean directory path, matching filepath.Join
// without the cost of cleaning the result on every entry.
func join(dir, name string) string {
	if dir == "." {
		return name
	}
	return dir + string(filepath.Separator) + name
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
