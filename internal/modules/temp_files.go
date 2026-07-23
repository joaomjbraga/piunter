package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type TempFilesModule struct {
	BaseModule
	paths []string
}

func NewTempFilesModule() *TempFilesModule {
	return &TempFilesModule{
		BaseModule: BaseModule{
			id:          "temp-files",
			name:        "Arquivos Temporários",
			description: "Limpa arquivos temporários antigos em /tmp e /var/tmp",
		},
		paths: []string{"/tmp", "/var/tmp"},
	}
}

func (m *TempFilesModule) IsAvailable() bool {
	return true
}

func (m *TempFilesModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{Module: m.id, Items: []types.CleanableItem{}, TotalSize: 0}

	cutoff := time.Now().Add(-24 * time.Hour)
	for _, path := range m.paths {
		if !utils.FileExists(path) {
			continue
		}
		filepath.Walk(path, func(current string, info os.FileInfo, err error) error {
			if err != nil || info == nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if info.ModTime().Before(cutoff) {
				result.Items = append(result.Items, types.CleanableItem{
					Path:        current,
					Size:        info.Size(),
					Type:        "temp-file",
					Description: fmt.Sprintf("Arquivo temporário antigo: %s", filepath.Base(current)),
				})
				result.TotalSize += info.Size()
			}
			return nil
		})
	}

	return result, nil
}

func (m *TempFilesModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %d arquivos temporários antigos", len(analysis.Items)))
		return result, nil
	}

	cutoff := time.Now().Add(-24 * time.Hour)
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
		utils.Item(m.Name(), fmt.Sprintf("%d arquivos temporários removidos", result.ItemsRemoved))
	}
	return result, nil
}
