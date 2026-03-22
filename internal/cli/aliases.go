package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	compatsvc "github.com/paolo/flare-edge-cli/internal/service/compat"
	logssvc "github.com/paolo/flare-edge-cli/internal/service/logs"
	projectsvc "github.com/paolo/flare-edge-cli/internal/service/project"
	releasesvc "github.com/paolo/flare-edge-cli/internal/service/release"
)

func newInitAliasCommand(deps Dependencies) *cobra.Command {
	var options projectsvc.InitOptions
	cmd := &cobra.Command{
		Use:   "init <name>",
		Short: "Alias for project init",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.Name = args[0]
			result, err := deps.Services.Project.Init(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), false).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "cwd", ".", "Target parent directory")
	cmd.Flags().StringVar(&options.ModulePath, "module-path", "", "Go module path")
	cmd.Flags().StringVar(&options.PackageName, "package", "", "Package name")
	cmd.Flags().StringVar(&options.Template, "template", "edge-http", "Starter template: edge-http|edge-json|scheduled|kv-api|d1-api|r2-api|ai-text")
	cmd.Flags().StringVar(&options.CompatDate, "compat-date", "", "Compatibility date")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&options.UseJSONC, "use-jsonc", false, "Generate wrangler.jsonc")
	cmd.Flags().BoolVar(&options.WithGit, "with-git", true, "Generate .gitignore")
	cmd.Flags().BoolVar(&options.Yes, "yes", false, "Accept defaults without prompting")
	return cmd
}

func newInfoAliasCommand(deps Dependencies) *cobra.Command {
	var options projectsvc.InfoOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Alias for project info",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.Project.Info(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "cwd", ".", "Project directory")
	cmd.Flags().BoolVar(&options.ShowGenerated, "show-generated", false, "Include generated metadata")
	cmd.Flags().BoolVar(&options.ShowBindings, "show-bindings", false, "Include bindings")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}

func newCheckAliasCommand(deps Dependencies) *cobra.Command {
	var options compatsvc.CheckOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Alias for compat check",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.Compat.Check(context.Background(), options)
			if err != nil {
				return err
			}
			if jsonOutput {
				return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), true).Print(result)
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), false).PrintDiagnostics(result.Diagnostics)
		},
	}
	cmd.Flags().StringVar(&options.Path, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Entry, "entry", "", "Package or file entry selector")
	cmd.Flags().StringVar(&options.Profile, "profile", "worker-wasm", "Profile name")
	cmd.Flags().StringVar(&options.FailOn, "fail-on", "error", "Threshold")
	cmd.Flags().StringArrayVar(&options.Exclude, "exclude", nil, "Exclude globs")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}

func newTailAliasCommand(deps Dependencies) *cobra.Command {
	var options logssvc.Options
	cmd := &cobra.Command{
		Use:   "tail",
		Short: "Alias for logs tail",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return deps.Services.Logs.Tail(context.Background(), options, cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().StringVar(&options.Worker, "worker", "", "Worker name or route")
	cmd.Flags().StringVar(&options.Format, "format", "", "pretty|json")
	cmd.Flags().StringVar(&options.Search, "search", "", "Text filter")
	return cmd
}

func newRollbackAliasCommand(deps Dependencies) *cobra.Command {
	var options releasesvc.RollbackOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "rollback <version_id>",
		Short: "Alias for release rollback",
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
