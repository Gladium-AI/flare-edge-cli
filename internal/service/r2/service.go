package r2

import (
	"context"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	"github.com/paolo/flare-edge-cli/internal/service/shared"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

type Service struct {
	store    *configstore.Store
	fs       *fs.FileSystem
	wrangler *shared.WranglerExecutor
}

type BucketCreateOptions struct {
	Dir          string
	Binding      string
	Name         string
	Location     string
	StorageClass string
	Env          string
}

type ObjectPutOptions struct {
	Dir                string
	Binding            string
	Key                string
	File               string
	ContentType        string
	CacheControl       string
	ContentDisposition string
	Env                string
}

type ObjectGetOptions struct {
	Dir     string
	Binding string
	Key     string
	Out     string
	Env     string
}

type ObjectDeleteOptions struct {
	Dir     string
	Binding string
	Key     string
	Env     string
}

func NewService(store *configstore.Store, fs *fs.FileSystem, wrangler *shared.WranglerExecutor) *Service {
	return &Service{store: store, fs: fs, wrangler: wrangler}
}

func (s *Service) BucketCreate(ctx context.Context, options BucketCreateOptions) (shared.CommandResult, error) {
	name := options.Name
	if name == "" {
		name = options.Binding
	}
	command := []string{"r2", "bucket", "create", name}
	if options.Location != "" {
		command = append(command, "--location", options.Location)
	}
	if options.StorageClass != "" {
		command = append(command, "--storage-class", options.StorageClass)
	}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	project, wranglerCfg, err := shared.LoadProjectAndWrangler(options.Dir, s.store, s.fs)
	if err != nil {
		return shared.CommandResult{}, err
	}
	project.Bindings.R2 = append(project.Bindings.R2, config.R2Binding{Binding: options.Binding, BucketName: name, StorageClass: options.StorageClass})
	wranglerCfg.R2Buckets = config.UpsertR2(wranglerCfg.R2Buckets, config.WranglerR2Bucket{Binding: options.Binding, BucketName: name})
	if err := s.store.SaveProject(options.Dir, project); err != nil {
		return shared.CommandResult{}, err
	}
	if err := shared.SaveWrangler(options.Dir, project, wranglerCfg, s.store); err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}

func (s *Service) Put(ctx context.Context, options ObjectPutOptions) (shared.CommandResult, error) {
	bucket, err := s.resolveBucket(options.Dir, options.Binding)
	if err != nil {
		return shared.CommandResult{}, err
	}
	command := []string{"r2", "object", "put", bucket + "/" + options.Key, "--file", options.File, "--remote"}
	if options.ContentType != "" {
		command = append(command, "--content-type", options.ContentType)
	}
	if options.CacheControl != "" {
		command = append(command, "--cache-control", options.CacheControl)
	}
	if options.ContentDisposition != "" {
		command = append(command, "--content-disposition", options.ContentDisposition)
	}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}

func (s *Service) Get(ctx context.Context, options ObjectGetOptions) (shared.CommandResult, error) {
	bucket, err := s.resolveBucket(options.Dir, options.Binding)
	if err != nil {
		return shared.CommandResult{}, err
	}
	command := []string{"r2", "object", "get", bucket + "/" + options.Key, "--remote"}
	if options.Out != "" {
		command = append(command, "--file", options.Out)
	}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}

func (s *Service) Delete(ctx context.Context, options ObjectDeleteOptions) (shared.CommandResult, error) {
	bucket, err := s.resolveBucket(options.Dir, options.Binding)
	if err != nil {
		return shared.CommandResult{}, err
	}
	command := []string{"r2", "object", "delete", bucket + "/" + options.Key, "--remote"}
	raw, err := s.wrangler.Run(ctx, options.Dir, options.Env, command...)
	if err != nil {
		return shared.CommandResult{}, err
	}
	return shared.NewCommandResult(command, raw), nil
}

func (s *Service) resolveBucket(dir, binding string) (string, error) {
	project, _, err := shared.LoadProjectAndWrangler(dir, s.store, s.fs)
	if err != nil {
		return "", err
	}
	for _, bucket := range project.Bindings.R2 {
		if bucket.Binding == binding {
			return bucket.BucketName, nil
		}
	}
	return binding, nil
}
