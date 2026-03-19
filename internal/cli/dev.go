package cli

import (
	"context"

	"github.com/spf13/cobra"

	devsvc "github.com/paolo/flare-edge-cli/internal/service/dev"
)

func newDevCommand(deps Dependencies) *cobra.Command {
	var options devsvc.Options

	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Run local or remote Worker development",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return deps.Services.Dev.Run(context.Background(), options, cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}

	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().IntVar(&options.Port, "port", 0, "Local port")
	cmd.Flags().BoolVar(&options.Remote, "remote", false, "Run on the remote Cloudflare network")
	cmd.Flags().BoolVar(&options.Local, "local", false, "Prefer local execution semantics")
	cmd.Flags().BoolVar(&options.Persist, "persist", false, "Persist local dev state")
	cmd.Flags().IntVar(&options.InspectorPort, "inspector-port", 0, "Inspector port")
	cmd.Flags().Bool("open", false, "Reserved for browser integration")
	cmd.Flags().Bool("watch", true, "Reserved for explicit watch control")

	return cmd
}
