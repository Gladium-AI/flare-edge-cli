---
name: flare-edge-cli
description: Use this skill when the task is to scaffold, validate, build, develop, deploy, inspect, operate, or tear down Cloudflare Workers with flare-edge-cli. Trigger on requests about Cloudflare edge functions, Go/Wasm Workers, native JavaScript Workers, Workers AI in Go, or cleanup of ephemeral Cloudflare resources created through flare-edge-cli.
license: MIT
metadata:
  author: Gladium AI
  version: 1.1.0
  category: developer-tools
  tags:
    - cloudflare
    - workers
    - golang
    - javascript
    - wasm
    - workers-ai
    - deployment
---

# Flare Edge CLI

Use this skill to operate `flare-edge-cli` safely and consistently.

## Use This Skill For

- Creating a new Go-based or native JavaScript Cloudflare Worker project
- Validating Go-for-Wasm compatibility before build or deploy
- Building and running local dev flows for Workers
- Deploying or operating KV, D1, R2, routes, secrets, releases, and logs
- Creating or testing Go-based Workers AI projects
- Tearing down disposable Workers and attached side effects

## Core Rules

- Prefer `flare-edge-cli` over raw `wrangler` when the task fits the CLI surface.
- Prefer `--json` output when another agent or program will consume the result.
- Keep project work scoped to an explicit `--path` when you are not already inside the generated project directory.
- Do not run `auth logout --all` unless the user explicitly asks to clear Cloudflare auth.
- For disposable test environments, finish with `flare-edge-cli teardown` so remote side effects are removed.
- If a task involves real Cloudflare AI usage, note that local dev still uses remote Workers AI and may incur charges.

## Quick Workflow

1. Verify prerequisites.
2. Initialize or inspect the project.
3. Run compatibility and build checks when the runtime needs them.
4. Use `dev` for local validation.
5. Use `deploy` for live rollout.
6. Use service-specific commands for KV, D1, R2, secrets, routes, logs, and releases.
7. Use `teardown` for cleanup when the environment is temporary.

## Prerequisites

- `flare-edge-cli` available on `PATH`, or use the repo-local binary/build path.
- Go installed for Go/Wasm projects.
- Wrangler installed.
- Cloudflare auth already configured.

For the latest install flow and release-binary behavior, see [references/install-and-release.md](references/install-and-release.md).

## Standard Command Path

For a standard Go Worker:

```bash
flare-edge-cli doctor --json
flare-edge-cli project init my-worker --template edge-http
flare-edge-cli compat check --path ./my-worker --json
flare-edge-cli build --path ./my-worker --json
flare-edge-cli dev --path ./my-worker --local
flare-edge-cli deploy --path ./my-worker --json
```

For a standard JavaScript Worker:

```bash
flare-edge-cli doctor --json
flare-edge-cli project init my-js-worker --runtime js
flare-edge-cli build --path ./my-js-worker --json
flare-edge-cli dev --path ./my-js-worker --local
flare-edge-cli deploy --path ./my-js-worker --json
```

If the user needs Cloudflare's Node.js compatibility layer, scaffold with:

```bash
flare-edge-cli project init my-js-worker --runtime js --node-compat
```

For an AI Worker:

```bash
flare-edge-cli project init my-ai-worker --template ai-chat
flare-edge-cli build --path ./my-ai-worker --json
flare-edge-cli dev --path ./my-ai-worker --local
flare-edge-cli deploy --path ./my-ai-worker --json
```

Load [references/ai-workers.md](references/ai-workers.md) when the task is about AI templates, current model defaults, or how to test Workers AI locally.

## Cleanup

If you create temporary infrastructure or disposable test projects, tear them down explicitly:

```bash
flare-edge-cli teardown --path ./my-worker --json
```

Use `--keep-bindings` only when the user wants to preserve KV, D1, or R2 resources.

## Troubleshooting

- Start with `flare-edge-cli doctor`.
- If deployment succeeds but the Worker fails at runtime, use `flare-edge-cli logs tail`.
- If a command mutates Cloudflare resources, verify whether the target project path and Worker name are correct before rerunning it.
- If a project mixes manual Wrangler edits with CLI-managed config, inspect both `flare-edge.json` and `wrangler.jsonc`.
- `compat check` is only meaningful for Go/Wasm projects; JavaScript Worker projects intentionally skip Go static analysis.

For common command sequences and operational guidance, read [references/workflows.md](references/workflows.md).
