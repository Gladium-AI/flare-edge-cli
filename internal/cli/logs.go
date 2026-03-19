package cli

import (
	"context"

	"github.com/spf13/cobra"

	logssvc "github.com/paolo/flare-edge-cli/internal/service/logs"
)

func newLogsCommand(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{Use: "logs", Short: "Tail Worker logs"}
	cmd.AddCommand(newLogsTailCommand(deps))
	return cmd
}

func newLogsTailCommand(deps Dependencies) *cobra.Command {
	var options logssvc.Options
	cmd := &cobra.Command{
		Use:   "tail",
		Short: "Tail runtime logs for the Worker",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return deps.Services.Logs.Tail(context.Background(), options, cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().StringVar(&options.Worker, "worker", "", "Worker name or route")
	cmd.Flags().StringVar(&options.Format, "format", "", "Output format: pretty|json")
	cmd.Flags().StringVar(&options.Search, "search", "", "Text filter")
	cmd.Flags().StringArrayVar(&options.Status, "status", nil, "Invocation status filter")
	cmd.Flags().Float64Var(&options.Sampling, "sampling", 0, "Sampling rate between 0 and 1")
	cmd.Flags().Bool("debug", false, "Reserved for debug log verbosity")
	return cmd
}
