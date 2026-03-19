package shared

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

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
		return envWithDefaults(nil)
	}

	state, err := w.State.Load()
	if err != nil {
		return envWithDefaults(nil)
	}

	env := make([]string, 0, 3)
	if state.APIToken != "" {
		env = append(env, "CLOUDFLARE_API_TOKEN="+state.APIToken)
	}
	if state.AccountID != "" {
		env = append(env, "CLOUDFLARE_ACCOUNT_ID="+state.AccountID)
	}
	return envWithDefaults(env)
}

func (w *WranglerExecutor) APIToken() (string, error) {
	if token := os.Getenv("CLOUDFLARE_API_TOKEN"); token != "" {
		return token, nil
	}
	if w.State != nil {
		state, err := w.State.Load()
		if err == nil && state.APIToken != "" {
			return state.APIToken, nil
		}
	}
	return loadWranglerOAuthToken()
}

func (w *WranglerExecutor) AccountID() string {
	if value := os.Getenv("CLOUDFLARE_ACCOUNT_ID"); value != "" {
		return value
	}
	if w.State == nil {
		return ""
	}
	state, err := w.State.Load()
	if err != nil {
		return ""
	}
	return state.AccountID
}

func AppendGlobalArgs(args []string, envName string) []string {
	if envName == "" {
		return args
	}

	return append(args, "--env", envName)
}

func envWithDefaults(env []string) []string {
	if hasEnvPrefix(env, "CLOUDFLARE_API_TOKEN=") || os.Getenv("CLOUDFLARE_API_TOKEN") != "" {
		return append([]string{"CI=1"}, env...)
	}
	return env
}

func loadWranglerOAuthToken() (string, error) {
	for _, path := range wranglerConfigCandidates() {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		match := regexp.MustCompile(`(?m)^oauth_token\s*=\s*"([^"]+)"\s*$`).FindSubmatch(data)
		if len(match) == 2 {
			return string(match[1]), nil
		}
	}
	return "", fmt.Errorf("cloudflare api token not found in local state or wrangler config")
}

func wranglerConfigCandidates() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	return []string{
		filepath.Join(home, "Library", "Preferences", ".wrangler", "config", "default.toml"),
		filepath.Join(home, ".wrangler", "config", "default.toml"),
	}
}

func hasEnvPrefix(env []string, prefix string) bool {
	for _, item := range env {
		if len(item) >= len(prefix) && item[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}
