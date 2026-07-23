package modules

import "testing"

func TestNewDownloadsOldModule(t *testing.T) {
	module := NewDownloadsOldModule()

	if module.ID() != "downloads-old" {
		t.Fatalf("expected id 'downloads-old', got %s", module.ID())
	}

	if module.Name() != "Downloads Antigos" {
		t.Fatalf("expected name 'Downloads Antigos', got %s", module.Name())
	}
}
