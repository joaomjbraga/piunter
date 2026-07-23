package modules

import "testing"

func TestNewTempFilesModule(t *testing.T) {
	module := NewTempFilesModule()

	if module.ID() != "temp-files" {
		t.Fatalf("expected id 'temp-files', got %s", module.ID())
	}

	if module.Name() != "Arquivos Temporários" {
		t.Fatalf("expected name 'Arquivos Temporários', got %s", module.Name())
	}
}
