package release

import (
	"context"

	"github.com/paolo/flare-edge-cli/internal/service/shared"
)

type Service struct {
	wrangler *shared.WranglerExecutor
}

type ListOptions struct {
	Dir   string
	Env   string
	Name  string
	Limit int
}

type PromoteOptions struct {
	Dir       string
	Env       string
	VersionID string
	Message   string
	Yes       bool
}

type RollbackOptions struct {
	Dir       string
	Env       string
	VersionID string
	Yes       bool
	Message   string
}

func NewService(wrangler *shared.WranglerExecutor) *Service {
	return &Service{wrangler: wrangler}
}

func (s *Service) List(ctx context.Context, options ListOptions) (shared.CommandResult, error) {
	command := []string{"versions", "list", "--json"}
	if options.Name != "" {
		command = append(command, "--name", options.Name)
	}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}

func (s *Service) Promote(ctx context.Context, options PromoteOptions) (shared.CommandResult, error) {
	command := []string{"versions", "deploy", "--version-id", options.VersionID, "--percentage", "100"}
	if options.Message != "" {
		command = append(command, "--message", options.Message)
	}
	if options.Yes {
		command = append(command, "--yes")
	}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}

func (s *Service) Rollback(ctx context.Context, options RollbackOptions) (shared.CommandResult, error) {
	command := []string{"rollback", options.VersionID}
	if options.Message != "" {
		command = append(command, "--message", options.Message)
	}
	if options.Yes {
		command = append(command, "--yes")
	}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}
