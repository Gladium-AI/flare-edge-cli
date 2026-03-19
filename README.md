# flare-edge-cli

`flare-edge-cli` validates Go-for-Wasm compatibility, builds Workers-ready Wasm artifacts, generates the Worker shim/config, and delegates Cloudflare operations to Wrangler.

## Status

This repository implements the CLI surface for:

- auth
- project
- compat
- build
- dev
- deploy
- route
- secret
- kv
- d1
- r2
- logs
- release
- doctor

## Development

```bash
go test ./...
go vet ./...
```

