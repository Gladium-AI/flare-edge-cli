package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "flare-edge-cli",
		Short:         "Build and deploy Go/Wasm Workers with Cloudflare",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(newAuthCommand(deps))
	cmd.AddCommand(newBuildCommand(deps))
	cmd.AddCommand(newCompatCommand(deps))
	cmd.AddCommand(newProjectCommand(deps))

	return cmd
}
