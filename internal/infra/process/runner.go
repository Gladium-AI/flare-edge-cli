package process

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

type Command struct {
	Name string
	Args []string
	Dir  string
	Env  []string
}

type Result struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}

type Runner interface {
	Run(ctx context.Context, cmd Command) (Result, error)
}

type ExecRunner struct{}

func NewExecRunner() *ExecRunner {
	return &ExecRunner{}
}

func (r *ExecRunner) Run(ctx context.Context, cmd Command) (Result, error) {
	execCmd := exec.CommandContext(ctx, cmd.Name, cmd.Args...)
	execCmd.Dir = cmd.Dir
	execCmd.Env = cmd.Env

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr

	err := execCmd.Run()
	result := Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if execCmd.ProcessState != nil {
		result.ExitCode = execCmd.ProcessState.ExitCode()
	}

	if err != nil {
		return result, fmt.Errorf("run %s: %w", cmd.Name, err)
	}

	return result, nil
}

