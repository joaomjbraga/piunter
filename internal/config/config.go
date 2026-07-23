package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DefaultConfigFile = "piunter.json"
	EnvConfigFile     = "PIUNTER_CONFIG_FILE"
	EnvThreshold      = "PIUNTER_THRESHOLD"
	EnvAll            = "PIUNTER_ALL"
	EnvAnalyze        = "PIUNTER_ANALYZE"
	EnvDryRun         = "PIUNTER_DRY_RUN"
	EnvForce          = "PIUNTER_FORCE"
	EnvModules        = "PIUNTER_MODULES"
	EnvSkipUpdate     = "PIUNTER_SKIP_UPDATE_CHECK"
)

type Config struct {
	ThresholdMB     int      `json:"threshold_mb"`
	All             bool     `json:"all"`
	Analyze         bool     `json:"analyze"`
	DryRun          bool     `json:"dry_run"`
	Force           bool     `json:"force"`
	Verbose         bool     `json:"verbose"`
	Modules         []string `json:"modules"`
	SkipUpdateCheck bool     `json:"skip_update_check"`
}

func Default() Config {
	return Config{
		ThresholdMB: 100,
		All:         false,
		Analyze:     false,
		DryRun:      false,
		Force:       false,
		Verbose:     false,
		Modules:     nil,
		SkipUpdateCheck: false,
	}
}

func Load() Config {
	cfg := Default()

	if path := configFilePath(); path != "" {
		fileCfg, err := loadFile(path)
		if err == nil {
			cfg = mergeConfig(cfg, fileCfg)
		} else {
			fmt.Fprintf(os.Stderr, "[config] failed to load %s: %v\n", path, err)
		}
	}

	cfg = mergeEnv(cfg)
	return cfg
}

func configFilePath() string {
	if path := os.Getenv(EnvConfigFile); path != "" {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "piunter", DefaultConfigFile)
}

func loadFile(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func mergeConfig(base, override Config) Config {
	if override.ThresholdMB != 0 {
		base.ThresholdMB = override.ThresholdMB
	}
	if override.All {
		base.All = override.All
	}
	if override.Analyze {
		base.Analyze = override.Analyze
	}
	if override.DryRun {
		base.DryRun = override.DryRun
	}
	if override.Force {
		base.Force = override.Force
	}
	if override.Verbose {
		base.Verbose = override.Verbose
	}
	if len(override.Modules) > 0 {
		base.Modules = override.Modules
	}
	if override.SkipUpdateCheck {
		base.SkipUpdateCheck = override.SkipUpdateCheck
	}
	return base
}

func mergeEnv(base Config) Config {
	if value := os.Getenv(EnvThreshold); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			base.ThresholdMB = parsed
		}
	}

	if value := os.Getenv(EnvAll); value != "" {
		base.All = parseBool(value)
	}
	if value := os.Getenv(EnvAnalyze); value != "" {
		base.Analyze = parseBool(value)
	}
	if value := os.Getenv(EnvDryRun); value != "" {
		base.DryRun = parseBool(value)
	}
	if value := os.Getenv(EnvForce); value != "" {
		base.Force = parseBool(value)
	}
	if value := os.Getenv("PIUNTER_VERBOSE"); value != "" {
		base.Verbose = parseBool(value)
	}
	if value := os.Getenv(EnvModules); value != "" {
		base.Modules = splitModules(value)
	}
	if value := os.Getenv(EnvSkipUpdate); value != "" {
		base.SkipUpdateCheck = parseBool(value)
	}
	return base
}

func parseBool(value string) bool {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "1" || value == "true" || value == "yes" || value == "on" {
		return true
	}
	return false
}

func splitModules(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
