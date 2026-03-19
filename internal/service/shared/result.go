package shared

import (
	"strings"

	"github.com/paolo/flare-edge-cli/internal/infra/process"
)

type CommandResult struct {
	Command []string `json:"command"`
	Stdout  string   `json:"stdout,omitempty"`
	Stderr  string   `json:"stderr,omitempty"`
}

func NewCommandResult(command []string, result process.Result) CommandResult {
	return CommandResult{
		Command: command,
		Stdout:  strings.TrimSpace(result.Stdout),
		Stderr:  strings.TrimSpace(result.Stderr),
	}
}
