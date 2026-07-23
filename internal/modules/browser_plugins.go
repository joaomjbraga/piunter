package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type BrowserPluginsModule struct {
	BaseModule
	paths []string
}

func NewBrowserPluginsModule() *BrowserPluginsModule {
	return &BrowserPluginsModule{
		BaseModule: BaseModule{
			id:          "browser-plugins",
			name:        "Plugins de Navegador",
			description: "Limpa resíduos de extensões e plugins de navegadores",
		},
		paths: []string{
			".config/google-chrome/Default/Extensions",
			".config/google-chrome/Default/Service Worker/ScriptStorage",
			".mozilla/firefox/Crash Reports",
			".mozilla/firefox/.parentlock",
		},
	}
}

func (m *BrowserPluginsModule) IsAvailable() bool {
	return true
}

func (m *BrowserPluginsModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{Module: m.id, Items: []types.CleanableItem{}, TotalSize: 0}
	home := utils.GetHomeDir()

	for _, relative := range m.paths {
		path := filepath.Join(home, relative)
		if !utils.FileExists(path) {
			continue
		}
		size, err := utils.GetDirSize(path)
		if err != nil {
			continue
		}
		result.Items = append(result.Items, types.CleanableItem{
			Path:        path,
			Size:        size,
			Type:        "directory",
			Description: fmt.Sprintf("Resíduo de plugin: %s", filepath.Base(path)),
		})
		result.TotalSize += size
	}

	return result, nil
}

func (m *BrowserPluginsModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %d resíduos de plugins", len(analysis.Items)))
		return result, nil
	}

	for _, item := range analysis.Items {
		if err := os.RemoveAll(item.Path); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Falha ao remover %s: %s", item.Path, err.Error()))
			continue
		}
		result.SpaceFreed += item.Size
		result.ItemsRemoved++
	}

	if result.ItemsRemoved > 0 {
		utils.Item(m.Name(), fmt.Sprintf("%d resíduos removidos", result.ItemsRemoved))
	}
	if len(result.Errors) > 0 {
		result.Success = false
	}
	return result, nil
}
