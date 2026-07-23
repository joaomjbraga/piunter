package core

import (
	"strings"
	"testing"
)

func TestBuildConfirmationSummary(t *testing.T) {
	modules := []string{"cache", "trash"}
	summary := BuildConfirmationSummary(modules, true)

	if !strings.Contains(summary, "cache") || !strings.Contains(summary, "trash") {
		t.Fatalf("expected summary to mention selected modules, got %q", summary)
	}
	if !strings.Contains(summary, "dry-run") {
		t.Fatalf("expected summary to mention dry-run mode, got %q", summary)
	}
}
