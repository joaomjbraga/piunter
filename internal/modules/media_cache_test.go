package modules

import "testing"

func TestNewMediaCacheModule(t *testing.T) {
	module := NewMediaCacheModule()

	if module.ID() != "media-cache" {
		t.Fatalf("expected id 'media-cache', got %s", module.ID())
	}

	if module.Name() != "Cache de Mídia" {
		t.Fatalf("expected name 'Cache de Mídia', got %s", module.Name())
	}
}
