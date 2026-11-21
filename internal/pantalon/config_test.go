package pantalon

import (
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	minimalYamlDoc = `
---
apiVersion: pantalon.kallan.dev/v1alpha1
kind: TerraformConfiguration
metadata:
  name: hello-world
`
)

func TestUnmarshalTerraformConfiguration_Minimal(t *testing.T) {
	v, _, err := Parse([]byte(minimalYamlDoc))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello-world", v.Metadata.Name)
}

func TestUnmarshalTerraformConfiguration_ApiVersionKindMatch(t *testing.T) {
	v, _, err := Parse([]byte(minimalYamlDoc))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, TerraformKind, v.Kind)
	assert.Equal(t, PantalonVersion, v.APIVersion)
}

func TestUnmarshalTerraformConfiguration_WithContext(t *testing.T) {
	yamlDoc := `
---
apiVersion: pantalon.kallan.dev/v1alpha1
kind: TerraformConfiguration
metadata:
  name: hello-world
context:
  foo: "1"
  bar: b
  baz: 'true'
  qux: |
    multi
`
	tfCfg, _, err := Parse([]byte(yamlDoc))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "1", tfCfg.Context["foo"])
	assert.Equal(t, "b", tfCfg.Context["bar"])
	assert.Equal(t, "true", tfCfg.Context["baz"])
	assert.Equal(t, "multi\n", tfCfg.Context["qux"])
}

// Like Kubernetes labels and annotations, values must be YAML strings, so an
// unquoted 1.10 fails instead of silently becoming "1.1".
func TestUnmarshalTerraformConfiguration_RejectsNonStringValues(t *testing.T) {
	for name, tc := range map[string]struct{ yaml, wantErr string }{
		"label integer":    {"  labels:\n    tier: 1\n", `line 8: metadata.labels value for key "tier" must be a string; quote it`},
		"label float":      {"  labels:\n    version: 1.10\n", `line 8: metadata.labels value for key "version" must be a string; quote it`},
		"label bool":       {"  labels:\n    env: prod\n    legacy: true\n", `line 9: metadata.labels value for key "legacy" must be a string; quote it`},
		"label null":       {"  labels:\n    env:\n", `line 8: metadata.labels value for key "env" must be a string; quote it`},
		"label flow style": {"  labels: {tier: 1}\n", `line 7: metadata.labels value for key "tier" must be a string; quote it`},
		"context integer":  {"context:\n  project-number: 123\n", `line 8: context value for key "project-number" must be a string; quote it`},
		"context bool":     {"context:\n  enabled: false\n", `line 8: context value for key "enabled" must be a string; quote it`},
		"context null":     {"context:\n  sa: ~\n", `line 8: context value for key "sa" must be a string; quote it`},
	} {
		t.Run(name, func(t *testing.T) {
			_, _, err := Parse([]byte(minimalYamlDoc + tc.yaml))
			assert.EqualError(t, err, tc.wantErr)
		})
	}
}

func TestUnmarshalTerraformConfiguration_ApiVersionMissing(t *testing.T) {
	yamlDoc := `
---
kind: TerraformConfiguration
metadata:
  name: hello-world
`
	_, _, err := Parse([]byte(yamlDoc))
	assert.EqualError(t, err, "invalid version")
}

func TestUnmarshalTerraformConfiguration_InvalidMetadataNameSnakeCase(t *testing.T) {
	yamlDoc := `
---
kind: TerraformConfiguration
apiVersion: pantalon.kallan.dev/v1alpha1
metadata:
  name: hello_world
`
	_, _, err := Parse([]byte(yamlDoc))
	assert.ErrorContains(t, err, "invalid metadata.name")
}

func TestUnmarshalTerraformConfiguration_InvalidMetadataNameDots(t *testing.T) {
	yamlDoc := `
---
kind: TerraformConfiguration
apiVersion: pantalon.kallan.dev/v1alpha1
metadata:
  name: hello.world
`
	_, _, err := Parse([]byte(yamlDoc))
	assert.ErrorContains(t, err, "invalid metadata.name")
}

// Names are rejected, not lowercased, so "Dev" and "dev" can't both exist.
func TestUnmarshalTerraformConfiguration_InvalidMetadataNameUppercase(t *testing.T) {
	yamlDoc := `
---
kind: TerraformConfiguration
apiVersion: pantalon.kallan.dev/v1alpha1
metadata:
  name: Hello-world
`
	_, _, err := Parse([]byte(yamlDoc))
	assert.ErrorContains(t, err, `invalid metadata.name "Hello-world"`)
}

