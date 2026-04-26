package core

import (
	"fmt"
	"time"

	"github.com/joaomjbraga/piunter/internal/modules"
	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type Cleaner struct {
	modules  []modules.Module
	dryRun   bool
	parallel bool
}

func NewCleaner(moduleIds []string, dryRun bool) *Cleaner {
	mods := modules.GetModulesByIds(moduleIds)
	cfg, _ := utils.LoadConfig()
	return &Cleaner{
		modules:  mods,
		dryRun:   dryRun,
		parallel: cfg.Parallel,
	}
}

func (c *Cleaner) Clean() (*types.Report, error) {
	startTime := time.Now()
	var results []types.CleaningResult

	if c.parallel {
		results = c.cleanParallel()
	} else {
		results = c.cleanSequential()
	}

	endTime := time.Now()

	var totalSpaceFreed int64
	var totalItemsRemoved int
	var errors []string

	for _, r := range results {
		totalSpaceFreed += r.SpaceFreed
		totalItemsRemoved += r.ItemsRemoved
		errors = append(errors, r.Errors...)
	}

	return &types.Report{
		StartTime:        startTime.Format(time.RFC3339),
		EndTime:          endTime.Format(time.RFC3339),
		Modules:          results,
		TotalSpaceFreed:  totalSpaceFreed,
		TotalItemsRemoved: totalItemsRemoved,
		Errors:           errors,
	}, nil
}

func (c *Cleaner) cleanSequential() []types.CleaningResult {
	var results []types.CleaningResult

	for _, m := range c.modules {
		if !m.IsAvailable() {
			continue
		}

		result, err := m.Clean(c.dryRun)
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

func (c *Cleaner) cleanParallel() []types.CleaningResult {
	results := make([]types.CleaningResult, len(c.modules))
	resultChan := make(chan struct {
		index   int
		result types.CleaningResult
		err    error
	}, len(c.modules))

	workerCount := utils.GetOptimalWorkers(len(c.modules))

	jobChan := make(chan int, len(c.modules))
	for i := range c.modules {
		jobChan <- i
	}
	close(jobChan)

	for w := 0; w < workerCount; w++ {
		go func() {
			for idx := range jobChan {
				m := c.modules[idx]
				result, err := m.Clean(c.dryRun)
				if err != nil {
					resultChan <- struct {
						index   int
						result types.CleaningResult
						err    error
					}{idx, types.CleaningResult{
						Module:       m.ID(),
						Success:      false,
						SpaceFreed:   0,
						ItemsRemoved: 0,
						Errors:       []string{err.Error()},
					}, err}
				} else {
					resultChan <- struct {
						index   int
						result types.CleaningResult
						err    error
					}{idx, *result, nil}
				}
			}
		}()
	}

	for i := 0; i < len(c.modules); i++ {
		res := <-resultChan
		results[res.index] = res.result
	}

	return results
}

func (c *Cleaner) PrintReport(report *types.Report) {
	startTime, _ := time.Parse(time.RFC3339, report.StartTime)
	endTime, _ := time.Parse(time.RFC3339, report.EndTime)
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
	fmt.Printf("  \033[90m%s\033[0m\n", repeat("─", 40))

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
		fmt.Printf("  \033[1;31mErros:\033[0m\n")
		for _, err := range report.Errors {
			fmt.Printf("    \033[90m-\033[0m \033[31m%s\033[0m\n", err)
		}
	}

	utils.Space()

	if c.dryRun {
		fmt.Printf("  \033[33m!\033[0m Dry-run concluído\n")
		fmt.Printf("  \033[90mExecute sem --dry-run para aplicar\033[0m\n")
	} else {
		fmt.Printf("  \033[32m*\033[0m Limpeza concluída em %s\n", durationStr)
	}
	fmt.Println()
}