# pantalon

Pantalon is a simple and lightweight system to identify multiple Terraform configurations in a single repository.

## Why?

Terraform does not have a deterministic method to identify the different between a root module, and a child module. Technically you can run `terraform init && terraform apply` against any Terraform module, given you provide the right variables. A real world repository may contain a mix of root modules indented to represent instances of resources, and child modules intended to be consumed by other modules.

### Goals

Pantalon intends to:

- Create a machine readable output of the Terraform configurations (root modules) within a repository.
- Allow selection of configurations based on labels (filtering).
- Allow specification of metadata needed to executing the Terraform configuration such as Service Accounts.
- Allow filtering by files changed in the git commit (such as a push to `main` or introduced by a pull request).

### Anti-Goals

Pantalon does not intend to:

- Manage the execution of Terraform. Either directly or as an orchestrator.
- Interact with the Terraform state or external systems and backends.

## Usage

### Configuration File

Within each root Terraform module, create a `pantalon.yaml` file with the following content:

```yaml
apiVersion: pantalon.kallan.dev/v1alpha1
kind: TerraformConfiguration
metadata:
  name: my-configuration
  labels:
    environment: prod
    team: platform
```

The `metadata.name` field is used to identify the configuration. It MUST be unique within the repository and a valid DNS label (lowercase `kebab-case`; uppercase is rejected, not converted). Pantalon fails if two configurations share a name.

The optional `metadata.labels` field is a map of key/value pairs used to select configurations, following [Kubernetes label syntax](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set): keys have an optional DNS-subdomain prefix and `/` followed by a name of up to 63 characters, and values are empty or up to 63 characters of alphanumerics, `-`, `_` and `.`. Labels are included in the output as `labels`.

The filename must strictly be `pantalon.yaml`. `pantalon.yml` or `pantalon.json` is not supported.

The file must be plain, explicit YAML: anchors (`&a`), aliases (`*a`) and merge keys (`<<`) are rejected with an error naming the file and line.

Unknown fields produce a warning on stderr rather than an error, so older versions of pantalon can read configurations that use newer fields. Under GitHub Actions (`GITHUB_ACTIONS=true`) warnings are emitted as `::warning` annotations on the offending file and line.

### Listing Configurations

Pantalon can list the configurations within a repository.

The following command will list all configurations within the current directory:

```shell
.
└── terraform
    ├── compute
    │   ├── environments
    │   │   ├── dev
    │   │   │   ├── main.tf
    │   │   │   └── pantalon.yaml
    │   │   ├── prod
    │   │   │   ├── main.tf
    │   │   │   └── pantalon.yaml
    │   │   └── qa
    │   │       ├── main.tf
    │   │       └── pantalon.yaml
    │   └── modules
    │       ├── instance-group
    │       │   ├── v1
    │       │   │   └── main.tf
    │       │   └── v2
    │       │       └── main.tf
    │       └── template
    │           ├── v1
    │           │   └── main.tf
    │           └── v2
    │               └── main.tf
    ├── data
    │   ├── environments
    │   │   ├── dev
    │   │   │   ├── main.tf
    │   │   │   └── pantalon.yaml
    │   │   ├── prod
    │   │   │   ├── main.tf
    │   │   │   └── pantalon.yaml
    │   │   └── qa
    │   │       ├── main.tf
    │   │       └── pantalon.yaml
    │   └── modules
    │       ├── db-instance
    │       │   ├── v1
    │       │   │   └── main.tf
    │       │   └── v2
    │       │       └── main.tf
    │       └── iam
    │           ├── v1
    │           │   └── main.tf
    │           └── v2
    │               └── main.tf
    └── load-balancer
        ├── environments
        │   ├── dev
        │   │   ├── main.tf
        │   │   └── pantalon.yaml
        │   ├── prod
        │   │   ├── main.tf
        │   │   └── pantalon.yaml
        │   └── qa
        │       ├── main.tf
        │       └── pantalon.yaml
        └── modules
            ├── certificates
            │   ├── v1
            │   │   └── main.tf
            │   └── v2
            │       └── main.tf
            └── proxy
                ├── v1
                │   └── main.tf
                └── v2
                    └── main.tf
```

