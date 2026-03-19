package d1

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	"github.com/paolo/flare-edge-cli/internal/service/shared"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

type Service struct {
	store    *configstore.Store
	fs       *fs.FileSystem
	wrangler *shared.WranglerExecutor
}

type CreateOptions struct {
	Dir     string
	Binding string
	Name    string
	Env     string
}

type ExecuteOptions struct {
	Dir     string
	Binding string
	SQL     string
	File    string
	Remote  bool
	Local   bool
	JSON    bool
	Env     string
}

type MigrationNewOptions struct {
	Dir  string
	Name string
	Path string
}

type MigrationApplyOptions struct {
	Dir     string
	Binding string
	Remote  bool
	Local   bool
	Env     string
}

type CreateResult struct {
	Command shared.CommandResult `json:"command"`
	Binding string               `json:"binding"`
	ID      string               `json:"id,omitempty"`
}

func NewService(store *configstore.Store, fs *fs.FileSystem, wrangler *shared.WranglerExecutor) *Service {
	return &Service{store: store, fs: fs, wrangler: wrangler}
}

func (s *Service) Create(ctx context.Context, options CreateOptions) (CreateResult, error) {
	name := options.Name
	if name == "" {
		name = options.Binding
	}
	command := []string{"d1", "create", name}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return CreateResult{}, err
	}
	id := shared.ExtractID(raw.Stdout)
	project, wranglerCfg, err := shared.LoadProjectAndWrangler(options.Dir, s.store, s.fs)
	if err != nil {
		return CreateResult{}, err
	}
	project.Bindings.D1 = append(project.Bindings.D1, config.D1Binding{Binding: options.Binding, DatabaseName: name, DatabaseID: id})
	wranglerCfg.D1Databases = config.UpsertD1(wranglerCfg.D1Databases, config.WranglerD1Database{Binding: options.Binding, DatabaseName: name, DatabaseID: id})
	if err := s.store.SaveProject(options.Dir, project); err != nil {
		return CreateResult{}, err
	}
	if err := shared.SaveWrangler(options.Dir, project, wranglerCfg, s.store); err != nil {
		return CreateResult{}, err
	}
	return CreateResult{Command: shared.NewCommandResult(command, raw), Binding: options.Binding, ID: id}, nil
}

func (s *Service) Execute(ctx context.Context, options ExecuteOptions) (shared.CommandResult, error) {
	command := []string{"d1", "execute", options.Binding}
	if options.SQL != "" {
		command = append(command, "--command", options.SQL)
	}
	if options.File != "" {
		command = append(command, "--file", options.File)
	}
	if options.Remote {
		command = append(command, "--remote")
	}
	if options.Local {
		command = append(command, "--local")
	}
	if options.JSON {
		command = append(command, "--json")
	}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}

func (s *Service) MigrationNew(_ context.Context, options MigrationNewOptions) (map[string]string, error) {
	dir := options.Path
	if dir == "" {
		dir = filepath.Join(options.Dir, "migrations")
	}
	filename := fmt.Sprintf("%s_%s.sql", time.Now().Format("20060102150405"), slug(options.Name))
	fullPath := filepath.Join(dir, filename)
	if err := s.fs.WriteFile(fullPath, []byte("-- Write migration here.\n"), 0o644); err != nil {
		return nil, err
	}
	return map[string]string{"path": fullPath}, nil
}

func (s *Service) MigrationApply(ctx context.Context, options MigrationApplyOptions) (shared.CommandResult, error) {
	command := []string{"d1", "migrations", "apply", options.Binding}
	if options.Remote {
		command = append(command, "--remote")
	}
	if options.Local {
		command = append(command, "--local")
	}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}

func slug(value string) string {
	pattern := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return pattern.ReplaceAllString(value, "_")
}
