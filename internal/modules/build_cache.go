package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type BuildCacheModule struct {
	BaseModule
	paths []string
}

func NewBuildCacheModule() *BuildCacheModule {
	return &BuildCacheModule{
		BaseModule: BaseModule{
			id:          "build-cache",
			name:        "Cache de Build",
			description: "Limpa caches de ferramentas de build e IA",
		},
		paths: []string{
			".gradle",
			".m2",
			".cache/torch",
			".cache/huggingface",
			".cache/keras",
		},
	}
}

func (m *BuildCacheModule) IsAvailable() bool {
	return true
}

func (m *BuildCacheModule) Analyze(threshold int) (*types.AnalysisResult, error) {
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
			Description: fmt.Sprintf("Cache de build/IA: %s", filepath.Base(path)),
		})
		result.TotalSize += size
	}

	return result, nil
}

func (m *BuildCacheModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %d caches de build/IA", len(analysis.Items)))
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
		utils.Item(m.Name(), fmt.Sprintf("%d caches removidos", result.ItemsRemoved))
	}
	if len(result.Errors) > 0 {
		result.Success = false
	}
	return result, nil
}
