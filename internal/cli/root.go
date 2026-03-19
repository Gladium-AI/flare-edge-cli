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
	cmd.AddCommand(newD1Command(deps))
	cmd.AddCommand(newDeployCommand(deps))
	cmd.AddCommand(newDevCommand(deps))
	cmd.AddCommand(newDoctorCommand(deps))
	cmd.AddCommand(newKVCommand(deps))
	cmd.AddCommand(newLogsCommand(deps))
	cmd.AddCommand(newProjectCommand(deps))
	cmd.AddCommand(newR2Command(deps))
	cmd.AddCommand(newReleaseCommand(deps))
	cmd.AddCommand(newRouteCommand(deps))
	cmd.AddCommand(newSecretCommand(deps))
	cmd.AddCommand(newInitAliasCommand(deps))
	cmd.AddCommand(newInfoAliasCommand(deps))
	cmd.AddCommand(newCheckAliasCommand(deps))
	cmd.AddCommand(newTailAliasCommand(deps))
	cmd.AddCommand(newRollbackAliasCommand(deps))

	return cmd
}
