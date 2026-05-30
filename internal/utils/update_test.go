package utils

import (
	"os"
	"testing"
	"time"
)

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		current string
		latest  string
		want    bool
	}{
		{"1.5.0", "1.5.0", false},
		{"1.5.0", "1.5.1", true},
		{"1.5.0", "2.0.0", true},
		{"2.0.0", "1.9.9", false},
		{"1.5.0", "1.6.0", true},
		{"1.6.0", "1.5.0", false},
		{"1.5.0", "1.5.0", false},
		{"0.1.0", "0.2.0", true},
		{"1.5.0", "v1.6.0", true},
		{"v1.5.0", "1.6.0", true},
		{"v1.5.0", "v1.5.0", false},
		{"1.5", "1.5.0", false},
		{"1.5", "1.6.0", true},
		{"1.5.0", "1.5", false},
	}

	for _, tt := range tests {
		got := isNewerVersion(tt.current, tt.latest)
		if got != tt.want {
			t.Errorf("isNewerVersion(%q, %q) = %v, want %v", tt.current, tt.latest, got, tt.want)
		}
	}
}

func TestIsCacheStale(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		cache VersionCache
		want  bool
	}{
		{
			name:  "empty cache",
			cache: VersionCache{},
			want:  true,
		},
		{
			name:  "recent check",
			cache: VersionCache{LastCheck: now.Unix()},
			want:  false,
		},
		{
			name:  "23 hours ago",
			cache: VersionCache{LastCheck: now.Add(-23 * time.Hour).Unix()},
			want:  false,
		},
		{
			name:  "25 hours ago",
			cache: VersionCache{LastCheck: now.Add(-25 * time.Hour).Unix()},
			want:  true,
		},
		{
			name:  "one week ago",
			cache: VersionCache{LastCheck: now.Add(-7 * 24 * time.Hour).Unix()},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCacheStale(tt.cache)
			if got != tt.want {
				t.Errorf("isCacheStale() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckForUpdate_SkipEnvVar(t *testing.T) {
	os.Setenv(skipUpdateCheckEnv, "1")
	defer os.Unsetenv(skipUpdateCheckEnv)

	latest, err := CheckForUpdate("1.5.0")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if latest != "" {
		t.Errorf("expected empty latest, got %q", latest)
	}
}

func TestIsNewerVersion_EdgeCases(t *testing.T) {
	if !isNewerVersion("", "1.0.0") {
		t.Error("empty current should treat latest as newer")
	}

	if !isNewerVersion("0.0.9", "0.0.10") {
		t.Error("0.0.10 should be newer than 0.0.9")
	}

	if isNewerVersion("1.0.0", "") {
		t.Error("empty latest should not be considered newer")
	}
}

func TestNotifiedVersion_SupressesDuplicate(t *testing.T) {
	if !isNewerVersion("1.5.0", "v1.6.0") {
		t.Fatal("expected v1.6.0 to be newer than 1.5.0")
	}

	cache := VersionCache{
		LastCheck:        time.Now().Unix(),
		LatestVersion:    "v1.6.0",
		NotifiedVersion:  "v1.6.0",
	}

	orig := getVersionCachePath
	customPath := t.TempDir() + "/piunter/version_cache.json"
	getVersionCachePath = func() string { return customPath }
	t.Cleanup(func() { getVersionCachePath = orig })

	saveVersionCache(cache)

	latest, err := CheckForUpdate("1.5.0")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if latest != "" {
		t.Errorf("expected empty (already notified), got %q", latest)
	}
}

func TestSaveAndLoadCache(t *testing.T) {
	orig := getVersionCachePath
	customPath := t.TempDir() + "/piunter/version_cache.json"

	getVersionCachePath = func() string {
		return customPath
	}
	t.Cleanup(func() { getVersionCachePath = orig })

	cache := VersionCache{
		LastCheck:     1000,
		LatestVersion: "v1.6.0",
	}
	saveVersionCache(cache)

	loaded := loadVersionCache()
	if loaded.LastCheck != cache.LastCheck {
		t.Errorf("LastCheck: got %d, want %d", loaded.LastCheck, cache.LastCheck)
	}
	if loaded.LatestVersion != cache.LatestVersion {
		t.Errorf("LatestVersion: got %q, want %q", loaded.LatestVersion, cache.LatestVersion)
	}
}
