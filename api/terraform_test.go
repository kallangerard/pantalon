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
					Path:     "/path/to/item1",
				},
			},
			expected: []ConfigurationItem{
				{
					Name:    "item1",
					Context: map[string]string{"key1": "value1"},
					Path:    "/path/to/item1",
				},
			},
		},
		{
			name: "Multiple items",
			input: []TerraformConfiguration{
				{
					Metadata: Metadata{Name: "item1"},
					Context:  map[string]string{"key1": "value1"},
					Path:     "/path/to/item1",
				},
				{
					Metadata: Metadata{Name: "item2"},
					Context:  map[string]string{"key2": "value2"},
					Path:     "/path/to/item2",
				},
			},
			expected: []ConfigurationItem{
				{
					Name:    "item1",
					Context: map[string]string{"key1": "value1"},
					Path:    "/path/to/item1",
				},
				{
					Name:    "item2",
					Context: map[string]string{"key2": "value2"},
					Path:    "/path/to/item2",
				},
			},
		},
		{
			name: "Empty context",
			input: []TerraformConfiguration{
				{
					Metadata: Metadata{Name: "item1"},
					Context:  map[string]string{},
					Path:     "/path/to/item1",
				},
			},
			expected: []ConfigurationItem{
				{
					Name:    "item1",
					Context: map[string]string{},
					Path:    "/path/to/item1",
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
