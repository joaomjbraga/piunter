package modules

import "testing"

func TestNewPackageCacheModule(t *testing.T) {
	module := NewPackageCacheModule()

	if module.ID() != "package-cache" {
		t.Fatalf("expected id 'package-cache', got %s", module.ID())
	}

	if module.Name() != "Cache de Pacotes" {
		t.Fatalf("expected name 'Cache de Pacotes', got %s", module.Name())
	}
}
