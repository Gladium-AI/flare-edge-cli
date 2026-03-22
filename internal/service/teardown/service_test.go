package teardown

import (
	"testing"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
)

func TestResetLocalState(t *testing.T) {
	t.Parallel()

	project := config.Project{
		OutDir: "dist",
		Bindings: config.ProjectBindings{
			Vars:    map[string]string{"KEEP": "1"},
			Secrets: []string{"SECRET"},
			AI:      &config.AIBinding{Binding: "AI", Remote: true},
			KV:      []config.KVBinding{{Binding: "KV", ID: "kv-id"}},
			D1:      []config.D1Binding{{Binding: "DB", DatabaseName: "db"}},
			R2:      []config.R2Binding{{Binding: "BUCKET", BucketName: "bucket"}},
		},
	}
	wrangler := config.WranglerConfig{
		AI:           &config.WranglerAIBinding{Binding: "AI", Remote: true},
		Routes:       []config.WranglerRoute{{Pattern: "example.com/*"}},
		KVNamespaces: []config.WranglerKVNamespace{{Binding: "KV", ID: "kv-id"}},
		D1Databases:  []config.WranglerD1Database{{Binding: "DB", DatabaseName: "db"}},
		R2Buckets:    []config.WranglerR2Bucket{{Binding: "BUCKET", BucketName: "bucket"}},
	}

	cleanProject, cleanWrangler := resetLocalState(project, wrangler, false)
	if len(cleanProject.Bindings.KV) != 0 || len(cleanProject.Bindings.D1) != 0 || len(cleanProject.Bindings.R2) != 0 {
		t.Fatalf("expected bindings to be cleared")
	}
	if len(cleanProject.Bindings.Secrets) != 0 {
		t.Fatalf("expected secrets to be cleared")
	}
	if cleanProject.Bindings.AI != nil {
		t.Fatalf("expected ai binding to be cleared")
	}
	if got := cleanProject.Bindings.Vars["KEEP"]; got != "1" {
		t.Fatalf("expected vars to be preserved, got %q", got)
	}
	if cleanWrangler.AI != nil || len(cleanWrangler.Routes) != 0 || len(cleanWrangler.KVNamespaces) != 0 || len(cleanWrangler.D1Databases) != 0 || len(cleanWrangler.R2Buckets) != 0 {
		t.Fatalf("expected wrangler resources to be cleared")
	}

	keptProject, keptWrangler := resetLocalState(project, wrangler, true)
	if len(keptProject.Bindings.KV) != 1 || len(keptProject.Bindings.D1) != 1 || len(keptProject.Bindings.R2) != 1 {
		t.Fatalf("expected bindings to be preserved when keep-bindings is true")
	}
	if len(keptProject.Bindings.Secrets) != 0 {
		t.Fatalf("expected secrets to be cleared even when keeping bindings")
	}
	if keptProject.Bindings.AI == nil {
		t.Fatalf("expected ai binding to be preserved when keep-bindings is true")
	}
	if len(keptWrangler.Routes) != 0 {
		t.Fatalf("expected routes to be cleared")
	}
	if keptWrangler.AI == nil {
		t.Fatalf("expected ai binding to be preserved when keep-bindings is true")
	}
	if len(keptWrangler.KVNamespaces) != 1 || len(keptWrangler.D1Databases) != 1 || len(keptWrangler.R2Buckets) != 1 {
		t.Fatalf("expected wrangler bindings to be preserved when keep-bindings is true")
	}
}
