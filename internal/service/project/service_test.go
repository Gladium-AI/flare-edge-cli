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

	tests := []struct {
		name     string
		template string
		model    string
		marker   string
	}{
		{name: "text", template: "ai-text", model: `defaultAIModel = "@cf/moonshotai/kimi-k2.5"`, marker: `messages.Call("push", user)`},
		{name: "chat", template: "ai-chat", model: `defaultAIModel = "@cf/moonshotai/kimi-k2.5"`, marker: `messages.Call("push", user)`},
		{name: "vision", template: "ai-vision", model: `defaultAIModel = "@cf/moonshotai/kimi-k2.5"`, marker: `image_url query parameter is required`},
		{name: "stt", template: "ai-stt", model: `defaultAIModel = "@cf/deepgram/nova-3"`, marker: `send audio bytes in the request body`},
		{name: "tts", template: "ai-tts", model: `defaultAIModel = "@cf/deepgram/aura-2-en"`, marker: `input.Set("speaker", speaker)`},
		{name: "image", template: "ai-image", model: `defaultAIModel = "@cf/black-forest-labs/flux-2-klein-9b"`, marker: `form.Call("append", "prompt", prompt)`},
		{name: "embeddings", template: "ai-embeddings", model: `defaultAIModel = "@cf/qwen/qwen3-embedding-0.6b"`, marker: `input.Set("text", text)`},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			filesystem := fs.New()
			store := configstore.New(filesystem)
			service := NewService(store, filesystem)
			root := t.TempDir()

			result, err := service.Init(context.Background(), InitOptions{
				Dir:        root,
				Name:       "ai-worker",
				ModulePath: "github.com/example/ai-worker",
				Template:   tt.template,
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
			if !strings.Contains(content, tt.model) {
				t.Fatalf("expected ai worker template to include default model %s, got %s", tt.model, content)
			}
			if !strings.Contains(content, `env.Get("AI")`) {
				t.Fatalf("expected ai worker template to read env.AI")
			}
			if !strings.Contains(content, tt.marker) {
				t.Fatalf("expected ai worker template to include marker %s", tt.marker)
			}
		})
	}
}

func TestInitJavaScriptTemplateWithNodeCompat(t *testing.T) {
	t.Parallel()

	filesystem := fs.New()
	store := configstore.New(filesystem)
	service := NewService(store, filesystem)
	root := t.TempDir()

	result, err := service.Init(context.Background(), InitOptions{
		Dir:        root,
		Name:       "js-worker",
		Runtime:    "js",
		NodeCompat: true,
		WithGit:    true,
	})
	if err != nil {
		t.Fatalf("init project: %v", err)
	}

	project, err := store.LoadProject(result.ProjectDir)
	if err != nil {
		t.Fatalf("load project: %v", err)
	}
	if project.EffectiveRuntime() != "js-worker" {
		t.Fatalf("unexpected runtime: %q", project.EffectiveRuntime())
	}
	if project.Template != "js-worker" {
		t.Fatalf("unexpected template: %q", project.Template)
	}
	if !project.NodeCompat {
		t.Fatalf("expected node compatibility to be enabled")
	}
	if project.MainPath() != "src/worker.mjs" {
		t.Fatalf("unexpected main path: %q", project.MainPath())
	}

	wranglerCfg, err := store.LoadWrangler(result.ProjectDir, project.WranglerConfig)
	if err != nil {
		t.Fatalf("load wrangler: %v", err)
	}
	if wranglerCfg.Main != "src/worker.mjs" {
		t.Fatalf("unexpected wrangler main: %q", wranglerCfg.Main)
	}
	if len(wranglerCfg.CompatibilityFlags) != 1 || wranglerCfg.CompatibilityFlags[0] != "nodejs_compat" {
		t.Fatalf("unexpected compatibility flags: %+v", wranglerCfg.CompatibilityFlags)
	}

	workerPath := filepath.Join(result.ProjectDir, "src", "worker.mjs")
	workerMain, err := filesystem.ReadFile(workerPath)
	if err != nil {
		t.Fatalf("read js worker: %v", err)
	}
	if !strings.Contains(string(workerMain), `export default`) {
		t.Fatalf("expected JavaScript worker module scaffold, got %s", string(workerMain))
	}
}
