package core

import (
	"testing"

	"github.com/joaomjbraga/piunter/pkg/types"
)

func TestSummarizeReport(t *testing.T) {
	report := &types.Report{
		TotalSpaceFreed:   120,
		TotalItemsRemoved: 3,
		Errors:            []string{"permission denied"},
		Modules: []types.CleaningResult{{Module: "cache", Success: true, SpaceFreed: 120, ItemsRemoved: 3}},
	}

	summary := summarizeReport(report)
	if summary == "" {
		t.Fatal("expected non-empty summary")
	}
	if len(summary) == 0 {
		t.Fatal("expected summary content")
	}
}
