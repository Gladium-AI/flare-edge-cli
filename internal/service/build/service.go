package build

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	"github.com/paolo/flare-edge-cli/internal/infra/process"
	"github.com/paolo/flare-edge-cli/internal/infra/toolchain"
	projectsvc "github.com/paolo/flare-edge-cli/internal/service/project"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
	"github.com/tetratelabs/wabin/binary"
	"github.com/tetratelabs/wabin/wasm"
)

type Service struct {
	store  *configstore.Store
	fs     *fs.FileSystem
	runner process.Runner
	goTool *toolchain.GoToolchain
}

type WasmOptions struct {
	Path     string
	Entry    string
	OutDir   string
	OutFile  string
	ShimOut  string
	Target   string
	Optimize string
	TinyGo   bool
	NoShim   bool
	Clean    bool
}

type WasmResult struct {
	Artifact string   `json:"artifact"`
	Shim     string   `json:"shim,omitempty"`
	Files    []string `json:"files"`
	Compiler string   `json:"compiler"`
}

type InspectOptions struct {
	Artifact string
	Size     bool
	Exports  bool
	Imports  bool
}

type InspectResult struct {
	Artifact string   `json:"artifact"`
	Size     int64    `json:"size,omitempty"`
	Exports  []string `json:"exports,omitempty"`
	Imports  []string `json:"imports,omitempty"`
}

func NewService(store *configstore.Store, fs *fs.FileSystem, runner process.Runner, goTool *toolchain.GoToolchain) *Service {
	return &Service{
		store:  store,
		fs:     fs,
		runner: runner,
		goTool: goTool,
	}
}

func (s *Service) Wasm(ctx context.Context, options WasmOptions) (WasmResult, error) {
	project, err := s.loadProject(options.Path)
	if err != nil {
		return WasmResult{}, err
	}

	outDir := defaultString(options.OutDir, filepath.Join(options.Path, project.OutDir))
	outFile := defaultString(options.OutFile, project.WasmFile)
	shimOut := defaultString(options.ShimOut, filepath.Join(outDir, project.ShimFile))
	artifact := filepath.Join(outDir, outFile)

	if options.Clean {
		if err := s.fs.RemoveAll(outDir); err != nil {
			return WasmResult{}, err
		}
	}
	if err := s.fs.MkdirAll(outDir, 0o755); err != nil {
		return WasmResult{}, err
	}
	compileArtifact, err := filepath.Abs(artifact)
	if err != nil {
		return WasmResult{}, fmt.Errorf("resolve artifact path: %w", err)
	}

	entry := defaultString(options.Entry, project.Entry)
	compiler := "go"
	if options.TinyGo {
		compiler = "tinygo"
	}
	if err := s.compile(ctx, compiler, options.Path, entry, compileArtifact, options.Optimize); err != nil {
		return WasmResult{}, err
	}
	exists, err := s.fs.Exists(artifact)
	if err != nil {
		return WasmResult{}, err
	}
	if !exists {
		return WasmResult{}, fmt.Errorf("wasm build reported success but artifact is missing: %s", artifact)
	}

	files := []string{artifact}
	if !options.NoShim {
		wasmExecPath, err := s.goTool.WasmExecPath(ctx, options.Path)
		if err != nil {
			return WasmResult{}, err
		}
		if err := s.fs.CopyFile(wasmExecPath, filepath.Join(outDir, "wasm_exec.js"), 0o644); err != nil {
			return WasmResult{}, err
		}

		content := projectsvc.WorkerShimTemplate(outFile)
		if err := s.fs.WriteFile(shimOut, []byte(content), 0o644); err != nil {
			return WasmResult{}, err
		}
		files = append(files, filepath.Join(outDir, "wasm_exec.js"), shimOut)
	}

	return WasmResult{
		Artifact: artifact,
		Shim:     shimOut,
		Files:    files,
		Compiler: compiler,
	}, nil
}

func (s *Service) Inspect(_ context.Context, options InspectOptions) (InspectResult, error) {
	data, err := s.fs.ReadFile(options.Artifact)
	if err != nil {
		return InspectResult{}, err
	}

	module, err := binary.DecodeModule(data, wasm.CoreFeaturesV2)
	if err != nil {
		return InspectResult{}, fmt.Errorf("decode wasm module: %w", err)
	}

	result := InspectResult{Artifact: options.Artifact}
	if options.Size || (!options.Exports && !options.Imports) {
		result.Size = int64(len(data))
	}
	if options.Exports || (!options.Size && !options.Imports) {
		for _, export := range module.ExportSection {
			result.Exports = append(result.Exports, fmt.Sprintf("%s:%s", export.Name, wasm.ExternTypeName(export.Type)))
		}
	}
	if options.Imports || (!options.Size && !options.Exports) {
		for _, importItem := range module.ImportSection {
			result.Imports = append(result.Imports, fmt.Sprintf("%s.%s:%s", importItem.Module, importItem.Name, wasm.ExternTypeName(importItem.Type)))
		}
	}

	return result, nil
}

func (s *Service) compile(ctx context.Context, compiler, dir, entry, artifact, optimize string) error {
	switch compiler {
	case "tinygo":
		args := []string{"build", "-target", "wasm", "-o", artifact, entry}
		_, err := s.runner.Run(ctx, process.Command{Name: "tinygo", Args: args, Dir: dir})
		return err
	default:
		args := []string{"build", "-trimpath"}
		if strings.EqualFold(optimize, "size") {
			args = append(args, "-ldflags=-s -w")
		}
		args = append(args, "-o", artifact, entry)
		_, err := s.runner.Run(ctx, process.Command{
			Name: "go",
			Args: args,
			Dir:  dir,
			Env:  []string{"GOOS=js", "GOARCH=wasm"},
		})
		return err
	}
}

func (s *Service) loadProject(dir string) (config.Project, error) {
	projectPath := filepath.Join(dir, config.DefaultProjectConfigFile)
	exists, err := s.fs.Exists(projectPath)
	if err != nil {
		return config.Project{}, err
	}
	if !exists {
		project := config.DefaultProject(filepath.Base(dir), "github.com/paolo/"+filepath.Base(dir), "main", "edge-http", config.DefaultCompatibilityDate, "")
		return project, nil
	}
	return s.store.LoadProject(dir)
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
