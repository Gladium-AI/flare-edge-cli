package route

import (
	"context"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	"github.com/paolo/flare-edge-cli/internal/service/shared"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

type Service struct {
	store *configstore.Store
	fs    *fs.FileSystem
}

type AttachOptions struct {
	Dir    string
	Route  string
	Zone   string
	Env    string
	Script string
}

type DomainOptions struct {
	Dir      string
	Hostname string
	Zone     string
	Env      string
	Script   string
}

type DetachOptions struct {
	Dir      string
	Route    string
	Hostname string
	Zone     string
	Env      string
}

type Result struct {
	Updated bool                   `json:"updated"`
	Routes  []config.WranglerRoute `json:"routes"`
}

func NewService(store *configstore.Store, fs *fs.FileSystem) *Service {
	return &Service{store: store, fs: fs}
}

func (s *Service) Attach(_ context.Context, options AttachOptions) (Result, error) {
	project, wranglerCfg, err := shared.LoadProjectAndWrangler(options.Dir, s.store, s.fs)
	if err != nil {
		return Result{}, err
	}
	if options.Script != "" {
		project.WorkerName = options.Script
		wranglerCfg.Name = options.Script
	}
	route := config.WranglerRoute{Pattern: options.Route}
	assignZone(&route, options.Zone)
	wranglerCfg.Routes = config.UpsertRoute(wranglerCfg.Routes, route)
	if err := s.store.SaveProject(options.Dir, project); err != nil {
		return Result{}, err
	}
	if err := shared.SaveWrangler(options.Dir, project, wranglerCfg, s.store); err != nil {
		return Result{}, err
	}
	return Result{Updated: true, Routes: wranglerCfg.Routes}, nil
}

func (s *Service) Domain(_ context.Context, options DomainOptions) (Result, error) {
	project, wranglerCfg, err := shared.LoadProjectAndWrangler(options.Dir, s.store, s.fs)
	if err != nil {
		return Result{}, err
	}
	if options.Script != "" {
		project.WorkerName = options.Script
		wranglerCfg.Name = options.Script
	}
	route := config.WranglerRoute{Pattern: options.Hostname, CustomDomain: true}
	assignZone(&route, options.Zone)
	wranglerCfg.Routes = config.UpsertRoute(wranglerCfg.Routes, route)
	if err := s.store.SaveProject(options.Dir, project); err != nil {
		return Result{}, err
	}
	if err := shared.SaveWrangler(options.Dir, project, wranglerCfg, s.store); err != nil {
		return Result{}, err
	}
	return Result{Updated: true, Routes: wranglerCfg.Routes}, nil
}

func (s *Service) Detach(_ context.Context, options DetachOptions) (Result, error) {
	project, wranglerCfg, err := shared.LoadProjectAndWrangler(options.Dir, s.store, s.fs)
	if err != nil {
		return Result{}, err
	}
	if options.Route != "" {
		wranglerCfg.Routes = config.RemoveRoute(wranglerCfg.Routes, options.Route, false)
	}
	if options.Hostname != "" {
		wranglerCfg.Routes = config.RemoveRoute(wranglerCfg.Routes, options.Hostname, true)
	}
	if err := shared.SaveWrangler(options.Dir, project, wranglerCfg, s.store); err != nil {
		return Result{}, err
	}
	return Result{Updated: true, Routes: wranglerCfg.Routes}, nil
}

func assignZone(route *config.WranglerRoute, zone string) {
	if zone == "" {
		return
	}
	if len(zone) > 16 && !contains(zone, ".") {
		route.ZoneID = zone
		return
	}
	route.ZoneName = zone
}

func contains(value, needle string) bool {
	for index := 0; index+len(needle) <= len(value); index++ {
		if value[index:index+len(needle)] == needle {
			return true
		}
	}
	return false
}
