package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type DownloadsOldModule struct {
	BaseModule
	paths []string
}

func NewDownloadsOldModule() *DownloadsOldModule {
	return &DownloadsOldModule{
		BaseModule: BaseModule{
			id:          "downloads-old",
			name:        "Downloads Antigos",
			description: "Limpa arquivos antigos em Downloads e pastas semelhantes",
		},
		paths: []string{"Downloads", ".mozilla/firefox", ".config/google-chrome/Default/Downloads"},
	}
}

func (m *DownloadsOldModule) IsAvailable() bool {
	return true
}

func (m *DownloadsOldModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{Module: m.id, Items: []types.CleanableItem{}, TotalSize: 0}
	home := utils.GetHomeDir()
	cutoff := time.Now().AddDate(0, 0, -30)

	for _, relative := range m.paths {
		path := filepath.Join(home, relative)
		if !utils.FileExists(path) {
			continue
		}
		filepath.Walk(path, func(current string, info os.FileInfo, err error) error {
			if err != nil || info == nil || info.IsDir() {
				return nil
			}
			if info.ModTime().Before(cutoff) {
				result.Items = append(result.Items, types.CleanableItem{
					Path:        current,
					Size:        info.Size(),
					Type:        "download-file",
					Description: fmt.Sprintf("Arquivo antigo: %s", filepath.Base(current)),
				})
				result.TotalSize += info.Size()
			}
			return nil
		})
	}

	return result, nil
}

func (m *DownloadsOldModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %d arquivos antigos em Downloads", len(analysis.Items)))
		return result, nil
	}

	cutoff := time.Now().AddDate(0, 0, -30)
	for _, item := range analysis.Items {
		info, err := os.Stat(item.Path)
		if err != nil {
			continue
		}
		if info.ModTime().After(cutoff) {
			continue
		}
		if err := os.Remove(item.Path); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Falha ao remover %s: %s", item.Path, err.Error()))
			continue
		}
		result.SpaceFreed += item.Size
		result.ItemsRemoved++
	}

	if result.ItemsRemoved > 0 {
		utils.Item(m.Name(), fmt.Sprintf("%d arquivos antigos removidos", result.ItemsRemoved))
	}
	if len(result.Errors) > 0 {
		result.Success = false
	}
	return result, nil
}
