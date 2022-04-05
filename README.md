# spini

SpIni (Spinnaker Initializer) - command line tool for managing [Spinnaker](https://spinnaker.io/) accounts, applications and pipelines.

## Get
```bash
git clone git@github.com:ealebed/spini.git
```

## Build
```bash
cd spini
```

```bash
go install github.com/ealebed/spini
```

or
```bash
make install
```

## Build docker image
```bash
make image
```

## Use

```bash
spini -h
```

---

## Syntax

Use the following syntax to run `spini` commands from your terminal window:

```bash
spini [command] [subcommand] [flags]
```

### Flags are

| flag | Description |
| ----------- | ------------ |
| `--config` | string; path to Spin CLI config file (default $HOME/.spin/config) |
| `--dry-run` | bool; print output / save generated files without real changing system configuration (default true) |
| `--gate-endpoint` | string; Gate (API server) endpoint (default "<http://localhost:8084>") |
| `-h`, `--help` | help for selected command |
| `--org` | string; GitHub source owner organization (default "ealebed") |
| `--version` | spini version |

### Commands are

| command | Description |
| ----------- | ------------ |
| `account`, `acc` | manage Spinnaker accounts (clusters) |
| `application`, `app` | manage Spinnaker application’s lifecycle |
| `help` | help about any command |
| `manifest` | manage Kubernetes manifests from remote repository |
| `pipeline`, `pipe` | manage Spinnaker pipelines |

### Account subcommands are

| subcommand | Description |
| ----------- | ------------ |
| `get` | returns the specified spinnaker account |
| `list`, `ls` | returns list of all spinnaker accounts |

### Application subcommands are

| subcommand | Description |
| ----------- | ------------ |
| `delete`, `del` | delete the specified application |
| `get` | returns the specified spinnaker application |
| `list`, `ls` | returns list of all spinnaker applications |
| `save`, `create` | save/update the provided spinnaker application |
| `save-all`, `create-all` | save/update all spinnaker applications from provided GitHub repository |

### Manifest subcommands are

| subcommand | Description |
| ----------- | ------------ |
| `delete`, `del` | delete yaml manifest(s) for provided application |
| `save`, `create`, `generate` | save/update yaml manifest(s) for provided application |
| `save-all`, `create-all`, `generate-all` | save/update yaml manifest(s) for for all applications from provided GitHub repository |

### Pipeline subcommands are

| subcommand | Description |
| ----------- | ------------ |
| `delete`, `del` | delete the provided pipeline from the provided spinnaker application |
| `delete-all`, `prune` | delete all pipelines in the provided spinnaker application |
| `disable`, `off` | disable pipelines in the provided spinnaker application |
| `disable-all` | disable all pipelines in the provided spinnaker account(cluster) |
| `enable`, `on` | enable pipelines in the provided spinnaker application |
| `enable-all` | enable all pipelines in the provided spinnaker account(cluster) |
| `execute`, `exec` | execute the provided pipeline in the provided spinnaker application |
| `execute-all` | execute all pipelines in the provided spinnaker application or kubernetes cluster |
| `get` | returns the pipeline with the provided name from the provided spinnaker application |
| `list`, `ls` | returns list of all pipelines for the provided spinnaker application |
| `save`, `create` | save/update pipeline(s) for the provided spinnaker application |
| `save-all`, `create-all` | save/update pipeline(s) for all spinnaker applications from provided GitHub repository |

## Examples: Common operations

### Manage Spinnaker applications

```bash
# Create a new (or update existing) Spinnaker application using the definition in configuration.json (from remote GitHub repository).
spini application save --name=spini-test-application --repo=test-k8s --local=false --dry-run=false

# Create a new (or update existing) Spinnaker application using the definition in configuration.json (from remote GitHub repository and custom branch).
spini application save --name=spini-test-application --repo=test-k8s --branch=custom --local=false --dry-run=false

# Create a new (or update existing) Spinnaker application using the definition in configuration.json (from local GitHub repository).
spini application save --name=spini-test-application --dry-run=false

# Create (or update if exist) all Spinnaker applications using the definitions in configuration.json (from remote GitHub repository).
spini application save-all --repo=test-k8s --local=false --dry-run=false

# Create (or update if exist) all Spinnaker applications using the definitions in configuration.json (from remote GitHub repository and custom branch).
spini application save-all --repo=test-k8s --branch=custom --local=false --dry-run=false

# Create (or update if exist) all Spinnaker applications using the definitions in configuration.json (from local GitHub).
spini application save-all --dry-run=false

# List all Spinnaker applications.
spini application list

# Retrieve a single Spinnaker application.
spini application get --name=spini-test-application

# Delete a single Spinnaker application.
spini application delete --name=spini-test-application --dry-run=false
```

### Manage Kubernetes manifests

```bash
# Create a new (or update existing) Kubernetes manifest(s) with custom commit message for provided application using the definitions in configuration.json (from remote GitHub repository).
spini manifest save --name=spini-test-application --commit-message="Custom message" --repo=test-k8s --local=false --dry-run=false

# Create a new (or update existing) Kubernetes manifest(s) with default commit message for provided application using the definitions in configuration.json (from local GitHub repository).
spini manifest save --name=spini-test-application --dry-run=false

# Create a new (or update existing) Kubernetes manifests for all applications from configuration.json (from remote GitHub repository and custom branch).
spini manifest save-all --repo=test-k8s --branch=custom --local=false --dry-run=false

# Delete Kubernetes manifest(s) for provided application using the definitions in configuration.json (from remote GitHub repository).
spini manifest delete --name=spini-test-application --repo=test-k8s --local=false --dry-run=false
```

### Manage Spinnaker pipelines

```bash
# Create a new (or update existing) Spinnaker pipeline(s) using the definition in configuration.json (from remote GitHub repository).
spini pipeline save --name=spini-test-application --repo=test-k8s --local=false --dry-run=false

# Create a new (or update existing) Spinnaker pipeline(s) using the definition in configuration.json (from local GitHub repository).
spini pipeline save --name=spini-test-application --dry-run=false

# Create a new (or update existing) pipeline(s) for all Spinnaker applications using the definition in configuration.json (from remote GitHub repository).
spini pipeline save-all --repo=test-k8s --local=false --dry-run=false

# Create a new (or update existing) pipeline(s) for all Spinnaker applications using the definition in configuration.json (from local GitHub repository).
spini pipeline save-all --dry-run=false

# List all pipelines in the provided Spinnaker application.
spini pipeline list

# Retrieve a single pipeline from the provided Spinnaker application.
spini pipeline get --name=spini-test-application --pipeline="deploy-gke1-dc(production)"

# Start a single pipeline execution from the provided Spinnaker application.
spini pipeline execute --name=spini-test-application --pipeline="deploy-gke1-dc(production)" --dry-run=false

# Start all pipelines execution from the provided Spinnaker application.
spini pipeline execute-all --name=spini-test-application --dry-run=false

# Start all pipelines execution in all Spinnaker applications from the provided Kubernetes cluster.
spini pipeline execute-all --account=gke1 --dry-run=false

# Enable pipelines in the provided Spinnaker application.
spini pipeline enable --name=spini-test-application --dry-run=false

# Enable all pipelines in the provided Spinnaker account(Kubernetes cluster).
spini pipeline enable-all --account=sgp1 --dry-run=false

# Disable pipelines in the provided Spinnaker application.
spini pipeline disable --name=spini-test-application --dry-run=false

# Disable all pipelines in the provided Spinnaker account(Kubernetes cluster).
spini pipeline disable-all --account=sgp1 --dry-run=false

# Delete a single pipeline from the provided Spinnaker application.
spini pipeline delete --name=spini-test-application --pipeline="deploy-gke1-dc(production)" --dry-run=false

# Delete all pipelines from the provided Spinnaker application.
spini pipeline delete-all --name=spini-test-application --dry-run=false
```

---
Sample definition application(s) properties are in `configuration.json` file repository

---

TODO:
- Configure colored/formatted output
- Add tests
- Configure CI/CD for PR
- Refactoring custom/hardcoded values to make tool more general
