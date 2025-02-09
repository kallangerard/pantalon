package api

import (
	"errors"
	"regexp"

	"github.com/goccy/go-yaml"
)

const (
	PantalonVersion = "pantalon.kallan.dev/v1alpha1"
	TerraformKind   = "TerraformConfiguration"
)

type PantalonConfig interface {
	New() config
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

func New() config {
	return config{}
}

func (c config) Unmarshal(yamlDoc []byte) (TerraformConfiguration, error) {
	cfg := TerraformConfiguration{}

	err := yaml.Unmarshal(yamlDoc, &cfg)
	if err != nil {
		return cfg, err
	}

	err = c.ValidateTerraform(cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (c config) ValidateTerraform(cfg TerraformConfiguration) error {
	if cfg.ApiVersion != PantalonVersion {
		return errors.New("invalid version")
	}

	if cfg.Kind != TerraformKind {
		return errors.New("invalid kind")
	}

	if !isValidSubdomainLabel(cfg.Metadata.Name) {
		return errors.New("invalid metadata.name")
	}
	return nil
}

// Must comply with RFC 1123 subdomain labels
//
// As described in https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names
func isValidSubdomainLabel(s string) bool {
	reg := regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

	if !reg.MatchString(s) {
		return false
	}

	if len(s) > 253 {
		return false
	}

	return true
}
