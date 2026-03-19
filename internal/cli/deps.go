package cli

import (
	authsvc "github.com/paolo/flare-edge-cli/internal/service/auth"
	projectsvc "github.com/paolo/flare-edge-cli/internal/service/project"
)

type Services struct {
	Auth    *authsvc.Service
	Project *projectsvc.Service
}

type Dependencies struct {
	Services Services
}
