package modules

import "testing"

func TestNewDevCacheModule(t *testing.T) {
	module := NewDevCacheModule()

	if module.ID() != "dev-cache" {
		t.Fatalf("expected id 'dev-cache', got %s", module.ID())
	}

	if module.Name() != "Cache de Desenvolvimento" {
		t.Fatalf("expected name 'Cache de Desenvolvimento', got %s", module.Name())
	}
}
