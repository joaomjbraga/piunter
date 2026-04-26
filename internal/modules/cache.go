package modules

import (
	"os"
	"path/filepath"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type CacheModule struct {
	BaseModule
}

func NewCacheModule() *CacheModule {
	return &CacheModule{
		BaseModule: BaseModule{
			id:          "cache",
			name:        "Cache do Usuário",
			description: "Limpa cache geral do usuário (~/.cache)",
		},
	}
}

func (m *CacheModule) IsAvailable() bool {
	return utils.FileExists(utils.GetCacheDir())
}

func (m *CacheModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	cacheDir := utils.GetCacheDir()
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:    []types.CleanableItem{},
		TotalSize: 0,
	}

	if !utils.FileExists(cacheDir) {
		return result, nil
	}

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return result, utils.NewAnalysisError(m.id, "falha ao ler diretório de cache", err)
	}

	skipDirs := map[string]bool{
		"thumbnails": true,
		"thumbnail":  true,
		"icon-cache": true,
	}

	for _, entry := range entries {
		fullPath := filepath.Join(cacheDir, entry.Name())
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		if skipDirs[entry.Name()] {
			continue
		}

		var size int64
		if info.IsDir() {
			size = utils.GetDirSizeAsync(fullPath)
		} else {
			size = info.Size()
		}

		itemType := "file"
		if info.IsDir() {
			itemType = "directory"
		}

		result.Items = append(result.Items, types.CleanableItem{
			Path:        fullPath,
			Size:        size,
			Type:        itemType,
			Description: "Diretório de cache: " + entry.Name(),
		})
		result.TotalSize += size
	}

	return result, nil
}

func (m *CacheModule) Clean(dryRun bool) (*types.CleaningResult, error) {
	analysis, err := m.Analyze(0)
	if err != nil {
		return &types.CleaningResult{
			Module: m.id,
			Success: false,
			Errors: []string{err.Error()},
		}, err
	}

	result := &types.CleaningResult{
		Module:       m.id,
		Success:      true,
		SpaceFreed:   0,
		ItemsRemoved: 0,
		Errors:       []string{},
	}

	skipDirs := map[string]bool{
		"thumbnails": true,
		"thumbnail":  true,
		"icon-cache": true,
	}

	for _, item := range analysis.Items {
		if skipDirs[filepath.Base(item.Path)] {
			continue
		}

		if dryRun {
			result.SpaceFreed += item.Size
			result.ItemsRemoved++
		} else {
			err := utils.RemovePath(item.Path, item.Type == "directory")
			if err != nil {
				result.Errors = append(result.Errors, "Falha ao remover "+item.Path+": "+err.Error())
			} else {
				result.SpaceFreed += item.Size
				result.ItemsRemoved++
			}
		}
	}

	if result.SpaceFreed > 0 && !dryRun {
		utils.Item(m.Name(), "Cache limpo")
	}

	return result, nil
}