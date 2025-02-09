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
```

The `metadata.name` field is used to identify the configuration. This should be unique within the repository, and must be a valid DNS label (i.e. `kebab-case`). This name will be cast to lowercase.

The filename must strictly be `pantalon.yaml`. `pantalon.yml` or `pantalon.json` is not supported.

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

```shell
pantalon --output-format=yaml .
```

```yaml
- name: pantalon-example-compute-dev
  context:
    gcp-service-account: infrastructure@pantalon-dev.iam.gserviceaccount.com
- name: pantalon-example-compute-prod
  context:
    gcp-service-account: infrastructure@pantalon-prod.iam.gserviceaccount.com
- name: pantalon-example-compute-qa
  context:
    gcp-service-account: infrastructure@pantalon-qa.iam.gserviceaccount.com
- name: pantalon-example-data-dev
  context:
    gcp-service-account: infrastructure@pantalon-dev.iam.gserviceaccount.com
- name: pantalon-example-data-prod
  context:
    gcp-service-account: infrastructure@pantalon-prod.iam.gserviceaccount.com
- name: pantalon-example-data-qa
  context:
    gcp-service-account: infrastructure@pantalon-qa.iam.gserviceaccount.com
- name: pantalon-example-lbl-dev
  context:
    gcp-service-account: infrastructure@pantalon-dev.iam.gserviceaccount.com
- name: pantalon-example-lbl-prod
  context:
    gcp-service-account: infrastructure@pantalon-prod.iam.gserviceaccount.com
- name: pantalon-example-lbl-qa
  context:
    gcp-service-account: infrastructure@pantalon-qa.iam.gserviceaccount.com
```

## Roadmap

- [ ] Support listing dependencies of a root module within the pantalon file.
- [ ] Allow filtering by label selectors.
- [ ] Allow filtering by path glob.
- [ ] Filter by the union of git files changed and directories detected
- [ ] Support other configuration use cases other than Terraform.
- [ ] Create a Docker release.
- [ ] Create a GitHub action.
- [ ] Allow the supply of arbitrary metadata.
