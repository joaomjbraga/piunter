package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version         string            `yaml:"version"`
	ThresholdMB     int               `yaml:"threshold_mb"`
	DisabledModules []string          `yaml:"disabled_modules"`
	ExcludePaths    []string          `yaml:"exclude_paths"`
	DryRunDefault   bool              `yaml:"dry_run_default"`
	DebugEnabled    bool              `yaml:"debug_enabled"`
	Parallel        bool              `yaml:"parallel"`
	PackageSizes    PackageSizeConfig `yaml:"package_sizes"`
}

type PackageSizeConfig struct {
	OrphanPackageMB  int64 `yaml:"orphan_package_mb"`
	FlatpakAppMB     int64 `yaml:"flatpak_app_mb"`
	SnapRevisionMB   int64 `yaml:"snap_revision_mb"`
}

var DefaultConfig = Config{
	Version:         "1.0",
	ThresholdMB:     100,
	DisabledModules: []string{},
	ExcludePaths:    []string{},
	DryRunDefault:   false,
	DebugEnabled:    false,
	Parallel:        false,
	PackageSizes: PackageSizeConfig{
		OrphanPackageMB:  10,
		FlatpakAppMB:     50,
		SnapRevisionMB:   200,
	},
}

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
)

var (
	configCache     Config
	configCacheOnce sync.Once
	configCacheErr  error
)

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "piunter", "config.yaml")
}

func LoadConfig() (Config, error) {
	configCacheOnce.Do(func() {
		configCache, configCacheErr = loadConfigFromFile()
	})
	return configCache, configCacheErr
}

func loadConfigFromFile() (Config, error) {
	configPath := GetConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

type ConfigValidator struct {
	cfg Config
}

func NewConfigValidator(cfg Config) *ConfigValidator {
	return &ConfigValidator{cfg: cfg}
}

func (v *ConfigValidator) Validate() error {
	if v.cfg.ThresholdMB < 1 || v.cfg.ThresholdMB > 100000 {
		return fmt.Errorf("threshold_mb must be between 1 and 100000")
	}

	for _, path := range v.cfg.ExcludePaths {
		if !filepath.IsAbs(path) {
			return fmt.Errorf("exclude_paths must be absolute paths: %s", path)
		}
	}

	return nil
}

func (v *ConfigValidator) IsModuleDisabled(moduleID string) bool {
	for _, disabled := range v.cfg.DisabledModules {
		if disabled == moduleID {
			return true
		}
	}
	return false
}

func (v *ConfigValidator) IsPathExcluded(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	for _, excluded := range v.cfg.ExcludePaths {
		if strings.Contains(absPath, excluded) {
			return true
		}
	}
	return false
}
