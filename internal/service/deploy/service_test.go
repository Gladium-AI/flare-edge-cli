package deploy

import (
	"testing"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
)

func TestSplitPair(t *testing.T) {
	key, value, ok := splitPair("FOO=bar")
	if !ok || key != "FOO" || value != "bar" {
		t.Fatalf("unexpected split result: %q %q %v", key, value, ok)
	}
}

func TestApplyVarsInitializesNilMaps(t *testing.T) {
	project := config.Project{}
	wranglerCfg := config.WranglerConfig{}

	applyVars(&project, &wranglerCfg, []string{"ORCHESTRATOR_URL=https://example.com", "BROKEN"})

	if got := wranglerCfg.Vars["ORCHESTRATOR_URL"]; got != "https://example.com" {
		t.Fatalf("unexpected wrangler var value: %q", got)
	}
	if got := project.Bindings.Vars["ORCHESTRATOR_URL"]; got != "https://example.com" {
		t.Fatalf("unexpected project var value: %q", got)
	}
	if _, exists := wranglerCfg.Vars["BROKEN"]; exists {
		t.Fatalf("invalid var entry should be ignored")
	}
}
