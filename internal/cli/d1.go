package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	d1svc "github.com/paolo/flare-edge-cli/internal/service/d1"
)

func newD1Command(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{Use: "d1", Short: "Manage D1 databases"}
	migrations := &cobra.Command{Use: "migrations", Short: "Manage D1 migrations"}
	migrations.AddCommand(newD1MigrationsNewCommand(deps), newD1MigrationsApplyCommand(deps))
	cmd.AddCommand(newD1CreateCommand(deps), newD1ExecuteCommand(deps), migrations)
	return cmd
}

func newD1CreateCommand(deps Dependencies) *cobra.Command {
	var options d1svc.CreateOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "create <binding>",
		Short: "Create a D1 database and update bindings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.Binding = args[0]
			result, err := deps.Services.D1.Create(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Name, "name", "", "Remote database name")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}

func newD1ExecuteCommand(deps Dependencies) *cobra.Command {
	var options d1svc.ExecuteOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "execute",
		Short: "Execute SQL against D1",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.D1.Execute(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput || options.JSON).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Binding, "binding", "", "D1 binding or database name")
	cmd.Flags().StringVar(&options.SQL, "sql", "", "SQL statement")
	cmd.Flags().StringVar(&options.File, "file", "", "SQL file path")
	cmd.Flags().BoolVar(&options.Remote, "remote", false, "Execute remotely")
	cmd.Flags().BoolVar(&options.Local, "local", false, "Execute locally")
	cmd.Flags().BoolVar(&options.JSON, "json", false, "Return clean JSON")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&jsonOutput, "result-json", false, "Emit wrapper JSON")
	_ = cmd.MarkFlagRequired("binding")
	return cmd
}

func newD1MigrationsNewCommand(deps Dependencies) *cobra.Command {
	var options d1svc.MigrationNewOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "new <name>",
		Short: "Create a new local SQL migration file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.Name = args[0]
			result, err := deps.Services.D1.MigrationNew(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Path, "dir", "", "Migration directory")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}

func newD1MigrationsApplyCommand(deps Dependencies) *cobra.Command {
	var options d1svc.MigrationApplyOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply D1 migrations",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.D1.MigrationApply(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Binding, "binding", "", "D1 binding or database name")
	cmd.Flags().BoolVar(&options.Remote, "remote", false, "Apply remotely")
	cmd.Flags().BoolVar(&options.Local, "local", false, "Apply locally")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().Bool("yes", false, "Reserved for non-interactive confirmation")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	_ = cmd.MarkFlagRequired("binding")
	return cmd
}
