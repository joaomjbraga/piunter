package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type AppLogsModule struct {
	BaseModule
	paths []string
}

func NewAppLogsModule() *AppLogsModule {
	return &AppLogsModule{
		BaseModule: BaseModule{
			id:          "app-logs",
			name:        "Logs de Aplicativos",
			description: "Limpa logs de aplicativos e shells",
		},
		paths: []string{
			".cache/spotify/logs",
			".cache/discord",
			".config/Code/logs",
			".config/Code - Insiders/logs",
			".local/state",
		},
	}
}

func (m *AppLogsModule) IsAvailable() bool {
	return true
}

func (m *AppLogsModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{Module: m.id, Items: []types.CleanableItem{}, TotalSize: 0}
	home := utils.GetHomeDir()

	for _, relative := range m.paths {
		path := filepath.Join(home, relative)
		if !utils.FileExists(path) {
			continue
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			fullPath := filepath.Join(path, entry.Name())
			if entry.IsDir() {
				continue
			}
			if !strings.HasSuffix(strings.ToLower(entry.Name()), ".log") && !strings.HasSuffix(strings.ToLower(entry.Name()), ".log.gz") {
				continue
			}
			info, err := os.Stat(fullPath)
			if err != nil {
				continue
			}
			result.Items = append(result.Items, types.CleanableItem{
				Path:        fullPath,
				Size:        info.Size(),
				Type:        "file",
				Description: fmt.Sprintf("Log: %s", entry.Name()),
			})
			result.TotalSize += info.Size()
		}
	}

	return result, nil
}

func (m *AppLogsModule) Clean(dryRun bool) (*types.CleaningResult, error) {
	result := &types.CleaningResult{Module: m.id, Success: true, SpaceFreed: 0, ItemsRemoved: 0, Errors: []string{}}
	analysis, err := m.Analyze(0)
	if err != nil {
		return result, err
	}

	if len(analysis.Items) == 0 {
		return result, nil
	}

	if dryRun {
		result.SpaceFreed = analysis.TotalSize
		result.ItemsRemoved = len(analysis.Items)
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %d logs de aplicativos", len(analysis.Items)))
		return result, nil
	}

	for _, item := range analysis.Items {
		if err := os.Remove(item.Path); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Falha ao remover %s: %s", item.Path, err.Error()))
			continue
		}
		result.SpaceFreed += item.Size
		result.ItemsRemoved++
	}

	if result.ItemsRemoved > 0 {
		utils.Item(m.Name(), fmt.Sprintf("%d logs removidos", result.ItemsRemoved))
	}
	if len(result.Errors) > 0 {
		result.Success = false
	}
	return result, nil
}
