# Copilot Coding Instructions

## Project Overview

Pantalon is a Go CLI tool that discovers and lists Terraform root module configurations within a repository. Each root module is identified by a `pantalon.yaml` file. The tool outputs a machine-readable list of configurations, optionally filtered by recently changed directories, suitable for use as a GitHub Actions job matrix.

## Repository Structure

```
.
├── api/           # Core domain types and logic (YAML parsing, validation, output marshaling)
├── cmd/pantalon/  # CLI entry point
├── file/          # File system operations (searching for pantalon.yaml files, changed-file filtering)
├── examples/      # Example GitHub Actions workflows showing intended usage
└── testdata/      # Test fixtures used by file package tests
```

## Language & Toolchain

- **Language**: Go (see `go.mod` for the exact version)
- **Dependencies**: `github.com/goccy/go-yaml` for YAML marshaling/unmarshaling, `github.com/stretchr/testify` for test assertions
- **Module path**: `github.com/kallangerard/pantalon`

## Building & Testing

```shell
# Run all tests
go test ./...

# Build the binary
go build ./cmd/pantalon
```

## Coding Conventions

- Follow standard Go formatting (`gofmt`).
- Package names are short and lowercase: `api`, `file`.
- Error handling uses `errors.New`, `fmt.Errorf`, and `errors.Join` (Go 1.20+) — avoid third-party error libraries.
- Prefer table-driven tests using `[]struct{ name string; ... }` as seen in `api/terraform_test.go`.
- Use `github.com/stretchr/testify/assert` for assertions in tests.
- The `api` package contains pure logic with no filesystem I/O; the `file` package handles all I/O.
- Configuration file must be named exactly `pantalon.yaml` — no other extensions are supported.
- `metadata.name` must be a valid RFC 1123 DNS subdomain label (lowercase, alphanumeric and hyphens, max 253 chars).
- The `context` field in `pantalon.yaml` is a free-form `map[string]string` for arbitrary metadata (e.g. GCP service account).

## Key Types

| Type | Package | Purpose |
|---|---|---|
| `TerraformConfiguration` | `api` | Parsed representation of a `pantalon.yaml` file |
| `ConfigurationItem` | `api` | Output-friendly representation with `name`, `path`, `dir`, and `context` |
| `Metadata` | `api` | Holds `name` from the YAML `metadata` block |

## CLI Flags

- `--output-format` — `json` (default) or `yaml`
- `--changed-dirs` — JSON array of directory paths (e.g. from `tj-actions/changed-files`) used to filter results to only configurations under those directories

## pantalon.yaml Format

```yaml
apiVersion: pantalon.kallan.dev/v1alpha1
kind: TerraformConfiguration
metadata:
  name: my-configuration   # kebab-case, unique within repo
context:                   # optional; arbitrary string key/value pairs
  gcp-service-account: sa@project.iam.gserviceaccount.com
```
