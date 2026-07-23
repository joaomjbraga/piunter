package modules

import "testing"

func TestNewGameCacheModule(t *testing.T) {
	module := NewGameCacheModule()

	if module.ID() != "game-cache" {
		t.Fatalf("expected id 'game-cache', got %s", module.ID())
	}

	if module.Name() != "Cache de Jogos" {
		t.Fatalf("expected name 'Cache de Jogos', got %s", module.Name())
	}
}
