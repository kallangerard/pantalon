package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	cfg := config{}
	v, err := cfg.Unmarshal([]byte(minimalYamlDoc))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello-world", v.Metadata.Name)
}

func TestUnmarshalTerraformConfiguration_ApiVersionKindMatch(t *testing.T) {
	cfg := config{}
	v, err := cfg.Unmarshal([]byte(minimalYamlDoc))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, TerraformKind, v.Kind)
	assert.Equal(t, PantalonVersion, v.ApiVersion)
}

func TestUnmarshalTerraformConfiguration_WithContext(t *testing.T) {
	yamlDoc := `
---
apiVersion: pantalon.kallan.dev/v1alpha1
kind: TerraformConfiguration
metadata:
  name: hello-world
context:
  foo: 1
  bar: b
  baz: true
`
	cfg := config{}
	tfCfg, err := cfg.Unmarshal([]byte(yamlDoc))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "1", tfCfg.Context["foo"])
	assert.Equal(t, "b", tfCfg.Context["bar"])
	assert.Equal(t, "true", tfCfg.Context["baz"])
}

func TestUnmarshalTerraformConfiguration_ApiVersionMissing(t *testing.T) {
	yamlDoc := `
---
kind: TerraformConfiguration
metadata:
  name: hello-world
`
	cfg := config{}
	_, err := cfg.Unmarshal([]byte(yamlDoc))
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
	cfg := config{}
	_, err := cfg.Unmarshal([]byte(yamlDoc))
	assert.EqualError(t, err, "invalid metadata.name")
}

func TestUnmarshalTerraformConfiguration_InvalidMetadataNameDots(t *testing.T) {
	yamlDoc := `
---
kind: TerraformConfiguration
apiVersion: pantalon.kallan.dev/v1alpha1
metadata:
  name: hello.world
`
	cfg := config{}
	_, err := cfg.Unmarshal([]byte(yamlDoc))
	assert.EqualError(t, err, "invalid metadata.name")
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
	tfCfg, err := config{}.Unmarshal([]byte(yamlDoc))
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
			_, err := config{}.Unmarshal([]byte(yamlDoc))
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
