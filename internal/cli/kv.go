package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/paolo/flare-edge-cli/internal/output"
	kvsvc "github.com/paolo/flare-edge-cli/internal/service/kv"
)

func newKVCommand(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{Use: "kv", Short: "Manage Workers KV"}
	namespace := &cobra.Command{Use: "namespace", Short: "Manage KV namespaces"}
	key := &cobra.Command{Use: "key", Short: "Manage KV entries"}
	namespace.AddCommand(newKVNamespaceCreateCommand(deps))
	put := newKVPutCommand(deps)
	get := newKVGetCommand(deps)
	list := newKVListCommand(deps)
	deleteCmd := newKVDeleteCommand(deps)
	key.AddCommand(newKVPutCommand(deps), newKVGetCommand(deps), newKVListCommand(deps), newKVDeleteCommand(deps))
	cmd.AddCommand(namespace, key, put, get, list, deleteCmd)
	return cmd
}

func newKVNamespaceCreateCommand(deps Dependencies) *cobra.Command {
	var options kvsvc.NamespaceCreateOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "create <binding>",
		Short: "Create a KV namespace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.Binding = args[0]
			result, err := deps.Services.KV.NamespaceCreate(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Title, "title", "", "Remote namespace title")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&options.Provision, "provision", false, "Update the local Worker binding")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	return cmd
}

func newKVPutCommand(deps Dependencies) *cobra.Command {
	var options kvsvc.KeyPutOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "put",
		Short: "Write a KV key",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.KV.Put(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Binding, "binding", "", "KV binding")
	cmd.Flags().StringVar(&options.Key, "key", "", "KV key")
	cmd.Flags().StringVar(&options.Value, "value", "", "KV value")
	cmd.Flags().StringVar(&options.FromFile, "from-file", "", "Load the value from a file")
	cmd.Flags().IntVar(&options.TTL, "ttl", 0, "Entry TTL in seconds")
	cmd.Flags().IntVar(&options.Expiration, "expiration", 0, "Expiration as a Unix timestamp")
	cmd.Flags().StringVar(&options.Metadata, "metadata", "", "Metadata JSON")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	_ = cmd.MarkFlagRequired("binding")
	_ = cmd.MarkFlagRequired("key")
	return cmd
}

func newKVGetCommand(deps Dependencies) *cobra.Command {
	var options kvsvc.KeyGetOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Read a KV value",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.KV.Get(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Binding, "binding", "", "KV binding")
	cmd.Flags().StringVar(&options.Key, "key", "", "KV key")
	cmd.Flags().BoolVar(&options.Text, "text", false, "Decode as UTF-8 text")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	_ = cmd.MarkFlagRequired("binding")
	_ = cmd.MarkFlagRequired("key")
	return cmd
}

func newKVListCommand(deps Dependencies) *cobra.Command {
	var options kvsvc.KeyListOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List KV keys",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.KV.List(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Binding, "binding", "", "KV binding")
	cmd.Flags().StringVar(&options.Prefix, "prefix", "", "Key prefix")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	_ = cmd.MarkFlagRequired("binding")
	return cmd
}

func newKVDeleteCommand(deps Dependencies) *cobra.Command {
	var options kvsvc.KeyDeleteOptions
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a KV key",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := deps.Services.KV.Delete(context.Background(), options)
			if err != nil {
				return err
			}
			return output.NewPrinter(cmd.OutOrStdout(), cmd.ErrOrStderr(), jsonOutput).Print(result)
		},
	}
	cmd.Flags().StringVar(&options.Dir, "path", ".", "Project path")
	cmd.Flags().StringVar(&options.Binding, "binding", "", "KV binding")
	cmd.Flags().StringVar(&options.Key, "key", "", "KV key")
	cmd.Flags().StringVar(&options.Env, "env", "", "Wrangler environment")
	cmd.Flags().Bool("yes", false, "Reserved for non-interactive confirmation")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Emit machine-readable JSON")
	_ = cmd.MarkFlagRequired("binding")
	_ = cmd.MarkFlagRequired("key")
	return cmd
}
