package core

import (
	"testing"
	"time"

	"github.com/joaomjbraga/piunter/pkg/types"
)

func TestBuildReportAggregatesResults(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Second)

	results := []types.CleaningResult{
		{Module: "cache", Success: true, SpaceFreed: 100, ItemsRemoved: 2, Errors: nil},
		{Module: "logs", Success: false, SpaceFreed: 0, ItemsRemoved: 0, Errors: []string{"failed"}},
	}

	report := buildReport(start, end, results)

	if report.TotalSpaceFreed != 100 {
		t.Fatalf("expected total space freed to be 100, got %d", report.TotalSpaceFreed)
	}

	if report.TotalItemsRemoved != 2 {
		t.Fatalf("expected total items removed to be 2, got %d", report.TotalItemsRemoved)
	}

	if len(report.Errors) != 1 || report.Errors[0] != "failed" {
		t.Fatalf("expected one aggregated error, got %v", report.Errors)
	}
}
