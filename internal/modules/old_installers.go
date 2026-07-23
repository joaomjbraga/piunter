package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type OldInstallersModule struct {
	BaseModule
	extensions []string
}

func NewOldInstallersModule() *OldInstallersModule {
	return &OldInstallersModule{
		BaseModule: BaseModule{
			id:          "old-installers",
			name:        "Instaladores Antigos",
			description: "Remove instaladores antigos em Downloads e pastas semelhantes",
		},
		extensions: []string{".deb", ".rpm", ".tar.gz", ".tgz", ".tar.xz", ".zip"},
	}
}

func (m *OldInstallersModule) IsAvailable() bool {
	return true
}

func (m *OldInstallersModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{Module: m.id, Items: []types.CleanableItem{}, TotalSize: 0}
	home := utils.GetHomeDir()
	paths := []string{"Downloads", "Desktop", "Public"}

	for _, relative := range paths {
		path := filepath.Join(home, relative)
		if !utils.FileExists(path) {
			continue
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if !hasInstallerExtension(name, m.extensions) {
				continue
			}
			fullPath := filepath.Join(path, name)
			info, err := os.Stat(fullPath)
			if err != nil {
				continue
			}
			result.Items = append(result.Items, types.CleanableItem{
				Path:        fullPath,
				Size:        info.Size(),
				Type:        "file",
				Description: fmt.Sprintf("Instalador antigo: %s", name),
			})
			result.TotalSize += info.Size()
		}
	}

	return result, nil
}

func hasInstallerExtension(name string, exts []string) bool {
	lower := strings.ToLower(name)
	for _, ext := range exts {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func (m *OldInstallersModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %d instaladores antigos", len(analysis.Items)))
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
		utils.Item(m.Name(), fmt.Sprintf("%d instaladores removidos", result.ItemsRemoved))
	}
	if len(result.Errors) > 0 {
		result.Success = false
	}
	return result, nil
}
