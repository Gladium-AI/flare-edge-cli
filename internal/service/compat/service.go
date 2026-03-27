package compat

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	domaincompat "github.com/paolo/flare-edge-cli/internal/domain/compat"
	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/domain/diagnostic"
	"golang.org/x/tools/go/packages"
)

type Service struct{}

type CheckOptions struct {
	Path    string
	Entry   string
	Profile string
	Strict  bool
	FailOn  string
	Exclude []string
}

type CheckResult struct {
	Profile     string            `json:"profile"`
	FailOn      string            `json:"fail_on"`
	Diagnostics []diagnostic.Item `json:"diagnostics"`
	ErrorCount  int               `json:"error_count"`
	WarnCount   int               `json:"warning_count"`
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Check(_ context.Context, options CheckOptions) (CheckResult, error) {
	if project, ok, err := loadProject(options.Path); err != nil {
		return CheckResult{}, err
	} else if ok && !project.RequiresBuild() {
		return CheckResult{
			Profile:     defaultProjectProfile(project, options.Profile),
			FailOn:      defaultString(options.FailOn, "error"),
			Diagnostics: nil,
		}, nil
	}

	loadMode := packages.NeedName | packages.NeedCompiledGoFiles | packages.NeedSyntax
	cfg := &packages.Config{
		Mode: loadMode,
		Dir:  defaultString(options.Path, "."),
	}

	pattern := "./..."
	if options.Entry != "" {
		pattern = options.Entry
	}

	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		return CheckResult{}, fmt.Errorf("load packages: %w", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		return CheckResult{}, fmt.Errorf("package loading reported errors")
	}

	var diagnosticsOut []diagnostic.Item
	for _, pkg := range pkgs {
		for index, file := range pkg.Syntax {
			filename := pkg.CompiledGoFiles[index]
			if excluded(filename, options.Exclude) {
				continue
			}

			imports := importAliases(file)
			ast.Inspect(file, func(node ast.Node) bool {
				switch typed := node.(type) {
				case *ast.ImportSpec:
					diagnosticsOut = append(diagnosticsOut, importDiagnostics(filename, pkg.Fset, typed)...)
				case *ast.GoStmt:
					diagnosticsOut = append(diagnosticsOut, newItem("FE006", filename, pkg.Fset, typed, "goroutine detected"))
				case *ast.CallExpr:
					diagnosticsOut = append(diagnosticsOut, callDiagnostics(filename, pkg.Fset, typed, imports)...)
				}
				return true
			})
		}
	}

	sort.Slice(diagnosticsOut, func(left, right int) bool {
		if diagnosticsOut[left].File != diagnosticsOut[right].File {
			return diagnosticsOut[left].File < diagnosticsOut[right].File
		}
		if diagnosticsOut[left].Line != diagnosticsOut[right].Line {
			return diagnosticsOut[left].Line < diagnosticsOut[right].Line
		}
		return diagnosticsOut[left].RuleID < diagnosticsOut[right].RuleID
	})

	result := CheckResult{
		Profile:     defaultString(options.Profile, "worker-wasm"),
		FailOn:      defaultString(options.FailOn, "error"),
		Diagnostics: diagnosticsOut,
	}
	for _, item := range diagnosticsOut {
		if item.Severity == diagnostic.SeverityError {
			result.ErrorCount++
		}
		if item.Severity == diagnostic.SeverityWarning {
			result.WarnCount++
		}
	}
	return result, nil
}

func loadProject(dir string) (config.Project, bool, error) {
	if dir == "" {
		dir = "."
	}
	path := filepath.Join(dir, config.DefaultProjectConfigFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config.Project{}, false, nil
		}
		return config.Project{}, false, fmt.Errorf("read %s: %w", path, err)
	}

	var project config.Project
	if err := json.Unmarshal(data, &project); err != nil {
		return config.Project{}, false, fmt.Errorf("decode %s: %w", path, err)
	}
	return project, true, nil
}

