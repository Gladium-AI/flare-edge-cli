package kv

import (
	"context"
	"strconv"

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

type NamespaceCreateOptions struct {
	Dir       string
	Binding   string
	Title     string
	Env       string
	Provision bool
}

type KeyPutOptions struct {
	Dir        string
	Binding    string
	Key        string
	Value      string
	FromFile   string
	TTL        int
	Expiration int
	Metadata   string
	Env        string
}

type KeyGetOptions struct {
	Dir     string
	Binding string
	Key     string
	Text    bool
	Env     string
}

type KeyListOptions struct {
	Dir     string
	Binding string
	Prefix  string
	Env     string
}

type KeyDeleteOptions struct {
	Dir     string
	Binding string
	Key     string
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

func (s *Service) NamespaceCreate(ctx context.Context, options NamespaceCreateOptions) (CreateResult, error) {
	title := options.Title
	if title == "" {
		title = options.Binding
	}
	command := []string{"kv", "namespace", "create", title}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return CreateResult{}, err
	}
	id := shared.ExtractID(raw.Stdout)
	if options.Provision {
		project, wranglerCfg, err := shared.LoadProjectAndWrangler(options.Dir, s.store, s.fs)
		if err != nil {
			return CreateResult{}, err
		}
		project.Bindings.KV = append(project.Bindings.KV, config.KVBinding{Binding: options.Binding, ID: id, Title: title})
		wranglerCfg.KVNamespaces = config.UpsertKV(wranglerCfg.KVNamespaces, config.WranglerKVNamespace{Binding: options.Binding, ID: id})
		if err := s.store.SaveProject(options.Dir, project); err != nil {
			return CreateResult{}, err
		}
		if err := shared.SaveWrangler(options.Dir, project, wranglerCfg, s.store); err != nil {
			return CreateResult{}, err
		}
	}
	return CreateResult{Command: shared.NewCommandResult(command, raw), Binding: options.Binding, ID: id}, nil
}

func (s *Service) Put(ctx context.Context, options KeyPutOptions) (shared.CommandResult, error) {
	command := []string{"kv", "key", "put", options.Key}
	if options.Value != "" {
		command = append(command, options.Value)
	}
	command = append(command, "--binding", options.Binding)
	if options.FromFile != "" {
		command = append(command, "--path", options.FromFile)
	}
	if options.TTL > 0 {
		command = append(command, "--ttl", itoa(options.TTL))
	}
	if options.Expiration > 0 {
		command = append(command, "--expiration", itoa(options.Expiration))
	}
	if options.Metadata != "" {
		command = append(command, "--metadata", options.Metadata)
	}
	command = append(command, "--remote")
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}

func (s *Service) Get(ctx context.Context, options KeyGetOptions) (shared.CommandResult, error) {
	command := []string{"kv", "key", "get", options.Key, "--binding", options.Binding}
	if options.Text {
		command = append(command, "--text")
	}
	command = append(command, "--remote")
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}

func (s *Service) List(ctx context.Context, options KeyListOptions) (shared.CommandResult, error) {
	command := []string{"kv", "key", "list", "--binding", options.Binding}
	if options.Prefix != "" {
		command = append(command, "--prefix", options.Prefix)
	}
	command = append(command, "--remote")
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}

func (s *Service) Delete(ctx context.Context, options KeyDeleteOptions) (shared.CommandResult, error) {
	command := []string{"kv", "key", "delete", options.Key, "--binding", options.Binding, "--remote"}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}

func itoa(value int) string {
	return strconv.Itoa(value)
}
