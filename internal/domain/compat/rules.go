package compat

import "github.com/paolo/flare-edge-cli/internal/domain/diagnostic"

type Rule struct {
	ID       string              `json:"id"`
	Severity diagnostic.Severity `json:"severity"`
	Summary  string              `json:"summary"`
	Why      string              `json:"why"`
	FixHint  string              `json:"fix_hint"`
}

func BuiltInRules() []Rule {
	return []Rule{
		{
			ID:       "FE001",
			Severity: diagnostic.SeverityError,
			Summary:  "cgo imports are not supported in Workers Wasm builds",
			Why:      "Cloudflare Workers does not expose native system libraries for Go cgo bindings.",
			FixHint:  "Remove import \"C\" and replace native integration points with pure Go or Web APIs.",
		},
		{
			ID:       "FE002",
			Severity: diagnostic.SeverityError,
			Summary:  "unsupported host packages are imported",
			Why:      "Workers does not support local process execution, plugins, or raw syscalls.",
			FixHint:  "Remove imports such as os/exec, plugin, or syscall and use Workers-compatible APIs.",
		},
		{
			ID:       "FE003",
			Severity: diagnostic.SeverityError,
			Summary:  "filesystem access is incompatible with Workers runtime constraints",
			Why:      "Workers does not provide a writable or readable local filesystem to Wasm modules.",
			FixHint:  "Move data access to KV, D1, R2, or request-bound state.",
		},
		{
			ID:       "FE004",
			Severity: diagnostic.SeverityError,
			Summary:  "server listeners are incompatible with the fetch-driven runtime model",
			Why:      "Workers handles requests through fetch events instead of opening network listeners.",
			FixHint:  "Replace http.ListenAndServe or net.Listen with a fetch handler exported to the Worker shim.",
		},
		{
			ID:       "FE005",
			Severity: diagnostic.SeverityError,
			Summary:  "process execution is unavailable in Workers",
			Why:      "Workers cannot spawn subprocesses from Go Wasm code.",
			FixHint:  "Move subprocess behavior to a separate service or replace it with an API call.",
		},
		{
			ID:       "FE006",
			Severity: diagnostic.SeverityWarning,
			Summary:  "goroutines can hide runtime assumptions in js/wasm builds",
			Why:      "Long-lived goroutines in js/wasm often translate into event-loop retention or scheduling surprises.",
			FixHint:  "Keep goroutines short-lived and ensure request handling remains event-driven.",
		},
		{
			ID:       "FE007",
			Severity: diagnostic.SeverityWarning,
			Summary:  "OS environment lookups do not map cleanly to Workers bindings",
			Why:      "Workers exposes configuration via env bindings, not process-wide environment variables.",
			FixHint:  "Thread env bindings through your handler instead of calling os.Getenv directly.",
		},
	}
}
