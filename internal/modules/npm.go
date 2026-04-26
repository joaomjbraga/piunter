package modules

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type PackageCacheModule struct {
	BaseModule
	cachePaths    []func() []string
	cleanCommands []string
}

func (m *PackageCacheModule) IsAvailable() bool {
	for _, getPaths := range m.cachePaths {
		for _, path := range getPaths() {
			if utils.FileExists(path) {
				return true
			}
		}
	}
	return false
}

func (m *PackageCacheModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.ID(),
		Items:     []types.CleanableItem{},
		TotalSize: 0,
	}

	for _, getPaths := range m.cachePaths {
		for _, path := range getPaths() {
			if !utils.FileExists(path) {
				continue
			}
			size := utils.GetDirSizeAsync(path)
			result.Items = append(result.Items, types.CleanableItem{
				Path:        path,
				Size:        size,
				Type:        "directory",
				Description: fmt.Sprintf("Cache do %s (%s)", m.Name(), path),
			})
			result.TotalSize += size
		}
	}

	return result, nil
}

func (m *PackageCacheModule) Clean(dryRun bool) (*types.CleaningResult, error) {
	analysis, err := m.Analyze(0)
	if err != nil {
		return &types.CleaningResult{
			Module:  m.ID(),
			Success: false,
			Errors:  []string{err.Error()},
		}, err
	}

	result := &types.CleaningResult{
		Module:       m.ID(),
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %s do cache %s", utils.FormatBytes(analysis.TotalSize), m.Name()))
		return result, nil
	}

	for _, cmd := range m.cleanCommands {
		parts := strings.Fields(cmd)
		execResult := utils.Exec(parts[0], parts[1:]...)
		if execResult.Success {
			result.SpaceFreed = analysis.TotalSize
			result.ItemsRemoved = 1
			utils.Item(m.Name(), "Cache limpo")
		} else {
			result.Success = false
			result.Errors = append(result.Errors, fmt.Sprintf("Falha ao limpar cache %s: %s", m.Name(), execResult.Stderr))
		}
	}

	return result, nil
}

type NpmModule struct {
	PackageCacheModule
}

func NewNpmModule() *NpmModule {
	home := utils.GetHomeDir()
	return &NpmModule{
		PackageCacheModule: PackageCacheModule{
			BaseModule: BaseModule{
				id:          "npm",
				name:        "NPM",
				description: "Limpa cache do npm",
			},
			cachePaths: []func() []string{
				func() []string { return []string{filepath.Join(home, ".npm")} },
			},
			cleanCommands: []string{"npm cache clean --force"},
		},
	}
}

type YarnModule struct {
	PackageCacheModule
}

func NewYarnModule() *YarnModule {
	home := utils.GetHomeDir()
	return &YarnModule{
		PackageCacheModule: PackageCacheModule{
			BaseModule: BaseModule{
				id:          "yarn",
				name:        "Yarn",
				description: "Limpa cache do Yarn",
			},
			cachePaths: []func() []string{
				func() []string { return []string{filepath.Join(home, ".cache", "yarn"), filepath.Join(home, ".yarn", "cache")} },
			},
			cleanCommands: []string{"yarn cache clean"},
		},
	}
}

type PnpmModule struct {
	PackageCacheModule
}

func NewPnpmModule() *PnpmModule {
	home := utils.GetHomeDir()
	return &PnpmModule{
		PackageCacheModule: PackageCacheModule{
			BaseModule: BaseModule{
				id:          "pnpm",
				name:        "PNPM",
				description: "Limpa cache do PNPM",
			},
			cachePaths: []func() []string{
				func() []string {
					paths := []string{filepath.Join(home, ".pnpm-store")}
					localPath := filepath.Join(home, ".local", "share", "pnpm", "cache")
					if utils.FileExists(localPath) {
						paths = append(paths, localPath)
					}
					return paths
				},
			},
			cleanCommands: []string{"pnpm store prune"},
		},
	}
}