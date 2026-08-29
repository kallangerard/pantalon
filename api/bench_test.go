package api

import (
	"fmt"
	"testing"
)

func BenchmarkUnmarshal(b *testing.B) {
	doc := []byte(`apiVersion: pantalon.kallan.dev/v1alpha1
kind: TerraformConfiguration
metadata:
  name: compute-dev
context:
  service_account: ci@example.iam.gserviceaccount.com
  project: example-dev
`)
	cfg := New()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := cfg.Unmarshal(doc); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIsValidSubdomainLabel(b *testing.B) {
	names := []string{"compute-dev", "a", "network-prod-ap-southeast-2", "Invalid_Name", "-leading"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, n := range names {
			isValidSubdomainLabel(n)
		}
	}
}

func BenchmarkMarshalItems(b *testing.B) {
	cfgs := make([]TerraformConfiguration, 200)
	for i := range cfgs {
		cfgs[i] = TerraformConfiguration{
			ApiVersion: PantalonVersion,
			Kind:       TerraformKind,
			Metadata:   Metadata{Name: fmt.Sprintf("stack-%d", i)},
			Path:       fmt.Sprintf("terraform/stack%d/environments/dev/pantalon.yaml", i),
		}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := MarshalItems(cfgs); err != nil {
			b.Fatal(err)
		}
	}
}
