package config

import (
	"fmt"
	"path/filepath"
)

const (
	DefaultProjectConfigFile      = "flare-edge.json"
	DefaultWranglerConfigFile     = "wrangler.jsonc"
	DefaultCompatibilityDate      = "2026-03-19"
	DefaultCompatibilityProfile   = "worker-wasm"
	DefaultJSCompatibilityProfile = "worker-js"

	RuntimeGoWasm     = "go-wasm"
	RuntimeJavaScript = "js-worker"
)

type Project struct {
	SchemaVersion        int                    `json:"schema_version" validate:"required,min=1"`
	ProjectName          string                 `json:"project_name" validate:"required"`
	Runtime              string                 `json:"runtime,omitempty" validate:"omitempty,oneof=go-wasm js-worker"`
	ModulePath           string                 `json:"module_path" validate:"required"`
	PackageName          string                 `json:"package_name" validate:"required"`
	Template             string                 `json:"template" validate:"required"`
	Entry                string                 `json:"entry" validate:"required"`
	Main                 string                 `json:"main,omitempty"`
	OutDir               string                 `json:"out_dir" validate:"required"`
	WasmFile             string                 `json:"wasm_file,omitempty"`
	ShimFile             string                 `json:"shim_file,omitempty"`
	WorkerName           string                 `json:"worker_name" validate:"required"`
	WranglerConfig       string                 `json:"wrangler_config" validate:"required"`
	CompatibilityDate    string                 `json:"compatibility_date" validate:"required"`
	CompatibilityProfile string                 `json:"compatibility_profile" validate:"required"`
	NodeCompat           bool                   `json:"nodejs_compat,omitempty"`
	Env                  string                 `json:"env,omitempty"`
	Bindings             ProjectBindings        `json:"bindings"`
	Environments         map[string]Environment `json:"environments,omitempty"`
	Generated            GeneratedArtifacts     `json:"generated"`
}

type ProjectBindings struct {
	Vars    map[string]string `json:"vars,omitempty"`
	Secrets []string          `json:"secrets,omitempty"`
	AI      *AIBinding        `json:"ai,omitempty"`
	KV      []KVBinding       `json:"kv,omitempty"`
	D1      []D1Binding       `json:"d1,omitempty"`
	R2      []R2Binding       `json:"r2,omitempty"`
}

type GeneratedArtifacts struct {
	ShimSource     string `json:"shim_source,omitempty"`
	WasmExecSource string `json:"wasm_exec_source,omitempty"`
}

type Environment struct {
	Name     string          `json:"name"`
	Bindings ProjectBindings `json:"bindings"`
}

type AIBinding struct {
	Binding string `json:"binding"`
	Remote  bool   `json:"remote,omitempty"`
}

type KVBinding struct {
	Binding     string `json:"binding"`
	ID          string `json:"id,omitempty"`
	PreviewID   string `json:"preview_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Environment string `json:"environment,omitempty"`
}

type D1Binding struct {
	Binding      string `json:"binding"`
	DatabaseName string `json:"database_name,omitempty"`
	DatabaseID   string `json:"database_id,omitempty"`
}

type R2Binding struct {
	Binding      string `json:"binding"`
	BucketName   string `json:"bucket_name,omitempty"`
	Jurisdiction string `json:"jurisdiction,omitempty"`
	StorageClass string `json:"storage_class,omitempty"`
}

func DefaultProject(name, modulePath, packageName, template, compatDate, env string) Project {
	return DefaultProjectWithRuntime(name, modulePath, packageName, template, compatDate, env, RuntimeGoWasm, false)
}

func DefaultProjectWithRuntime(name, modulePath, packageName, template, compatDate, env, runtime string, nodeCompat bool) Project {
	if compatDate == "" {
		compatDate = DefaultCompatibilityDate
	}

	project := Project{
		SchemaVersion:     1,
		ProjectName:       name,
		Runtime:           normalizeRuntime(runtime),
		ModulePath:        modulePath,
		PackageName:       packageName,
		Template:          template,
		WorkerName:        name,
		WranglerConfig:    DefaultWranglerConfigFile,
		CompatibilityDate: compatDate,
		NodeCompat:        nodeCompat,
		Env:               env,
		Bindings:          ProjectBindings{Vars: map[string]string{}},
		Environments:      map[string]Environment{},
	}
	switch project.Runtime {
	case RuntimeJavaScript:
		project.Entry = filepath.ToSlash(filepath.Join("src", "worker.mjs"))
		project.Main = project.Entry
		project.OutDir = "src"
		project.CompatibilityProfile = DefaultJSCompatibilityProfile
	default:
		project.Entry = "./cmd/worker"
		project.OutDir = "dist"
		project.WasmFile = "app.wasm"
		project.ShimFile = "worker.mjs"
		project.Main = filepath.ToSlash(filepath.Join(project.OutDir, project.ShimFile))
		project.CompatibilityProfile = DefaultCompatibilityProfile
		project.Generated = GeneratedArtifacts{
			ShimSource:     "internal/generated/worker_shim.mjs",
			WasmExecSource: "internal/generated/wasm_exec.js",
		}
	}
	if UsesAIBinding(template) {
		project.Bindings.AI = &AIBinding{Binding: "AI", Remote: true}
	}
	return project
}

func (p Project) ArtifactPath() string {
	if p.WasmFile == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", p.OutDir, p.WasmFile)
}

func (p Project) EffectiveRuntime() string {
	return normalizeRuntime(p.Runtime)
}

func (p Project) MainPath() string {
	if p.Main != "" {
		return filepath.ToSlash(p.Main)
	}
	if p.ShimFile != "" {
		return filepath.ToSlash(filepath.Join(p.OutDir, p.ShimFile))
	}
	return filepath.ToSlash(p.Entry)
}

func (p Project) RequiresBuild() bool {
	return p.EffectiveRuntime() == RuntimeGoWasm
}

func UsesAIBinding(template string) bool {
	switch template {
	case "ai-text", "ai-chat", "ai-vision", "ai-stt", "ai-tts", "ai-image", "ai-embeddings":
		return true
	default:
		return false
	}
}

func normalizeRuntime(runtime string) string {
	switch runtime {
	case RuntimeJavaScript:
		return RuntimeJavaScript
	default:
		return RuntimeGoWasm
	}
}
