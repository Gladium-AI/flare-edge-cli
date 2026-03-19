package config

type WranglerConfig struct {
	Name              string                       `json:"name"`
	Main              string                       `json:"main"`
	CompatibilityDate string                       `json:"compatibility_date"`
	Observability     *WranglerObservability       `json:"observability,omitempty"`
	Vars              map[string]string            `json:"vars,omitempty"`
	KVNamespaces      []WranglerKVNamespace        `json:"kv_namespaces,omitempty"`
	D1Databases       []WranglerD1Database         `json:"d1_databases,omitempty"`
	R2Buckets         []WranglerR2Bucket           `json:"r2_buckets,omitempty"`
	Routes            []WranglerRoute              `json:"routes,omitempty"`
	Env               map[string]WranglerEnvConfig `json:"env,omitempty"`
}

type WranglerObservability struct {
	Enabled bool `json:"enabled"`
}

type WranglerEnvConfig struct {
	Name         string                `json:"name,omitempty"`
	Vars         map[string]string     `json:"vars,omitempty"`
	KVNamespaces []WranglerKVNamespace `json:"kv_namespaces,omitempty"`
	D1Databases  []WranglerD1Database  `json:"d1_databases,omitempty"`
	R2Buckets    []WranglerR2Bucket    `json:"r2_buckets,omitempty"`
	Routes       []WranglerRoute       `json:"routes,omitempty"`
}

type WranglerKVNamespace struct {
	Binding   string `json:"binding"`
	ID        string `json:"id,omitempty"`
	PreviewID string `json:"preview_id,omitempty"`
}

type WranglerD1Database struct {
	Binding      string `json:"binding"`
	DatabaseName string `json:"database_name,omitempty"`
	DatabaseID   string `json:"database_id,omitempty"`
}

type WranglerR2Bucket struct {
	Binding      string `json:"binding"`
	BucketName   string `json:"bucket_name,omitempty"`
	Jurisdiction string `json:"jurisdiction,omitempty"`
}

type WranglerRoute struct {
	Pattern      string `json:"pattern"`
	ZoneName     string `json:"zone_name,omitempty"`
	ZoneID       string `json:"zone_id,omitempty"`
	CustomDomain bool   `json:"custom_domain,omitempty"`
}
