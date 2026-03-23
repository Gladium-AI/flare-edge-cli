package config

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWranglerConfigRoundTripPreservesUnknownFields(t *testing.T) {
	input := []byte(`{
  "name": "demo-worker",
  "main": "dist/worker.mjs",
  "compatibility_date": "2026-03-19",
  "ai": { "binding": "AI", "remote": true },
  "unsafe": { "bindings": [{ "name": "EXTRA" }] },
  "env": {
    "production": {
      "vars": { "FOO": "bar" },
      "durable_objects": {
        "bindings": [{ "name": "COUNTER", "class_name": "Counter" }]
      }
    }
  }
}`)

	var cfg WranglerConfig
	if err := json.Unmarshal(input, &cfg); err != nil {
		t.Fatalf("unmarshal wrangler config: %v", err)
	}

	cfg.Vars = map[string]string{"BAR": "baz"}

	output, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal wrangler config: %v", err)
	}

	text := string(output)
	if !strings.Contains(text, `"unsafe":{"bindings":[{"name":"EXTRA"}]}`) {
		t.Fatalf("unknown top-level field was not preserved: %s", text)
	}
	if !strings.Contains(text, `"durable_objects":{"bindings":[{"name":"COUNTER","class_name":"Counter"}]}`) {
		t.Fatalf("unknown env field was not preserved: %s", text)
	}
	if !strings.Contains(text, `"vars":{"BAR":"baz"}`) {
		t.Fatalf("known var update was not written: %s", text)
	}
}

func TestWranglerConfigMarshalDoesNotResurrectClearedKnownFields(t *testing.T) {
	input := []byte(`{
  "name": "demo-worker",
  "ai": { "binding": "AI", "remote": true },
  "unsafe": { "bindings": [{ "name": "EXTRA" }] }
}`)

	var cfg WranglerConfig
	if err := json.Unmarshal(input, &cfg); err != nil {
		t.Fatalf("unmarshal wrangler config: %v", err)
	}

	cfg.AI = nil

	output, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal wrangler config: %v", err)
	}

	text := string(output)
	if strings.Contains(text, `"ai":`) {
		t.Fatalf("cleared known field should not be preserved from extras: %s", text)
	}
	if !strings.Contains(text, `"unsafe":{"bindings":[{"name":"EXTRA"}]}`) {
		t.Fatalf("unknown field should still be preserved: %s", text)
	}
}
