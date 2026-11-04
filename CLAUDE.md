# Pantalon

A Go CLI tool that identifies Terraform root module configurations within a repository and produces machine-readable output for use in CI/CD pipelines (e.g. GitHub Actions job matrices).

## Project Structure

```
cmd/pantalon/    # CLI entrypoint (main.go: run(); flags.go) and golden CLI tests
internal/pantalon/  # All logic: config parsing, discovery, filtering, output
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
- **`metadata.name`** — must be a valid RFC 1123 DNS subdomain label (lowercase alphanumeric and hyphens, max 253 chars). Must be unique (case-insensitively) across the repository; `Discover` fails on duplicates before any filtering.
- Unknown fields warn on stderr (as `::warning` annotations under GitHub Actions) instead of failing, so older binaries can read newer configs. `Parse` returns them as `[]Warning`; `main` renders them.
- **`metadata.labels`** — optional Kubernetes-style labels (same key/value rules), emitted as `labels` and matched by `--selector`/`-l`.
- **`dependencies.paths`** — optional GitHub-Actions-style glob patterns (root-relative, `!` negation, last match wins) for directories outside the config's own dir; a changed dir matching them also selects the config under `--changed-dirs`.
- **`context`** — arbitrary key/value map for metadata like GCP service accounts, passed through to output.
- The tool always runs from the **repository root** via `Discover(os.DirFS("."))`.
- When `--changed-dirs` is supplied, output is filtered to configurations whose directory contains (or is contained by) a changed directory, or whose `dependencies.paths` match one.

## Architecture

- `Discover(fs.FS)` walks the filesystem looking for `pantalon.yaml` files, stopping descent into a directory once a match is found (`fs.SkipDir`).
- `MarshalItems()` converts `[]TerraformConfiguration` into the flat `[]ConfigurationItem` output struct.
- `Filter(items, FilterOptions)` is the single filtering entry point: `ParseChangedDirs` + `ChangedFiles` (changed dirs), `GlobFilter`, then `ParseSelector` + `SelectorFilter` (kubectl-style selector: `=`, `==`, `!=`, `in`, `notin`, `key`, `!key`, comma = AND).
- `Write(w, items, format)` renders JSON or YAML.
- Goldens in `testdata/golden/` pin flags, stdout, stderr and exit codes; regenerate only deliberately with `go test ./cmd/pantalon -update`.
