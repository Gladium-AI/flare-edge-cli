package process

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Name  string
	Args  []string
	Dir   string
	Env   []string
	Stdin string
}

type Result struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}

type Runner interface {
	Run(ctx context.Context, cmd Command) (Result, error)
	Stream(ctx context.Context, cmd Command, stdout, stderr io.Writer) error
}

type ExecRunner struct{}

func NewExecRunner() *ExecRunner {
	return &ExecRunner{}
}

func (r *ExecRunner) Run(ctx context.Context, cmd Command) (Result, error) {
	execCmd := exec.CommandContext(ctx, cmd.Name, cmd.Args...)
	execCmd.Dir = cmd.Dir
	if len(cmd.Env) > 0 {
		execCmd.Env = append(os.Environ(), cmd.Env...)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr
	if cmd.Stdin != "" {
		execCmd.Stdin = strings.NewReader(cmd.Stdin)
	}

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

func (r *ExecRunner) Stream(ctx context.Context, cmd Command, stdout, stderr io.Writer) error {
	execCmd := exec.CommandContext(ctx, cmd.Name, cmd.Args...)
	execCmd.Dir = cmd.Dir
	if len(cmd.Env) > 0 {
		execCmd.Env = append(os.Environ(), cmd.Env...)
	}
	execCmd.Stdout = stdout
	execCmd.Stderr = stderr
	execCmd.Stdin = os.Stdin
	if cmd.Stdin != "" {
		execCmd.Stdin = strings.NewReader(cmd.Stdin)
	}
	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("stream %s: %w", cmd.Name, err)
	}
	return nil
}
