package compat

import (
	"context"
	"os"
	"path/filepath"
	"testing"
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
