package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type ShellHistoryModule struct {
	BaseModule
	files []string
}

func NewShellHistoryModule() *ShellHistoryModule {
	return &ShellHistoryModule{
		BaseModule: BaseModule{
			id:          "shell-history",
			name:        "Histórico de Shell",
			description: "Limpa arquivos de histórico de shell do usuário",
		},
		files: []string{
			".bash_history",
			".zsh_history",
			".local/share/fish/fish_history",
		},
	}
}

func (m *ShellHistoryModule) IsAvailable() bool {
	return true
}

func (m *ShellHistoryModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{Module: m.id, Items: []types.CleanableItem{}, TotalSize: 0}
	home := utils.GetHomeDir()

	for _, relative := range m.files {
		path := filepath.Join(home, relative)
		if !utils.FileExists(path) {
			continue
		}
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		result.Items = append(result.Items, types.CleanableItem{
			Path:        path,
			Size:        info.Size(),
			Type:        "history-file",
			Description: fmt.Sprintf("Histórico de shell: %s", filepath.Base(path)),
		})
		result.TotalSize += info.Size()
	}

	return result, nil
}

func (m *ShellHistoryModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %d arquivos de histórico de shell", len(analysis.Items)))
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
		utils.Item(m.Name(), fmt.Sprintf("%d históricos removidos", result.ItemsRemoved))
	}
	if len(result.Errors) > 0 {
		result.Success = false
	}
	return result, nil
}

func init() {
	_ = strings.TrimSpace
}
