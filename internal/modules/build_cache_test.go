package modules

import "testing"

func TestNewBuildCacheModule(t *testing.T) {
	module := NewBuildCacheModule()

	if module.ID() != "build-cache" {
		t.Fatalf("expected id 'build-cache', got %s", module.ID())
	}

	if module.Name() != "Cache de Build" {
		t.Fatalf("expected name 'Cache de Build', got %s", module.Name())
	}
}
