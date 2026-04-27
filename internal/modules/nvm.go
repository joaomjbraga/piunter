package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type NvmModule struct {
	BaseModule
}

func NewNvmModule() *NvmModule {
	return &NvmModule{
		BaseModule: BaseModule{
			id:          "nvm",
			name:        "NVM",
			description: "Limpa cache do NVM (Node Version Manager)",
		},
	}
}

func (m *NvmModule) IsAvailable() bool {
	home := utils.GetHomeDir()
	nvmDir := filepath.Join(home, ".nvm")
	return utils.FileExists(nvmDir)
}

func (m *NvmModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:    []types.CleanableItem{},
		TotalSize: 0,
	}

	home := utils.GetHomeDir()
	nvmDir := filepath.Join(home, ".nvm")
	cacheDir := filepath.Join(nvmDir, ".cache")
	sourceDir := filepath.Join(nvmDir, ".source")
	versionsDir := filepath.Join(nvmDir, "versions", "node")

	dirsToCheck := []string{cacheDir, sourceDir}
	for _, dir := range dirsToCheck {
		if !utils.FileExists(dir) {
			continue
		}
		size, err := utils.GetDirSize(dir)
		if err != nil {
			continue
		}
		result.Items = append(result.Items, types.CleanableItem{
			Path:        dir,
			Size:        size,
			Type:        "directory",
			Description: "NVM " + filepath.Base(dir) + " directory",
		})
		result.TotalSize += size
	}

	if utils.FileExists(versionsDir) {
		entries, err := os.ReadDir(versionsDir)
		if err == nil {
			inactiveVersions := 0
			var inactiveSize int64
			currentVersion := m.getCurrentNvmVersion()
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				versionDir := filepath.Join(versionsDir, entry.Name())
				if entry.Name() != currentVersion && entry.Name() != "node" {
					inactiveVersions++
					size, _ := utils.GetDirSize(versionDir)
					inactiveSize += size
					result.Items = append(result.Items, types.CleanableItem{
						Path:        versionDir,
						Size:        size,
						Type:        "directory",
						Description: "Versão inativa do Node: " + entry.Name(),
					})
					result.TotalSize += size
				}
			}
		}
	}

	return result, nil
}

func (m *NvmModule) getCurrentNvmVersion() string {
	home := utils.GetHomeDir()
	nvmDir := filepath.Join(home, ".nvm", "versions")
	if !utils.FileExists(nvmDir) {
		return ""
	}

	entries, err := os.ReadDir(nvmDir)
	if err != nil {
		return ""
	}

	var latestVersion string
	var latestTime int64

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "node" {
			continue
		}
		versionDir := filepath.Join(nvmDir, entry.Name())
		info, err := os.Stat(versionDir)
		if err != nil {
			continue
		}
		if info.ModTime().Unix() > latestTime {
			latestTime = info.ModTime().Unix()
			latestVersion = entry.Name()
		}
	}

	return latestVersion
}

func (m *NvmModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		Errors:      []string{},
	}

	if analysis.TotalSize == 0 {
		return result, nil
	}

	if dryRun {
		result.SpaceFreed = analysis.TotalSize
		result.ItemsRemoved = len(analysis.Items)
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %s do NVM", utils.FormatBytes(analysis.TotalSize)))
		return result, nil
	}

	executor := utils.GetExecutor()
	execResult := executor.Exec("nvm", "cache", "clear")
	if !execResult.Success {
		home := utils.GetHomeDir()
		cacheDir := filepath.Join(home, ".nvm", ".cache")
		if utils.FileExists(cacheDir) {
			os.RemoveAll(cacheDir)
			result.ItemsRemoved++
		}
		sourceDir := filepath.Join(home, ".nvm", "source")
		if utils.FileExists(sourceDir) {
			os.RemoveAll(sourceDir)
			result.ItemsRemoved++
		}
	}

	result.SpaceFreed = analysis.TotalSize
	result.ItemsRemoved = len(analysis.Items)
	utils.Item(m.Name(), "Cache NVM limpo")

	return result, nil
}

