package api

import "github.com/goccy/go-yaml"

const (
	PantalonVersion = "pantalon.kallan.dev/v1alpha1"
	TerraformKind   = "TerraformConfiguration"
)

type PantalonConfig interface {
	Unmarshal([]byte) (TerraformConfiguration, error)
}

type config struct {
}

type TerraformConfiguration struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   Metadata          `yaml:"metadata"`
	Context    map[string]string `yaml:"context,omitempty"`
}

type Metadata struct {
	Name string `yaml:"name"`
}

func (c config) Unmarshal(yamlDoc []byte) (TerraformConfiguration, error) {
	cfg := TerraformConfiguration{}

	err := yaml.Unmarshal(yamlDoc, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
