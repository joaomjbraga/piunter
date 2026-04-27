package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type SdkmanModule struct {
	BaseModule
}

func NewSdkmanModule() *SdkmanModule {
	return &SdkmanModule{
		BaseModule: BaseModule{
			id:          "sdkman",
			name:        "SDKMAN",
			description: "Limpa cache do SDKMAN (Software Development Kit Manager)",
		},
	}
}

func (m *SdkmanModule) IsAvailable() bool {
	home := utils.GetHomeDir()
	sdkmanDir := filepath.Join(home, ".sdkman")
	return utils.FileExists(sdkmanDir)
}

func (m *SdkmanModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:    []types.CleanableItem{},
		TotalSize: 0,
	}

	home := utils.GetHomeDir()
	sdkmanDir := filepath.Join(home, ".sdkman")

	dirsToCheck := map[string]string{
		"archives":     filepath.Join(sdkmanDir, "archives"),
		"tmp":         filepath.Join(sdkmanDir, "tmp"),
		"downloads":    filepath.Join(sdkmanDir, "var", "downloads"),
		"temp":        filepath.Join(sdkmanDir, "tmp"),
	}

	for name, dir := range dirsToCheck {
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
			Description: "SDKMAN " + name + " directory",
		})
		result.TotalSize += size
	}

	candidatesDir := filepath.Join(sdkmanDir, "candidates")
	if utils.FileExists(candidatesDir) {
		inactiveSize, inactiveCount := m.getInactiveCandidatesSize(candidatesDir)
		if inactiveSize > 0 {
			result.Items = append(result.Items, types.CleanableItem{
				Path:        candidatesDir,
				Size:        inactiveSize,
				Type:        "directory",
				Description: fmt.Sprintf("SDKMAN candidatos inativos (%d versões)", inactiveCount),
			})
			result.TotalSize += inactiveSize
		}
	}

	return result, nil
}

func (m *SdkmanModule) getInactiveCandidatesSize(candidatesDir string) (int64, int) {
	var totalSize int64
	var count int

	entries, _ := os.ReadDir(candidatesDir)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidateDir := filepath.Join(candidatesDir, entry.Name())
		currentLink := filepath.Join(candidateDir, "current")
		if !utils.FileExists(currentLink) {
			continue
		}

		linkTarget, err := os.Readlink(currentLink)
		if err != nil {
			continue
		}

		subDirs, _ := os.ReadDir(candidateDir)
		for _, sub := range subDirs {
			if sub.Name() == "current" {
				continue
			}
			subDir := filepath.Join(candidateDir, sub.Name())
			size, _ := utils.GetDirSize(subDir)
			if !strings.Contains(linkTarget, sub.Name()) {
				count++
				totalSize += size
			}
		}
	}

	return totalSize, count
}

func (m *SdkmanModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %s do SDKMAN", utils.FormatBytes(analysis.TotalSize)))
		return result, nil
	}

	executor := utils.GetExecutor()
	execResult := executor.Exec("sdk", "flush", "temp")
	if !execResult.Success {
		home := utils.GetHomeDir()
		sdkmanDir := filepath.Join(home, ".sdkman")

		dirsToClean := []string{
			filepath.Join(sdkmanDir, "archives"),
			filepath.Join(sdkmanDir, "tmp"),
			filepath.Join(sdkmanDir, "var", "downloads"),
		}

		for _, dir := range dirsToClean {
			if utils.FileExists(dir) {
				os.RemoveAll(dir)
				result.ItemsRemoved++
			}
		}
	}

	result.SpaceFreed = analysis.TotalSize
	result.ItemsRemoved = len(analysis.Items)
	utils.Item(m.Name(), "Cache SDKMAN limpo")

	return result, nil
}

func (m *SdkmanModule) CleanCache(dryRun bool) (*types.CleaningResult, error) {
	home := utils.GetHomeDir()
	sdkmanDir := filepath.Join(home, ".sdkman")

	result := &types.CleaningResult{
		Module:       m.id,
		Success:      true,
		SpaceFreed:   0,
		ItemsRemoved: 0,
		Errors:      []string{},
	}

	cacheDirs := map[string]string{
		"archives": filepath.Join(sdkmanDir, "archives"),
		"tmp":     filepath.Join(sdkmanDir, "tmp"),
	}

	var totalSize int64
	for name, dir := range cacheDirs {
		if utils.FileExists(dir) {
			size, _ := utils.GetDirSize(dir)
			totalSize += size
			result.ItemsRemoved++
			if !dryRun {
				os.RemoveAll(dir)
				utils.Info(fmt.Sprintf("Removido: SDKMAN %s", name))
			}
		}
	}

	result.SpaceFreed = totalSize

	return result, nil
}

func (m *SdkmanModule) CleanOldVersions(dryRun bool) (*types.CleaningResult, error) {
	home := utils.GetHomeDir()
	sdkmanDir := filepath.Join(home, ".sdkman")
	candidatesDir := filepath.Join(sdkmanDir, "candidates")

	result := &types.CleaningResult{
		Module:       m.id,
		Success:      true,
		SpaceFreed:   0,
		ItemsRemoved: 0,
		Errors:      []string{},
	}

	if !utils.FileExists(candidatesDir) {
		return result, nil
	}

	entries, _ := os.ReadDir(candidatesDir)
	var totalSize int64

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidateDir := filepath.Join(candidatesDir, entry.Name())
		currentLink := filepath.Join(candidateDir, "current")
		if !utils.FileExists(currentLink) {
			continue
		}

		linkTarget, _ := os.Readlink(currentLink)

		subDirs, _ := os.ReadDir(candidateDir)
		for _, sub := range subDirs {
			if sub.Name() == "current" {
				continue
			}
			subDir := filepath.Join(candidateDir, sub.Name())
			if !strings.Contains(linkTarget, sub.Name()) {
				size, _ := utils.GetDirSize(subDir)
				totalSize += size
				result.ItemsRemoved++
				if !dryRun {
					os.RemoveAll(subDir)
					utils.Info(fmt.Sprintf("Removido: SDKMAN %s %s", entry.Name(), sub.Name()))
				}
			}
		}
	}

	result.SpaceFreed = totalSize

	return result, nil
}

func (m *SdkmanModule) GetSDKMANDir() string {
	home := utils.GetHomeDir()
	return filepath.Join(home, ".sdkman")
}

func (m *SdkmanModule) GetArchivesSize() int64 {
	home := utils.GetHomeDir()
	archivesDir := filepath.Join(home, ".sdkman", "archives")
	if !utils.FileExists(archivesDir) {
		return 0
	}
	size, _ := utils.GetDirSize(archivesDir)
	return size
}

func (m *SdkmanModule) GetTmpSize() int64 {
	home := utils.GetHomeDir()
	tmpDir := filepath.Join(home, ".sdkman", "tmp")
	if !utils.FileExists(tmpDir) {
		return 0
	}
	size, _ := utils.GetDirSize(tmpDir)
	return size
}

func (m *SdkmanModule) IsPathSafe(path string) bool {
	home := utils.GetHomeDir()
	sdkmanDir := filepath.Join(home, ".sdkman")
	return strings.HasPrefix(path, sdkmanDir)
}