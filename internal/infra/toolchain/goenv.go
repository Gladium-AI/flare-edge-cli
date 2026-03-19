package toolchain

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/paolo/flare-edge-cli/internal/infra/process"
)

type GoEnv struct {
	GOROOT string `json:"GOROOT"`
	GOPATH string `json:"GOPATH"`
}

type GoToolchain struct {
	runner process.Runner
}

func NewGoToolchain(runner process.Runner) *GoToolchain {
	return &GoToolchain{runner: runner}
}

func (g *GoToolchain) Version(ctx context.Context, dir string) (string, error) {
	result, err := g.runner.Run(ctx, process.Command{Name: "go", Args: []string{"version"}, Dir: dir})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result.Stdout), nil
}

func (g *GoToolchain) Env(ctx context.Context, dir string) (GoEnv, error) {
	result, err := g.runner.Run(ctx, process.Command{Name: "go", Args: []string{"env", "-json", "GOROOT", "GOPATH"}, Dir: dir})
	if err != nil {
		return GoEnv{}, err
	}

	var env GoEnv
	if err := json.Unmarshal([]byte(result.Stdout), &env); err != nil {
		return GoEnv{}, fmt.Errorf("decode go env: %w", err)
	}

	return env, nil
}

func (g *GoToolchain) WasmExecPath(ctx context.Context, dir string) (string, error) {
	env, err := g.Env(ctx, dir)
	if err != nil {
		return "", err
	}

	return filepath.Join(env.GOROOT, "lib", "wasm", "wasm_exec.js"), nil
}

