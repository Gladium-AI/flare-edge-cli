package auth

import (
	"context"
	"fmt"

	"github.com/paolo/flare-edge-cli/internal/infra/wrangler"
)

type Service struct {
	client *wrangler.Client
	state  *StateStore
}

type LoginOptions struct {
	Dir            string
	UseWrangler    bool
	Browser        bool
	APIToken       string
	AccountID      string
	Persist        bool
	NonInteractive bool
}

type LoginResult struct {
	Method    string `json:"method"`
	Persisted bool   `json:"persisted"`
	AccountID string `json:"account_id,omitempty"`
	Validated bool   `json:"validated"`
	Wrangler  bool   `json:"wrangler"`
}

type WhoAmIResult struct {
	Email          string `json:"email,omitempty"`
	AccountID      string `json:"account_id,omitempty"`
	WranglerHealth string `json:"wrangler_health"`
	LocalToken     bool   `json:"local_token"`
	Raw            string `json:"raw,omitempty"`
}

type LogoutOptions struct {
	Dir       string
	All       bool
	LocalOnly bool
}

func NewService(client *wrangler.Client, state *StateStore) *Service {
	return &Service{
		client: client,
		state:  state,
	}
}

func (s *Service) Login(ctx context.Context, options LoginOptions) (LoginResult, error) {
	if options.APIToken != "" {
		if options.Persist {
			if err := s.state.Save(State{
				APIToken:  options.APIToken,
				AccountID: options.AccountID,
			}); err != nil {
				return LoginResult{}, err
			}
		}

		identity, err := s.client.WhoAmI(ctx, options.Dir, authEnv(State{
			APIToken:  options.APIToken,
			AccountID: options.AccountID,
		}))
		validated := err == nil
		return LoginResult{
			Method:    "api-token",
			Persisted: options.Persist,
			AccountID: firstNonEmpty(identity.AccountID, options.AccountID),
			Validated: validated,
			Wrangler:  false,
		}, nil
	}

	args := make([]string, 0, 2)
	if !options.Browser {
		args = append(args, "--browser=false")
	}

	if _, err := s.client.Login(ctx, options.Dir, authEnvFromStore(s.state), args...); err != nil {
		return LoginResult{}, fmt.Errorf("wrangler login: %w", err)
	}

	return LoginResult{
		Method:    "wrangler-oauth",
		Persisted: true,
		Validated: true,
		Wrangler:  true,
	}, nil
}

func (s *Service) WhoAmI(ctx context.Context, dir string) (WhoAmIResult, error) {
	state, err := s.state.Load()
	if err != nil {
		return WhoAmIResult{}, err
	}

	identity, err := s.client.WhoAmI(ctx, dir, authEnv(state))
	if err != nil {
		return WhoAmIResult{
			AccountID:      state.AccountID,
			WranglerHealth: "error",
			LocalToken:     state.APIToken != "",
		}, nil
	}

	return WhoAmIResult{
		Email:          identity.Email,
		AccountID:      firstNonEmpty(identity.AccountID, state.AccountID),
		WranglerHealth: "ok",
		LocalToken:     state.APIToken != "",
		Raw:            identity.Raw,
	}, nil
}

func (s *Service) Logout(ctx context.Context, options LogoutOptions) error {
	if !options.LocalOnly {
		if _, err := s.client.Logout(ctx, options.Dir, authEnvFromStore(s.state)); err != nil && !options.All {
			return fmt.Errorf("wrangler logout: %w", err)
		}
	}

	return s.state.Delete()
}

func authEnvFromStore(store *StateStore) []string {
	state, err := store.Load()
	if err != nil {
		return nil
	}
	return authEnv(state)
}

func authEnv(state State) []string {
	env := make([]string, 0, 2)
	if state.APIToken != "" {
		env = append(env, "CLOUDFLARE_API_TOKEN="+state.APIToken)
	}
	if state.AccountID != "" {
		env = append(env, "CLOUDFLARE_ACCOUNT_ID="+state.AccountID)
	}
	return env
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
