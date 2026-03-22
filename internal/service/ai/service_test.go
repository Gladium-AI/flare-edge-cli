package ai

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

func TestSetAndClearBinding(t *testing.T) {
	t.Parallel()

	filesystem := fs.New()
	store := configstore.New(filesystem)
	service := NewService(store, filesystem)
	dir := t.TempDir()

	project := config.DefaultProject("ai-worker", "github.com/example/ai-worker", "main", "edge-http", config.DefaultCompatibilityDate, "")
	wranglerCfg := config.WranglerConfig{
		Name:              project.WorkerName,
		Main:              filepath.ToSlash(filepath.Join(project.OutDir, project.ShimFile)),
		CompatibilityDate: project.CompatibilityDate,
	}
	if err := store.SaveProject(dir, project); err != nil {
		t.Fatalf("save project: %v", err)
	}
	if err := store.SaveWrangler(dir, project.WranglerConfig, wranglerCfg); err != nil {
		t.Fatalf("save wrangler: %v", err)
	}

	result, err := service.SetBinding(context.Background(), SetBindingOptions{
		Dir:     dir,
		Binding: "AI",
		Remote:  true,
	})
	if err != nil {
		t.Fatalf("set binding: %v", err)
	}
	if result.Binding != "AI" || !result.Remote {
		t.Fatalf("unexpected set result: %+v", result)
	}

	project, err = store.LoadProject(dir)
	if err != nil {
		t.Fatalf("load project after set: %v", err)
	}
	if project.Bindings.AI == nil || project.Bindings.AI.Binding != "AI" || !project.Bindings.AI.Remote {
		t.Fatalf("project ai binding not persisted: %+v", project.Bindings.AI)
	}

	wranglerCfg, err = store.LoadWrangler(dir, project.WranglerConfig)
	if err != nil {
		t.Fatalf("load wrangler after set: %v", err)
	}
	if wranglerCfg.AI == nil || wranglerCfg.AI.Binding != "AI" || !wranglerCfg.AI.Remote {
		t.Fatalf("wrangler ai binding not persisted: %+v", wranglerCfg.AI)
	}

	clearResult, err := service.ClearBinding(context.Background(), ClearBindingOptions{Dir: dir})
	if err != nil {
		t.Fatalf("clear binding: %v", err)
	}
	if !clearResult.Cleared {
		t.Fatalf("expected cleared result, got %+v", clearResult)
	}

	project, err = store.LoadProject(dir)
	if err != nil {
		t.Fatalf("load project after clear: %v", err)
	}
	if project.Bindings.AI != nil {
		t.Fatalf("expected project ai binding to be cleared")
	}

	wranglerCfg, err = store.LoadWrangler(dir, project.WranglerConfig)
	if err != nil {
		t.Fatalf("load wrangler after clear: %v", err)
	}
	if wranglerCfg.AI != nil {
		t.Fatalf("expected wrangler ai binding to be cleared")
	}
}
