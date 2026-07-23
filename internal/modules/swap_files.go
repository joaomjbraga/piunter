package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type SwapFilesModule struct {
	BaseModule
	extensions []string
}

func NewSwapFilesModule() *SwapFilesModule {
	return &SwapFilesModule{
		BaseModule: BaseModule{
			id:          "swap-files",
			name:        "Arquivos Swap",
			description: "Remove arquivos swap e temporários de editores",
		},
		extensions: []string{".swp", ".swo", ".tmp"},
	}
}

func (m *SwapFilesModule) IsAvailable() bool {
	return true
}

func (m *SwapFilesModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{Module: m.id, Items: []types.CleanableItem{}, TotalSize: 0}
	home := utils.GetHomeDir()

	err := filepath.Walk(home, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		name := info.Name()
		if !hasSwapExtension(name, m.extensions) {
			return nil
		}
		result.Items = append(result.Items, types.CleanableItem{
			Path:        path,
			Size:        info.Size(),
			Type:        "file",
			Description: fmt.Sprintf("Arquivo temporário: %s", name),
		})
		result.TotalSize += info.Size()
		return nil
	})
	if err != nil {
		return result, nil
	}
	return result, nil
}

func hasSwapExtension(name string, exts []string) bool {
	lower := strings.ToLower(name)
	for _, ext := range exts {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func (m *SwapFilesModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %d arquivos swap/temporários", len(analysis.Items)))
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
		utils.Item(m.Name(), fmt.Sprintf("%d arquivos removidos", result.ItemsRemoved))
	}
	if len(result.Errors) > 0 {
		result.Success = false
	}
	return result, nil
}
