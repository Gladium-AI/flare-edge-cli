package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	teardownsvc "github.com/paolo/flare-edge-cli/internal/service/teardown"
)

func newTeardownCommand(deps Dependencies) *cobra.Command {
	var options teardownsvc.Options
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "teardown",
		Short: "Delete the Worker and associated Cloudflare resources",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.Teardown.Run(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().StringVar(&options.Name, "name", "", "Worker name override")
	cmd.Flags().BoolVar(&options.KeepBindings, "keep-bindings", false, "Keep KV, D1, and R2 resources")
	cmd.Flags().BoolVar(&options.KeepArtifacts, "keep-artifacts", false, "Keep local dist/ and .wrangler state")
	cmd.Flags().BoolVar(&options.DeleteProject, "delete-project", false, "Delete the entire local project directory after teardown")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}
