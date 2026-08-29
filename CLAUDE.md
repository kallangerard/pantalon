# Pantalon

A Go CLI tool that identifies Terraform root module configurations within a repository and produces machine-readable output for use in CI/CD pipelines (e.g. GitHub Actions job matrices).

## Project Structure

```
cmd/pantalon/    # CLI entrypoint (main.go)
api/             # Core types and logic (TerraformConfiguration, ConfigurationItem, marshaling)
file/            # Filesystem search and changed-file filtering
testdata/        # Fixtures used by tests
examples/        # Example GitHub Actions workflows
```

## Common Commands

```bash
# Run all tests
go test ./...

# Run linter
golangci-lint run ./...

# Build the binary
go build ./cmd/pantalon/

# Run the binary from the repo root
./pantalon --output-format=yaml
./pantalon --output-format=json --changed-dirs='["terraform/compute/environments/dev"]'
```

## Key Concepts

- **`pantalon.yaml`** — marker file placed in each Terraform root module. Must use `apiVersion: pantalon.kallan.dev/v1alpha1` and `kind: TerraformConfiguration`.
- **`metadata.name`** — must be a valid RFC 1123 DNS subdomain label (lowercase alphanumeric and hyphens, max 253 chars).
- **`context`** — arbitrary key/value map for metadata like GCP service accounts, passed through to output.
- The tool always runs from the **repository root**, walking outwards from `.`.
- When `--changed-dirs` is supplied, output is filtered to configurations whose directory is a prefix of (or equal to) a changed directory.

## Architecture

- `file.Search()` walks the filesystem looking for `pantalon.yaml` files, stopping descent into a directory once a match is found. The walk uses `os.ReadDir` (the entry list already says whether a marker file is present, so no per-directory `os.Stat` is needed), runs subdirectories across a bounded pool of goroutines, and concatenates per-directory results so output stays in `filepath.WalkDir` lexical order. `.git` and `.terraform` are skipped (`file.skippedDirs`).
- Marker files are read and parsed across the same bounded pool.
- Benchmarks live in `api/bench_test.go` and `file/bench_test.go`; run them with `go test -bench=. -benchmem ./...`. `go.mod` targets Go 1.23, so use `for i := 0; i < b.N; i++` rather than `b.Loop()`.
- `api.MarshalItems()` converts `[]TerraformConfiguration` into the flat `[]ConfigurationItem` output struct.
- `api.UnmarshalChangedFileJson()` parses the JSON array from `--changed-dirs`.
- `file.ChangedFiles()` filters items by matching against changed directories.
- Output is either JSON or YAML via `github.com/goccy/go-yaml`.
