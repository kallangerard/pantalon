package api

import (
	"testing"
)

func TestUnmarshalTerraformConfiguration(t *testing.T) {
	yamlDoc := `
---
apiVersion: pantalon.kallan.dev/v1alpha1
kind: TerraformConfiguration
metadata:
  name: hello-world
`

	cfg := config{}
	v, err := cfg.Unmarshal([]byte(yamlDoc))
	if err != nil {
		t.Fatal(err)
	}

	expected := TerraformConfiguration{
		ApiVersion: PantalonVersion,
		Kind:       TerraformKind,
		Metadata: Metadata{
			Name: "hello-world",
		},
	}

	if v != expected {
		t.Fatalf("expected %v, got %v", expected, v)
	}
}
