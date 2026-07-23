package modules

import "testing"

func TestNewEditorCacheModule(t *testing.T) {
	module := NewEditorCacheModule()

	if module.ID() != "editor-cache" {
		t.Fatalf("expected id 'editor-cache', got %s", module.ID())
	}

	if module.Name() != "Cache de Editores" {
		t.Fatalf("expected name 'Cache de Editores', got %s", module.Name())
	}
}
