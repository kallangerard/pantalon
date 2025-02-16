package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/goccy/go-yaml"

	"github.com/kallangerard/pantalon/api"
	"github.com/kallangerard/pantalon/file"
)

func main() {
	outputFormat := flag.String("output-format", "json", "Output format (json)")
	flag.Parse()

	configurations, err := file.Search()
	if err != nil {
		log.Fatalf("Error listing configurations: %v", err)
	}

	items, err := api.MarshalItems(configurations)
	if err != nil {
		log.Fatalf("Error marshaling items: %v", err)
	}

	switch *outputFormat {
	case "json":
		outputJson(items)
	case "json-compact":
		outputJsonCompact(items)
	case "yaml":
		outputYaml(items)
	default:
		log.Fatalf("Unsupported output format: %s", *outputFormat)
	}
}

func outputJson(configurations []api.ConfigurationItem) {
	data, err := json.MarshalIndent(configurations, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling json: %v", err)
	}
	fmt.Println(string(data))
}

func outputJsonCompact(configurations []api.ConfigurationItem) {
	data, err := json.Marshal(configurations)
	if err != nil {
		log.Fatalf("Error marshaling json-compact: %v", err)
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
