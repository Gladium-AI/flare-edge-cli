package config

import "strings"

func UpsertRoute(routes []WranglerRoute, route WranglerRoute) []WranglerRoute {
	for index, existing := range routes {
		if existing.Pattern == route.Pattern {
			routes[index] = route
			return routes
		}
	}
	return append(routes, route)
}

func RemoveRoute(routes []WranglerRoute, pattern string, customDomain bool) []WranglerRoute {
	filtered := make([]WranglerRoute, 0, len(routes))
	for _, route := range routes {
		if route.Pattern == pattern && route.CustomDomain == customDomain {
			continue
		}
		filtered = append(filtered, route)
	}
	return filtered
}

func UpsertKV(namespaces []WranglerKVNamespace, namespace WranglerKVNamespace) []WranglerKVNamespace {
	for index, existing := range namespaces {
		if strings.EqualFold(existing.Binding, namespace.Binding) {
			namespaces[index] = namespace
			return namespaces
		}
	}
	return append(namespaces, namespace)
}

func UpsertD1(databases []WranglerD1Database, database WranglerD1Database) []WranglerD1Database {
	for index, existing := range databases {
		if strings.EqualFold(existing.Binding, database.Binding) {
			databases[index] = database
			return databases
		}
	}
	return append(databases, database)
}

func UpsertR2(buckets []WranglerR2Bucket, bucket WranglerR2Bucket) []WranglerR2Bucket {
	for index, existing := range buckets {
		if strings.EqualFold(existing.Binding, bucket.Binding) {
			buckets[index] = bucket
			return buckets
		}
	}
	return append(buckets, bucket)
}

func EnsureEnv(cfg *WranglerConfig, name string) *WranglerEnvConfig {
	if cfg.Env == nil {
		cfg.Env = map[string]WranglerEnvConfig{}
	}
	envCfg := cfg.Env[name]
	return &envCfg
}

func SetEnv(cfg *WranglerConfig, name string, envCfg WranglerEnvConfig) {
	if cfg.Env == nil {
		cfg.Env = map[string]WranglerEnvConfig{}
	}
	cfg.Env[name] = envCfg
}