func TestMarshalItems(t *testing.T) {
	tests := []struct {
		name     string
		input    []TerraformConfiguration
		expected []ConfigurationItem
	}{
		{
			name: "Single item",
			input: []TerraformConfiguration{
				{
					Metadata: Metadata{Name: "item1"},
					Context:  map[string]string{"key1": "value1"},
					Path:     "/path/to/item1/pantalon.yaml",
				},
			},
			expected: []ConfigurationItem{
				{
					Name:    "item1",
					Context: map[string]string{"key1": "value1"},
					Path:    "/path/to/item1/pantalon.yaml",
					Dir:     "/path/to/item1",
				},
			},
		},
		{
			name: "Multiple items",
			input: []TerraformConfiguration{
				{
					Metadata: Metadata{Name: "item1"},
					Context:  map[string]string{"key1": "value1"},
					Path:     "/path/to/item1/pantalon.yaml",
				},
				{
					Metadata: Metadata{Name: "item2"},
					Context:  map[string]string{"key2": "value2"},
					Path:     "/path/to/item2/pantalon.yaml",
				},
			},
			expected: []ConfigurationItem{
				{
					Name:    "item1",
					Context: map[string]string{"key1": "value1"},
					Path:    "/path/to/item1/pantalon.yaml",
					Dir:     "/path/to/item1",
				},
				{
					Name:    "item2",
					Context: map[string]string{"key2": "value2"},
					Path:    "/path/to/item2/pantalon.yaml",
					Dir:     "/path/to/item2",
				},
			},
		},
		{
			name: "Empty context",
			input: []TerraformConfiguration{
				{
					Metadata: Metadata{Name: "item1"},
					Context:  map[string]string{},
					Path:     "/path/to/item1/pantalon.yaml",
				},
			},
			expected: []ConfigurationItem{
				{
					Name:    "item1",
					Context: map[string]string{},
					Path:    "/path/to/item1/pantalon.yaml",
					Dir:     "/path/to/item1",
				},
			},
		},
		{
			name:     "No items",
			input:    []TerraformConfiguration{},
			expected: []ConfigurationItem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MarshalItems(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnmarshalTerraformConfiguration_WithLabels(t *testing.T) {
	yamlDoc := `
---
apiVersion: pantalon.kallan.dev/v1alpha1
kind: TerraformConfiguration
metadata:
  name: hello-world
  labels:
    environment: prod
    example.com/tier: "1"
`
	tfCfg, _, err := Parse([]byte(yamlDoc))
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"environment": "prod", "example.com/tier": "1"}, tfCfg.Metadata.Labels)
}

func TestUnmarshalTerraformConfiguration_InvalidLabels(t *testing.T) {
	for name, labels := range map[string]string{
		"key":   "bad key: prod",
		"value": "environment: bad value",
	} {
		t.Run(name, func(t *testing.T) {
			yamlDoc := minimalYamlDoc + "  labels:\n    " + labels + "\n"
			_, _, err := Parse([]byte(yamlDoc))
			assert.Error(t, err)
		})
	}
}

func TestMarshalItems_CarriesLabels(t *testing.T) {
	items, err := MarshalItems([]TerraformConfiguration{{
		Metadata: Metadata{Name: "a", Labels: map[string]string{"env": "dev"}},
		Path:     "a/pantalon.yaml",
	}})
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"env": "dev"}, items[0].Labels)
}

func TestUnmarshalTerraformConfiguration_WithDependencies(t *testing.T) {
	yamlDoc := minimalYamlDoc + `dependencies:
  paths:
    - terraform/compute/modules/**
    - "!terraform/compute/modules/*/v1"
    - terraform/shared/providers
`
	tfCfg, _, err := Parse([]byte(yamlDoc))
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"terraform/compute/modules/**",
		"!terraform/compute/modules/*/v1",
		"terraform/shared/providers",
	}, tfCfg.Dependencies.Paths)
}

func TestUnmarshalTerraformConfiguration_InvalidDependencies(t *testing.T) {
	for name, paths := range map[string]string{
		"first negated":  `["!a/**"]`,
		"empty":          `["a", ""]`,
		"empty negation": `["a", "!"]`,
		"absolute":       `["/a/b"]`,
		"parent dir":     `["../a"]`,
		"nested parent":  `["a/../../b"]`,
		"negated parent": `["a", "!../b"]`,
		"malformed glob": `["a/[b"]`,
	} {
		t.Run(name, func(t *testing.T) {
			yamlDoc := minimalYamlDoc + "dependencies:\n  paths: " + paths + "\n"
			_, _, err := Parse([]byte(yamlDoc))
			assert.ErrorContains(t, err, "invalid dependencies.paths")
		})
	}
}

