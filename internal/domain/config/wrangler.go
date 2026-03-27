package config

import (
	"encoding/json"
	"reflect"
)

type WranglerConfig struct {
	Name               string                       `json:"name"`
	Main               string                       `json:"main"`
	CompatibilityDate  string                       `json:"compatibility_date"`
	CompatibilityFlags []string                     `json:"compatibility_flags,omitempty"`
	Observability      *WranglerObservability       `json:"observability,omitempty"`
	AI                 *WranglerAIBinding           `json:"ai,omitempty"`
	Vars               map[string]string            `json:"vars,omitempty"`
	KVNamespaces       []WranglerKVNamespace        `json:"kv_namespaces,omitempty"`
	D1Databases        []WranglerD1Database         `json:"d1_databases,omitempty"`
	R2Buckets          []WranglerR2Bucket           `json:"r2_buckets,omitempty"`
	Routes             []WranglerRoute              `json:"routes,omitempty"`
	Env                map[string]WranglerEnvConfig `json:"env,omitempty"`
	Extra              map[string]json.RawMessage   `json:"-"`
}

type WranglerObservability struct {
	Enabled bool `json:"enabled"`
}

type WranglerEnvConfig struct {
	Name               string                     `json:"name,omitempty"`
	AI                 *WranglerAIBinding         `json:"ai,omitempty"`
	Vars               map[string]string          `json:"vars,omitempty"`
	CompatibilityFlags []string                   `json:"compatibility_flags,omitempty"`
	KVNamespaces       []WranglerKVNamespace      `json:"kv_namespaces,omitempty"`
	D1Databases        []WranglerD1Database       `json:"d1_databases,omitempty"`
	R2Buckets          []WranglerR2Bucket         `json:"r2_buckets,omitempty"`
	Routes             []WranglerRoute            `json:"routes,omitempty"`
	Extra              map[string]json.RawMessage `json:"-"`
}

type WranglerAIBinding struct {
	Binding string `json:"binding"`
	Remote  bool   `json:"remote,omitempty"`
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

func (w *WranglerConfig) UnmarshalJSON(data []byte) error {
	type alias WranglerConfig
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	deleteKnownWranglerKeys(raw)

	*w = WranglerConfig(decoded)
	if len(raw) > 0 {
		w.Extra = raw
	} else {
		w.Extra = nil
	}
	return nil
}

func (w WranglerConfig) MarshalJSON() ([]byte, error) {
	raw := cloneRawMap(w.Extra)
	deleteKnownWranglerKeys(raw)

	putString(raw, "name", w.Name)
	putString(raw, "main", w.Main)
	putValue(raw, "compatibility_date", w.CompatibilityDate)
	putValue(raw, "compatibility_flags", w.CompatibilityFlags)
	putValue(raw, "observability", w.Observability)
	putValue(raw, "ai", w.AI)
	putValue(raw, "vars", w.Vars)
	putValue(raw, "kv_namespaces", w.KVNamespaces)
	putValue(raw, "d1_databases", w.D1Databases)
	putValue(raw, "r2_buckets", w.R2Buckets)
	putValue(raw, "routes", w.Routes)
	putValue(raw, "env", w.Env)

	return json.Marshal(raw)
}

func (w *WranglerEnvConfig) UnmarshalJSON(data []byte) error {
	type alias WranglerEnvConfig
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	deleteKnownWranglerEnvKeys(raw)

	*w = WranglerEnvConfig(decoded)
	if len(raw) > 0 {
		w.Extra = raw
	} else {
		w.Extra = nil
	}
	return nil
}

func (w WranglerEnvConfig) MarshalJSON() ([]byte, error) {
	raw := cloneRawMap(w.Extra)
	deleteKnownWranglerEnvKeys(raw)

	putString(raw, "name", w.Name)
	putValue(raw, "ai", w.AI)
	putValue(raw, "vars", w.Vars)
	putValue(raw, "compatibility_flags", w.CompatibilityFlags)
	putValue(raw, "kv_namespaces", w.KVNamespaces)
	putValue(raw, "d1_databases", w.D1Databases)
	putValue(raw, "r2_buckets", w.R2Buckets)
	putValue(raw, "routes", w.Routes)

	return json.Marshal(raw)
}

func cloneRawMap(source map[string]json.RawMessage) map[string]json.RawMessage {
	if len(source) == 0 {
		return map[string]json.RawMessage{}
	}
	cloned := make(map[string]json.RawMessage, len(source))
	for key, value := range source {
		cloned[key] = append(json.RawMessage(nil), value...)
	}
	return cloned
}

func deleteKnownWranglerKeys(raw map[string]json.RawMessage) {
	delete(raw, "name")
	delete(raw, "main")
	delete(raw, "compatibility_date")
	delete(raw, "compatibility_flags")
	delete(raw, "observability")
	delete(raw, "ai")
	delete(raw, "vars")
	delete(raw, "compatibility_flags")
	delete(raw, "kv_namespaces")
	delete(raw, "d1_databases")
	delete(raw, "r2_buckets")
	delete(raw, "routes")
	delete(raw, "env")
}

func deleteKnownWranglerEnvKeys(raw map[string]json.RawMessage) {
	delete(raw, "name")
	delete(raw, "ai")
	delete(raw, "vars")
	delete(raw, "kv_namespaces")
	delete(raw, "d1_databases")
	delete(raw, "r2_buckets")
	delete(raw, "routes")
}

func putString(raw map[string]json.RawMessage, key, value string) {
	if value == "" {
		delete(raw, key)
		return
	}
	putValue(raw, key, value)
}

func putValue(raw map[string]json.RawMessage, key string, value any) {
	if isEmptyJSONValue(value) {
		delete(raw, key)
		return
	}
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	raw[key] = data
}

func isEmptyJSONValue(value any) bool {
	if value == nil {
		return true
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice:
		if rv.IsNil() {
			return true
		}
	}

	switch typed := value.(type) {
	case string:
		return typed == ""
	case map[string]string:
		return len(typed) == 0
	case map[string]WranglerEnvConfig:
		return len(typed) == 0
	case []string:
		return len(typed) == 0
	case []WranglerKVNamespace:
		return len(typed) == 0
	case []WranglerD1Database:
		return len(typed) == 0
	case []WranglerR2Bucket:
		return len(typed) == 0
	case []WranglerRoute:
		return len(typed) == 0
	default:
		return false
	}
}
