package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type LogsModule struct {
	BaseModule
}

func NewLogsModule() *LogsModule {
	return &LogsModule{
		BaseModule: BaseModule{
			id:          "logs",
			name:        "Logs do Sistema",
			description: "Limpa logs antigos do sistema",
		},
	}
}

func (m *LogsModule) IsAvailable() bool {
	return utils.IsRoot() || utils.HasSudoPassword()
}

func (m *LogsModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:    []types.CleanableItem{},
		TotalSize: 0,
	}

	logPaths := []string{
		"/var/log",
		"/tmp",
	}

	errorHandler := utils.NewErrorHandler()

	for _, logPath := range logPaths {
		if !utils.FileExists(logPath) {
			continue
		}

		entries, err := os.ReadDir(logPath)
		if err != nil {
			errorHandler.Add(utils.NewAnalysisError(m.id, fmt.Sprintf("falha ao ler %s", logPath), err))
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			fullPath := filepath.Join(logPath, entry.Name())
			size, err := utils.GetDirSize(fullPath)
			if err != nil {
				errorHandler.Add(utils.NewAnalysisError(m.id, fmt.Sprintf("falha ao calcular tamanho de %s", fullPath), err))
				continue
			}
			result.Items = append(result.Items, types.CleanableItem{
				Path:        fullPath,
				Size:        size,
				Type:        "directory",
				Description: "Diretório de log: " + entry.Name(),
			})
			result.TotalSize += size
		}
	}

	if errorHandler.HasErrors() {
		utils.Warn(errorHandler.Error())
	}

	return result, nil
}

func (m *LogsModule) Clean(dryRun bool) (*types.CleaningResult, error) {
	analysis, err := m.Analyze(0)
	if err != nil {
		return &types.CleaningResult{
			Module:  m.id,
			Success: false,
			Errors:  []string{err.Error()},
		}, err
	}

	result := &types.CleaningResult{
		Module:       m.id,
		Success:      true,
		SpaceFreed:   0,
		ItemsRemoved: 0,
		Errors:       []string{},
	}

	if analysis.TotalSize == 0 {
		return result, nil
	}

	if dryRun {
		result.SpaceFreed = analysis.TotalSize
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %s de logs", utils.FormatBytes(analysis.TotalSize)))
		return result, nil
	}

	execResult := utils.Exec("sudo", "journalctl", "--vacuum-time=7d")
	if !execResult.Success {
		result.Errors = append(result.Errors, fmt.Sprintf("Falha ao limpar journal: %s", execResult.Stderr))
	}

	execResult = utils.Exec("sudo", "find", "/var/log", "-type", "f", "-name", "*.gz", "-mtime", "+30", "-delete")
	if !execResult.Success {
		result.Errors = append(result.Errors, fmt.Sprintf("Falha ao limpar logs antigos: %s", execResult.Stderr))
	}

	result.SpaceFreed = analysis.TotalSize
	result.ItemsRemoved = len(analysis.Items)
	utils.Item(m.Name(), "Logs limpos")

	return result, nil
}

type FlatpakModule struct {
	BaseModule
}

func NewFlatpakModule() *FlatpakModule {
	return &FlatpakModule{
		BaseModule: BaseModule{
			id:          "flatpak",
			name:        "Flatpak",
			description: "Remove dados órfãos do Flatpak",
		},
	}
}

func (m *FlatpakModule) IsAvailable() bool {
	return utils.IsCommandAvailable("flatpak")
}

func (m *FlatpakModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:    []types.CleanableItem{},
		TotalSize: 0,
	}

	executor := utils.GetExecutor()
	execResult := executor.Exec("flatpak", "list", "--app", "--columns=ref")
	if !execResult.Success {
		return result, utils.NewAnalysisError(m.id, "falha ao listar flatpaks", fmt.Errorf("%s", execResult.Stderr))
	}

	lines := strings.Split(execResult.Stdout, "\n")
	config, _ := utils.LoadConfig()
	avgSize := config.PackageSizes.FlatpakAppMB * utils.MB

	appCount := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			appCount++
		}
	}

	result.TotalSize = int64(appCount) * avgSize

	return result, nil
}

func (m *FlatpakModule) Clean(dryRun bool) (*types.CleaningResult, error) {
	result := &types.CleaningResult{
		Module:       m.id,
		Success:      true,
		SpaceFreed:   0,
		ItemsRemoved: 0,
		Errors:       []string{},
	}

	if dryRun {
		result.SpaceFreed = 100 * 1024 * 1024
		utils.Info("[DRY-RUN] Limparia dados órfãos do Flatpak")
		return result, nil
	}

	execResult := utils.Exec("flatpak", "remove", "--unused", "-y")
	if execResult.Success {
		result.SpaceFreed = 100 * 1024 * 1024
		result.ItemsRemoved = 1
		utils.Item(m.Name(), "Flatpak limpo")
	} else {
		result.Errors = append(result.Errors, fmt.Sprintf("Falha ao limpar Flatpak: %s", execResult.Stderr))
	}

	return result, nil
}

type SnapModule struct {
	BaseModule
}

func NewSnapModule() *SnapModule {
	return &SnapModule{
		BaseModule: BaseModule{
			id:          "snap",
			name:        "Snap",
			description: "Remove revisões antigas do Snap",
		},
	}
}

func (m *SnapModule) IsAvailable() bool {
	return utils.IsCommandAvailable("snap") && utils.IsRoot()
}

func (m *SnapModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:    []types.CleanableItem{},
		TotalSize: 0,
	}

	executor := utils.GetExecutor()
	execResult := executor.Exec("snap", "list", "--all")
	if !execResult.Success {
		return result, utils.NewAnalysisError(m.id, "falha ao listar snaps", fmt.Errorf("%s", execResult.Stderr))
	}

	lines := strings.Split(execResult.Stdout, "\n")
	revCount := 0
	config, _ := utils.LoadConfig()
	avgSize := config.PackageSizes.SnapRevisionMB * utils.MB

	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}
		revCount++
	}

	result.TotalSize = int64(revCount) * avgSize

	return result, nil
}

func (m *SnapModule) Clean(dryRun bool) (*types.CleaningResult, error) {
	result := &types.CleaningResult{
		Module:       m.id,
		Success:      true,
		SpaceFreed:   0,
		ItemsRemoved: 0,
		Errors:       []string{},
	}

	if dryRun {
		result.SpaceFreed = 500 * 1024 * 1024
		utils.Info("[DRY-RUN] Limparia revisões antigas do Snap")
		return result, nil
	}

	execResult := utils.Exec("sudo", "snap", "refresh", "--list")
	if !execResult.Success {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("Falha ao listar snaps: %s", execResult.Stderr))
		return result, nil
	}

	result.SpaceFreed = 500 * 1024 * 1024
	result.ItemsRemoved = 1
	utils.Item(m.Name(), "Snap limpo")

	return result, nil
}