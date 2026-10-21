package api

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"

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
	Labels  map[string]string `yaml:"labels"`
	Context map[string]string `yaml:"context"`
}

type Metadata struct {
	Name   string            `yaml:"name"`
	Labels map[string]string `yaml:"labels,omitempty"`
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

	for k, v := range cfg.Metadata.Labels {
		if !isValidLabelKey(k) {
			return fmt.Errorf("invalid metadata.labels key %q", k)
		}
		if !isValidLabelValue(v) {
			return fmt.Errorf("invalid metadata.labels value %q for key %q", v, k)
		}
	}
	return nil
}

func MarshalItems(cfgs []TerraformConfiguration) ([]ConfigurationItem, error) {

	items := make([]ConfigurationItem, 0)

	for _, cfg := range cfgs {
		item := ConfigurationItem{
			Name:    cfg.Metadata.Name,
			Labels:  cfg.Metadata.Labels,
			Context: cfg.Context,
			Path:    cfg.Path,
			Dir:     path.Dir(cfg.Path),
		}
		items = append(items, item)
	}

	return items, nil
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

var labelNameRegexp = regexp.MustCompile(`^[A-Za-z0-9]([-A-Za-z0-9_.]*[A-Za-z0-9])?$`)

// isValidLabelKey follows Kubernetes label keys: an optional DNS subdomain
// prefix and "/", then a name of at most 63 characters.
//
// As described in https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set
func isValidLabelKey(s string) bool {
	name := s
	if prefix, n, found := strings.Cut(s, "/"); found {
		if len(prefix) > 253 || !isValidDNSSubdomain(prefix) {
			return false
		}
		name = n
	}

	return len(name) <= 63 && labelNameRegexp.MatchString(name)
}

// isValidLabelValue follows Kubernetes label values: empty, or at most 63
// characters matching the label name syntax.
func isValidLabelValue(s string) bool {
	return s == "" || (len(s) <= 63 && labelNameRegexp.MatchString(s))
}

var dnsSubdomainRegexp = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)

func isValidDNSSubdomain(s string) bool {
	return dnsSubdomainRegexp.MatchString(s)
}
