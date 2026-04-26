package core

import (
	"fmt"

	"github.com/joaomjbraga/piunter/internal/modules"
	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type Analyzer struct {
	modules   []modules.Module
	threshold int
}

func NewAnalyzer(moduleIds []string, threshold int) *Analyzer {
	var mods []modules.Module
	if len(moduleIds) == 0 {
		mods = modules.GetAvailableModules()
	} else {
		mods = modules.GetModulesByIds(moduleIds)
	}
	return &Analyzer{
		modules:   mods,
		threshold: threshold,
	}
}

func (a *Analyzer) Analyze() ([]*types.AnalysisResult, error) {
	var results []*types.AnalysisResult
	for _, m := range a.modules {
		if !m.IsAvailable() {
			continue
		}
		result, err := m.Analyze(a.threshold)
		if err != nil {
			utils.Debug(fmt.Sprintf("%s: %s", m.Name(), err.Error()))
			continue
		}
		results = append(results, result)
	}
	return results, nil
}

func (a *Analyzer) GetSummary(results []*types.AnalysisResult) struct {
	TotalSize   int64
	TotalItems  int
	ByModule    map[string]struct{ Size int64; Items int }
} {
	byModule := make(map[string]struct{ Size int64; Items int })
	var totalSize int64
	var totalItems int

	for _, result := range results {
		byModule[result.Module] = struct{ Size int64; Items int }{
			Size:  result.TotalSize,
			Items: len(result.Items),
		}
		totalSize += result.TotalSize
		totalItems += len(result.Items)
	}

	return struct {
		TotalSize   int64
		TotalItems  int
		ByModule    map[string]struct{ Size int64; Items int }
	}{
		TotalSize:  totalSize,
		TotalItems: totalItems,
		ByModule:   byModule,
	}
}

func (a *Analyzer) PrintAnalysis(results []*types.AnalysisResult) {
	summary := a.GetSummary(results)

	fmt.Printf("  \033[1mAnálise de espaço recuperável\033[0m\n\n")

	for _, result := range results {
		size := utils.FormatBytes(result.TotalSize)
		count := ""
		if len(result.Items) > 0 {
			count = fmt.Sprintf("(%d itens)", len(result.Items))
		}

		if len(result.Items) > 0 {
			fmt.Printf("    \033[90m-\033[0m %-20s \033[36m%s\033[0m %s\n", result.Module, size, count)
		} else {
			fmt.Printf("    \033[90m-\033[0m %-20s \033[90m0 B\033[0m\n", result.Module)
		}
	}

	utils.Space()
	fmt.Printf("  \033[90m%s\033[0m\n", repeat("─", 40))

	totalSize := utils.FormatBytes(summary.TotalSize)

	fmt.Println()
	fmt.Printf("  \033[1mTotal\033[0m\n")
	fmt.Printf("    \033[90m-\033[0m \033[37mEspaço:\033[0m \033[32m%s\033[0m\n", totalSize)
	fmt.Printf("    \033[90m-\033[0m \033[37mItens:\033[0m \033[36m%d\033[0m\n", summary.TotalItems)
	fmt.Println()
}

func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}