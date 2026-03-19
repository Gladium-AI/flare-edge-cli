package deploy

import "testing"

func TestSplitPair(t *testing.T) {
	key, value, ok := splitPair("FOO=bar")
	if !ok || key != "FOO" || value != "bar" {
		t.Fatalf("unexpected split result: %q %q %v", key, value, ok)
	}
}
