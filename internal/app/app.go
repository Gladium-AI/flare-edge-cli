package app

import (
	"fmt"
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
	d1svc "github.com/paolo/flare-edge-cli/internal/service/d1"
	deploysvc "github.com/paolo/flare-edge-cli/internal/service/deploy"
	devsvc "github.com/paolo/flare-edge-cli/internal/service/dev"
	doctorsvc "github.com/paolo/flare-edge-cli/internal/service/doctor"
	kvsvc "github.com/paolo/flare-edge-cli/internal/service/kv"
	logssvc "github.com/paolo/flare-edge-cli/internal/service/logs"
	projectsvc "github.com/paolo/flare-edge-cli/internal/service/project"
	r2svc "github.com/paolo/flare-edge-cli/internal/service/r2"
	releasesvc "github.com/paolo/flare-edge-cli/internal/service/release"
	routesvc "github.com/paolo/flare-edge-cli/internal/service/route"
	secretsvc "github.com/paolo/flare-edge-cli/internal/service/secret"
	"github.com/paolo/flare-edge-cli/internal/service/shared"
	teardownsvc "github.com/paolo/flare-edge-cli/internal/service/teardown"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

func Run() int {
	filesystem := fs.New()
	runner := process.NewExecRunner()
	store := configstore.New(filesystem)
	wranglerClient := wrangler.NewClient(runner)
	authState := authsvc.NewStateStore(filesystem)
	goTool := toolchain.NewGoToolchain(runner)
	_ = logging.New(os.Stderr)
	wranglerExec := &shared.WranglerExecutor{Client: wranglerClient, State: authState}
	buildService := buildsvc.NewService(store, filesystem, runner, goTool)

	deps := cli.Dependencies{
		Services: cli.Services{
			Auth:     authsvc.NewService(wranglerClient, authState),
			Build:    buildService,
			Compat:   compatsvc.NewService(),
			D1:       d1svc.NewService(store, filesystem, wranglerExec),
			Deploy:   deploysvc.NewService(store, filesystem, buildService, compatsvc.NewService(), wranglerExec),
			Dev:      devsvc.NewService(runner, wranglerExec),
			Doctor:   doctorsvc.NewService(store, filesystem, runner, goTool, wranglerClient, authState, buildService),
			KV:       kvsvc.NewService(store, filesystem, wranglerExec),
			Logs:     logssvc.NewService(runner, wranglerExec),
			Project:  projectsvc.NewService(store, filesystem),
			R2:       r2svc.NewService(store, filesystem, wranglerExec),
			Release:  releasesvc.NewService(wranglerExec),
			Route:    routesvc.NewService(store, filesystem, buildService, wranglerExec),
			Secret:   secretsvc.NewService(runner, filesystem, wranglerExec),
			Teardown: teardownsvc.NewService(store, filesystem, wranglerExec),
		},
	}

	root := cli.NewRootCommand(deps)
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return exitcode.RuntimeError
	}
	return exitcode.Success
}
