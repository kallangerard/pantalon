package api

import (
	"errors"

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
	Path       string
}

type ConfigurationItem struct {
	Name    string            `yaml:"name"`
	Path    string            `yaml:"path"`
	Dir     string            `yaml:"dir"`
	Context map[string]string `yaml:"context"`
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

	err = c.validateTerraform(cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (c config) validateTerraform(cfg TerraformConfiguration) error {
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

func MarshalItems(cfgs []TerraformConfiguration) ([]ConfigurationItem, error) {

	items := make([]ConfigurationItem, 0, len(cfgs))

	for _, cfg := range cfgs {
		item := ConfigurationItem{
			Name:    cfg.Metadata.Name,
			Context: cfg.Context,
			Path:    cfg.Path,
			Dir:     dir(cfg.Path),
		}
		items = append(items, item)
	}

	return items, nil
}

// dir returns the directory of a slash separated path. It is equivalent to
// path.Dir for the already-clean paths the walk produces, without re-running
// path.Clean over every path.
func dir(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' {
			if i == 0 {
				return "/"
			}
			return p[:i]
		}
	}
	// No separator: path.Dir cleans the empty prefix to ".".
	return "."
}

// Must comply with RFC 1123 subdomain labels
//
// As described in https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names
// It is equivalent to the regular expression
// `^[a-z0-9]([a-z0-9-]*[a-z0-9])?$` with a 253 character limit, hand rolled to
// avoid compiling a regular expression on every configuration read.
func isValidSubdomainLabel(s string) bool {
	if len(s) == 0 || len(s) > 253 {
		return false
	}

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z', c >= '0' && c <= '9':
		case c == '-':
			// A hyphen may not lead or trail the label.
			if i == 0 || i == len(s)-1 {
				return false
			}
		default:
			return false
		}
	}

	return true
}
