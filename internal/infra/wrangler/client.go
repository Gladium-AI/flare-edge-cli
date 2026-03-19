package wrangler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/paolo/flare-edge-cli/internal/infra/process"
)

type Client struct {
	runner process.Runner
}

type WhoAmI struct {
	Email     string `json:"email,omitempty"`
	AccountID string `json:"account_id,omitempty"`
	Raw       string `json:"raw,omitempty"`
}

func NewClient(runner process.Runner) *Client {
	return &Client{runner: runner}
}

func (c *Client) Run(ctx context.Context, dir string, args ...string) (process.Result, error) {
	return c.RunWithEnv(ctx, dir, nil, args...)
}

func (c *Client) RunWithEnv(ctx context.Context, dir string, env []string, args ...string) (process.Result, error) {
	return c.runner.Run(ctx, process.Command{
		Name: "wrangler",
		Args: args,
		Dir:  dir,
		Env:  env,
	})
}

func (c *Client) Login(ctx context.Context, dir string, env []string, args ...string) (process.Result, error) {
	return c.RunWithEnv(ctx, dir, env, append([]string{"login"}, args...)...)
}

func (c *Client) Logout(ctx context.Context, dir string, env []string, args ...string) (process.Result, error) {
	return c.RunWithEnv(ctx, dir, env, append([]string{"logout"}, args...)...)
}

func (c *Client) WhoAmI(ctx context.Context, dir string, env []string) (WhoAmI, error) {
	result, err := c.RunWithEnv(ctx, dir, env, "whoami", "--json")
	if err == nil && strings.TrimSpace(result.Stdout) != "" {
		var payload map[string]any
		if decodeErr := json.Unmarshal([]byte(result.Stdout), &payload); decodeErr == nil {
			identity := WhoAmI{Raw: strings.TrimSpace(result.Stdout)}
			if email, ok := payload["email"].(string); ok {
				identity.Email = email
			}
			if accounts, ok := payload["accounts"].([]any); ok && len(accounts) > 0 {
				if first, ok := accounts[0].(map[string]any); ok {
					if id, ok := first["id"].(string); ok {
						identity.AccountID = id
					}
				}
			}
			return identity, nil
		}
	}

	result, err = c.RunWithEnv(ctx, dir, env, "whoami")
	if err != nil {
		return WhoAmI{}, fmt.Errorf("wrangler whoami: %w", err)
	}

	return WhoAmI{Raw: strings.TrimSpace(result.Stdout)}, nil
}
