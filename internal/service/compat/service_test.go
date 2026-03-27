package compat

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
)

func TestCheckFindsUnsupportedPatterns(t *testing.T) {
	dir := t.TempDir()
	source := `package main

import (
	"net"
	"os"
)

func main() {
	_, _ = os.Open("file.txt")
	_, _ = net.Listen("tcp", ":8080")
}`
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.26.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := NewService().Check(context.Background(), CheckOptions{Path: dir})
	if err != nil {
		t.Fatal(err)
	}
	if result.ErrorCount < 2 {
		t.Fatalf("expected at least 2 compatibility errors, got %d", result.ErrorCount)
	}
	for _, item := range result.Diagnostics {
		if item.Line == 0 {
			t.Fatalf("expected diagnostic line to be set")
		}
	}
}

func TestCheckSkipsStaticAnalysisForJavaScriptWorkers(t *testing.T) {
	dir := t.TempDir()
	project := []byte(`{
  "schema_version": 1,
  "project_name": "js-worker",
  "runtime": "js-worker",
  "module_path": "github.com/example/js-worker",
  "package_name": "main",
  "template": "js-worker",
  "entry": "src/worker.mjs",
  "main": "src/worker.mjs",
  "out_dir": "src",
  "worker_name": "js-worker",
  "wrangler_config": "wrangler.jsonc",
  "compatibility_date": "2026-03-19",
  "compatibility_profile": "worker-js",
  "nodejs_compat": true,
  "bindings": {},
  "generated": {}
}`)
	if err := os.WriteFile(filepath.Join(dir, config.DefaultProjectConfigFile), project, 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := NewService().Check(context.Background(), CheckOptions{Path: dir})
	if err != nil {
		t.Fatal(err)
	}
	if result.Profile != config.DefaultJSCompatibilityProfile {
		t.Fatalf("unexpected profile: %q", result.Profile)
	}
	if result.ErrorCount != 0 || result.WarnCount != 0 || len(result.Diagnostics) != 0 {
		t.Fatalf("expected JavaScript worker compatibility check to be skipped, got %+v", result)
	}
}
