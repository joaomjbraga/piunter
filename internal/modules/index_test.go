package modules

import "testing"

func TestGetModuleIDs(t *testing.T) {
	ids := GetModuleIDs()

	expected := []string{
		"packages",
		"package-cache",
		"temp-files",
		"shell-history",
		"dev-cache",
		"browser-cache",
		"editor-cache",
		"media-cache",
		"game-cache",
		"container-cache",
		"build-cache",
		"ides-cache",
		"browser-plugins",
		"old-installers",
		"swap-files",
		"app-logs",
		"downloads-old",
		"cache",
		"flatpak",
		"snap",
		"docker",
		"logs",
		"large-files",
		"appimage",
		"thumbs",
		"recent",
		"trash",
	}

	if len(ids) != len(expected) {
		t.Fatalf("expected %d modules, got %d", len(expected), len(ids))
	}

	for i, id := range expected {
		if ids[i] != id {
			t.Fatalf("expected module %q at position %d, got %q", id, i, ids[i])
		}
	}
}

func TestGetAllModuleInfosIncludesRegisteredModules(t *testing.T) {
	infos := GetAllModuleInfos()
	if len(infos) == 0 {
		t.Fatal("expected module infos to be populated")
	}

	seen := make(map[string]bool)
	for _, info := range infos {
		seen[info.ID] = true
	}

	for _, id := range GetModuleIDs() {
		if !seen[id] {
			t.Fatalf("expected module %q to be present in module infos", id)
		}
	}
}
