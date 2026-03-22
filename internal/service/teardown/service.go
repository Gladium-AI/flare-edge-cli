package teardown

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

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

type Options struct {
	Dir           string
	Env           string
	Name          string
	KeepBindings  bool
	KeepArtifacts bool
	DeleteProject bool
}

type Result struct {
	Worker               string   `json:"worker"`
	WorkerDeleted        bool     `json:"worker_deleted"`
	RoutesDeleted        []string `json:"routes_deleted,omitempty"`
	CustomDomainsDeleted []string `json:"custom_domains_deleted,omitempty"`
	KVDeleted            []string `json:"kv_deleted,omitempty"`
	D1Deleted            []string `json:"d1_deleted,omitempty"`
	R2Deleted            []string `json:"r2_deleted,omitempty"`
	LocalRemoved         []string `json:"local_removed,omitempty"`
	ConfigUpdated        bool     `json:"config_updated"`
}

func NewService(store *configstore.Store, fs *fs.FileSystem, wrangler *shared.WranglerExecutor) *Service {
	return &Service{store: store, fs: fs, wrangler: wrangler}
}

func (s *Service) Run(ctx context.Context, options Options) (Result, error) {
	project, wranglerCfg, err := shared.LoadProjectAndWrangler(options.Dir, s.store, s.fs)
	if err != nil {
		return Result{}, err
	}

	result := Result{Worker: firstNonEmpty(options.Name, wranglerCfg.Name, project.WorkerName)}
	var failures []string

	client, accountID, err := s.cloudflare(ctx)
	if err != nil {
		failures = append(failures, "cloudflare client: "+err.Error())
	}

	if client != nil {
		routesDeleted, customDeleted, routeFailures := s.deleteTriggers(ctx, client, accountID, wranglerCfg.Routes)
		result.RoutesDeleted = routesDeleted
		result.CustomDomainsDeleted = customDeleted
		failures = append(failures, routeFailures...)
	}

	if result.Worker != "" {
		if deleteErr := s.deleteWorker(ctx, options.Dir, options.Env, result.Worker); deleteErr != nil {
			failures = append(failures, "worker delete: "+deleteErr.Error())
		} else {
			result.WorkerDeleted = true
		}
	}

	if !options.KeepBindings {
		kvDeleted, d1Deleted, r2Deleted, bindingFailures := s.deleteBindings(ctx, options, project, wranglerCfg)
		result.KVDeleted = kvDeleted
		result.D1Deleted = d1Deleted
		result.R2Deleted = r2Deleted
		failures = append(failures, bindingFailures...)
	}

	if options.DeleteProject {
		if err := s.fs.RemoveAll(options.Dir); err != nil {
			failures = append(failures, "remove project: "+err.Error())
		} else {
			result.LocalRemoved = append(result.LocalRemoved, options.Dir)
		}
	} else {
		updatedProject, updatedWrangler := resetLocalState(project, wranglerCfg, options.KeepBindings)
		if err := s.store.SaveProject(options.Dir, updatedProject); err != nil {
			failures = append(failures, "save project config: "+err.Error())
		} else if err := shared.SaveWrangler(options.Dir, updatedProject, updatedWrangler, s.store); err != nil {
			failures = append(failures, "save wrangler config: "+err.Error())
		} else {
			result.ConfigUpdated = true
		}
		if !options.KeepArtifacts {
			paths := []string{
				filepath.Join(options.Dir, updatedProject.OutDir),
				filepath.Join(options.Dir, ".wrangler"),
			}
			for _, path := range paths {
				if err := s.fs.RemoveAll(path); err != nil {
					failures = append(failures, "remove "+path+": "+err.Error())
					continue
				}
				result.LocalRemoved = append(result.LocalRemoved, path)
			}
		}
	}

	if len(failures) > 0 {
		return result, fmt.Errorf("teardown completed with errors: %s", strings.Join(failures, "; "))
	}
	return result, nil
}

func (s *Service) cloudflare(ctx context.Context) (*cloudflare.Client, string, error) {
	token, err := s.wrangler.APIToken()
	if err != nil {
		return nil, "", err
	}
	return cloudflare.NewClient(token), s.wrangler.AccountID(), nil
}

func (s *Service) deleteTriggers(ctx context.Context, client *cloudflare.Client, accountID string, routes []config.WranglerRoute) ([]string, []string, []string) {
	var routesDeleted []string
	var customDeleted []string
	var failures []string
	for _, route := range routes {
		if route.CustomDomain {
			deleted, err := s.deleteCustomDomain(ctx, client, accountID, route.Pattern)
			if err != nil {
				failures = append(failures, "custom domain "+route.Pattern+": "+err.Error())
				continue
			}
			if deleted {
				customDeleted = append(customDeleted, route.Pattern)
			}
			continue
		}
		deleted, err := s.deleteRoute(ctx, client, accountID, route)
		if err != nil {
			failures = append(failures, "route "+route.Pattern+": "+err.Error())
			continue
		}
		if deleted {
			routesDeleted = append(routesDeleted, route.Pattern)
		}
	}
	return routesDeleted, customDeleted, failures
}