Pantalon is always ran from the current working directory, and should always be ran from the root of the git repository. This is important to ensure simple integration with tools like changed-files as well as passing correct paths to GitHub Workflow job matrixes.

```shell
pantalon --output-format=yaml
```

```yaml
- name: compute-dev
  path: terraform/compute/environments/dev/pantalon.yaml
  dir: terraform/compute/environments/dev
  labels:
    component: compute
    environment: dev
  dependencies:
    paths:
    - terraform/compute/modules/**
  context:
    gcp-service-account: infrastructure@pantalon-dev.iam.gserviceaccount.com
- name: compute-prod
  path: terraform/compute/environments/prod/pantalon.yaml
  dir: terraform/compute/environments/prod
  labels:
    component: compute
    environment: prod
  dependencies:
    paths:
    - terraform/compute/modules/**
    - "!terraform/compute/modules/*/v2"
  context:
    gcp-service-account: infrastructure@pantalon-prod.iam.gserviceaccount.com
- name: compute-qa
  path: terraform/compute/environments/qa/pantalon.yaml
  dir: terraform/compute/environments/qa
  labels:
    component: compute
    environment: qa
  dependencies:
    paths:
    - terraform/compute/modules/**
  context:
    gcp-service-account: infrastructure@pantalon-qa.iam.gserviceaccount.com
- name: data-dev
  path: terraform/data/environments/dev/pantalon.yaml
  dir: terraform/data/environments/dev
  labels:
    component: data
    environment: dev
  context:
    gcp-service-account: infrastructure@pantalon-dev.iam.gserviceaccount.com
- name: data-prod
  path: terraform/data/environments/prod/pantalon.yaml
  dir: terraform/data/environments/prod
  labels:
    component: data
    environment: prod
  context:
    gcp-service-account: infrastructure@pantalon-prod.iam.gserviceaccount.com
- name: data-qa
  path: terraform/data/environments/qa/pantalon.yaml
  dir: terraform/data/environments/qa
  labels:
    component: data
    environment: qa
  context:
    gcp-service-account: infrastructure@pantalon-qa.iam.gserviceaccount.com
- name: lbl-dev
  path: terraform/load-balancer/environments/dev/pantalon.yaml
  dir: terraform/load-balancer/environments/dev
  labels:
    component: load-balancer
    environment: dev
  context:
    gcp-service-account: infrastructure@pantalon-dev.iam.gserviceaccount.com
- name: lbl-prod
  path: terraform/load-balancer/environments/prod/pantalon.yaml
  dir: terraform/load-balancer/environments/prod
  labels:
    component: load-balancer
    environment: prod
  context:
    gcp-service-account: infrastructure@pantalon-prod.iam.gserviceaccount.com
- name: lbl-qa
  path: terraform/load-balancer/environments/qa/pantalon.yaml
  dir: terraform/load-balancer/environments/qa
  labels:
    component: load-balancer
    environment: qa
  context:
    gcp-service-account: infrastructure@pantalon-qa.iam.gserviceaccount.com
```

### Changed Directories

Pantalon can filter configurations based on the directories changed in the git commit.

```shell
pantalon --output-format=yaml --changed-dirs='["terraform/compute/environments/dev","terraform/data/environments/dev"]'
```

