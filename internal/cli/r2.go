package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	r2svc "github.com/paolo/flare-edge-cli/internal/service/r2"
)

func newR2Command(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{Use: "r2", Short: "Manage R2 buckets and objects"}
	bucket := &cobra.Command{Use: "bucket", Short: "Manage R2 buckets"}
	object := &cobra.Command{Use: "object", Short: "Manage R2 objects"}
	bucket.AddCommand(newR2BucketCreateCommand(deps))
	object.AddCommand(newR2ObjectPutCommand(deps), newR2ObjectGetCommand(deps), newR2ObjectDeleteCommand(deps))
	cmd.AddCommand(bucket, object)
	return cmd
}

func newR2BucketCreateCommand(deps Dependencies) *cobra.Command {
	var options r2svc.BucketCreateOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "create <binding>",
		Short: "Create an R2 bucket and update bindings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.Binding = args[0]
			result, err := deps.Services.R2.BucketCreate(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Name, "name", "", "Bucket name")
	cmd.Flags().StringVar(&options.Location, "location", "", "Bucket location")
	cmd.Flags().StringVar(&options.StorageClass, "storage-class", "", "Storage class")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}

func newR2ObjectPutCommand(deps Dependencies) *cobra.Command {
	var options r2svc.ObjectPutOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "put",
		Short: "Upload an object to R2",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.R2.Put(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Binding, "binding", "", "R2 bucket or binding")
	cmd.Flags().StringVar(&options.Key, "key", "", "Object key")
	cmd.Flags().StringVar(&options.File, "file", "", "Source file")
	cmd.Flags().StringVar(&options.ContentType, "content-type", "", "Content-Type header")
	cmd.Flags().StringVar(&options.CacheControl, "cache-control", "", "Cache-Control header")
	cmd.Flags().StringVar(&options.ContentDisposition, "content-disposition", "", "Content-Disposition header")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	_ = cmd.MarkFlagRequired("binding")
	_ = cmd.MarkFlagRequired("key")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newR2ObjectGetCommand(deps Dependencies) *cobra.Command {
	var options r2svc.ObjectGetOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Download an object from R2",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.R2.Get(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Binding, "binding", "", "R2 bucket or binding")
	cmd.Flags().StringVar(&options.Key, "key", "", "Object key")
	cmd.Flags().StringVar(&options.Out, "out", "", "Destination file")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	_ = cmd.MarkFlagRequired("binding")
	_ = cmd.MarkFlagRequired("key")
	return cmd
}

func newR2ObjectDeleteCommand(deps Dependencies) *cobra.Command {
	var options r2svc.ObjectDeleteOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an object from R2",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.R2.Delete(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Binding, "binding", "", "R2 bucket or binding")
	cmd.Flags().StringVar(&options.Key, "key", "", "Object key")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().Bool("yes", false, "Reserved for non-interactive confirmation")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	_ = cmd.MarkFlagRequired("binding")
	_ = cmd.MarkFlagRequired("key")
	return cmd
}