func (m *NvmModule) GetCacheSize() int64 {
	home := utils.GetHomeDir()
	nvmDir := filepath.Join(home, ".nvm")
	cacheDir := filepath.Join(nvmDir, ".cache")
	sourceDir := filepath.Join(nvmDir, "source")

	var totalSize int64

	if utils.FileExists(cacheDir) {
		size, _ := utils.GetDirSize(cacheDir)
		totalSize += size
	}

	if utils.FileExists(sourceDir) {
		size, _ := utils.GetDirSize(sourceDir)
		totalSize += size
	}

	return totalSize
}

func (m *NvmModule) GetInactiveVersions() ([]string, []int64) {
	home := utils.GetHomeDir()
	versionsDir := filepath.Join(home, ".nvm", "versions", "node")

	if !utils.FileExists(versionsDir) {
		return nil, nil
	}

	currentVersion := m.getCurrentNvmVersion()
	var versions []string
	var sizes []int64

	entries, _ := os.ReadDir(versionsDir)
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "node" {
			continue
		}
		if entry.Name() != currentVersion {
			versionDir := filepath.Join(versionsDir, entry.Name())
			size, _ := utils.GetDirSize(versionDir)
			versions = append(versions, entry.Name())
			sizes = append(sizes, size)
		}
	}

	return versions, sizes
}

func (m *NvmModule) CleanInactiveVersions(dryRun bool) (*types.CleaningResult, error) {
	home := utils.GetHomeDir()
	versionsDir := filepath.Join(home, ".nvm", "versions", "node")

	result := &types.CleaningResult{
		Module:       m.id,
		Success:      true,
		SpaceFreed:   0,
		ItemsRemoved: 0,
		Errors:      []string{},
	}

	if !utils.FileExists(versionsDir) {
		return result, nil
	}

	currentVersion := m.getCurrentNvmVersion()
	entries, _ := os.ReadDir(versionsDir)

	var totalSize int64
	var removed int

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "node" {
			continue
		}
		if entry.Name() != currentVersion {
			versionDir := filepath.Join(versionsDir, entry.Name())
			size, _ := utils.GetDirSize(versionDir)
			totalSize += size
			removed++

			if !dryRun {
				if err := os.RemoveAll(versionDir); err != nil {
					result.Errors = append(result.Errors, fmt.Sprintf("falha ao remover %s: %v", entry.Name(), err))
				}
			}
		}
	}

	result.SpaceFreed = totalSize
	result.ItemsRemoved = removed

	return result, nil
}

func (m *NvmModule) CleanCache(dryRun bool) (*types.CleaningResult, error) {
	home := utils.GetHomeDir()
	nvmDir := filepath.Join(home, ".nvm")
	cacheDir := filepath.Join(nvmDir, ".cache")
	sourceDir := filepath.Join(nvmDir, "source")

	result := &types.CleaningResult{
		Module:       m.id,
		Success:      true,
		SpaceFreed:   0,
		ItemsRemoved: 0,
		Errors:      []string{},
	}

	var totalSize int64

	if utils.FileExists(cacheDir) {
		size, _ := utils.GetDirSize(cacheDir)
		totalSize += size
		if !dryRun {
			os.RemoveAll(cacheDir)
		}
		result.ItemsRemoved++
	}

	if utils.FileExists(sourceDir) {
		size, _ := utils.GetDirSize(sourceDir)
		totalSize += size
		if !dryRun {
			os.RemoveAll(sourceDir)
		}
		result.ItemsRemoved++
	}

	result.SpaceFreed = totalSize

	return result, nil
}

func (m *NvmModule) IsPathSafe(path string) bool {
	home := utils.GetHomeDir()
	nvmDir := filepath.Join(home, ".nvm")
	return strings.HasPrefix(path, nvmDir)
}