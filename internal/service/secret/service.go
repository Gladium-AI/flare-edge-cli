package secret

import (
	"context"
	"fmt"

	"github.com/paolo/flare-edge-cli/internal/infra/process"
	"github.com/paolo/flare-edge-cli/internal/service/shared"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

type Service struct {
	runner   process.Runner
	fs       *fs.FileSystem
	wrangler *shared.WranglerExecutor
}

type PutOptions struct {
	Dir       string
	Key       string
	Value     string
	FromFile  string
	Env       string
	Versioned bool
}

type ListOptions struct {
	Dir string
	Env string
}

type DeleteOptions struct {
	Dir       string
	Key       string
	Env       string
	Versioned bool
}

func NewService(runner process.Runner, fs *fs.FileSystem, wrangler *shared.WranglerExecutor) *Service {
	return &Service{runner: runner, fs: fs, wrangler: wrangler}
}

func (s *Service) Put(ctx context.Context, options PutOptions) (shared.CommandResult, error) {
	payload := options.Value
	if options.FromFile != "" {
		data, err := s.fs.ReadFile(options.FromFile)
		if err != nil {
			return shared.CommandResult{}, err
		}
		payload = string(data)
	}
	if payload == "" {
		return shared.CommandResult{}, fmt.Errorf("secret value is required")
	}

	command := []string{"secret", "put", options.Key}
	if options.Versioned {
		command = []string{"versions", "secret", "put", options.Key}
	}
	args := shared.AppendGlobalArgs(command, options.Env)
	run, err := s.runner.Run(ctx, process.Command{
		Name:  "wrangler",
		Args:  args,
		Dir:   options.Dir,
		Env:   s.wrangler.EnvVars(),
		Stdin: payload,
	})
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(args, run), nil
}

func (s *Service) List(ctx context.Context, options ListOptions) (shared.CommandResult, error) {
	command := []string{"secret", "list"}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}

func (s *Service) Delete(ctx context.Context, options DeleteOptions) (shared.CommandResult, error) {
	command := []string{"secret", "delete", options.Key}
	if options.Versioned {
		command = []string{"versions", "secret", "delete", options.Key}
	}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}
