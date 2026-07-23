package modules

import "testing"

func TestNewOldInstallersModule(t *testing.T) {
	module := NewOldInstallersModule()

	if module.ID() != "old-installers" {
		t.Fatalf("expected id 'old-installers', got %s", module.ID())
	}

	if module.Name() != "Instaladores Antigos" {
		t.Fatalf("expected name 'Instaladores Antigos', got %s", module.Name())
	}
}
