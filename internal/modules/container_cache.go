package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type ContainerCacheModule struct {
	BaseModule
	paths []string
}

func NewContainerCacheModule() *ContainerCacheModule {
	return &ContainerCacheModule{
		BaseModule: BaseModule{
			id:          "container-cache",
			name:        "Cache de Containers",
			description: "Limpa caches e artefatos de containers e VMs",
		},
		paths: []string{
			".cache/podman",
			".local/share/containers",
			".local/share/libvirt",
			".cache/virt-manager",
		},
	}
}

func (m *ContainerCacheModule) IsAvailable() bool {
	return true
}

func (m *ContainerCacheModule) Analyze(threshold int) (*types.AnalysisResult, error) {
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
			Description: fmt.Sprintf("Cache de container/VM: %s", filepath.Base(path)),
		})
		result.TotalSize += size
	}

	return result, nil
}

func (m *ContainerCacheModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %d caches de containers/VMs", len(analysis.Items)))
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