func defaultProjectProfile(project config.Project, profile string) string {
	if profile != "" {
		return profile
	}
	if project.CompatibilityProfile != "" {
		return project.CompatibilityProfile
	}
	return defaultString(profile, "worker-wasm")
}

func (s *Service) Rules(severity string) []domaincompat.Rule {
	all := domaincompat.BuiltInRules()
	if severity == "" {
		return all
	}
	filtered := make([]domaincompat.Rule, 0, len(all))
	for _, rule := range all {
		if string(rule.Severity) == severity {
			filtered = append(filtered, rule)
		}
	}
	return filtered
}

func importAliases(file *ast.File) map[string]string {
	aliases := map[string]string{}
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		name := filepath.Base(path)
		if spec.Name != nil && spec.Name.Name != "" {
			name = spec.Name.Name
		}
		aliases[name] = path
	}
	return aliases
}

func importDiagnostics(filename string, fset *token.FileSet, spec *ast.ImportSpec) []diagnostic.Item {
	path, err := strconv.Unquote(spec.Path.Value)
	if err != nil {
		return nil
	}

	switch path {
	case "C":
		return []diagnostic.Item{newItem("FE001", filename, fset, spec, "import \"C\" is incompatible with Workers Wasm builds")}
	case "os/exec", "plugin", "syscall":
		return []diagnostic.Item{newItem("FE002", filename, fset, spec, "unsupported host package imported: "+path)}
	}

	return nil
}

func callDiagnostics(filename string, fset *token.FileSet, call *ast.CallExpr, aliases map[string]string) []diagnostic.Item {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return nil
	}

	importPath := aliases[ident.Name]
	switch importPath {
	case "os":
		if contains(selector.Sel.Name, "Open", "OpenFile", "Create", "ReadFile", "WriteFile", "Remove", "RemoveAll", "Mkdir", "MkdirAll", "Rename", "Stat", "Lstat") {
			return []diagnostic.Item{newItem("FE003", filename, fset, call, "filesystem access is unavailable in Workers")}
		}
		if contains(selector.Sel.Name, "Getenv", "LookupEnv") {
			return []diagnostic.Item{newItem("FE007", filename, fset, call, "use Workers env bindings instead of os.Getenv")}
		}
	case "net":
		if contains(selector.Sel.Name, "Listen", "ListenPacket") {
			return []diagnostic.Item{newItem("FE004", filename, fset, call, "network listeners are incompatible with the Workers runtime")}
		}
	case "net/http":
		if contains(selector.Sel.Name, "ListenAndServe", "ListenAndServeTLS", "Serve") {
			return []diagnostic.Item{newItem("FE004", filename, fset, call, "HTTP servers should be replaced with a fetch handler")}
		}
	case "os/exec":
		if contains(selector.Sel.Name, "Command", "CommandContext") {
			return []diagnostic.Item{newItem("FE005", filename, fset, call, "process execution is not available in Workers")}
		}
	}

	return nil
}

func newItem(ruleID string, filename string, fset *token.FileSet, node ast.Node, message string) diagnostic.Item {
	rule := ruleByID(ruleID)
	position := fset.Position(node.Pos())
	return diagnostic.Item{
		RuleID:   rule.ID,
		Severity: rule.Severity,
		File:     filename,
		Line:     position.Line,
		Message:  message,
		Why:      rule.Why,
		FixHint:  rule.FixHint,
	}
}

func ruleByID(ruleID string) domaincompat.Rule {
	for _, rule := range domaincompat.BuiltInRules() {
		if rule.ID == ruleID {
			return rule
		}
	}
	return domaincompat.Rule{ID: ruleID, Severity: diagnostic.SeverityInfo}
}

func excluded(filename string, globs []string) bool {
	for _, pattern := range globs {
		match, err := filepath.Match(pattern, filename)
		if err == nil && match {
			return true
		}
		if strings.Contains(filename, pattern) {
			return true
		}
	}
	return false
}

func contains(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
