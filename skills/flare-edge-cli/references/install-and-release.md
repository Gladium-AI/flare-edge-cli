# Install And Release

## Install Options

One-line installer:

```bash
curl -fsSL https://raw.githubusercontent.com/Gladium-AI/flare-edge-cli/main/install.sh | sh
```

Source build:

```bash
make build
make install
```

## Release Assets

GitHub releases publish binaries for:

- `linux/amd64`
- `linux/arm64`
- `darwin/amd64`
- `darwin/arm64`
- `windows/amd64`
- `windows/arm64`

Release assets are attached automatically by the GitHub Actions release workflow when a GitHub release is published.
