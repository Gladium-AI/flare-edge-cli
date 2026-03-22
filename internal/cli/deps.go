package cli

import (
	aisvc "github.com/paolo/flare-edge-cli/internal/service/ai"
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
	teardownsvc "github.com/paolo/flare-edge-cli/internal/service/teardown"
)

type Services struct {
	AI       *aisvc.Service
	Auth     *authsvc.Service
	Build    *buildsvc.Service
	Compat   *compatsvc.Service
	D1       *d1svc.Service
	Deploy   *deploysvc.Service
	Dev      *devsvc.Service
	Doctor   *doctorsvc.Service
	KV       *kvsvc.Service
	Logs     *logssvc.Service
	Project  *projectsvc.Service
	R2       *r2svc.Service
	Release  *releasesvc.Service
	Route    *routesvc.Service
	Secret   *secretsvc.Service
	Teardown *teardownsvc.Service
}

type Dependencies struct {
	Services Services
}
