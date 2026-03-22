package project

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

func TestInitAITemplate(t *testing.T) {
	t.Parallel()

	filesystem := fs.New()
	store := configstore.New(filesystem)
	service := NewService(store, filesystem)
	root := t.TempDir()

	result, err := service.Init(context.Background(), InitOptions{
		Dir:        root,
		Name:       "ai-worker",
		ModulePath: "github.com/example/ai-worker",
		Template:   "ai-text",
		WithGit:    true,
	})
	if err != nil {
		t.Fatalf("init project: %v", err)
	}

	project, err := store.LoadProject(result.ProjectDir)
	if err != nil {
		t.Fatalf("load project: %v", err)
	}
	if project.Bindings.AI == nil || project.Bindings.AI.Binding != "AI" || !project.Bindings.AI.Remote {
		t.Fatalf("expected project ai binding to be scaffolded, got %+v", project.Bindings.AI)
	}

	wranglerCfg, err := store.LoadWrangler(result.ProjectDir, project.WranglerConfig)
	if err != nil {
		t.Fatalf("load wrangler: %v", err)
	}
	if wranglerCfg.AI == nil || wranglerCfg.AI.Binding != "AI" || !wranglerCfg.AI.Remote {
		t.Fatalf("expected wrangler ai binding to be scaffolded, got %+v", wranglerCfg.AI)
	}

	workerMainPath := filepath.Join(result.ProjectDir, "cmd/worker/main.go")
	workerMain, err := filesystem.ReadFile(workerMainPath)
	if err != nil {
		t.Fatalf("read worker main: %v", err)
	}
	content := string(workerMain)
	if !strings.Contains(content, `defaultAIModel = "@cf/meta/llama-3.1-8b-instruct"`) {
		t.Fatalf("expected ai worker template to include default model, got %s", content)
	}
	if !strings.Contains(content, `env.Get("AI")`) {
		t.Fatalf("expected ai worker template to read env.AI")
	}
}
