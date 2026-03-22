package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	aisvc "github.com/paolo/flare-edge-cli/internal/service/ai"
)

func newAICommand(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{Use: "ai", Short: "Manage Workers AI bindings"}
	binding := &cobra.Command{Use: "binding", Short: "Manage Workers AI binding config"}
	binding.AddCommand(newAIBindingSetCommand(deps), newAIBindingClearCommand(deps))
	cmd.AddCommand(binding)
	return cmd
}

func newAIBindingSetCommand(deps Dependencies) *cobra.Command {
	var options aisvc.SetBindingOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Configure the Workers AI binding",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.AI.SetBinding(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Binding, "binding", "AI", "Workers AI binding name")
	cmd.Flags().BoolVar(&options.Remote, "remote", true, "Use the remote Workers AI binding")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}

func newAIBindingClearCommand(deps Dependencies) *cobra.Command {
	var options aisvc.ClearBindingOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Remove the Workers AI binding from local config",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.AI.ClearBinding(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}
