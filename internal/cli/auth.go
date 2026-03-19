package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	authsvc "github.com/paolo/flare-edge-cli/internal/service/auth"
)

func newAuthCommand(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with Cloudflare",
	}

	cmd.AddCommand(newAuthLoginCommand(deps))
	cmd.AddCommand(newAuthWhoAmICommand(deps))
	cmd.AddCommand(newAuthLogoutCommand(deps))

	return cmd
}

func newAuthLoginCommand(deps Dependencies) *cobra.Command {
	var options authsvc.LoginOptions

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate the local machine against Cloudflare",
		RunE: func(cmd *cobra.Command, _ []string) error {
			options.Dir, _ = cmd.Flags().GetString("cwd")
			if options.Dir == "" {
				options.Dir = "."
			}
			result, err := deps.Services.Auth.Login(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), false).Print(result)
		},
	}

	cmd.Flags().StringVar(&options.APIToken, "api-token", "", "API token to validate or persist locally")
	cmd.Flags().StringVar(&options.AccountID, "account-id", "", "Default Cloudflare account ID")
	cmd.Flags().BoolVar(&options.UseWrangler, "wrangler", false, "Prefer Wrangler-managed OAuth login")
	cmd.Flags().BoolVar(&options.Browser, "browser", true, "Open the browser during Wrangler login")
	cmd.Flags().BoolVar(&options.Persist, "persist", false, "Persist API token metadata for future commands")
	cmd.Flags().BoolVar(&options.NonInteractive, "non-interactive", false, "Avoid interactive prompts where possible")
	cmd.Flags().String("cwd", ".", "Working directory")

	return cmd
}

func newAuthWhoAmICommand(deps Dependencies) *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Show the effective Cloudflare identity",
		RunE: func(cmd *cobra.Command, _ []string) error {
			dir, _ := cmd.Flags().GetString("cwd")
			result, err := deps.Services.Auth.WhoAmI(context.Background(), dir)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	cmd.Flags().String("cwd", ".", "Working directory")

	return cmd
}

func newAuthLogoutCommand(deps Dependencies) *cobra.Command {
	var options authsvc.LogoutOptions

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Clear local credentials and optionally Wrangler auth",
		RunE: func(cmd *cobra.Command, _ []string) error {
			options.Dir, _ = cmd.Flags().GetString("cwd")
			if err := deps.Services.Auth.Logout(context.Background(), options); err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), false).Print("logged out")
		},
	}

	cmd.Flags().BoolVar(&options.All, "all", false, "Remove both local and Wrangler-managed authentication")
	cmd.Flags().BoolVar(&options.LocalOnly, "local-only", false, "Remove only flare-edge-cli local auth state")
	cmd.Flags().String("cwd", ".", "Working directory")

	return cmd
}
