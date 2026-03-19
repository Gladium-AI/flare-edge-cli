package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	buildsvc "github.com/paolo/flare-edge-cli/internal/service/build"
)

func newBuildCommand(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build",
		Short:   "Build Wasm artifacts for Workers",
		Aliases: []string{"compile"},
	}

	cmd.AddCommand(newBuildWasmCommand(deps))
	cmd.AddCommand(newBuildInspectCommand(deps))
	return cmd
}

func newBuildWasmCommand(deps Dependencies) *cobra.Command {
	var options buildsvc.WasmOptions
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "wasm",
		Short: "Compile Go code to Wasm and stage the Worker shim",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.Build.Wasm(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}

	cmd.Flags().StringVar(&options.Path, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Entry, "entry", "", "Package or file to build")
	cmd.Flags().StringVar(&options.OutDir, "out-dir", "", "Output directory")
	cmd.Flags().StringVar(&options.OutFile, "out-file", "", "Output Wasm filename")
	cmd.Flags().StringVar(&options.ShimOut, "shim-out", "", "Output Worker shim path")
	cmd.Flags().StringVar(&options.Target, "target", "js/wasm", "Compilation target")
	cmd.Flags().StringVar(&options.Optimize, "optimize", "size", "Optimization preference: size|speed")
	cmd.Flags().BoolVar(&options.TinyGo, "tinygo", false, "Compile with TinyGo instead of the standard Go toolchain")
	cmd.Flags().BoolVar(&options.NoShim, "no-shim", false, "Skip Worker shim generation")
	cmd.Flags().BoolVar(&options.Clean, "clean", false, "Remove the output directory before building")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")

	return cmd
}

func newBuildInspectCommand(deps Dependencies) *cobra.Command {
	var options buildsvc.InspectOptions
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect the generated Wasm artifact",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.Build.Inspect(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}

	cmd.Flags().StringVar(&options.Artifact, "artifact", "", "Path to the Wasm artifact")
	cmd.Flags().BoolVar(&options.Size, "size", false, "Show artifact size")
	cmd.Flags().BoolVar(&options.Exports, "exports", false, "Show exported symbols")
	cmd.Flags().BoolVar(&options.Imports, "imports", false, "Show imported symbols")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	_ = cmd.MarkFlagRequired("artifact")

	return cmd
}
