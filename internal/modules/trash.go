package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type TrashModule struct {
	BaseModule
}

func NewTrashModule() *TrashModule {
	return &TrashModule{
		BaseModule: BaseModule{
			id:          "trash",
			name:        "Lixeira",
			description: "Esvazia a lixeira do sistema",
		},
	}
}

func (m *TrashModule) IsAvailable() bool {
	home := utils.GetHomeDir()
	trashPath := filepath.Join(home, ".local", "share", "Trash")
	return utils.FileExists(trashPath)
}

func (m *TrashModule) getTrashPaths() []string {
	home := utils.GetHomeDir()
	return []string{
		filepath.Join(home, ".local", "share", "Trash", "files"),
		filepath.Join(home, ".local", "share", "Trash", "info"),
	}
}

func (m *TrashModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:     []types.CleanableItem{},
		TotalSize: 0,
	}

	for _, trashPath := range m.getTrashPaths() {
		if !utils.FileExists(trashPath) {
			continue
		}

		entries, err := os.ReadDir(trashPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			fullPath := filepath.Join(trashPath, entry.Name())
			info, err := os.Stat(fullPath)
			if err != nil {
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
				Description: "Lixeira: " + entry.Name(),
			})
			result.TotalSize += size
		}
	}

	return result, nil
}

func (m *TrashModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		result.ItemsRemoved = len(analysis.Items)
		utils.Info(fmt.Sprintf("[DRY-RUN] Esvaziaria %s da lixeira", utils.FormatBytes(analysis.TotalSize)))
		return result, nil
	}

	for _, item := range analysis.Items {
		err := utils.RemovePath(item.Path, item.Type == "directory")
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Falha ao remover %s: %s", item.Path, err.Error()))
		} else {
			result.SpaceFreed += item.Size
			result.ItemsRemoved++
		}
	}

	if result.ItemsRemoved > 0 {
		utils.Item(m.Name(), fmt.Sprintf("%d itens removidos", result.ItemsRemoved))
	}

	return result, nil
}