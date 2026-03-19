package route

import (
	"context"
	"fmt"

	"github.com/paolo/flare-edge-cli/internal/domain/config"
	"github.com/paolo/flare-edge-cli/internal/infra/cloudflare"
	"github.com/paolo/flare-edge-cli/internal/infra/configstore"
	"github.com/paolo/flare-edge-cli/internal/service/shared"
	"github.com/paolo/flare-edge-cli/internal/support/fs"
)

type Service struct {
	store    *configstore.Store
	fs       *fs.FileSystem
	wrangler *shared.WranglerExecutor
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
	Applied bool                   `json:"applied"`
	Routes  []config.WranglerRoute `json:"routes"`
	Deploy  *shared.CommandResult  `json:"deploy,omitempty"`
}

func NewService(store *configstore.Store, fs *fs.FileSystem, wrangler *shared.WranglerExecutor) *Service {
	return &Service{store: store, fs: fs, wrangler: wrangler}
}

func (s *Service) Attach(ctx context.Context, options AttachOptions) (Result, error) {
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
	result, err := s.apply(ctx, options.Dir, options.Env)
	if err != nil {
		return Result{}, err
	}
	result.Updated = true
	result.Routes = wranglerCfg.Routes
	return result, nil
}

func (s *Service) Domain(ctx context.Context, options DomainOptions) (Result, error) {
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
	result, err := s.apply(ctx, options.Dir, options.Env)
	if err != nil {
		return Result{}, err
	}
	result.Updated = true
	result.Routes = wranglerCfg.Routes
	return result, nil
}

func (s *Service) Detach(ctx context.Context, options DetachOptions) (Result, error) {
	project, wranglerCfg, err := shared.LoadProjectAndWrangler(options.Dir, s.store, s.fs)
	if err != nil {
		return Result{}, err
	}
	existingRoutes := append([]config.WranglerRoute(nil), wranglerCfg.Routes...)
	if options.Route != "" {
		wranglerCfg.Routes = config.RemoveRoute(wranglerCfg.Routes, options.Route, false)
	}
	if options.Hostname != "" {
		wranglerCfg.Routes = config.RemoveRoute(wranglerCfg.Routes, options.Hostname, true)
	}
	if err := shared.SaveWrangler(options.Dir, project, wranglerCfg, s.store); err != nil {
		return Result{}, err
	}

	client, err := s.cloudflareClient()
	if err != nil {
		return Result{}, err
	}
	if options.Route != "" {
		if err := s.deleteRouteTrigger(ctx, client, existingRoutes, options); err != nil {
			return Result{}, err
		}
	}
	if options.Hostname != "" {
		if err := s.deleteCustomDomain(ctx, client, options); err != nil {
			return Result{}, err
		}
	}
	return Result{Updated: true, Applied: options.Route != "" || options.Hostname != "", Routes: wranglerCfg.Routes}, nil
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

func (s *Service) apply(ctx context.Context, dir, env string) (Result, error) {
	raw, err := s.wrangler.Run(ctx, dir, env, "deploy")
	if err != nil {
		return Result{}, err
	}
	return Result{
		Applied: true,
		Deploy:  commandResultPtr(shared.NewCommandResult([]string{"deploy"}, raw)),
	}, nil
}

func (s *Service) cloudflareClient() (*cloudflare.Client, error) {
	token, err := s.wrangler.APIToken()
	if err != nil {
		return nil, err
	}
	return cloudflare.NewClient(token), nil
}

func (s *Service) deleteRouteTrigger(ctx context.Context, client *cloudflare.Client, existing []config.WranglerRoute, options DetachOptions) error {
	zoneRef := firstZoneReference(options.Zone, existing, options.Route, false)
	if zoneRef == "" {
		return fmt.Errorf("zone is required to detach route %q", options.Route)
	}
	zoneID, err := resolveZoneID(ctx, client, s.wrangler.AccountID(), zoneRef)
	if err != nil {
		return err
	}
	routes, err := client.ListRoutes(ctx, zoneID)
	if err != nil {
		return err
	}
	for _, route := range routes {
		if route.Pattern == options.Route {
			return client.DeleteRoute(ctx, zoneID, route.ID)
		}
	}
	return nil
}

func (s *Service) deleteCustomDomain(ctx context.Context, client *cloudflare.Client, options DetachOptions) error {
	accountID := s.wrangler.AccountID()
	if accountID == "" {
		return fmt.Errorf("account id is required to detach custom domain %q", options.Hostname)
	}
	records, err := client.ListDomainRecords(ctx, accountID, options.Hostname)
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Hostname == options.Hostname {
			return client.DeleteDomainRecord(ctx, accountID, record.ID)
		}
	}
	return nil
}

func resolveZoneID(ctx context.Context, client *cloudflare.Client, accountID, zoneRef string) (string, error) {
	if len(zoneRef) == 32 && !contains(zoneRef, ".") {
		return zoneRef, nil
	}
	return client.FindZoneID(ctx, accountID, zoneRef)
}

func firstZoneReference(explicit string, routes []config.WranglerRoute, pattern string, customDomain bool) string {
	if explicit != "" {
		return explicit
	}
	for _, route := range routes {
		if route.Pattern != pattern || route.CustomDomain != customDomain {
			continue
		}
		if route.ZoneID != "" {
			return route.ZoneID
		}
		if route.ZoneName != "" {
			return route.ZoneName
		}
	}
	return ""
}

func commandResultPtr(value shared.CommandResult) *shared.CommandResult {
	return &value
}
