package project

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

type Service struct {
	store *configstore.Store
	fs    *fs.FileSystem
}

type InitOptions struct {
	Dir         string
	Name        string
	Runtime     string
	ModulePath  string
	PackageName string
	Template    string
	CompatDate  string
	NodeCompat  bool
	Env         string
	UseJSONC    bool
	WithGit     bool
	Yes         bool
}

type InitResult struct {
	ProjectDir    string `json:"project_dir"`
	ProjectConfig string `json:"project_config"`
	Wrangler      string `json:"wrangler_config"`
}

type InfoOptions struct {
	Dir           string
	ShowGenerated bool
	ShowBindings  bool
}

type InfoResult struct {
	Name              string                     `json:"name"`
	Runtime           string                     `json:"runtime"`
	Entrypoint        string                     `json:"entrypoint"`
	Main              string                     `json:"main"`
	WorkerName        string                     `json:"worker_name"`
	CompatibilityDate string                     `json:"compatibility_date"`
	CompatibilityMode string                     `json:"compatibility_profile"`
	NodeCompat        bool                       `json:"nodejs_compat,omitempty"`
	OutputDir         string                     `json:"output_dir"`
	WasmFile          string                     `json:"wasm_file"`
	ShimFile          string                     `json:"shim_file"`
	WranglerConfig    string                     `json:"wrangler_config"`
	Bindings          *config.ProjectBindings    `json:"bindings,omitempty"`
	Generated         *config.GeneratedArtifacts `json:"generated,omitempty"`
}

func NewService(store *configstore.Store, fs *fs.FileSystem) *Service {
	return &Service{
		store: store,
		fs:    fs,
	}
}

func (s *Service) Init(_ context.Context, options InitOptions) (InitResult, error) {
	projectDir := filepath.Join(options.Dir, options.Name)
	exists, err := s.fs.Exists(projectDir)
	if err != nil {
		return InitResult{}, err
	}
	if exists {
		return InitResult{}, fmt.Errorf("project directory already exists: %s", projectDir)
	}

	runtime := normalizeRuntime(options.Runtime)
	template, err := resolveTemplate(runtime, options.Template)
	if err != nil {
		return InitResult{}, err
	}

	project := config.DefaultProject(
		options.Name,
		defaultString(options.ModulePath, "github.com/paolo/"+options.Name),
		defaultString(options.PackageName, "main"),
		template,
		options.CompatDate,
		options.Env,
	)
	project = config.DefaultProjectWithRuntime(
		options.Name,
		project.ModulePath,
		project.PackageName,
		template,
		options.CompatDate,
		options.Env,
		runtime,
		options.NodeCompat,
	)

	wrangler := config.WranglerConfig{
		Name:              project.WorkerName,
		Main:              project.MainPath(),
		CompatibilityDate: project.CompatibilityDate,
		Observability:     &config.WranglerObservability{Enabled: true},
	}
	if project.NodeCompat {
		wrangler.CompatibilityFlags = append(wrangler.CompatibilityFlags, "nodejs_compat")
	}
	if project.Bindings.AI != nil {
		wrangler.AI = &config.WranglerAIBinding{
			Binding: project.Bindings.AI.Binding,
			Remote:  project.Bindings.AI.Remote,
		}
	}

	files := scaffoldFiles(project)
	for path, content := range files {
		if err := s.fs.WriteFile(filepath.Join(projectDir, path), []byte(content), 0o644); err != nil {
			return InitResult{}, err
		}
	}

	if options.WithGit {
		if err := s.fs.WriteFile(filepath.Join(projectDir, ".gitignore"), []byte(defaultGitignore), 0o644); err != nil {
			return InitResult{}, err
		}
	}

	if err := s.store.SaveProject(projectDir, project); err != nil {
		return InitResult{}, err
	}
	if err := s.store.SaveWrangler(projectDir, project.WranglerConfig, wrangler); err != nil {
		return InitResult{}, err
	}

	return InitResult{
		ProjectDir:    projectDir,
		ProjectConfig: filepath.Join(projectDir, config.DefaultProjectConfigFile),
		Wrangler:      filepath.Join(projectDir, config.DefaultWranglerConfigFile),
	}, nil
}

func (s *Service) Info(_ context.Context, options InfoOptions) (InfoResult, error) {
	project, err := s.store.LoadProject(options.Dir)
	if err != nil {
		return InfoResult{}, err
	}

	result := InfoResult{
		Name:              project.ProjectName,
		Runtime:           project.EffectiveRuntime(),
		Entrypoint:        project.Entry,
		Main:              project.MainPath(),
		WorkerName:        project.WorkerName,
		CompatibilityDate: project.CompatibilityDate,
		CompatibilityMode: project.CompatibilityProfile,
		NodeCompat:        project.NodeCompat,
		OutputDir:         project.OutDir,
		WasmFile:          project.WasmFile,
		ShimFile:          project.ShimFile,
		WranglerConfig:    project.WranglerConfig,
	}
	if options.ShowBindings {
		result.Bindings = &project.Bindings
	}
	if options.ShowGenerated {
		result.Generated = &project.Generated
	}

	return result, nil
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func normalizeRuntime(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "", "go", "go-wasm", "wasm":
		return config.RuntimeGoWasm
	case "js", "javascript", "js-worker", "node", "nodejs":
		return config.RuntimeJavaScript
	default:
		return value
	}
}

func resolveTemplate(runtime, template string) (string, error) {
	selected := strings.TrimSpace(template)
	if selected == "" {
		if runtime == config.RuntimeJavaScript {
			return "js-worker", nil
		}
		return "edge-http", nil
	}

	if runtime == config.RuntimeJavaScript {
		if selected != "js-worker" {
			return "", fmt.Errorf("runtime %q only supports template %q", runtime, "js-worker")
		}
		return selected, nil
	}

	if selected == "js-worker" {
		return "", fmt.Errorf("template %q requires --runtime js", selected)
	}
	return selected, nil
}
