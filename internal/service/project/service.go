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
	ModulePath  string
	PackageName string
	Template    string
	CompatDate  string
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
	Entrypoint        string                     `json:"entrypoint"`
	WorkerName        string                     `json:"worker_name"`
	CompatibilityDate string                     `json:"compatibility_date"`
	CompatibilityMode string                     `json:"compatibility_profile"`
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

	project := config.DefaultProject(
		options.Name,
		defaultString(options.ModulePath, "github.com/paolo/"+options.Name),
		defaultString(options.PackageName, "main"),
		defaultString(options.Template, "edge-http"),
		options.CompatDate,
		options.Env,
	)

	wrangler := config.WranglerConfig{
		Name:              project.WorkerName,
		Main:              filepath.ToSlash(filepath.Join(project.OutDir, project.ShimFile)),
		CompatibilityDate: project.CompatibilityDate,
		Observability:     &config.WranglerObservability{Enabled: true},
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
		Entrypoint:        project.Entry,
		WorkerName:        project.WorkerName,
		CompatibilityDate: project.CompatibilityDate,
		CompatibilityMode: project.CompatibilityProfile,
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
