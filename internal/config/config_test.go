package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMergesFileAndEnvironment(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	content := `{
		"threshold_mb": 120,
		"all": false,
		"analyze": false,
		"dry_run": true,
		"force": false,
		"modules": ["cache", "docker"],
		"skip_update_check": true
	}`

	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	t.Setenv("PIUNTER_CONFIG_FILE", configPath)
	t.Setenv("PIUNTER_THRESHOLD", "250")
	t.Setenv("PIUNTER_DRY_RUN", "0")
	t.Setenv("PIUNTER_FORCE", "1")
	t.Setenv("PIUNTER_MODULES", "trash,logs")

	cfg := Load()

	if cfg.ThresholdMB != 250 {
		t.Fatalf("expected threshold 250, got %d", cfg.ThresholdMB)
	}

	if cfg.DryRun {
		t.Fatalf("expected dry_run to be false from environment override")
	}

	if !cfg.Force {
		t.Fatalf("expected force to be enabled from environment")
	}

	if cfg.SkipUpdateCheck != true {
		t.Fatalf("expected skip_update_check to be true from file")
	}

	if len(cfg.Modules) != 2 || cfg.Modules[0] != "trash" || cfg.Modules[1] != "logs" {
		t.Fatalf("expected modules to be overridden from env, got %v", cfg.Modules)
	}
}
