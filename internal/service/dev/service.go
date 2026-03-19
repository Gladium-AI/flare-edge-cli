package dev

import (
	"context"
	"io"
	"strconv"

	"github.com/paolo/flare-edge-cli/internal/infra/process"
	"github.com/paolo/flare-edge-cli/internal/service/shared"
)

type Service struct {
	runner   process.Runner
	wrangler *shared.WranglerExecutor
}

type Options struct {
	Dir           string
	Env           string
	Port          int
	Remote        bool
	Local         bool
	Persist       bool
	InspectorPort int
}

func NewService(runner process.Runner, wrangler *shared.WranglerExecutor) *Service {
	return &Service{runner: runner, wrangler: wrangler}
}

func (s *Service) Run(ctx context.Context, options Options, stdout, stderr io.Writer) error {
	args := []string{"dev"}
	if options.Port > 0 {
		args = append(args, "--port", strconv.Itoa(options.Port))
	}
	if options.Remote {
		args = append(args, "--remote")
	}
	if options.Local {
		args = append(args, "--local-protocol", "http")
	}
	if options.Persist {
		args = append(args, "--persist-to", ".wrangler/state")
	}
	if options.InspectorPort > 0 {
		args = append(args, "--inspector-port", strconv.Itoa(options.InspectorPort))
	}
	args = shared.AppendGlobalArgs(args, options.Env)
	return s.runner.Stream(ctx, process.Command{
		Name: "wrangler",
		Args: args,
		Dir:  options.Dir,
		Env:  s.wrangler.EnvVars(),
	}, stdout, stderr)
}
