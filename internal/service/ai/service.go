package ai

import (
	"context"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

type Service struct {
	store *configstore.Store
	fs    *fs.FileSystem
}

type SetBindingOptions struct {
	Dir     string
	Binding string
	Remote  bool
}

type ClearBindingOptions struct {
	Dir string
}

type BindingResult struct {
	Binding string `json:"binding,omitempty"`
	Remote  bool   `json:"remote,omitempty"`
	Cleared bool   `json:"cleared,omitempty"`
}

func NewService(store *configstore.Store, fs *fs.FileSystem) *Service {
	return &Service{store: store, fs: fs}
}

func (s *Service) SetBinding(_ context.Context, options SetBindingOptions) (BindingResult, error) {
	project, wranglerCfg, err := loadProjectAndWrangler(options.Dir, s.store)
	if err != nil {
		return BindingResult{}, err
	}

	project.Bindings.AI = &config.AIBinding{
		Binding: options.Binding,
		Remote:  options.Remote,
	}
	wranglerCfg.AI = &config.WranglerAIBinding{
		Binding: options.Binding,
		Remote:  options.Remote,
	}

	if err := s.store.SaveProject(options.Dir, project); err != nil {
		return BindingResult{}, err
	}
	if err := s.store.SaveWrangler(options.Dir, project.WranglerConfig, wranglerCfg); err != nil {
		return BindingResult{}, err
	}

	return BindingResult{
		Binding: options.Binding,
		Remote:  options.Remote,
	}, nil
}

func (s *Service) ClearBinding(_ context.Context, options ClearBindingOptions) (BindingResult, error) {
	project, wranglerCfg, err := loadProjectAndWrangler(options.Dir, s.store)
	if err != nil {
		return BindingResult{}, err
	}

	project.Bindings.AI = nil
	wranglerCfg.AI = nil

	if err := s.store.SaveProject(options.Dir, project); err != nil {
		return BindingResult{}, err
	}
	if err := s.store.SaveWrangler(options.Dir, project.WranglerConfig, wranglerCfg); err != nil {
		return BindingResult{}, err
	}

	return BindingResult{Cleared: true}, nil
}

func loadProjectAndWrangler(dir string, store *configstore.Store) (config.Project, config.WranglerConfig, error) {
	project, err := store.LoadProject(dir)
	if err != nil {
		return config.Project{}, config.WranglerConfig{}, err
	}
	wranglerCfg, err := store.LoadWrangler(dir, project.WranglerConfig)
	if err != nil {
		return config.Project{}, config.WranglerConfig{}, err
	}
	return project, wranglerCfg, nil
}
