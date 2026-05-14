package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/goccy/go-yaml"

	"github.com/kallangerard/pantalon/api"
	"github.com/kallangerard/pantalon/file"
)

// pathGlobs implements flag.Value for a repeatable --path-glob flag.
type pathGlobs []string

func (p *pathGlobs) String() string { return fmt.Sprintf("%v", *p) }
func (p *pathGlobs) Set(v string) error {
	*p = append(*p, v)
	return nil
}

func main() {
	outputFormat := flag.String("output-format", "json", "Output format (json)")
	changedDirsJson := flag.String("changed-dirs", "", `[".", "foo", "foo/bar"]`)
	var globs pathGlobs
	flag.Var(&globs, "path-glob", "Doublestar glob pattern to filter configurations by directory path (repeatable, OR logic)")
	flag.Parse()

	configurations, err := file.Search()
	if err != nil {
		log.Fatalf("Error listing configurations: %v", err)
	}

	unfilteredItems, err := api.MarshalItems(configurations)
	if err != nil {
		log.Fatalf("Error marshaling items: %v", err)
	}

	var items = []api.ConfigurationItem{}

	if changedDirsJson != nil {
		changedDirs, err := api.UnmarshalChangedFileJson([]byte(*changedDirsJson))
		if err != nil {
			log.Fatalf("Error unmarshaling changed dirs: %v", err)
		}
		items, err = file.ChangedFiles(unfilteredItems, changedDirs)
		if err != nil {
			log.Fatalf("Error filtering changed files: %v", err)
		}
	} else {
		items = unfilteredItems
	}

	if len(globs) > 0 {
		items, err = file.GlobFilter(items, globs)
		if err != nil {
			log.Fatalf("Error filtering by path glob: %v", err)
		}
	}

	switch *outputFormat {
	case "json":
		outputJson(items)
	case "yaml":
		outputYaml(items)
	default:
		log.Fatalf("Unsupported output format: %s", *outputFormat)
	}
}

func outputJson(configurations []api.ConfigurationItem) {
	data, err := yaml.MarshalWithOptions(configurations,
		yaml.JSON(),
	)
	if err != nil {
		log.Fatalf("Error marshaling json: %v", err)
	}
	fmt.Println(string(data))
}

func outputYaml(configurations []api.ConfigurationItem) {
	data, err := yaml.Marshal(configurations)
	if err != nil {
		log.Fatalf("Error marshaling yaml: %v", err)
	}
	fmt.Println(string(data))
}
