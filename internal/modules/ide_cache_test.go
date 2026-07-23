package modules

import "testing"

func TestNewIdesCacheModule(t *testing.T) {
	module := NewIdesCacheModule()

	if module.ID() != "ides-cache" {
		t.Fatalf("expected id 'ides-cache', got %s", module.ID())
	}

	if module.Name() != "Cache de IDEs" {
		t.Fatalf("expected name 'Cache de IDEs', got %s", module.Name())
	}
}
