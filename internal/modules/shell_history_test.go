package modules

import "testing"

func TestNewShellHistoryModule(t *testing.T) {
	module := NewShellHistoryModule()

	if module.ID() != "shell-history" {
		t.Fatalf("expected id 'shell-history', got %s", module.ID())
	}

	if module.Name() != "Histórico de Shell" {
		t.Fatalf("expected name 'Histórico de Shell', got %s", module.Name())
	}
}
