package build

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	"github.com/paolo/flare-edge-cli/internal/infra/process"
	"github.com/paolo/flare-edge-cli/internal/infra/toolchain"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

type fakeRunner struct {
	last process.Command
}

func (r *fakeRunner) Run(_ context.Context, cmd process.Command) (process.Result, error) {
	r.last = cmd

	output := ""
	for index := 0; index+1 < len(cmd.Args); index++ {
		if cmd.Args[index] == "-o" {
			output = cmd.Args[index+1]
			break
		}
	}
	if output == "" {
		return process.Result{}, nil
	}
	if !filepath.IsAbs(output) {
		output = filepath.Join(cmd.Dir, output)
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return process.Result{}, err
	}
	if err := os.WriteFile(output, []byte("wasm"), 0o644); err != nil {
		return process.Result{}, err
	}
	return process.Result{}, nil
}

func (r *fakeRunner) Stream(_ context.Context, _ process.Command, _, _ io.Writer) error {
	return nil
}

func TestWasmUsesStableArtifactPathAndRequiresArtifact(t *testing.T) {
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "services", "orchestrator")
	filesystem := fs.New()
	store := configstore.New(filesystem)

	project := config.DefaultProject("orchestrator", "github.com/paolo/orchestrator", "main", "edge-http", config.DefaultCompatibilityDate, "")
	if err := store.SaveProject(projectDir, project); err != nil {
		t.Fatalf("save project: %v", err)
	}

	runner := &fakeRunner{}
	service := NewService(store, filesystem, runner, &toolchain.GoToolchain{})

	result, err := service.Wasm(context.Background(), WasmOptions{
		Path:   projectDir,
		NoShim: true,
	})
	if err != nil {
		t.Fatalf("build wasm: %v", err)
	}

	expectedArtifact := filepath.Join(projectDir, "dist", "app.wasm")
	if result.Artifact != expectedArtifact {
		t.Fatalf("unexpected artifact path: %q", result.Artifact)
	}
	outputIndex := -1
	for index := 0; index+1 < len(runner.last.Args); index++ {
		if runner.last.Args[index] == "-o" {
			outputIndex = index + 1
			break
		}
	}
	if outputIndex == -1 {
		t.Fatalf("expected -o flag in args: %v", runner.last.Args)
	}
	if !filepath.IsAbs(runner.last.Args[outputIndex]) {
		t.Fatalf("compiler output path should be absolute, got %q", runner.last.Args[outputIndex])
	}
	if runner.last.Args[outputIndex] != expectedArtifact {
		t.Fatalf("unexpected compiler output path: %q", runner.last.Args[outputIndex])
	}
	exists, err := filesystem.Exists(expectedArtifact)
	if err != nil {
		t.Fatalf("stat artifact: %v", err)
	}
	if !exists {
		t.Fatalf("expected artifact to exist at %s", expectedArtifact)
	}
}
