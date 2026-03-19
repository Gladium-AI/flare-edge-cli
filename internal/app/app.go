package app

import (
	"os"

	"github.com/paolo/flare-edge-cli/internal/cli"
	"github.com/paolo/flare-edge-cli/internal/domain/exitcode"
	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	"github.com/paolo/flare-edge-cli/internal/infra/process"
	"github.com/paolo/flare-edge-cli/internal/infra/toolchain"
	"github.com/paolo/flare-edge-cli/internal/infra/wrangler"
	"github.com/paolo/flare-edge-cli/internal/logging"
	authsvc "github.com/paolo/flare-edge-cli/internal/service/auth"
	buildsvc "github.com/paolo/flare-edge-cli/internal/service/build"
	compatsvc "github.com/paolo/flare-edge-cli/internal/service/compat"
	projectsvc "github.com/paolo/flare-edge-cli/internal/service/project"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

func Run() int {
	filesystem := fs.New()
	runner := process.NewExecRunner()
	store := configstore.New(filesystem)
	wranglerClient := wrangler.NewClient(runner)
	authState := authsvc.NewStateStore(filesystem)
	_ = toolchain.NewGoToolchain(runner)
	_ = logging.New(os.Stderr)

	deps := cli.Dependencies{
		Services: cli.Services{
			Auth:    authsvc.NewService(wranglerClient, authState),
			Build:   buildsvc.NewService(store, filesystem, runner, toolchain.NewGoToolchain(runner)),
			Compat:  compatsvc.NewService(),
			Project: projectsvc.NewService(store, filesystem),
		},
	}

	root := cli.NewRootCommand(deps)
	if err := root.Execute(); err != nil {
		return exitcode.RuntimeError
	}
	return exitcode.Success
}
