package shared

import (
	"path/filepath"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

func LoadProjectAndWrangler(dir string, store *configstore.Store, filesystem *fs.FileSystem) (config.Project, config.WranglerConfig, error) {
	project, err := store.LoadProject(dir)
	if err != nil {
		return config.Project{}, config.WranglerConfig{}, err
	}
	wrangler, err := store.LoadWrangler(dir, project.WranglerConfig)
	if err != nil {
		return config.Project{}, config.WranglerConfig{}, err
	}
	return project, wrangler, nil
}

func SaveWrangler(dir string, project config.Project, wrangler config.WranglerConfig, store *configstore.Store) error {
	return store.SaveWrangler(dir, project.WranglerConfig, wrangler)
}

func ProjectConfigPath(dir string) string {
	return filepath.Join(dir, config.DefaultProjectConfigFile)
}
