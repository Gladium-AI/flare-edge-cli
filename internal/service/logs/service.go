package logs

import (
	"context"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/paolo/flare-edge-cli/internal/infra/process"
	"github.com/paolo/flare-edge-cli/internal/service/shared"
)

type Service struct {
	runner   process.Runner
	wrangler *shared.WranglerExecutor
}

type Options struct {
	Dir      string
	Env      string
	Worker   string
	Format   string
	Search   string
	Status   []string
	Sampling float64
}

func NewService(runner process.Runner, wrangler *shared.WranglerExecutor) *Service {
	return &Service{runner: runner, wrangler: wrangler}
}

func (s *Service) Tail(ctx context.Context, options Options, stdout, stderr io.Writer) error {
	args := []string{"tail"}
	if target := normalizeWorkerTarget(options.Worker); target != "" {
		args = append(args, target)
	}
	if options.Format != "" {
		args = append(args, "--format", options.Format)
	}
	if options.Search != "" {
		args = append(args, "--search", options.Search)
	}
	for _, status := range options.Status {
		args = append(args, "--status", status)
	}
	if options.Sampling > 0 {
		args = append(args, "--sampling-rate", strconv.FormatFloat(options.Sampling, 'f', -1, 64))
	}
	args = shared.AppendGlobalArgs(args, options.Env)
	return s.runner.Stream(ctx, process.Command{
		Name: "wrangler",
		Args: args,
		Dir:  options.Dir,
		Env:  s.wrangler.EnvVars(),
	}, stdout, stderr)
}

func normalizeWorkerTarget(value string) string {
	if value == "" {
		return ""
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return value
	}
	host := parsed.Hostname()
	if strings.HasSuffix(host, ".workers.dev") {
		labels := strings.Split(host, ".")
		if len(labels) > 0 && labels[0] != "" {
			return labels[0]
		}
	}
	path := strings.TrimSuffix(parsed.EscapedPath(), "/")
	if path == "" {
		return host
	}
	return host + path
}
