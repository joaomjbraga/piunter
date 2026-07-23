package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/joaomjbraga/piunter/internal/modules"
	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type Cleaner struct {
	modules []modules.Module
	dryRun  bool
	verbose bool
}

func NewCleaner(moduleIds []string, dryRun bool) *Cleaner {
	mods := modules.GetModulesByIds(moduleIds)
	return &Cleaner{
		modules: mods,
		dryRun:  dryRun,
		verbose: false,
	}
}

func NewCleanerWithOptions(moduleIds []string, dryRun bool, verbose bool) *Cleaner {
	mods := modules.GetModulesByIds(moduleIds)
	return &Cleaner{
		modules: mods,
		dryRun:  dryRun,
		verbose: verbose,
	}
}

func (c *Cleaner) Clean() (*types.Report, error) {
	startTime := time.Now()
	results := c.cleanSequential()
	endTime := time.Now()

	return buildReport(startTime, endTime, results), nil
}

func isVerboseEnabled(verboseFlag bool, configVerbose bool) bool {
	return verboseFlag || configVerbose
}

func BuildConfirmationSummary(moduleIDs []string, dryRun bool) string {
	if len(moduleIDs) == 0 {
		return "Nenhum módulo selecionado."
	}

	base := fmt.Sprintf("Os seguintes módulos serão processados: %s.", strings.Join(moduleIDs, ", "))
	if dryRun {
		return base + " Execução em modo dry-run, nenhuma alteração será aplicada."
	}
	return base + " Confirmação necessária para aplicar as alterações."
}

func buildReport(startTime, endTime time.Time, results []types.CleaningResult) *types.Report {
	var totalSpaceFreed int64
	var totalItemsRemoved int
	var errors []string

	for _, r := range results {
		totalSpaceFreed += r.SpaceFreed
		totalItemsRemoved += r.ItemsRemoved
		errors = append(errors, r.Errors...)
	}

	return &types.Report{
		StartTime:         startTime.Format(time.RFC3339),
		EndTime:           endTime.Format(time.RFC3339),
		Modules:           results,
		TotalSpaceFreed:   totalSpaceFreed,
		TotalItemsRemoved: totalItemsRemoved,
		Errors:            errors,
	}
}

func (c *Cleaner) cleanSequential() []types.CleaningResult {
	var results []types.CleaningResult

	for _, m := range c.modules {
		if !m.IsAvailable() {
			continue
		}

		result, err := m.Clean(c.dryRun)
		if c.verbose {
			utils.Info(fmt.Sprintf("[verbose] módulo %s concluído", m.ID()))
		}
		if err != nil {
			results = append(results, types.CleaningResult{
				Module:       m.ID(),
				Success:      false,
				SpaceFreed:   0,
				ItemsRemoved: 0,
				Errors:       []string{fmt.Sprintf("%s: %s", m.Name(), err.Error())},
			})
			continue
		}
		results = append(results, *result)
	}

	return results
}

func summarizeReport(report *types.Report) string {
	if len(report.Errors) == 0 {
		return fmt.Sprintf("Limpeza concluída com sucesso. Espaço liberado: %s", utils.FormatBytes(report.TotalSpaceFreed))
	}
	return fmt.Sprintf("Limpeza concluída com avisos. Espaço liberado: %s. Erros: %d", utils.FormatBytes(report.TotalSpaceFreed), len(report.Errors))
}

func (c *Cleaner) PrintReport(report *types.Report) {
	startTime, err := time.Parse(time.RFC3339, report.StartTime)
	if err != nil {
		startTime = time.Now()
	}
	endTime, err := time.Parse(time.RFC3339, report.EndTime)
	if err != nil {
		endTime = time.Now()
	}
	duration := endTime.Sub(startTime)

	var durationStr string
	if duration.Seconds() < 60 {
		durationStr = fmt.Sprintf("%.1fs", duration.Seconds())
	} else {
		durationStr = fmt.Sprintf("%.1fmin", duration.Minutes())
	}

	for _, r := range report.Modules {
		if r.Success {
			utils.Item(r.Module, utils.FormatBytes(r.SpaceFreed))
		} else {
			utils.Item(r.Module, "\033[31merro\033[0m")
		}
	}

	utils.Space()
	fmt.Printf("  \033[90m%s\033[0m\n", strings.Repeat("─", 40))

	totalSize := utils.FormatBytes(report.TotalSpaceFreed)
	totalItems := fmt.Sprintf("%d", report.TotalItemsRemoved)
	totalErrors := fmt.Sprintf("%d", len(report.Errors))

	fmt.Println()
	fmt.Printf("  \033[1mResumo\033[0m\n")
	fmt.Printf("    \033[90m-\033[0m \033[37mEspaço liberado:\033[0m \033[32m%s\033[0m\n", totalSize)
	fmt.Printf("    \033[90m-\033[0m \033[37mItens removidos:\033[0m \033[36m%s\033[0m\n", totalItems)

	errorColor := "\033[32m"
	if len(report.Errors) > 0 {
		errorColor = "\033[31m"
	}
	fmt.Printf("    \033[90m-\033[0m \033[37mErros:\033[0m %s%s\033[0m\n", errorColor, totalErrors)

	if len(report.Errors) > 0 {
		utils.Space()
		fmt.Printf("  \033[33mAvisos:\033[0m\n")
		for _, err := range report.Errors {
			fmt.Printf("    \033[90m-\033[0m \033[33m%s\033[0m\n", err)
		}
	}

	utils.Space()

	fmt.Printf("  \033[32m*\033[0m %s\n", summarizeReport(report))
	if c.dryRun {
		fmt.Printf("  \033[90mExecute sem --dry-run para aplicar\033[0m\n")
	} else {
		fmt.Printf("  \033[90mTempo total: %s\033[0m\n", durationStr)
	}
	fmt.Println()
}