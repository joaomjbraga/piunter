package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Version         string            `yaml:"version"`
	ThresholdMB    int               `yaml:"threshold_mb"`
	DisabledModules []string         `yaml:"disabled_modules"`
	ExcludePaths   []string         `yaml:"exclude_paths"`
	DryRunDefault  bool             `yaml:"dry_run_default"`
	DebugEnabled  bool             `yaml:"debug_enabled"`
	Parallel      bool             `yaml:"parallel"`
	PackageSizes  PackageSizeConfig `yaml:"package_sizes"`
}

type PackageSizeConfig struct {
	OrphanPackageMB     int64 `yaml:"orphan_package_mb"`
	FlatpakAppMB      int64 `yaml:"flatpak_app_mb"`
	SnapRevisionMB    int64 `yaml:"snap_revision_mb"`
}

var DefaultConfig = Config{
	Version:         "1.0",
	ThresholdMB:     100,
	DisabledModules: []string{},
	ExcludePaths:    []string{},
	DryRunDefault:   false,
	DebugEnabled:   false,
	Parallel:      false,
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

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "piunter", "config.yaml")
}

func LoadConfig() (Config, error) {
	configPath := GetConfigPath()
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig, nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig, fmt.Errorf("failed to read config: %w", err)
	}
	
	var cfg Config
	if err := parseConfig(string(data), &cfg); err != nil {
		return DefaultConfig, fmt.Errorf("failed to parse config: %w", err)
	}
	
	return cfg, nil
}

func SaveConfig(cfg Config) error {
	configPath := GetConfigPath()
	dir := filepath.Dir(configPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	data := marshalConfig(cfg)

	if err := os.WriteFile(configPath, []byte(data), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func parseConfig(data string, cfg *Config) error {
	*cfg = DefaultConfig
	lines := parseConfigLines(data)
	
	for _, line := range lines {
		if line == "" || contains(line, "#") {
			continue
		}
		parts := splitConfigLine(line)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		switch key {
		case "threshold_mb":
			fmt.Sscanf(value, "%d", &cfg.ThresholdMB)
		case "dry_run_default":
			cfg.DryRunDefault = value == "true"
		case "debug_enabled":
			cfg.DebugEnabled = value == "true"
		case "parallel":
			cfg.Parallel = value == "true"
		case "orphan_package_mb":
			fmt.Sscanf(value, "%d", &cfg.PackageSizes.OrphanPackageMB)
		case "flatpak_app_mb":
			fmt.Sscanf(value, "%d", &cfg.PackageSizes.FlatpakAppMB)
		case "snap_revision_mb":
			fmt.Sscanf(value, "%d", &cfg.PackageSizes.SnapRevisionMB)
		}
	}
	
	return nil
}

func marshalConfig(cfg Config) string {
	var lines []string
	lines = append(lines, "# Configuração do Piunter")
	lines = append(lines, fmt.Sprintf("version: %s", cfg.Version))
	lines = append(lines, fmt.Sprintf("threshold_mb: %d", cfg.ThresholdMB))
	if cfg.DryRunDefault {
		lines = append(lines, "dry_run_default: true")
	}
	if cfg.DebugEnabled {
		lines = append(lines, "debug_enabled: true")
	}
	if cfg.Parallel {
		lines = append(lines, "parallel: true")
	}
	if len(cfg.DisabledModules) > 0 {
		lines = append(lines, "disabled_modules:")
		for _, m := range cfg.DisabledModules {
			lines = append(lines, fmt.Sprintf("  - %s", m))
		}
	}
	if len(cfg.ExcludePaths) > 0 {
		lines = append(lines, "exclude_paths:")
		for _, p := range cfg.ExcludePaths {
			lines = append(lines, fmt.Sprintf("  - %s", p))
		}
	}
	lines = append(lines, "# Tamanhos estimados por item (MB)")
	lines = append(lines, fmt.Sprintf("orphan_package_mb: %d", cfg.PackageSizes.OrphanPackageMB))
	lines = append(lines, fmt.Sprintf("flatpak_app_mb: %d", cfg.PackageSizes.FlatpakAppMB))
	lines = append(lines, fmt.Sprintf("snap_revision_mb: %d", cfg.PackageSizes.SnapRevisionMB))
	return stringsJoin(lines, "\n")
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func parseConfigLines(data string) []string {
	var lines []string
	var current string
	for _, r := range data {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func splitConfigLine(line string) []string {
	var result []string
	var current string
	var inList bool
	var listContent string
	
	for _, r := range line {
		if r == ':' && !inList {
			result = append(result, current)
			current = ""
			continue
		}
		if r == '-' && inList {
			listContent += string(r)
			continue
		}
		if r == '\n' {
			continue
		}
		current += string(r)
	}
	result = append(result, current)
	
	if len(result) > 1 && contains(result[1], "-") {
		inList = true
	}
	
	return result
}

func stringsJoin(slice []string, sep string) string {
	var result string
	for i, s := range slice {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
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
		if contains(absPath, excluded) {
			return true
		}
	}
	return false
}

type ConfigManager struct {
	cfg Config
}

func NewConfigManager() (*ConfigManager, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return &ConfigManager{cfg: cfg}, nil
}

func (m *ConfigManager) Get() Config {
	return m.cfg
}

func (m *ConfigManager) Update(updateFn func(*Config)) {
	updateFn(&m.cfg)
}

func (m *ConfigManager) Save() error {
	return SaveConfig(m.cfg)
}

func (m *ConfigManager) GetThresholdMB() int {
	if m.cfg.ThresholdMB == 0 {
		return DefaultConfig.ThresholdMB
	}
	return m.cfg.ThresholdMB
}

func (m *ConfigManager) IsModuleEnabled(moduleID string) bool {
	validator := NewConfigValidator(m.cfg)
	return !validator.IsModuleDisabled(moduleID)
}

func (m *ConfigManager) ShouldRunParallel() bool {
	return m.cfg.Parallel
}

type ConfigLogger struct {
	enabled bool
}

func (l *ConfigLogger) SetEnabled(enabled bool) {
	l.enabled = enabled
}

func (l *ConfigLogger) IsEnabled() bool {
	return l.enabled
}

func (l *ConfigLogger) Log(msg string) {
	if l.enabled {
		fmt.Printf("\033[36m[CONFIG %s]\033[0m %s\n", time.Now().Format("15:04:05"), msg)
	}
}