The changed-dirs argument has been designed to be supplied by [tj-actions/changed-files](https://github.com/tj-actions/changed-files) using the `dir_names` option.

```yaml
    - name: Run changed-files with dir_names
      id: changed-files-dir-names
      uses: tj-actions/changed-files@v45
      with:
        dir_names: "true"
        json: "true"
        escape_json: "false"
```

`json` makes changed-files emit a JSON array rather than a space-separated list, and `escape_json: "false"` keeps the quotes unescaped so Pantalon can parse it.

Entries are normalised before matching, so `a/dev/`, `./a/dev` and `a/dev` are equivalent, and an empty entry is treated as the repository root (matching every configuration). Absolute paths and paths that leave the repository (`../`) are rejected with an error.

Each matching configuration is only emitted once, even if multiple entries in `--changed-dirs` fall within the same configuration's directory (for example, separate files changed in two different submodules of the same root module). This avoids generating duplicate matrix jobs for the same `pantalon.yaml`.

### Dependencies

A configuration often depends on directories outside its own, such as shared child modules. List them under `dependencies.paths` so that changes there also select the configuration when using `--changed-dirs`:

```yaml
apiVersion: pantalon.kallan.dev/v1alpha1
kind: TerraformConfiguration
metadata:
  name: compute-prod
dependencies:
  paths:
    - terraform/compute/modules/**
    # prod has not been upgraded to v2 modules yet
    - "!terraform/compute/modules/*/v2"
    - terraform/shared/providers
```

Patterns follow the style of GitHub Actions [`on.<event>.paths`](https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions#onpushpull_requestpull_request_targetpathspaths-ignore) filters, using [doublestar](https://github.com/bmatcuk/doublestar) globs:

- Patterns are relative to the repository root, like `dir` and `--changed-dirs`. Absolute paths and `..` are rejected.
- A pattern matches a directory and everything below it, so `terraform/shared/providers`, `terraform/shared/providers/**` and `terraform/compute/modules/*/v2` all match changes in subdirectories too.
- A leading `!` excludes directories matched by earlier patterns. Patterns are evaluated in order and the last match wins, so the first pattern cannot be negated. Quote negated patterns, as `!` is special in YAML.
- Exclusions only apply to dependencies. Changes within the configuration's own directory always select it.
- `--changed-dirs` contains directories, not files, so file patterns such as `**/*.tf` will not match them.

Changing `terraform/compute/modules/template/v2` would then select every compute configuration except `compute-prod`:

```shell
pantalon --output-format=yaml --changed-dirs='["terraform/compute/modules/template/v2"]'
```

```yaml
- name: compute-dev
  path: terraform/compute/environments/dev/pantalon.yaml
  dir: terraform/compute/environments/dev
  labels:
    component: compute
    environment: dev
  dependencies:
    paths:
    - terraform/compute/modules/**
  context:
    gcp-service-account: infrastructure@pantalon-dev.iam.gserviceaccount.com
- name: compute-qa
  ...
```

`dependencies` is included in the output when set, and omitted otherwise. `--path-glob` and `--selector` match only the configuration's own `dir` and labels; dependencies affect `--changed-dirs` only.

### Path Glob Filtering

Pantalon can filter configurations by directory path using [doublestar](https://github.com/bmatcuk/doublestar) glob patterns. Pass `--path-glob` one or more times; a configuration is included if its directory matches **any** of the supplied patterns (OR logic).

```shell
# All configurations under terraform/compute
pantalon --output-format=yaml --path-glob='terraform/compute/**'
```

```yaml
- name: compute-dev
  path: terraform/compute/environments/dev/pantalon.yaml
  dir: terraform/compute/environments/dev
  labels:
    component: compute
    environment: dev
  dependencies:
    paths:
    - terraform/compute/modules/**
  context:
    gcp-service-account: infrastructure@pantalon-dev.iam.gserviceaccount.com
- name: compute-prod
  path: terraform/compute/environments/prod/pantalon.yaml
  dir: terraform/compute/environments/prod
  labels:
    component: compute
    environment: prod
  dependencies:
    paths:
    - terraform/compute/modules/**
    - "!terraform/compute/modules/*/v2"
  context:
    gcp-service-account: infrastructure@pantalon-prod.iam.gserviceaccount.com
- name: compute-qa
  path: terraform/compute/environments/qa/pantalon.yaml
  dir: terraform/compute/environments/qa
  labels:
    component: compute
    environment: qa
  dependencies:
    paths:
    - terraform/compute/modules/**
  context:
    gcp-service-account: infrastructure@pantalon-qa.iam.gserviceaccount.com
```

Use `*` to match a single path segment and `**` to match any number of segments:

| Pattern | Matches |
|---|---|
| `terraform/compute/**` | Every configuration nested under `terraform/compute/` |
| `**/environments/prod` | Every `prod` environment regardless of component |
| `terraform/*/environments/dev` | Every `dev` environment, one component level deep |
| `terraform/compute/environments/dev` | Exact directory match |

Multiple `--path-glob` flags are combined with OR logic — a configuration is included if it matches at least one pattern:

```shell
# Production environments for compute and data only
pantalon --output-format=yaml \
  --path-glob='terraform/compute/environments/prod' \
  --path-glob='terraform/data/environments/prod'
```

`--path-glob` can be combined with `--changed-dirs`. Changed-dirs filtering is applied first, and glob filtering is applied to the result:

```shell
pantalon --output-format=yaml \
  --changed-dirs='["terraform/compute/environments/dev","terraform/data/environments/dev"]' \
  --path-glob='terraform/compute/**'
```

This returns only changed configurations that also match the glob — in this case, just `compute-dev`.

See the [path glob example](examples/.github/workflows/terraform-plan-prod.yaml) for a GitHub Actions workflow that uses `--path-glob` to target production environments.

### Label Selectors

Pantalon can filter configurations by label using `--selector` (or `-l`), with the same syntax as `kubectl -l`:

```shell
pantalon --output-format=yaml -l 'environment=prod,component in (compute,data)'
```

| Selector | Matches |
|---|---|
| `environment=prod`, `environment==prod` | Label equals value |
| `environment!=prod` | Label differs from value, or is absent |
| `environment in (dev,qa)` | Label is one of the values |
| `environment notin (dev,qa)` | Label is none of the values, or is absent |
| `team` | Label is present |
| `!team` | Label is absent |

Requirements are separated by commas and combined with AND. Unlike kubectl, only one `--selector` is accepted (repeating it, or passing an empty value, is an error) and comparison operators (`>`, `<`) are not supported.

> **Caution:** negative selectors (`!=`, `notin`, `!key`) also match configurations that lack the label entirely, so a typo such as `enviroment: prod` would be selected by `-l 'environment!=prod'`. For gating jobs, prefer positive selectors like `-l 'environment in (dev,qa)'`.

> **Note:** quote label values in YAML (`version: "1.10"`). Unquoted scalars are converted to strings by the YAML parser, so `1.10` becomes `1.1`.

`--selector` can be combined with `--changed-dirs` and `--path-glob`; a configuration must satisfy all of the supplied filters.

### Matrix

The primary intent is to  use Pantalon to generate a matrix of configurations to be executed by a GitHub Actions.

```yaml
...
jobs:
  ...
  plan:
    name: ${{ matrix.configs.name }}
    needs:
      - define-matrix
    runs-on: ubuntu-latest
    strategy:
      matrix:
        configs: ${{ fromJSON(needs.define-matrix.outputs.configs) }}
    defaults:
      run:
        working-directory: ${{ matrix.configs.dir }}
    ...
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - uses: google-github-actions/auth@v2
        with:
          service_account: ${{ matrix.configs.context.gcp-service-account }}
          workload_identity_provider: projects/12345678912345/locations/global/workloadIdentityPools/my-identity-pool/providers/my-provider
      ...
      - name: Terraform Plan
        run: |
          terraform init
          terraform plan
```

See the [example](examples/.github/workflows/terraform-plan.yaml) for the full workflow.

## Roadmap

- [x] Support listing dependencies of a root module within the pantalon file.
- [x] Allow filtering by label selectors.
- [x] Allow filtering by path glob.
- [x] Filter by the union of git files changed and directories detected
- [ ] Support other configuration use cases other than Terraform.
- [ ] Create a Docker release.
- [x] Allow the supply of arbitrary metadata.
