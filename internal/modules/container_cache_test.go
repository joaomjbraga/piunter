package modules

import "testing"

func TestNewContainerCacheModule(t *testing.T) {
	module := NewContainerCacheModule()

	if module.ID() != "container-cache" {
		t.Fatalf("expected id 'container-cache', got %s", module.ID())
	}

	if module.Name() != "Cache de Containers" {
		t.Fatalf("expected name 'Cache de Containers', got %s", module.Name())
	}
}
