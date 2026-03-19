package deploy

import (
	"context"
	"fmt"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	buildsvc "github.com/paolo/flare-edge-cli/internal/service/build"
	compatsvc "github.com/paolo/flare-edge-cli/internal/service/compat"
	"github.com/paolo/flare-edge-cli/internal/service/shared"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

type Service struct {
	store    *configstore.Store
	fs       *fs.FileSystem
	build    *buildsvc.Service
	compat   *compatsvc.Service
	wrangler *shared.WranglerExecutor
}

type Options struct {
	Dir          string
	Env          string
	Name         string
	CompatDate   string
	Route        []string
	CustomDomain []string
	WorkersDev   bool
	DryRun       bool
	UploadOnly   bool
	Message      string
	Var          []string
	KeepVars     bool
	Minify       bool
	Latest       bool
}

type Result struct {
	Compatibility compatsvc.CheckResult `json:"compatibility"`
	Build         buildsvc.WasmResult   `json:"build"`
	Deploy        shared.CommandResult  `json:"deploy"`
}

func NewService(store *configstore.Store, fs *fs.FileSystem, build *buildsvc.Service, compat *compatsvc.Service, wrangler *shared.WranglerExecutor) *Service {
	return &Service{store: store, fs: fs, build: build, compat: compat, wrangler: wrangler}
}

func (s *Service) Deploy(ctx context.Context, options Options) (Result, error) {
	compatibility, err := s.compat.Check(ctx, compatsvc.CheckOptions{Path: options.Dir, Profile: "worker-wasm", FailOn: "error"})
	if err != nil {
		return Result{}, err
	}
	if compatibility.ErrorCount > 0 {
		return Result{}, fmt.Errorf("compatibility check failed with %d error(s)", compatibility.ErrorCount)
	}

	buildResult, err := s.build.Wasm(ctx, buildsvc.WasmOptions{Path: options.Dir})
	if err != nil {
		return Result{}, err
	}

	project, wranglerCfg, err := shared.LoadProjectAndWrangler(options.Dir, s.store, s.fs)
	if err != nil {
		return Result{}, err
	}

	if options.Name != "" {
		project.WorkerName = options.Name
		wranglerCfg.Name = options.Name
	}
	if options.CompatDate != "" {
		project.CompatibilityDate = options.CompatDate
		wranglerCfg.CompatibilityDate = options.CompatDate
	}
	if len(options.Var) > 0 {
		if wranglerCfg.Vars == nil {
			wranglerCfg.Vars = map[string]string{}
		}
		for _, item := range options.Var {
			key, value, ok := splitPair(item)
			if ok {
				wranglerCfg.Vars[key] = value
				project.Bindings.Vars[key] = value
			}
		}
	}
	for _, route := range options.Route {
		wranglerCfg.Routes = config.UpsertRoute(wranglerCfg.Routes, config.WranglerRoute{Pattern: route})
	}
	for _, hostname := range options.CustomDomain {
		wranglerCfg.Routes = config.UpsertRoute(wranglerCfg.Routes, config.WranglerRoute{Pattern: hostname, CustomDomain: true})
	}

	if err := s.store.SaveProject(options.Dir, project); err != nil {
		return Result{}, err
	}
	if err := shared.SaveWrangler(options.Dir, project, wranglerCfg, s.store); err != nil {
		return Result{}, err
	}

	command := []string{"deploy"}
	if options.UploadOnly {
		command = []string{"versions", "upload"}
	}
	if options.Name != "" {
		command = append(command, "--name", options.Name)
	}
	if options.CompatDate != "" {
		command = append(command, "--compatibility-date", options.CompatDate)
	}
	for _, item := range options.Var {
		command = append(command, "--var", item)
	}
	for _, route := range options.Route {
		command = append(command, "--route", route)
	}
	if options.KeepVars {
		command = append(command, "--keep-vars")
	}
	if options.Minify {
		command = append(command, "--minify")
	}
	if options.Latest {
		command = append(command, "--latest")
	}
	if options.DryRun {
		command = append(command, "--dry-run")
	}
	if options.Message != "" {
		command = append(command, "--message", options.Message)
	}

	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Compatibility: compatibility,
		Build:         buildResult,
		Deploy:        shared.NewCommandResult(command, raw),
	}, nil
}

func splitPair(value string) (string, string, bool) {
	for index := 0; index < len(value); index++ {
		if value[index] == '=' {
			return value[:index], value[index+1:], true
		}
	}
	return "", "", false
}
