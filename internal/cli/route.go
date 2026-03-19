package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	routesvc "github.com/paolo/flare-edge-cli/internal/service/route"
)

func newRouteCommand(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route",
		Short: "Manage route and custom-domain config",
	}
	cmd.AddCommand(newRouteAttachCommand(deps))
	cmd.AddCommand(newRouteDomainCommand(deps))
	cmd.AddCommand(newRouteDetachCommand(deps))
	return cmd
}

func newRouteAttachCommand(deps Dependencies) *cobra.Command {
	var options routesvc.AttachOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "attach",
		Short: "Attach a Worker route in wrangler.jsonc",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.Route.Attach(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Route, "route", "", "Route pattern")
	cmd.Flags().StringVar(&options.Zone, "zone", "", "Zone name or ID")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().StringVar(&options.Script, "script", "", "Worker name override")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	_ = cmd.MarkFlagRequired("route")
	return cmd
}

func newRouteDomainCommand(deps Dependencies) *cobra.Command {
	var options routesvc.DomainOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "Attach a custom domain in wrangler.jsonc",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.Route.Domain(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Hostname, "hostname", "", "Custom domain hostname")
	cmd.Flags().StringVar(&options.Zone, "zone", "", "Zone name or ID")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().StringVar(&options.Script, "script", "", "Worker name override")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	_ = cmd.MarkFlagRequired("hostname")
	return cmd
}

func newRouteDetachCommand(deps Dependencies) *cobra.Command {
	var options routesvc.DetachOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "detach",
		Short: "Detach routes or custom domains from wrangler.jsonc",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.Route.Detach(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Route, "route", "", "Route pattern")
	cmd.Flags().StringVar(&options.Hostname, "hostname", "", "Custom domain hostname")
	cmd.Flags().StringVar(&options.Zone, "zone", "", "Zone name or ID")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}
