package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type MiseModule struct {
	BaseModule
}

func NewMiseModule() *MiseModule {
	return &MiseModule{
		BaseModule: BaseModule{
			id:          "mise",
			name:        "Mise",
			description: "Limpa cache do Mise (gerenciador de runtimes e ferramentas)",
		},
	}
}

func (m *MiseModule) IsAvailable() bool {
	home := utils.GetHomeDir()
	dataDir := filepath.Join(home, ".local", "share", "mise")
	cacheDir := filepath.Join(home, ".cache", "mise")
	return utils.FileExists(dataDir) || utils.FileExists(cacheDir) || utils.IsCommandAvailable("mise")
}

func (m *MiseModule) getDataDir() string {
	home := utils.GetHomeDir()
	return filepath.Join(home, ".local", "share", "mise")
}

func (m *MiseModule) getCacheDir() string {
	home := utils.GetHomeDir()
	return filepath.Join(home, ".cache", "mise")
}

func (m *MiseModule) getDownloadsDir() string {
	return filepath.Join(m.getDataDir(), "downloads")
}

func (m *MiseModule) getInstallsDir() string {
	return filepath.Join(m.getDataDir(), "installs")
}

func (m *MiseModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:    []types.CleanableItem{},
		TotalSize: 0,
	}

	cacheDir := m.getCacheDir()
	if utils.FileExists(cacheDir) {
		size, err := utils.GetDirSize(cacheDir)
		if err == nil && size > 0 {
			result.Items = append(result.Items, types.CleanableItem{
				Path:        cacheDir,
				Size:        size,
				Type:        "directory",
				Description: "Cache do Mise",
			})
			result.TotalSize += size
		}
	}

	downloadsDir := m.getDownloadsDir()
	if utils.FileExists(downloadsDir) {
		size, err := utils.GetDirSize(downloadsDir)
		if err == nil && size > 0 {
			result.Items = append(result.Items, types.CleanableItem{
				Path:        downloadsDir,
				Size:        size,
				Type:        "directory",
				Description: "Downloads do Mise",
			})
			result.TotalSize += size
		}
	}

	installsDir := m.getInstallsDir()
	if utils.FileExists(installsDir) {
		inactiveSize, inactiveCount := m.getInactiveToolsSize(installsDir)
		if inactiveSize > 0 {
			result.Items = append(result.Items, types.CleanableItem{
				Path:        installsDir,
				Size:        inactiveSize,
				Type:        "directory",
				Description: fmt.Sprintf("Versões inativas de ferramentas (%d ferramentas)", inactiveCount),
			})
			result.TotalSize += inactiveSize
		}
	}

	return result, nil
}

func (m *MiseModule) getInactiveToolsSize(installsDir string) (int64, int) {
	var totalSize int64
	var count int

	entries, err := os.ReadDir(installsDir)
	if err != nil {
		return 0, 0
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		toolDir := filepath.Join(installsDir, entry.Name())
		dirSize, dirCount := m.getInactiveVersionsForTool(toolDir)
		totalSize += dirSize
		count += dirCount
	}

	return totalSize, count
}

func (m *MiseModule) getInactiveVersionsForTool(toolDir string) (int64, int) {
	var totalSize int64
	var count int

	currentLink := filepath.Join(toolDir, "current")
	if !utils.FileExists(currentLink) {
		return 0, 0
	}

	linkTarget, err := os.Readlink(currentLink)
	if err != nil {
		return 0, 0
	}
	activeVersion := filepath.Base(linkTarget)

	entries, err := os.ReadDir(toolDir)
	if err != nil {
		return 0, 0
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if entry.Name() == activeVersion {
			continue
		}
		versionDir := filepath.Join(toolDir, entry.Name())
		size, _ := utils.GetDirSize(versionDir)
		if size > 0 {
			totalSize += size
			count++
		}
	}

	return totalSize, count
}

func (m *MiseModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %s do Mise", utils.FormatBytes(analysis.TotalSize)))
		return result, nil
	}

	if utils.IsCommandAvailable("mise") {
		executor := utils.GetExecutor()
		execResult := executor.Exec("mise", "cache", "clear")
		if execResult.Success {
			utils.Item(m.Name(), "Cache limpo via mise cache clear")
		} else {
			m.cleanCacheDir(result)
		}
	} else {
		m.cleanCacheDir(result)
	}

	m.cleanDownloadsDir(result)
	m.cleanInactiveVersions(result)

	if result.ItemsRemoved > 0 {
		utils.Item(m.Name(), fmt.Sprintf("%d itens removidos", result.ItemsRemoved))
	}

	return result, nil
}

func (m *MiseModule) cleanCacheDir(result *types.CleaningResult) {
	cacheDir := m.getCacheDir()
	if !utils.FileExists(cacheDir) {
		return
	}
	size, _ := utils.GetDirSize(cacheDir)
	if size > 0 {
		result.SpaceFreed += size
		result.ItemsRemoved++
	}
	os.RemoveAll(cacheDir)
	utils.Info("Cache do Mise removido")
}

func (m *MiseModule) cleanDownloadsDir(result *types.CleaningResult) {
	downloadsDir := m.getDownloadsDir()
	if !utils.FileExists(downloadsDir) {
		return
	}
	size, _ := utils.GetDirSize(downloadsDir)
	if size > 0 {
		result.SpaceFreed += size
		result.ItemsRemoved++
	}
	os.RemoveAll(downloadsDir)
	utils.Info("Downloads do Mise removidos")
}

func (m *MiseModule) cleanInactiveVersions(result *types.CleaningResult) {
	installsDir := m.getInstallsDir()
	if !utils.FileExists(installsDir) {
		return
	}

	entries, err := os.ReadDir(installsDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		toolDir := filepath.Join(installsDir, entry.Name())
		m.cleanInactiveToolVersions(toolDir, entry.Name(), result)
	}
}

func (m *MiseModule) cleanInactiveToolVersions(toolDir, toolName string, result *types.CleaningResult) {
	currentLink := filepath.Join(toolDir, "current")
	if !utils.FileExists(currentLink) {
		return
	}

	linkTarget, err := os.Readlink(currentLink)
	if err != nil {
		return
	}
	activeVersion := filepath.Base(linkTarget)

	entries, err := os.ReadDir(toolDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if entry.Name() == activeVersion {
			continue
		}
		versionDir := filepath.Join(toolDir, entry.Name())
		size, _ := utils.GetDirSize(versionDir)
		if size > 0 {
			result.SpaceFreed += size
			result.ItemsRemoved++
		}
		if err := os.RemoveAll(versionDir); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("falha ao remover %s %s: %v", toolName, entry.Name(), err))
		}
	}
}

func (m *MiseModule) IsPathSafe(path string) bool {
	home := utils.GetHomeDir()
	dataDir := filepath.Join(home, ".local", "share", "mise")
	cacheDir := filepath.Join(home, ".cache", "mise")
	return strings.HasPrefix(path, dataDir) || strings.HasPrefix(path, cacheDir)
}
