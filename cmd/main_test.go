package main

import "testing"

func TestGetInteractiveSuggestionsMatchesPartialInput(t *testing.T) {
	infos := []moduleSelectionInfo{
		{ID: "cache", Name: "cache", Available: true},
		{ID: "package-cache", Name: "package cache", Available: true},
		{ID: "packages", Name: "packages", Available: true},
		{ID: "trash", Name: "trash", Available: false},
	}

	suggestions := getInteractiveSuggestions("pack", infos)
	if len(suggestions) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(suggestions))
	}

	if suggestions[0] != "package-cache" || suggestions[1] != "packages" {
		t.Fatalf("unexpected suggestions: %#v", suggestions)
	}
}

func TestGetModuleCategory(t *testing.T) {
	if got := getModuleCategory("cache"); got != "Cache e temporários" {
		t.Fatalf("expected cache category, got %q", got)
	}

	if got := getModuleCategory("docker"); got != "Containers e ambientes" {
		t.Fatalf("expected container category, got %q", got)
	}
}

func TestParseInteractiveSelectionSupportsNamesAndNumbers(t *testing.T) {
	infos := []moduleSelectionInfo{
		{ID: "cache", Name: "cache", Available: true},
		{ID: "packages", Name: "packages", Available: true},
		{ID: "trash", Name: "trash", Available: false},
	}

	selected := parseInteractiveSelection("1 packages cache", infos)
	if len(selected) != 2 {
		t.Fatalf("expected 2 selections, got %d", len(selected))
	}

	if selected[0] != "cache" || selected[1] != "packages" {
		t.Fatalf("unexpected selection order: %#v", selected)
	}
}

func TestParseInteractiveSelectionSupportsNormalizedNames(t *testing.T) {
	infos := []moduleSelectionInfo{
		{ID: "package-cache", Name: "package cache", Available: true},
	}

	selected := parseInteractiveSelection("package-cache", infos)
	if len(selected) != 1 || selected[0] != "package-cache" {
		t.Fatalf("expected normalized name selection, got %#v", selected)
	}
}

func TestRootCommandHasCompletionSubcommand(t *testing.T) {
	var found bool
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "completion" {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("expected root command to expose a completion subcommand")
	}
}

func TestFormatModuleList(t *testing.T) {
	cases := []struct {
		name      string
		moduleIDs []string
		want      string
	}{
		{
			name:      "single module",
			moduleIDs: []string{"cache"},
			want:      "cache",
		},
		{
			name:      "multiple modules",
			moduleIDs: []string{"cache", "docker", "trash"},
			want:      "cache, docker, trash",
		},
		{
			name:      "empty selection",
			moduleIDs: nil,
			want:      "nenhum módulo selecionado",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := formatModuleList(tc.moduleIDs); got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}
