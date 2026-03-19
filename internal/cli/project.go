package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	projectsvc "github.com/paolo/flare-edge-cli/internal/service/project"
)

func newProjectCommand(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Project scaffolding and inspection",
	}

	cmd.AddCommand(newProjectInitCommand(deps))
	cmd.AddCommand(newProjectInfoCommand(deps))

	return cmd
}

func newProjectInitCommand(deps Dependencies) *cobra.Command {
	var options projectsvc.InitOptions

	cmd := &cobra.Command{
		Use:   "init <name>",
		Short: "Scaffold a Workers-ready Go/Wasm project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.Name = args[0]
			options.Dir, _ = cmd.Flags().GetString("cwd")
			result, err := deps.Services.Project.Init(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), false).Print(result)
		},
	}

	cmd.Flags().StringVar(&options.ModulePath, "module-path", "", "Go module path to initialize")
	cmd.Flags().StringVar(&options.PackageName, "package", "", "Package name for generated Go entrypoint")
	cmd.Flags().StringVar(&options.Template, "template", "edge-http", "Starter template: edge-http|edge-json|scheduled|kv-api|d1-api|r2-api")
	cmd.Flags().StringVar(&options.CompatDate, "compat-date", "", "Cloudflare compatibility date")
	cmd.Flags().StringVar(&options.Env, "env", "", "Default Wrangler environment")
	cmd.Flags().BoolVar(&options.UseJSONC, "use-jsonc", false, "Generate wrangler.jsonc output")
	cmd.Flags().BoolVar(&options.WithGit, "with-git", true, "Generate .gitignore in the scaffolded project")
	cmd.Flags().BoolVar(&options.Yes, "yes", false, "Accept defaults without prompting")
	cmd.Flags().String("cwd", ".", "Target parent directory")

	return cmd
}

func newProjectInfoCommand(deps Dependencies) *cobra.Command {
	var options projectsvc.InfoOptions
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show project metadata inferred from local config",
		RunE: func(cmd *cobra.Command, _ []string) error {
			options.Dir, _ = cmd.Flags().GetString("cwd")
			result, err := deps.Services.Project.Info(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	cmd.Flags().BoolVar(&options.ShowGenerated, "show-generated", false, "Include generated artifact metadata")
	cmd.Flags().BoolVar(&options.ShowBindings, "show-bindings", false, "Include binding metadata")
	cmd.Flags().String("cwd", ".", "Project directory")

	return cmd
}
