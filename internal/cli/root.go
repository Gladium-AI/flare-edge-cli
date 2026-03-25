package cli

import (
	"os"

	"github.com/spf13/cobra"
)

func NewRootCommand(deps Dependencies) *cobra.Command {
	var accountID string

	cmd := &cobra.Command{
		Use:           "flare-edge-cli",
		Short:         "Build and deploy Cloudflare Workers",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if accountID != "" {
				return os.Setenv("CLOUDFLARE_ACCOUNT_ID", accountID)
			}
			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&accountID, "account-id", "", "Cloudflare account ID override for all commands")

	cmd.AddCommand(newAICommand(deps))
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
	cmd.AddCommand(newTeardownCommand(deps))
	cmd.AddCommand(newInitAliasCommand(deps))
	cmd.AddCommand(newInfoAliasCommand(deps))
	cmd.AddCommand(newCheckAliasCommand(deps))
	cmd.AddCommand(newTailAliasCommand(deps))
	cmd.AddCommand(newRollbackAliasCommand(deps))

	return cmd
}
