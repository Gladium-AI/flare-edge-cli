package config

import "fmt"

const (
	DefaultProjectConfigFile    = "flare-edge.json"
	DefaultWranglerConfigFile   = "wrangler.jsonc"
	DefaultCompatibilityDate    = "2026-03-19"
	DefaultCompatibilityProfile = "worker-wasm"
)

type Project struct {
	SchemaVersion        int                    `json:"schema_version" validate:"required,min=1"`
	ProjectName          string                 `json:"project_name" validate:"required"`
	ModulePath           string                 `json:"module_path" validate:"required"`
	PackageName          string                 `json:"package_name" validate:"required"`
	Template             string                 `json:"template" validate:"required"`
	Entry                string                 `json:"entry" validate:"required"`
	OutDir               string                 `json:"out_dir" validate:"required"`
	WasmFile             string                 `json:"wasm_file" validate:"required"`
	ShimFile             string                 `json:"shim_file" validate:"required"`
	WorkerName           string                 `json:"worker_name" validate:"required"`
	WranglerConfig       string                 `json:"wrangler_config" validate:"required"`
	CompatibilityDate    string                 `json:"compatibility_date" validate:"required"`
	CompatibilityProfile string                 `json:"compatibility_profile" validate:"required"`
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
	ShimSource     string `json:"shim_source" validate:"required"`
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
	if compatDate == "" {
		compatDate = DefaultCompatibilityDate
	}

	entry := "./cmd/worker"
	project := Project{
		SchemaVersion:        1,
		ProjectName:          name,
		ModulePath:           modulePath,
		PackageName:          packageName,
		Template:             template,
		Entry:                entry,
		OutDir:               "dist",
		WasmFile:             "app.wasm",
		ShimFile:             "worker.mjs",
		WorkerName:           name,
		WranglerConfig:       DefaultWranglerConfigFile,
		CompatibilityDate:    compatDate,
		CompatibilityProfile: DefaultCompatibilityProfile,
		Env:                  env,
		Bindings:             ProjectBindings{Vars: map[string]string{}},
		Environments:         map[string]Environment{},
		Generated: GeneratedArtifacts{
			ShimSource:     "internal/generated/worker_shim.mjs",
			WasmExecSource: "internal/generated/wasm_exec.js",
		},
	}
	if template == "ai-text" {
		project.Bindings.AI = &AIBinding{Binding: "AI", Remote: true}
	}
	return project
}

func (p Project) ArtifactPath() string {
	return fmt.Sprintf("%s/%s", p.OutDir, p.WasmFile)
}
