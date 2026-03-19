package shared

import (
	"context"

	"github.com/paolo/flare-edge-cli/internal/infra/process"
	"github.com/paolo/flare-edge-cli/internal/infra/wrangler"
	"github.com/paolo/flare-edge-cli/internal/service/auth"
)

type WranglerExecutor struct {
	Client *wrangler.Client
	State  *auth.StateStore
}

func (w *WranglerExecutor) Run(ctx context.Context, dir string, envName string, args ...string) (process.Result, error) {
	args = AppendGlobalArgs(args, envName)
	return w.Client.RunWithEnv(ctx, dir, w.EnvVars(), args...)
}

func (w *WranglerExecutor) EnvVars() []string {
	if w.State == nil {
		return nil
	}

	state, err := w.State.Load()
	if err != nil {
		return nil
	}

	env := make([]string, 0, 2)
	env = append(env, "CI=1")
	if state.APIToken != "" {
		env = append(env, "CLOUDFLARE_API_TOKEN="+state.APIToken)
	}
	if state.AccountID != "" {
		env = append(env, "CLOUDFLARE_ACCOUNT_ID="+state.AccountID)
	}
	return env
}

func AppendGlobalArgs(args []string, envName string) []string {
	if envName == "" {
		return args
	}

	return append(args, "--env", envName)
}
