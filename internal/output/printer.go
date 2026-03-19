package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/paolo/flare-edge-cli/internal/domain/diagnostic"
)

type Printer struct {
	stdout io.Writer
	stderr io.Writer
	json   bool
}

func NewPrinter(stdout, stderr io.Writer, json bool) *Printer {
	return &Printer{
		stdout: stdout,
		stderr: stderr,
		json:   json,
	}
}

func (p *Printer) Print(value any) error {
	if p.json {
		payload, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(p.stdout, string(payload))
		return err
	}

	_, err := fmt.Fprintln(p.stdout, value)
	return err
}

func (p *Printer) PrintDiagnostics(items []diagnostic.Item) error {
	if p.json {
		return p.Print(items)
	}

	if len(items) == 0 {
		_, err := fmt.Fprintln(p.stdout, "No diagnostics.")
		return err
	}

	var lines []string
	for _, item := range items {
		location := item.File
		if item.Line > 0 {
			location = fmt.Sprintf("%s:%d", location, item.Line)
		}
		if location == "" {
			location = "-"
		}
		lines = append(lines, fmt.Sprintf("[%s] %s %s", strings.ToUpper(string(item.Severity)), item.RuleID, location))
		lines = append(lines, fmt.Sprintf("  %s", item.Message))
		if item.Why != "" {
			lines = append(lines, fmt.Sprintf("  why: %s", item.Why))
		}
		if item.FixHint != "" {
			lines = append(lines, fmt.Sprintf("  fix: %s", item.FixHint))
		}
	}

	_, err := fmt.Fprintln(p.stdout, strings.Join(lines, "\n"))
	return err
}

func (p *Printer) Errorf(format string, args ...any) error {
	_, err := fmt.Fprintf(p.stderr, format+"\n", args...)
	return err
}
