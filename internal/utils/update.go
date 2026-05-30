package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	versionCheckURL   = "https://api.github.com/repos/joaomjbraga/piunter/releases/latest"
	versionCacheTTL   = 24 * time.Hour
	httpTimeout       = 5 * time.Second
	skipUpdateCheckEnv = "PIUNTER_SKIP_UPDATE_CHECK"
)

type VersionCache struct {
	LastCheck        int64  `json:"last_check"`
	LatestVersion    string `json:"latest_version"`
	NotifiedVersion  string `json:"notified_version,omitempty"`
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

var getVersionCachePath = func() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "piunter", "version_cache.json")
}

func loadVersionCache() VersionCache {
	var cache VersionCache
	data, err := os.ReadFile(getVersionCachePath())
	if err != nil {
		return cache
	}
	json.Unmarshal(data, &cache)
	return cache
}

func saveVersionCache(cache VersionCache) {
	path := getVersionCachePath()
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	data, err := json.Marshal(cache)
	if err != nil {
		return
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		Debug(fmt.Sprintf("falha ao salvar cache de versão: %s", err))
	}
}

func isCacheStale(cache VersionCache) bool {
	if cache.LastCheck == 0 {
		return true
	}
	lastCheck := time.Unix(cache.LastCheck, 0)
	return time.Since(lastCheck) > versionCacheTTL
}

func fetchLatestVersion(currentVersion string) (string, error) {
	client := &http.Client{Timeout: httpTimeout}
	req, err := http.NewRequest("GET", versionCheckURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "piunter/"+currentVersion)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("json decode failed: %w", err)
	}

	if release.TagName == "" {
		return "", fmt.Errorf("empty tag_name")
	}

	return release.TagName, nil
}

func isNewerVersion(current, latest string) bool {
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	partsCur := strings.Split(current, ".")
	partsLat := strings.Split(latest, ".")

	for i := 0; i < 3; i++ {
		var vCur, vLat int
		if i < len(partsCur) {
			fmt.Sscanf(partsCur[i], "%d", &vCur)
		}
		if i < len(partsLat) {
			fmt.Sscanf(partsLat[i], "%d", &vLat)
		}
		if vCur != vLat {
			return vLat > vCur
		}
	}
	return false
}

func CheckForUpdate(currentVersion string) (string, error) {
	if os.Getenv(skipUpdateCheckEnv) != "" {
		return "", nil
	}

	cache := loadVersionCache()

	if !isCacheStale(cache) {
		if cache.LatestVersion != "" && isNewerVersion(currentVersion, cache.LatestVersion) {
			if cache.LatestVersion == cache.NotifiedVersion {
				return "", nil
			}
			return cache.LatestVersion, nil
		}
		return "", nil
	}

	latest, err := fetchLatestVersion(currentVersion)
	if err != nil {
		return "", err
	}

	saveVersionCache(VersionCache{
		LastCheck:     time.Now().Unix(),
		LatestVersion: latest,
	})

	if isNewerVersion(currentVersion, latest) {
		if latest == cache.NotifiedVersion {
			return "", nil
		}
		cache.NotifiedVersion = latest
		saveVersionCache(cache)
		return latest, nil
	}

	return "", nil
}
