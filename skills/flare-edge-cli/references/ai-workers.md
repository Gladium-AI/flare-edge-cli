# AI Workers With Flare Edge CLI

`flare-edge-cli` supports first-class Workers AI scaffolds for Go/Wasm Workers.

## Supported AI Templates

- `ai-text`
- `ai-chat`
- `ai-vision`
- `ai-stt`
- `ai-tts`
- `ai-image`
- `ai-embeddings`

## Current Default Models

- `ai-text` / `ai-chat`: `@cf/moonshotai/kimi-k2.5`
- `ai-vision`: `@cf/moonshotai/kimi-k2.5`
- `ai-stt`: `@cf/deepgram/nova-3`
- `ai-tts`: `@cf/deepgram/aura-2-en`
- `ai-image`: `@cf/black-forest-labs/flux-2-klein-9b`
- `ai-embeddings`: `@cf/qwen/qwen3-embedding-0.6b`

## Local Testing Guidance

- `flare-edge-cli dev --local` still uses the remote Cloudflare AI binding.
- Expect real account usage and potential charges during local testing.
- First request in a fresh `wrangler dev` session may prompt for Cloudflare account selection.

## Typical AI Workflow

```bash
flare-edge-cli project init demo-ai --template ai-chat
flare-edge-cli build --path ./demo-ai --json
flare-edge-cli dev --path ./demo-ai --local
```

Then test the local endpoint with a simple request, for example:

```bash
curl 'http://127.0.0.1:8787/?prompt=Reply%20with%20OK'
```

## Cleanup

If the AI Worker was created for a temporary run, prefer:

```bash
flare-edge-cli teardown --path ./demo-ai --json
```
