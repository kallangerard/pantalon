package api

import (
	"errors"
	"fmt"
	"maps"
	"path"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
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
	ApiVersion   string            `yaml:"apiVersion"`
	Kind         string            `yaml:"kind"`
	Metadata     Metadata          `yaml:"metadata"`
	Dependencies Dependencies      `yaml:"dependencies,omitempty"`
	Context      map[string]string `yaml:"context,omitempty"`
	Path         string            `yaml:"-"`
}

type ConfigurationItem struct {
	Name         string            `yaml:"name"`
	Path         string            `yaml:"path"`
	Dir          string            `yaml:"dir"`
	Labels       map[string]string `yaml:"labels"`
	Dependencies *Dependencies     `yaml:"dependencies,omitempty"`
	Context      map[string]string `yaml:"context"`
}

// Dependencies lists what a configuration depends on outside its own
// directory, so that changes there also select it.
type Dependencies struct {
	// Paths are doublestar glob patterns relative to the repository root, in
	// the style of GitHub Actions on.<event>.paths. A pattern matches a
	// directory and everything below it. A leading "!" excludes directories
	// matched by earlier patterns; the last matching pattern wins.
	Paths []string `yaml:"paths,omitempty"`
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

// Warning is a problem in a pantalon.yaml that doesn't stop it being read.
type Warning struct {
	Line    int
	Message string
}

// Warnings reports fields this version of pantalon doesn't recognise.
// Unmarshal ignores them rather than rejecting them so that older binaries
// can read configs using newer fields.
func Warnings(yamlDoc []byte) []Warning {
	f, err := parser.ParseBytes(yamlDoc, 0)
	if err != nil || len(f.Docs) == 0 || f.Docs[0].Body == nil {
		return nil // Unmarshal reports syntax errors
	}
	return unknownFields(f.Docs[0].Body, reflect.TypeOf(TerraformConfiguration{}), "")
}

// unknownFields walks node against the yaml tags of struct type t.
func unknownFields(node ast.Node, t reflect.Type, prefix string) []Warning {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil // maps such as context accept any key
	}

	var values []*ast.MappingValueNode
	switch n := node.(type) {
	case *ast.MappingNode:
		values = n.Values
	case *ast.MappingValueNode:
		values = []*ast.MappingValueNode{n}
	case *ast.AnchorNode:
		return unknownFields(n.Value, t, prefix)
	case *ast.TagNode:
		return unknownFields(n.Value, t, prefix)
	}

	fields := make(map[string]reflect.Type)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name, _, _ := strings.Cut(f.Tag.Get("yaml"), ",")
		if name != "" && name != "-" {
			fields[name] = f.Type
		}
	}

	var warnings []Warning
	for _, v := range values {
		tk := v.Key.GetToken()
		key := tk.Value
		if key == "<<" {
			continue // YAML merge key
		}
		ft, ok := fields[key]
		if !ok {
			warnings = append(warnings, Warning{tk.Position.Line, fmt.Sprintf("unknown field %q", prefix+key)})
			continue
		}
		warnings = append(warnings, unknownFields(v.Value, ft, prefix+key+".")...)
	}
	return warnings
}

func (c config) validateTerraform(cfg TerraformConfiguration) error {
	if cfg.ApiVersion != PantalonVersion {
		return errors.New("invalid version")
	}

	if cfg.Kind != TerraformKind {
		return errors.New("invalid kind")
	}

	if !isValidSubdomainLabel(cfg.Metadata.Name) {
		return fmt.Errorf("invalid metadata.name %q", cfg.Metadata.Name)
	}

	keys := slices.Sorted(maps.Keys(cfg.Metadata.Labels))
	for _, k := range keys {
		v := cfg.Metadata.Labels[k]
		if !isValidLabelKey(k) {
			return fmt.Errorf("invalid metadata.labels key %q", k)
		}
		if !isValidLabelValue(v) {
			return fmt.Errorf("invalid metadata.labels value %q for key %q", v, k)
		}
	}
	for i, p := range cfg.Dependencies.Paths {
		if err := validateDependencyPath(p, i == 0); err != nil {
			return fmt.Errorf("invalid dependencies.paths[%d] %q: %w", i, p, err)
		}
	}

	return nil
}

func validateDependencyPath(p string, first bool) error {
	pattern, negated := strings.CutPrefix(p, "!")
	if negated && first {
		return errors.New("first pattern cannot be negated, as there is nothing to exclude")
	}
	if pattern == "" {
		return errors.New("pattern is empty")
	}
	if path.IsAbs(pattern) || slices.Contains(strings.Split(pattern, "/"), "..") {
		return errors.New("pattern must be relative to the repository root")
	}
	if !doublestar.ValidatePattern(pattern) {
		return errors.New("pattern is not a valid glob")
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
		if len(cfg.Dependencies.Paths) > 0 {
			item.Dependencies = &cfg.Dependencies
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
