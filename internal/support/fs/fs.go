package fs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type FileSystem struct {
	fs afero.Fs
}

func New() *FileSystem {
	return &FileSystem{fs: afero.NewOsFs()}
}

func (f *FileSystem) MkdirAll(path string, perm os.FileMode) error {
	return f.fs.MkdirAll(path, perm)
}

func (f *FileSystem) Exists(path string) (bool, error) {
	exists, err := afero.Exists(f.fs, path)
	if err != nil {
		return false, fmt.Errorf("stat %s: %w", path, err)
	}

	return exists, nil
}

func (f *FileSystem) ReadFile(path string) ([]byte, error) {
	data, err := afero.ReadFile(f.fs, path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return data, nil
}

func (f *FileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	if err := f.fs.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}
	if err := afero.WriteFile(f.fs, path, data, perm); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func (f *FileSystem) Remove(path string) error {
	if err := f.fs.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove %s: %w", path, err)
	}
	return nil
}

func (f *FileSystem) RemoveAll(path string) error {
	if err := f.fs.RemoveAll(path); err != nil {
		return fmt.Errorf("remove all %s: %w", path, err)
	}
	return nil
}

func (f *FileSystem) CopyFile(src, dst string, perm os.FileMode) error {
	data, err := f.ReadFile(src)
	if err != nil {
		return err
	}
	return f.WriteFile(dst, data, perm)
}

func (f *FileSystem) WriteJSON(path string, value any, perm os.FileMode) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	data = append(data, '\n')
	return f.WriteFile(path, data, perm)
}