func (s *Service) deleteRoute(ctx context.Context, client *cloudflare.Client, accountID string, route config.WranglerRoute) (bool, error) {
	zoneID, err := resolveRouteZoneID(ctx, client, accountID, route)
	if err != nil {
		return false, err
	}
	items, err := client.ListRoutes(ctx, zoneID)
	if err != nil {
		return false, err
	}
	for _, item := range items {
		if item.Pattern == route.Pattern {
			if err := client.DeleteRoute(ctx, zoneID, item.ID); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) deleteCustomDomain(ctx context.Context, client *cloudflare.Client, accountID, hostname string) (bool, error) {
	if accountID == "" {
		return false, fmt.Errorf("account id is required")
	}
	items, err := client.ListDomainRecords(ctx, accountID, hostname)
	if err != nil {
		return false, err
	}
	for _, item := range items {
		if item.Hostname == hostname {
			if err := client.DeleteDomainRecord(ctx, accountID, item.ID); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) deleteWorker(ctx context.Context, dir, env, name string) error {
	_, err := s.wrangler.Run(ctx, dir, env, "delete", "--name", name)
	return err
}

func (s *Service) deleteBindings(ctx context.Context, options Options, project config.Project, wranglerCfg config.WranglerConfig) ([]string, []string, []string, []string) {
	var kvDeleted []string
	var d1Deleted []string
	var r2Deleted []string
	var failures []string

	for _, namespace := range uniqueKVNamespaces(project, wranglerCfg) {
		if namespace.ID == "" {
			continue
		}
		if _, err := s.wrangler.Run(ctx, options.Dir, options.Env, "kv", "namespace", "delete", "--namespace-id", namespace.ID); err != nil {
			failures = append(failures, "kv "+namespace.Binding+": "+err.Error())
			continue
		}
		kvDeleted = append(kvDeleted, firstNonEmpty(namespace.Binding, namespace.ID))
	}

	for _, database := range uniqueD1Databases(project, wranglerCfg) {
		name := firstNonEmpty(database.DatabaseName, database.Binding)
		if name == "" {
			continue
		}
		if _, err := s.wrangler.Run(ctx, options.Dir, options.Env, "d1", "delete", name, "--skip-confirmation"); err != nil {
			failures = append(failures, "d1 "+name+": "+err.Error())
			continue
		}
		d1Deleted = append(d1Deleted, name)
	}

	for _, bucket := range uniqueR2Buckets(project, wranglerCfg) {
		name := firstNonEmpty(bucket.BucketName, bucket.Binding)
		if name == "" {
			continue
		}
		if _, err := s.wrangler.Run(ctx, options.Dir, options.Env, "r2", "bucket", "delete", name); err != nil {
			failures = append(failures, "r2 "+name+": "+err.Error())
			continue
		}
		r2Deleted = append(r2Deleted, name)
	}

	return kvDeleted, d1Deleted, r2Deleted, failures
}

func resetLocalState(project config.Project, wranglerCfg config.WranglerConfig, keepBindings bool) (config.Project, config.WranglerConfig) {
	project.Bindings.Secrets = nil
	wranglerCfg.Routes = nil
	if keepBindings {
		return project, wranglerCfg
	}
	project.Bindings.AI = nil
	project.Bindings.KV = nil
	project.Bindings.D1 = nil
	project.Bindings.R2 = nil
	wranglerCfg.AI = nil
	wranglerCfg.KVNamespaces = nil
	wranglerCfg.D1Databases = nil
	wranglerCfg.R2Buckets = nil
	return project, wranglerCfg
}

func resolveRouteZoneID(ctx context.Context, client *cloudflare.Client, accountID string, route config.WranglerRoute) (string, error) {
	if route.ZoneID != "" {
		return route.ZoneID, nil
	}
	if route.ZoneName != "" {
		return client.FindZoneID(ctx, accountID, route.ZoneName)
	}
	host := routeHost(route.Pattern)
	zone, err := client.FindZoneByHostname(ctx, accountID, host)
	if err != nil {
		return "", err
	}
	return zone.ID, nil
}

func routeHost(pattern string) string {
	host := pattern
	if slash := strings.Index(host, "/"); slash >= 0 {
		host = host[:slash]
	}
	return strings.TrimSpace(host)
}

func uniqueKVNamespaces(project config.Project, wranglerCfg config.WranglerConfig) []config.WranglerKVNamespace {
	seen := map[string]bool{}
	var items []config.WranglerKVNamespace
	for _, item := range wranglerCfg.KVNamespaces {
		if item.ID == "" || seen[item.ID] {
			continue
		}
		seen[item.ID] = true
		items = append(items, item)
	}
	for _, item := range project.Bindings.KV {
		if item.ID == "" || seen[item.ID] {
			continue
		}
		seen[item.ID] = true
		items = append(items, config.WranglerKVNamespace{Binding: item.Binding, ID: item.ID})
	}
	return items
}

func uniqueD1Databases(project config.Project, wranglerCfg config.WranglerConfig) []config.WranglerD1Database {
	seen := map[string]bool{}
	var items []config.WranglerD1Database
	for _, item := range wranglerCfg.D1Databases {
		key := firstNonEmpty(item.DatabaseName, item.DatabaseID, item.Binding)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		items = append(items, item)
	}
	for _, item := range project.Bindings.D1 {
		key := firstNonEmpty(item.DatabaseName, item.DatabaseID, item.Binding)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		items = append(items, config.WranglerD1Database{Binding: item.Binding, DatabaseName: item.DatabaseName, DatabaseID: item.DatabaseID})
	}
	return items
}

func uniqueR2Buckets(project config.Project, wranglerCfg config.WranglerConfig) []config.WranglerR2Bucket {
	seen := map[string]bool{}
	var items []config.WranglerR2Bucket
	for _, item := range wranglerCfg.R2Buckets {
		key := firstNonEmpty(item.BucketName, item.Binding)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		items = append(items, item)
	}
	for _, item := range project.Bindings.R2 {
		key := firstNonEmpty(item.BucketName, item.Binding)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		items = append(items, config.WranglerR2Bucket{Binding: item.Binding, BucketName: item.BucketName})
	}
	return items
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
