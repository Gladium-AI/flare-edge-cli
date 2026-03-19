package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
	baseArgs := make([]string, 0, 16)
	if options.Name != "" {
		baseArgs = append(baseArgs, "--name", options.Name)
	}
	if options.CompatDate != "" {
		baseArgs = append(baseArgs, "--compatibility-date", options.CompatDate)
	}
	for _, item := range options.Var {
		baseArgs = append(baseArgs, "--var", item)
	}
	for _, route := range options.Route {
		baseArgs = append(baseArgs, "--route", route)
	}
	if options.KeepVars {
		baseArgs = append(baseArgs, "--keep-vars")
	}
	if options.Minify {
		baseArgs = append(baseArgs, "--minify")
	}
	if options.Latest {
		baseArgs = append(baseArgs, "--latest")
	}
	if options.DryRun {
		baseArgs = append(baseArgs, "--dry-run")
	}

	if options.UploadOnly || options.Message != "" {
		uploadArgs := append([]string{"versions", "upload"}, baseArgs...)
		if options.Message != "" {
			uploadArgs = append(uploadArgs, "--message", options.Message)
		}
		uploadRaw, err := s.wrangler.Run(ctx, options.Dir, options.Env, uploadArgs...)
		if err != nil {
			if canFallbackToPlainDeploy(options, err) {
				return s.plainDeploy(ctx, options, compatibility, buildResult, baseArgs)
			}
			return Result{}, err
		}
		if options.UploadOnly {
			return Result{
				Compatibility: compatibility,
				Build:         buildResult,
				Deploy:        shared.NewCommandResult(uploadArgs, uploadRaw),
			}, nil
		}

		versionID, err := s.latestVersionID(ctx, options.Dir, options.Env)
		if err != nil {
			if canFallbackToPlainDeploy(options, err) {
				return s.plainDeploy(ctx, options, compatibility, buildResult, baseArgs)
			}
			return Result{}, err
		}
		command = []string{"versions", "deploy", "--version-id", versionID, "--percentage", "100"}
		if options.Message != "" {
			command = append(command, "--message", options.Message)
		}
		raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
		if err != nil {
			if canFallbackToPlainDeploy(options, err) {
				return s.plainDeploy(ctx, options, compatibility, buildResult, baseArgs)
			}
			return Result{}, err
		}
		return Result{
			Compatibility: compatibility,
			Build:         buildResult,
			Deploy:        shared.NewCommandResult(command, raw),
		}, nil
	}

	command = append(command, baseArgs...)
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

func (s *Service) plainDeploy(ctx context.Context, options Options, compatibility compatsvc.CheckResult, buildResult buildsvc.WasmResult, baseArgs []string) (Result, error) {
	command := append([]string{"deploy"}, baseArgs...)
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

func (s *Service) latestVersionID(ctx context.Context, dir, env string) (string, error) {
	raw, err := s.wrangler.Run(ctx, dir, env, "versions", "list", "--json")
	if err != nil {
		return "", err
	}

	var versions []struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(raw.Stdout), &versions); err != nil {
		return "", fmt.Errorf("decode versions list: %w", err)
	}
	if len(versions) == 0 {
		return "", fmt.Errorf("no worker versions available after upload")
	}
	return versions[len(versions)-1].ID, nil
}

func splitPair(value string) (string, string, bool) {
	for index := 0; index < len(value); index++ {
		if value[index] == '=' {
			return value[:index], value[index+1:], true
		}
	}
	return "", "", false
}

func canFallbackToPlainDeploy(options Options, err error) bool {
	if options.UploadOnly || options.Message == "" {
		return false
	}
	message := err.Error()
	return strings.Contains(message, "This Worker does not exist on your account") || strings.Contains(message, "TTY initialization failed")
}
