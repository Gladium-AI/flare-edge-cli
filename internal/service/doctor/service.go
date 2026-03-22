package doctor

import (
	"context"
	"path/filepath"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	"github.com/paolo/flare-edge-cli/internal/infra/process"
	"github.com/paolo/flare-edge-cli/internal/infra/toolchain"
	"github.com/paolo/flare-edge-cli/internal/infra/wrangler"
	authsvc "github.com/paolo/flare-edge-cli/internal/service/auth"
	buildsvc "github.com/paolo/flare-edge-cli/internal/service/build"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

type Service struct {
	store    *configstore.Store
	fs       *fs.FileSystem
	runner   process.Runner
	goTool   *toolchain.GoToolchain
	wrangler *wrangler.Client
	auth     *authsvc.StateStore
	build    *buildsvc.Service
}

type Options struct {
	Dir     string
	Verbose bool
}

type Check struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Details string `json:"details"`
}

type Result struct {
	Checks []Check `json:"checks"`
}

func NewService(store *configstore.Store, fs *fs.FileSystem, runner process.Runner, goTool *toolchain.GoToolchain, wrangler *wrangler.Client, auth *authsvc.StateStore, build *buildsvc.Service) *Service {
	return &Service{store: store, fs: fs, runner: runner, goTool: goTool, wrangler: wrangler, auth: auth, build: build}
}

func (s *Service) Run(ctx context.Context, options Options) (Result, error) {
	var checks []Check

	if version, err := s.goTool.Version(ctx, options.Dir); err == nil {
		checks = append(checks, Check{Name: "go", Status: "ok", Details: version})
	} else {
		checks = append(checks, Check{Name: "go", Status: "error", Details: err.Error()})
	}

	if result, err := s.runner.Run(ctx, process.Command{Name: "wrangler", Args: []string{"--version"}, Dir: options.Dir}); err == nil {
		checks = append(checks, Check{Name: "wrangler", Status: "ok", Details: result.Stdout})
	} else {
		checks = append(checks, Check{Name: "wrangler", Status: "error", Details: err.Error()})
	}

	state, err := s.auth.Load()
	if err == nil && state.APIToken != "" {
		checks = append(checks, Check{Name: "auth", Status: "ok", Details: "local API token present"})
	} else {
		if _, whoErr := s.wrangler.WhoAmI(ctx, options.Dir, nil); whoErr == nil {
			checks = append(checks, Check{Name: "auth", Status: "ok", Details: "Wrangler auth available"})
		} else {
			checks = append(checks, Check{Name: "auth", Status: "warning", Details: "Wrangler auth not configured"})
		}
	}

	projectPath := filepath.Join(options.Dir, config.DefaultProjectConfigFile)
	if exists, _ := s.fs.Exists(projectPath); exists {
		project, err := s.store.LoadProject(options.Dir)
		if err != nil {
			checks = append(checks, Check{Name: "project-config", Status: "error", Details: err.Error()})
		} else {
			checks = append(checks, Check{Name: "project-config", Status: "ok", Details: project.WorkerName})
			if project.CompatibilityDate == "" {
				checks = append(checks, Check{Name: "compatibility-date", Status: "error", Details: "compatibility date missing"})
			} else {
				checks = append(checks, Check{Name: "compatibility-date", Status: "ok", Details: project.CompatibilityDate})
			}
			if project.Bindings.AI != nil {
				if project.Bindings.AI.Remote {
					checks = append(checks, Check{Name: "ai-binding", Status: "ok", Details: project.Bindings.AI.Binding + " (remote)"})
				} else {
					checks = append(checks, Check{Name: "ai-binding", Status: "warning", Details: "Workers AI should use remote=true for development"})
				}
			}
			if _, err := s.build.Wasm(ctx, buildsvc.WasmOptions{Path: options.Dir, NoShim: true, OutDir: filepath.Join(options.Dir, ".doctor-dist"), Clean: true}); err != nil {
				checks = append(checks, Check{Name: "wasm-build", Status: "error", Details: err.Error()})
			} else {
				checks = append(checks, Check{Name: "wasm-build", Status: "ok", Details: "Wasm target buildable"})
			}
		}
	} else {
		checks = append(checks, Check{Name: "project-config", Status: "warning", Details: "flare-edge.json not found"})
	}

	return Result{Checks: checks}, nil
}
