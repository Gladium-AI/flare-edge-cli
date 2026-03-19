package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	releasesvc "github.com/paolo/flare-edge-cli/internal/service/release"
)

func newReleaseCommand(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{Use: "release", Short: "Manage Worker releases"}
	cmd.AddCommand(newReleaseListCommand(deps), newReleasePromoteCommand(deps), newReleaseRollbackCommand(deps))
	return cmd
}

func newReleaseListCommand(deps Dependencies) *cobra.Command {
	var options releasesvc.ListOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List deployed Worker versions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.Release.List(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().StringVar(&options.Name, "name", "", "Worker name")
	cmd.Flags().IntVar(&options.Limit, "limit", 10, "Reserved for local post-filtering")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}

func newReleasePromoteCommand(deps Dependencies) *cobra.Command {
	var options releasesvc.PromoteOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "promote <version_id>",
		Short: "Promote a previously uploaded version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.VersionID = args[0]
			result, err := deps.Services.Release.Promote(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().StringVar(&options.Message, "message", "", "Deployment message")
	cmd.Flags().BoolVar(&options.Yes, "yes", false, "Accept defaults without prompting")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}

func newReleaseRollbackCommand(deps Dependencies) *cobra.Command {
	var options releasesvc.RollbackOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "rollback <version_id>",
		Short: "Roll back to a previous Worker version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.VersionID = args[0]
			result, err := deps.Services.Release.Rollback(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&options.Yes, "yes", false, "Accept defaults without prompting")
	cmd.Flags().StringVar(&options.Message, "message", "", "Rollback message")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}