func TestMarshalItems_CarriesDependencies(t *testing.T) {
	items, err := MarshalItems([]TerraformConfiguration{
		{
			Metadata:     Metadata{Name: "a"},
			Dependencies: Dependencies{Paths: []string{"modules/**"}},
			Path:         "a/pantalon.yaml",
		},
		{
			Metadata: Metadata{Name: "b"},
			Path:     "b/pantalon.yaml",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, &Dependencies{Paths: []string{"modules/**"}}, items[0].Dependencies)
	assert.Nil(t, items[1].Dependencies)
}

// Configurations without dependencies must produce the same output as before
// dependencies existed.
func TestMarshalItems_DependenciesOutput(t *testing.T) {
	items, err := MarshalItems([]TerraformConfiguration{
		{
			Metadata:     Metadata{Name: "a"},
			Dependencies: Dependencies{Paths: []string{"modules/**"}},
			Path:         "a/pantalon.yaml",
		},
		{
			Metadata: Metadata{Name: "b"},
			Path:     "b/pantalon.yaml",
		},
	})
	assert.NoError(t, err)

	out, err := yaml.Marshal(items)
	assert.NoError(t, err)
	assert.Equal(t, `- name: a
  path: a/pantalon.yaml
  dir: a
  labels: {}
  dependencies:
    paths:
    - modules/**
  context: {}
- name: b
  path: b/pantalon.yaml
  dir: b
  labels: {}
  context: {}
`, string(out))
}

func TestUnknownFieldWarnings(t *testing.T) {
	for name, tc := range map[string]struct {
		yamlDoc string
		want    []Warning
	}{
		"all known fields": {minimalYamlDoc + `  labels:
    env: dev
dependencies:
  paths: [a]
context:
  anything: goes
`, nil},
		"every unknown field, with its path": {minimalYamlDoc + `  lables:
    env: dev
dependecies:
  paths: [a]
`, []Warning{
			{Line: 7, Message: `unknown field "metadata.lables"`},
			{Line: 9, Message: `unknown field "dependecies"`},
		}},
		"path is set by pantalon, not the user": {minimalYamlDoc + `path: foo
`, []Warning{{Line: 7, Message: `unknown field "path"`}}},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, unknownFieldWarnings([]byte(tc.yamlDoc)))
		})
	}
}

// Older binaries must still read configs that use fields added later.
func TestUnmarshalTerraformConfiguration_IgnoresUnknownFields(t *testing.T) {
	_, _, err := Parse([]byte(minimalYamlDoc + "fromTheFuture: true\n"))
	assert.NoError(t, err)
}

func TestUnmarshalTerraformConfiguration_RejectsAnchors(t *testing.T) {
	for name, doc := range map[string]string{
		"anchor": minimalYamlDoc + "context:\n  a: &x 1\n",
		"alias":  minimalYamlDoc + "context:\n  b: *x\n",
		"merge":  minimalYamlDoc + "context:\n  <<: {a: 1}\n",
	} {
		t.Run(name, func(t *testing.T) {
			_, _, err := Parse([]byte(doc))
			assert.ErrorContains(t, err, "line 8: YAML anchors, aliases and merge keys are not supported")
		})
	}
}

func TestUnmarshalTerraformConfiguration_AllowsAnchorLookalikes(t *testing.T) {
	doc := minimalYamlDoc + `context:
  a: "&x"
  b: 'it *is*'
  c: http://x?a=1&b=2
  "<<": quoted  # &x *x << in a comment
---
later: &ignored 1
`
	_, _, err := Parse([]byte(doc))
	assert.NoError(t, err)
}

func TestUnknownFieldWarnings_NoneWhenAnchored(t *testing.T) {
	assert.Nil(t, unknownFieldWarnings([]byte(minimalYamlDoc+"typo: &x 1\n")))
}

func TestIsValidDNSLabel(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{"single char", "a", true},
		{"kebab case", "compute-dev", true},
		{"digits", "0-9", true},
		{"63 chars is the maximum", strings.Repeat("a", 63), true},
		{"64 chars is too long", strings.Repeat("a", 64), false},
		{"254 chars is too long", strings.Repeat("a", 254), false},
		{"dot", "compute.dev", false},
		{"leading dot", ".compute", false},
		{"trailing dot", "compute.", false},
		{"dotted 63 chars", strings.Repeat("a", 31) + "." + strings.Repeat("b", 31), false},
		{"empty", "", false},
		{"leading hyphen", "-a", false},
		{"trailing hyphen", "a-", false},
		{"uppercase", "A", false},
		{"underscore", "a_b", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, isValidDNSLabel(c.input), "isValidDNSLabel(%q)", c.input)
		})
	}
}

func TestUnmarshalTerraformConfiguration_InvalidMetadataNameTooLong(t *testing.T) {
	name := strings.Repeat("a", 64)
	yamlDoc := nameHeader + "metadata:\n  name: " + name + "\n"
	_, _, err := Parse([]byte(yamlDoc))
	assert.ErrorContains(t, err, `invalid metadata.name "`+name+`"`)
}

func TestUnmarshalTerraformConfiguration_MetadataNameAtMaxLength(t *testing.T) {
	name := strings.Repeat("a", 63)
	cfg, _, err := Parse([]byte(nameHeader + "metadata:\n  name: " + name + "\n"))
	require.NoError(t, err)
	assert.Equal(t, name, cfg.Metadata.Name)
}

const nameHeader = "apiVersion: pantalon.kallan.dev/v1alpha1\nkind: TerraformConfiguration\n"
