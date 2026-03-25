# Flare Edge CLI Workflows

Use these patterns when operating `flare-edge-cli` projects.

## 1. Bootstrap a Standard Go Worker

```bash
flare-edge-cli doctor --json
flare-edge-cli project init hello-edge --template edge-http
flare-edge-cli compat check --path ./hello-edge --json
flare-edge-cli build --path ./hello-edge --json
flare-edge-cli dev --path ./hello-edge --local
flare-edge-cli deploy --path ./hello-edge --json
```

## 2. Bootstrap a Standard JavaScript Worker

```bash
flare-edge-cli doctor --json
flare-edge-cli project init hello-js --runtime js
flare-edge-cli build --path ./hello-js --json
flare-edge-cli dev --path ./hello-js --local
flare-edge-cli deploy --path ./hello-js --json
```

To enable Cloudflare's Node.js compatibility layer for the generated Worker:

```bash
flare-edge-cli project init hello-js --runtime js --node-compat
```

## 3. Bootstrap an AI Worker

```bash
flare-edge-cli project init hello-ai --template ai-chat
flare-edge-cli build --path ./hello-ai --json
flare-edge-cli dev --path ./hello-ai --local
flare-edge-cli deploy --path ./hello-ai --json
```

## 4. Operate an Existing Project

```bash
flare-edge-cli project info --path ./hello-edge --json
flare-edge-cli compat check --path ./hello-edge --json
flare-edge-cli build --path ./hello-edge --json
flare-edge-cli logs tail --path ./hello-edge
```

For a JavaScript project, skip `compat check`:

```bash
flare-edge-cli project info --path ./hello-js --json
flare-edge-cli build --path ./hello-js --json
flare-edge-cli logs tail --path ./hello-js
```

## 5. Provision Data Resources

```bash
flare-edge-cli kv namespace create CACHE --path ./hello-edge --json
flare-edge-cli d1 create DB --path ./hello-edge --json
flare-edge-cli r2 bucket create FILES --path ./hello-edge --json
```

## 6. Safe Cleanup

For temporary resources created during tests or agent workflows:

```bash
flare-edge-cli teardown --path ./hello-edge --json
```

If the local project should also be removed:

```bash
flare-edge-cli teardown --path ./hello-edge --delete-project --json
```

## Notes

- Prefer `--json` for machine consumption.
- Prefer `--path` over relying on the current working directory when the target project is not obvious.
- Use `doctor` early if auth or tooling health is uncertain.
