package modules

import "testing"

func TestNewBrowserPluginsModule(t *testing.T) {
	module := NewBrowserPluginsModule()

	if module.ID() != "browser-plugins" {
		t.Fatalf("expected id 'browser-plugins', got %s", module.ID())
	}

	if module.Name() != "Plugins de Navegador" {
		t.Fatalf("expected name 'Plugins de Navegador', got %s", module.Name())
	}
}
