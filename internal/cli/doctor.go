package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	doctorsvc "github.com/paolo/flare-edge-cli/internal/service/doctor"
)

func newDoctorCommand(deps Dependencies) *cobra.Command {
	var options doctorsvc.Options
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Run environment diagnostics",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.Doctor.Run(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	cmd.Flags().Bool("fix", false, "Reserved for future automated fixes")
	cmd.Flags().BoolVar(&options.Verbose, "verbose", false, "Include more diagnostic detail")
	return cmd
}
