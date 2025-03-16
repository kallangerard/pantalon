package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/goccy/go-yaml"

	"github.com/kallangerard/pantalon/api"
	"github.com/kallangerard/pantalon/file"
)

func main() {
	outputFormat := flag.String("output-format", "json", "Output format (json)")
	changedDirsJson := flag.String("changed-dirs", "", `[".", "foo", "foo/bar"]`)
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
