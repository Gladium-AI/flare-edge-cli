package config

import "testing"

func TestUpsertRouteReplacesExistingPattern(t *testing.T) {
	routes := []WranglerRoute{{Pattern: "example.com/*", ZoneName: "example.com"}}
	updated := UpsertRoute(routes, WranglerRoute{Pattern: "example.com/*", CustomDomain: true})
	if len(updated) != 1 {
		t.Fatalf("expected 1 route, got %d", len(updated))
	}
	if !updated[0].CustomDomain {
		t.Fatalf("expected route to be replaced with custom domain")
	}
}
