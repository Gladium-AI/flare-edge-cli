# flare-edge-cli

[![Go Version](https://img.shields.io/badge/Go-1.26.0-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

`flare-edge-cli` scaffolds, validates, builds, deploys, tails, and tears down Go-for-Wasm Cloudflare Workers projects. It generates the Worker shim and Wrangler config, runs compatibility checks against a Workers/Wasm profile, and delegates Cloudflare operations to Wrangler and Cloudflare APIs where appropriate.

This CLI is designed first for AI agents. Humans can use it directly, but the primary goal is to give coding agents a stable, scriptable control surface for creating and operating Cloudflare edge functions and lightweight microservices in Go.

## Why This Tool Exists

`flare-edge-cli` exists to standardize the end-to-end agent workflow for Go on Cloudflare:

- scaffold a deployable project with deterministic structure
- validate whether the Go code fits Workers/Wasm constraints
- build the `.wasm` artifact and Worker shim correctly every time
- provision and manage Cloudflare resources through one consistent interface
- emit machine-readable output that agents can inspect and chain into later actions
- tear down remote and local side effects when an ephemeral environment is no longer needed

The intended operator is usually an AI agent acting on behalf of a developer. Because of that, the interface is optimized around automation-friendly behavior:

- stable command names and flag semantics
- explicit config files
- predictable project layout
- JSON output for diagnostics and follow-up steps
- clear separation between command, service, and infrastructure logic
- safe cleanup for disposable environments created during agent workflows

## Primary Use Case

The primary use case is simple: an AI agent needs a standard way to create and deploy an edge function or small microservice on Cloudflare using Go without rebuilding the same scaffolding, compatibility analysis, build orchestration, deployment logic, and cleanup flow for every task.

This means `flare-edge-cli` is not just a deploy wrapper. It is an agent-oriented execution surface for:

- project generation
- compatibility analysis
- build orchestration
- deployment orchestration
- Cloudflare resource provisioning
- release control
- environment cleanup

## What It Does

- Scaffolds Go Worker projects with a reproducible layout
- Validates Go code against a Workers/Wasm compatibility profile
- Builds `.wasm` artifacts and the JavaScript Worker shim
- Runs local or remote dev sessions through Wrangler
- Deploys versioned Workers and manages routes, secrets, KV, D1, R2, and releases
- Tears down Workers and optional side-effect resources cleanly

## Design Principles

The implementation is intentionally biased toward agent use:

- machine-readable first, with human-readable output as a secondary mode
- deterministic scaffolding and build output
- explicit config mutation instead of hidden side effects
- thin CLI handlers with testable service-layer logic
- idempotent or safely repeatable operations where Cloudflare semantics allow it
- clean separation between local project state and remote Cloudflare state
- strong teardown support for disposable agent-created environments

## Requirements

- Go `1.26.0`
- [Wrangler](https://developers.cloudflare.com/workers/wrangler/) installed and available on `PATH`
- Cloudflare authentication already configured through Wrangler or an API token
- Optional: TinyGo for `--tinygo` builds

## Install

Install the Claude Code skill:

```bash
curl -fsSL https://raw.githubusercontent.com/Gladium-AI/flare-edge-cli/main/install-skill.sh | sh
```

That installs the skill to `~/.claude/skills/flare-edge-cli` by default. Override the target directory with `SKILLS_DIR=/path/to/skills`.

Build the binary from the repository root:

```bash
make build
```

That produces a local `./bin/flare-edge-cli` binary.

Install it into a user-local bin directory:

```bash
make install
```

The install target prefers `XDG_BIN_HOME` when set, then an existing `~/.local/bin`, then `~/bin`, and otherwise creates `~/.local/bin`. You can override the destination explicitly with `make install INSTALL_DIR=/path/to/bin`.

You can also run it directly during development:

```bash
go run ./cmd/flare-edge-cli --help
```

## Typical Workflow

Initialize a project:

```bash
./flare-edge-cli project init my-worker --module-path github.com/example/my-worker
```

Check compatibility:

```bash
./flare-edge-cli compat check --path ./my-worker
```

Build the Wasm bundle:

```bash
./flare-edge-cli build --path ./my-worker
```

Run local development:

```bash
./flare-edge-cli dev --path ./my-worker --local
```

Deploy:

```bash
./flare-edge-cli deploy --path ./my-worker
```

Clean everything up later:

```bash
./flare-edge-cli teardown --path ./my-worker
```

## Generated Project Layout

`project init` creates a Go Worker project with these important files:

```text
<project>/
  flare-edge.json
  wrangler.jsonc
  go.mod
  cmd/worker/main.go
  internal/generated/worker_shim.mjs
  README.md
  .gitignore
```

Important generated paths:

- `flare-edge.json`: CLI-owned project metadata
- `wrangler.jsonc`: deployable Wrangler configuration
- `cmd/worker/main.go`: Go/Wasm entrypoint
- `dist/app.wasm`: compiled artifact after build
- `dist/worker.mjs`: Worker shim after build

## Configuration Files

### `flare-edge.json`

This is the CLI’s typed project metadata file. It tracks:

- project and module names
- Go entrypoint
- output directory and artifact names
- Worker name
- compatibility date and profile
- bindings for KV, D1, R2, vars, and secrets
- generated shim metadata

### `wrangler.jsonc`

This is the file deployed by Wrangler. The CLI updates it when you:

- build
- deploy
- attach or detach routes
- add KV, D1, or R2 bindings
- configure custom domains
- run teardown

## Output Modes

Most commands support two output styles:

- human-readable output for interactive use
- `--json` for machine-readable automation

Compatibility checks also support SARIF:

```bash
./flare-edge-cli compat check --path ./my-worker --sarif
```

## Global Flag

Every command accepts:

```text
--account-id <id>
```

This overrides `CLOUDFLARE_ACCOUNT_ID` for the current invocation.

## Command Reference

### `auth`

Authenticate against Cloudflare and inspect local auth state.

```bash
./flare-edge-cli auth login [--wrangler] [--browser] [--api-token <token>] [--account-id <id>] [--persist] [--non-interactive]
./flare-edge-cli auth whoami [--json]
./flare-edge-cli auth logout [--all] [--local-only]
```

Notes:

- `auth login --wrangler` uses Wrangler-managed OAuth
- `auth login --api-token` validates and optionally persists token metadata
- `auth logout --local-only` removes only CLI-local state
- `auth logout --all` also clears Wrangler-managed auth

### `project`

Create and inspect projects.

```bash
./flare-edge-cli project init <name> [--cwd <dir>] [--module-path <go_module>] [--package <pkg>] [--template <name>] [--compat-date <YYYY-MM-DD>] [--env <name>] [--use-jsonc] [--with-git] [--yes]
./flare-edge-cli project info [--cwd <dir>] [--json] [--show-generated] [--show-bindings]
```

Supported templates:

- `edge-http`
- `edge-json`
- `scheduled`
- `kv-api`
- `d1-api`
- `r2-api`

### `compat`

Run the Go/Wasm compatibility analyzer or inspect built-in rules.

```bash
./flare-edge-cli compat check [--path <dir>] [--entry <pkg-or-file>] [--profile worker-wasm] [--strict] [--json] [--sarif] [--fail-on warning|error] [--exclude <glob>]
./flare-edge-cli compat rules [--json] [--severity error|warning|info]
```

Diagnostics include structured fields such as rule ID, severity, file, line, message, why, and fix hint.

### `build`

Compile Go to Wasm and inspect artifacts.

```bash
./flare-edge-cli build [--path <dir>] [--entry <pkg-or-file>] [--out-dir <dir>] [--out-file <file.wasm>] [--shim-out <file>] [--target js/wasm] [--optimize size|speed] [--tinygo] [--no-shim] [--clean] [--json]
./flare-edge-cli build wasm [same flags as build]
./flare-edge-cli build inspect [--artifact <path>] [--size] [--exports] [--imports] [--json]
```

The build writes:

- `dist/app.wasm`
- `dist/worker.mjs`
- `dist/wasm_exec.js`

### `dev`

Start a Wrangler-powered development session.

```bash
./flare-edge-cli dev [--path <dir>] [--env <name>] [--port <port>] [--remote] [--local] [--persist] [--inspector-port <port>]
```

Flags `--open` and `--watch` exist as reserved compatibility flags.

### `deploy`

Validate, build, and deploy the Worker.

```bash
./flare-edge-cli deploy [--path <dir>] [--env <name>] [--name <worker>] [--compat-date <YYYY-MM-DD>] [--route <pattern>] [--custom-domain <hostname>] [--workers-dev] [--dry-run] [--upload-only] [--message <text>] [--var <KEY=VALUE>] [--keep-vars] [--minify] [--latest] [--json]
```

Behavior:

- runs compatibility checks first
- builds the Wasm artifact and Worker shim
- updates Wrangler config
- deploys through Wrangler
- can attach routes and custom domains during deploy

### `route`

Attach and detach remote routes and custom domains. These commands also update local config and will build artifacts automatically if needed.

```bash
./flare-edge-cli route attach --route <pattern> [--path <dir>] [--zone <zone>] [--env <name>] [--script <worker>] [--json]
./flare-edge-cli route domain --hostname <hostname> [--path <dir>] [--zone <zone>] [--env <name>] [--script <worker>] [--json]
./flare-edge-cli route detach [--route <pattern>] [--hostname <hostname>] [--path <dir>] [--zone <zone>] [--env <name>] [--json]
```

### `secret`

Manage Worker secrets.

```bash
./flare-edge-cli secret put <KEY> [--path <dir>] [--value <value>] [--from-file <path>] [--env <name>] [--versioned] [--json]
./flare-edge-cli secret list [--path <dir>] [--env <name>] [--json]
./flare-edge-cli secret delete <KEY> [--path <dir>] [--env <name>] [--versioned] [--json]
```

### `kv`

Manage KV namespaces and entries.

```bash
./flare-edge-cli kv namespace create <binding> [--path <dir>] [--title <title>] [--env <name>] [--provision] [--json]
./flare-edge-cli kv put --binding <binding> --key <key> [--path <dir>] [--value <value>] [--from-file <path>] [--ttl <seconds>] [--expiration <unix>] [--metadata <json>] [--env <name>] [--json]
./flare-edge-cli kv get --binding <binding> --key <key> [--path <dir>] [--text] [--env <name>] [--json]
./flare-edge-cli kv list --binding <binding> [--path <dir>] [--prefix <prefix>] [--env <name>] [--json]
./flare-edge-cli kv delete --binding <binding> --key <key> [--path <dir>] [--env <name>] [--json]
```

### `d1`

Manage D1 databases and migrations.

```bash
./flare-edge-cli d1 create <binding> [--path <dir>] [--name <db_name>] [--env <name>] [--json]
./flare-edge-cli d1 execute --binding <binding> [--path <dir>] [--sql <statement>] [--file <path.sql>] [--remote] [--local] [--json] [--result-json] [--env <name>]
./flare-edge-cli d1 migrations new <name> [--path <dir>] [--dir <migration_dir>] [--json]
./flare-edge-cli d1 migrations apply --binding <binding> [--path <dir>] [--remote] [--local] [--env <name>] [--json]
```

Notes:

- `--json` on `d1 execute` requests clean D1 JSON from Wrangler
- `--result-json` wraps the CLI response itself in JSON

### `r2`

Manage R2 buckets and objects.

```bash
./flare-edge-cli r2 bucket create <binding> [--path <dir>] [--name <bucket>] [--location <region>] [--storage-class <class>] [--env <name>] [--json]
./flare-edge-cli r2 object put --binding <binding> --key <key> --file <path> [--path <dir>] [--content-type <mime>] [--cache-control <value>] [--content-disposition <value>] [--env <name>] [--json]
./flare-edge-cli r2 object get --binding <binding> --key <key> [--path <dir>] [--out <path>] [--env <name>] [--json]
./flare-edge-cli r2 object delete --binding <binding> --key <key> [--path <dir>] [--env <name>] [--json]
```

### `logs`

Tail runtime logs through Wrangler.

```bash
./flare-edge-cli logs tail [--path <dir>] [--env <name>] [--worker <name-or-route>] [--format pretty|json] [--search <text>] [--status <value>] [--sampling <0-1>]
```

### `release`

Inspect and manage versioned releases.

```bash
./flare-edge-cli release list [--path <dir>] [--env <name>] [--name <worker>] [--limit <n>] [--json]
./flare-edge-cli release promote <version_id> [--path <dir>] [--env <name>] [--message <text>] [--yes] [--json]
./flare-edge-cli release rollback <version_id> [--path <dir>] [--env <name>] [--message <text>] [--yes] [--json]
```

### `doctor`

Check whether the local environment and project are deployable.

```bash
./flare-edge-cli doctor [--path <dir>] [--json] [--verbose]
```

`doctor` checks:

- Go installation
- Wrangler installation
- auth health
- project config validity
- compatibility date presence
- Wasm build readiness
- binding/config sanity

### `teardown`

Delete a Worker and optionally remove related Cloudflare resources and local artifacts.

```bash
./flare-edge-cli teardown [--path <dir>] [--env <name>] [--name <worker>] [--keep-bindings] [--keep-artifacts] [--delete-project] [--json]
```

Default teardown behavior:

- deletes the deployed Worker
- deletes attached routes and custom domains
- deletes bound KV namespaces, D1 databases, and R2 buckets
- removes `dist/` and `.wrangler/`
- scrubs local binding and routing config from `flare-edge.json` and `wrangler.jsonc`

Use `--keep-bindings` if you want to keep KV, D1, and R2 resources.

## Aliases

These top-level aliases are available:

```text
flare-edge-cli init       -> flare-edge-cli project init
flare-edge-cli info       -> flare-edge-cli project info
flare-edge-cli check      -> flare-edge-cli compat check
flare-edge-cli tail       -> flare-edge-cli logs tail
flare-edge-cli rollback   -> flare-edge-cli release rollback
```

The top-level `build` command is itself the primary Wasm build entrypoint, with `build wasm` available as an explicit subcommand.

## Examples

Scaffold a JSON worker:

```bash
./flare-edge-cli init test-project --module-path github.com/example/test-project --template edge-json
```

Check compatibility and emit JSON:

```bash
./flare-edge-cli check --path ./test-project --json
```

Inspect the built artifact:

```bash
./flare-edge-cli build inspect --artifact ./test-project/dist/app.wasm --size --exports
```

Create a KV namespace and provision the local binding:

```bash
./flare-edge-cli kv namespace create CACHE --path ./test-project --title test-project-cache --provision
```

Create a D1 database and apply migrations remotely:

```bash
./flare-edge-cli d1 create DB --path ./test-project --name test-project-db
./flare-edge-cli d1 migrations new init --path ./test-project
./flare-edge-cli d1 migrations apply --path ./test-project --binding DB --remote
```

Upload and read an R2 object:

```bash
./flare-edge-cli r2 bucket create ASSETS --path ./test-project --name test-project-assets
./flare-edge-cli r2 object put --path ./test-project --binding ASSETS --key hello.txt --file ./README.md
./flare-edge-cli r2 object get --path ./test-project --binding ASSETS --key hello.txt --out /tmp/hello.txt
```

Attach a route or custom domain:

```bash
./flare-edge-cli route attach --path ./test-project --route example.com/api/*
./flare-edge-cli route domain --path ./test-project --hostname worker.example.com
```

## Development

Useful local commands:

```bash
make build
make test
make install
go test ./...
go vet ./...
go test -race ./...
go run ./cmd/flare-edge-cli --help
```

GitHub automation shipped with this repository:

- `.github/workflows/ci.yml`: `gofmt`, `go vet`, race-tested unit tests, and cross-platform builds
- `.github/workflows/security.yml`: dependency review on pull requests, `govulncheck`, and CodeQL analysis
- `.github/dependabot.yml`: weekly Go module and GitHub Actions dependency updates

## Notes

- Human-readable output is intended for operators; use `--json` when scripting
- The CLI keeps command handlers thin and pushes work into service and infrastructure layers
- `flare-edge-cli` depends on Wrangler for Cloudflare deployment flows rather than replacing Wrangler entirely
