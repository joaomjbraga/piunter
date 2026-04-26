package modules

import (
	"testing"
)

func TestNewNpmModule(t *testing.T) {
	module := NewNpmModule()

	if module.ID() != "npm" {
		t.Errorf("expected id 'npm', got %s", module.ID())
	}

	if module.Name() != "NPM" {
		t.Errorf("expected name 'NPM', got %s", module.Name())
	}
}

func TestNewYarnModule(t *testing.T) {
	module := NewYarnModule()

	if module.ID() != "yarn" {
		t.Errorf("expected id 'yarn', got %s", module.ID())
	}
}

func TestNewPnpmModule(t *testing.T) {
	module := NewPnpmModule()

	if module.ID() != "pnpm" {
		t.Errorf("expected id 'pnpm', got %s", module.ID())
	}
}

func TestNewDockerModule(t *testing.T) {
	module := NewDockerModule()

	if module.ID() != "docker" {
		t.Errorf("expected id 'docker', got %s", module.ID())
	}

	if module.Name() != "Docker" {
		t.Errorf("expected name 'Docker', got %s", module.Name())
	}
}

func TestNewPackagesModule(t *testing.T) {
	module := NewPackagesModule()

	if module.ID() != "packages" {
		t.Errorf("expected id 'packages', got %s", module.ID())
	}
}

func TestNewLogsModule(t *testing.T) {
	module := NewLogsModule()

	if module.ID() != "logs" {
		t.Errorf("expected id 'logs', got %s", module.ID())
	}
}

func TestNewCacheModule(t *testing.T) {
	module := NewCacheModule()

	if module.ID() != "cache" {
		t.Errorf("expected id 'cache', got %s", module.ID())
	}
}

func TestNewFlatpakModule(t *testing.T) {
	module := NewFlatpakModule()

	if module.ID() != "flatpak" {
		t.Errorf("expected id 'flatpak', got %s", module.ID())
	}
}

func TestNewSnapModule(t *testing.T) {
	module := NewSnapModule()

	if module.ID() != "snap" {
		t.Errorf("expected id 'snap', got %s", module.ID())
	}
}