package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	secretsvc "github.com/paolo/flare-edge-cli/internal/service/secret"
)

func newSecretCommand(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Manage Worker secrets",
	}
	cmd.AddCommand(newSecretPutCommand(deps))
	cmd.AddCommand(newSecretListCommand(deps))
	cmd.AddCommand(newSecretDeleteCommand(deps))
	return cmd
}

func newSecretPutCommand(deps Dependencies) *cobra.Command {
	var options secretsvc.PutOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "put <KEY>",
		Short: "Create or update a Worker secret",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.Key = args[0]
			result, err := deps.Services.Secret.Put(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Value, "value", "", "Secret value")
	cmd.Flags().StringVar(&options.FromFile, "from-file", "", "Load secret value from file")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&options.Versioned, "versioned", false, "Use versioned secrets workflow")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}

func newSecretListCommand(deps Dependencies) *cobra.Command {
	var options secretsvc.ListOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Worker secrets",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.Secret.List(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}

func newSecretDeleteCommand(deps Dependencies) *cobra.Command {
	var options secretsvc.DeleteOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "delete <KEY>",
		Short: "Delete a Worker secret",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.Key = args[0]
			result, err := deps.Services.Secret.Delete(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&options.Versioned, "versioned", false, "Use versioned secrets workflow")
	cmd.Flags().Bool("yes", false, "Reserved for non-interactive confirmation")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}
