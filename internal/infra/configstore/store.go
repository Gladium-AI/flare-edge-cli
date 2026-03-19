package configstore

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

type Store struct {
	fs        *fs.FileSystem
	validator *validator.Validate
}

func New(fs *fs.FileSystem) *Store {
	return &Store{
		fs:        fs,
		validator: validator.New(),
	}
}

func (s *Store) LoadProject(dir string) (config.Project, error) {
	path := filepath.Join(dir, config.DefaultProjectConfigFile)
	data, err := s.fs.ReadFile(path)
	if err != nil {
		return config.Project{}, err
	}

	var project config.Project
	if err := json.Unmarshal(data, &project); err != nil {
		return config.Project{}, fmt.Errorf("decode %s: %w", path, err)
	}

	if err := s.validator.Struct(project); err != nil {
		return config.Project{}, fmt.Errorf("validate %s: %w", path, err)
	}

	return project, nil
}

func (s *Store) SaveProject(dir string, project config.Project) error {
	if err := s.validator.Struct(project); err != nil {
		return fmt.Errorf("validate project: %w", err)
	}

	return s.fs.WriteJSON(filepath.Join(dir, config.DefaultProjectConfigFile), project, 0o644)
}

func (s *Store) LoadWrangler(dir string, filename string) (config.WranglerConfig, error) {
	path := filepath.Join(dir, filename)
	data, err := s.fs.ReadFile(path)
	if err != nil {
		return config.WranglerConfig{}, err
	}

	var wrangler config.WranglerConfig
	if err := json.Unmarshal(data, &wrangler); err != nil {
		return config.WranglerConfig{}, fmt.Errorf("decode %s: %w", path, err)
	}

	return wrangler, nil
}

func (s *Store) SaveWrangler(dir string, filename string, wrangler config.WranglerConfig) error {
	return s.fs.WriteJSON(filepath.Join(dir, filename), wrangler, 0o644)
}
