package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	deploysvc "github.com/paolo/flare-edge-cli/internal/service/deploy"
)

func newDeployCommand(deps Dependencies) *cobra.Command {
	var options deploysvc.Options
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Validate, build, and deploy the Worker",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.Deploy.Deploy(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}

	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().StringVar(&options.Name, "name", "", "Worker name override")
	cmd.Flags().StringVar(&options.CompatDate, "compat-date", "", "Compatibility date override")
	cmd.Flags().StringArrayVar(&options.Route, "route", nil, "Route to attach during deploy")
	cmd.Flags().StringArrayVar(&options.CustomDomain, "custom-domain", nil, "Custom domain to attach during deploy")
	cmd.Flags().BoolVar(&options.WorkersDev, "workers-dev", false, "Reserved for workers.dev deployments")
	cmd.Flags().BoolVar(&options.DryRun, "dry-run", false, "Run without applying the deployment")
	cmd.Flags().BoolVar(&options.UploadOnly, "upload-only", false, "Upload a version without promoting it")
	cmd.Flags().StringVar(&options.Message, "message", "", "Deployment message")
	cmd.Flags().StringArrayVar(&options.Var, "var", nil, "Variable in KEY=VALUE format")
	cmd.Flags().BoolVar(&options.KeepVars, "keep-vars", false, "Keep dashboard-managed variables")
	cmd.Flags().BoolVar(&options.Minify, "minify", false, "Minify the Worker bundle")
	cmd.Flags().BoolVar(&options.Latest, "latest", false, "Use the latest runtime")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")

	return cmd
}
