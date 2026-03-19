package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/domain/diagnostic"
	"github.com/paolo/flare-edge-cli/internal/output"
	compatsvc "github.com/paolo/flare-edge-cli/internal/service/compat"
)

func newCompatCommand(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compat",
		Short: "Go/Wasm compatibility analysis for Workers",
	}

	cmd.AddCommand(newCompatCheckCommand(deps))
	cmd.AddCommand(newCompatRulesCommand(deps))
	return cmd
}

func newCompatCheckCommand(deps Dependencies) *cobra.Command {
	var options compatsvc.CheckOptions
	var jsonOutput bool
	var sarif bool

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Run the static Workers/Wasm compatibility checker",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.Compat.Check(context.Background(), options)
			if err != nil {
				return err
			}

			printer := output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput || sarif)
			if sarif {
				return printer.Print(diagnostic.SARIF("flare-edge-cli compat", result.Diagnostics))
			}
			if jsonOutput {
				return printer.Print(result)
			}
			if err := printer.PrintDiagnostics(result.Diagnostics); err != nil {
				return err
			}
			return printer.Print(fmt.Sprintf("errors=%d warnings=%d", result.ErrorCount, result.WarnCount))
		},
	}

	cmd.Flags().StringVar(&options.Path, "path", ".", "Project path to analyze")
	cmd.Flags().StringVar(&options.Entry, "entry", "", "Package or file entry selector")
	cmd.Flags().StringVar(&options.Profile, "profile", "worker-wasm", "Compatibility rule profile")
	cmd.Flags().BoolVar(&options.Strict, "strict", false, "Enable stricter compatibility checks")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	cmd.Flags().BoolVar(&sarif, "sarif", false, "Emit SARIF JSON")
	cmd.Flags().StringVar(&options.FailOn, "fail-on", "error", "Threshold: warning|error")
	cmd.Flags().StringArrayVar(&options.Exclude, "exclude", nil, "Glob or substring to exclude from analysis")

	return cmd
}

func newCompatRulesCommand(deps Dependencies) *cobra.Command {
	var jsonOutput bool
	var severity string

	cmd := &cobra.Command{
		Use:   "rules",
		Short: "List built-in compatibility rules",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(deps.Services.Compat.Rules(severity))
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	cmd.Flags().StringVar(&severity, "severity", "", "Filter rules by severity: error|warning|info")

	return cmd
}
