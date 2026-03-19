package cli

import (
	authsvc "github.com/paolo/flare-edge-cli/internal/service/auth"
	buildsvc "github.com/paolo/flare-edge-cli/internal/service/build"
	compatsvc "github.com/paolo/flare-edge-cli/internal/service/compat"
	projectsvc "github.com/paolo/flare-edge-cli/internal/service/project"
)

type Services struct {
	Auth    *authsvc.Service
	Build   *buildsvc.Service
	Compat  *compatsvc.Service
	Project *projectsvc.Service
}

type Dependencies struct {
	Services Services
}
