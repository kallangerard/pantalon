package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-yaml"

	"github.com/kallangerard/pantalon/api"
	"github.com/kallangerard/pantalon/file"
)

func main() {
	outputFormat := flag.String("output-format", "json", "Output format (json)")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatalf("Usage: %s <directory>", os.Args[0])
	}

	root := flag.Arg(0)
	configurations, err := file.Search(root)
	if err != nil {
		log.Fatalf("Error listing configurations: %v", err)
	}

	switch *outputFormat {
	case "json":
		outputJson(configurations)
	case "json-compact":
		outputJsonCompact(configurations)
	case "yaml":
		outputYaml(configurations)
	default:
		log.Fatalf("Unsupported output format: %s", *outputFormat)
	}
}

func outputJson(configurations []api.TerraformConfiguration) {
	data, err := json.MarshalIndent(configurations, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling json: %v", err)
	}
	fmt.Println(string(data))
}

func outputJsonCompact(configurations []api.TerraformConfiguration) {
	data, err := json.Marshal(configurations)
	if err != nil {
		log.Fatalf("Error marshaling json-compact: %v", err)
	}
	fmt.Println(string(data))
}

func outputYaml(configurations []api.TerraformConfiguration) {
	data, err := yaml.Marshal(configurations)
	if err != nil {
		log.Fatalf("Error marshaling yaml: %v", err)
	}
	fmt.Println(string(data))
}
