package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultConfig(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.ThresholdMB != 100 {
		t.Errorf("expected default threshold 100, got %d", cfg.ThresholdMB)
	}

	if cfg.Version != "1.0" {
		t.Errorf("expected version 1.0, got %s", cfg.Version)
	}
}

func TestGetConfigPath(t *testing.T) {
	configPath := GetConfigPath()
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".config", "piunter", "config.yaml")

	if configPath != expected {
		t.Errorf("expected %s, got %s", expected, configPath)
	}
}

func TestConfigValidator_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg    Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: Config{
				ThresholdMB:  100,
				ExcludePaths: []string{},
			},
			wantErr: false,
		},
		{
			name: "threshold too low",
			cfg: Config{
				ThresholdMB:  0,
			},
			wantErr: true,
		},
		{
			name: "threshold too high",
			cfg: Config{
				ThresholdMB: 200000,
			},
			wantErr: true,
		},
		{
			name: "invalid exclude path",
			cfg: Config{
				ThresholdMB: 100,
				ExcludePaths: []string{"relative/path"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewConfigValidator(tt.cfg)
			err := validator.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigValidator_IsModuleDisabled(t *testing.T) {
	validator := NewConfigValidator(Config{
		DisabledModules: []string{"npm", "docker"},
	})

	tests := []struct {
		moduleID string
		want    bool
	}{
		{"npm", true},
		{"docker", true},
		{"cache", false},
		{"packages", false},
	}

	for _, tt := range tests {
		t.Run(tt.moduleID, func(t *testing.T) {
			if got := validator.IsModuleDisabled(tt.moduleID); got != tt.want {
				t.Errorf("IsModuleDisabled(%s) = %v, want %v", tt.moduleID, got, tt.want)
			}
		})
	}
}

func TestConfigValidator_IsPathExcluded(t *testing.T) {
	validator := NewConfigValidator(Config{
		ExcludePaths: []string{"/home/user/.cache", "/tmp"},
	})

	tests := []struct {
		path  string
		want  bool
	}{
		{"/home/user/.cache/npm", true},
		{"/tmp/piunter", true},
		{"/home/user/documents", false},
		{"/var/log", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := validator.IsPathExcluded(tt.path); got != tt.want {
				t.Errorf("IsPathExcluded(%s) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestConfigManager_GetThresholdMB(t *testing.T) {
	tests := []struct {
		name        string
		cfg         Config
		want        int
	}{
		{
			name:        "custom threshold",
			cfg:         Config{ThresholdMB: 500},
			want:        500,
		},
		{
			name:        "zero uses default",
			cfg:         Config{ThresholdMB: 0},
			want:        100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ConfigManager{cfg: tt.cfg}
			if got := m.GetThresholdMB(); got != tt.want {
				t.Errorf("GetThresholdMB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigManager_IsModuleEnabled(t *testing.T) {
	m := &ConfigManager{
		cfg: Config{
			DisabledModules: []string{"npm"},
		},
	}

	if !m.IsModuleEnabled("docker") {
		t.Error("expected docker to be enabled")
	}

	if m.IsModuleEnabled("npm") {
		t.Error("expected npm to be disabled")
	}
}

func TestPackageSizeConstants(t *testing.T) {
	if KB != 1024 {
		t.Errorf("expected KB = 1024, got %d", KB)
	}
	if MB != 1024*1024 {
		t.Errorf("expected MB = 1024*1024, got %d", MB)
	}
	if GB != 1024*1024*1024 {
		t.Errorf("expected GB = 1024*1024*1024, got %d", GB)
	}
}