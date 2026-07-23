package modules

import "testing"

func TestNewSwapFilesModule(t *testing.T) {
	module := NewSwapFilesModule()

	if module.ID() != "swap-files" {
		t.Fatalf("expected id 'swap-files', got %s", module.ID())
	}

	if module.Name() != "Arquivos Swap" {
		t.Fatalf("expected name 'Arquivos Swap', got %s", module.Name())
	}
}
