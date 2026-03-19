package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

const stateFilename = "state.json"

type State struct {
	APIToken  string `json:"api_token,omitempty"`
	AccountID string `json:"account_id,omitempty"`
}

type StateStore struct {
	fs *fs.FileSystem
}

func NewStateStore(fs *fs.FileSystem) *StateStore {
	return &StateStore{fs: fs}
}

func (s *StateStore) path() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}
	return filepath.Join(configDir, "flare-edge-cli", stateFilename), nil
}

func (s *StateStore) Load() (State, error) {
	path, err := s.path()
	if err != nil {
		return State{}, err
	}

	exists, err := s.fs.Exists(path)
	if err != nil {
		return State{}, err
	}
	if !exists {
		return State{}, nil
	}

	data, err := s.fs.ReadFile(path)
	if err != nil {
		return State{}, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, fmt.Errorf("decode auth state: %w", err)
	}
	return state, nil
}

func (s *StateStore) Save(state State) error {
	path, err := s.path()
	if err != nil {
		return err
	}

	return s.fs.WriteJSON(path, state, 0o600)
}

func (s *StateStore) Delete() error {
	path, err := s.path()
	if err != nil {
		return err
	}

	return s.fs.Remove(path)
}
