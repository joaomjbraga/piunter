package modules

import "testing"

func TestNewBrowserCacheModule(t *testing.T) {
	module := NewBrowserCacheModule()

	if module.ID() != "browser-cache" {
		t.Fatalf("expected id 'browser-cache', got %s", module.ID())
	}

	if module.Name() != "Cache de Navegador" {
		t.Fatalf("expected name 'Cache de Navegador', got %s", module.Name())
	}
}
