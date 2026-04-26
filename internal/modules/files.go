package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type LargeFilesModule struct {
	BaseModule
	threshold int
}

func NewLargeFilesModule() *LargeFilesModule {
	return &LargeFilesModule{
		BaseModule: BaseModule{
			id:          "large-files",
			name:        "Arquivos Grandes",
			description: "Encontra arquivos grandes no sistema",
		},
		threshold: 100,
	}
}

func (m *LargeFilesModule) IsAvailable() bool {
	return true
}

func (m *LargeFilesModule) SetThreshold(t int) {
	m.threshold = t
}

func (m *LargeFilesModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	if threshold > 0 {
		m.threshold = threshold
	}

	result := &types.AnalysisResult{
		Module:    m.id,
		Items:     []types.CleanableItem{},
		TotalSize: 0,
	}

	home := utils.GetHomeDir()
	thresholdBytes := int64(m.threshold) * 1024 * 1024

	err := filepath.Walk(home, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if info.Size() >= thresholdBytes {
			result.Items = append(result.Items, types.CleanableItem{
				Path:        path,
				Size:        info.Size(),
				Type:        "file",
				Description: fmt.Sprintf("Arquivo grande: %s", filepath.Base(path)),
			})
			result.TotalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		return result, nil
	}

	return result, nil
}

func (m *LargeFilesModule) Clean(dryRun bool) (*types.CleaningResult, error) {
	result := &types.CleaningResult{
		Module:       m.id,
		Success:      true,
		SpaceFreed:   0,
		ItemsRemoved: 0,
		Errors:       []string{},
	}

	utils.Info("Use --analyze para encontrar arquivos grandes primeiro")
	return result, nil
}

type ThumbsModule struct {
	BaseModule
}

func NewThumbsModule() *ThumbsModule {
	return &ThumbsModule{
		BaseModule: BaseModule{
			id:          "thumbs",
			name:        "Miniaturas",
			description: "Limpa miniaturas em cache do usuário",
		},
	}
}

func (m *ThumbsModule) IsAvailable() bool {
	home := utils.GetHomeDir()
	thumbsPath := filepath.Join(home, ".cache", "thumbnails")
	return utils.FileExists(thumbsPath)
}

func (m *ThumbsModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	home := utils.GetHomeDir()
	thumbsPath := filepath.Join(home, ".cache", "thumbnails")

	result := &types.AnalysisResult{
		Module:    m.id,
		Items:     []types.CleanableItem{},
		TotalSize: 0,
	}

	if !utils.FileExists(thumbsPath) {
		return result, nil
	}

	entries, err := os.ReadDir(thumbsPath)
	if err != nil {
		return result, nil
	}

	for _, entry := range entries {
		fullPath := filepath.Join(thumbsPath, entry.Name())
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		size := info.Size()
		result.Items = append(result.Items, types.CleanableItem{
			Path:        fullPath,
			Size:        size,
			Type:        "file",
			Description: "Miniatura: " + entry.Name(),
		})
		result.TotalSize += size
	}

	return result, nil
}

func (m *ThumbsModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %d miniaturas", len(analysis.Items)))
		return result, nil
	}

	for _, item := range analysis.Items {
		err := os.Remove(item.Path)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Falha ao remover %s: %s", item.Path, err.Error()))
		} else {
			result.SpaceFreed += item.Size
			result.ItemsRemoved++
		}
	}

	if result.ItemsRemoved > 0 {
		utils.Item(m.Name(), fmt.Sprintf("%d miniaturas removidas", result.ItemsRemoved))
	}

	return result, nil
}

type RecentFilesModule struct {
	BaseModule
	days int
}

func NewRecentFilesModule() *RecentFilesModule {
	return &RecentFilesModule{
		BaseModule: BaseModule{
			id:          "recent",
			name:        "Arquivos Recentes",
			description: "Lista arquivos modificados recentemente",
		},
		days: 7,
	}
}

func (m *RecentFilesModule) IsAvailable() bool {
	return true
}

func (m *RecentFilesModule) SetDays(d int) {
	m.days = d
}

func (m *RecentFilesModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:     []types.CleanableItem{},
		TotalSize: 0,
	}

	home := utils.GetHomeDir()
	cutoff := time.Now().AddDate(0, 0, -m.days)

	filepath.Walk(home, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if info.ModTime().After(cutoff) && info.Size() > 1024*1024 {
			result.Items = append(result.Items, types.CleanableItem{
				Path:        path,
				Size:        info.Size(),
				Type:        "file",
				Description: fmt.Sprintf("Arquivo recente: %s", filepath.Base(path)),
			})
			result.TotalSize += info.Size()
		}
		return nil
	})

	return result, nil
}

func (m *RecentFilesModule) Clean(dryRun bool) (*types.CleaningResult, error) {
	result := &types.CleaningResult{
		Module:       m.id,
		Success:      true,
		SpaceFreed:   0,
		ItemsRemoved: 0,
		Errors:       []string{},
	}

	utils.Info("Use --analyze para ver arquivos recentes primeiro")
	return result, nil
}

type AppImageModule struct {
	BaseModule
}

func NewAppImageModule() *AppImageModule {
	return &AppImageModule{
		BaseModule: BaseModule{
			id:          "appimage",
			name:        "AppImage",
			description: "Remove arquivos AppImage antigos",
		},
	}
}

func (m *AppImageModule) IsAvailable() bool {
	home := utils.GetHomeDir()
	downloadsPath := filepath.Join(home, "Downloads")
	return utils.FileExists(downloadsPath)
}

func (m *AppImageModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	home := utils.GetHomeDir()
	downloadsPath := filepath.Join(home, "Downloads")

	result := &types.AnalysisResult{
		Module:    m.id,
		Items:     []types.CleanableItem{},
		TotalSize: 0,
	}

	if !utils.FileExists(downloadsPath) {
		return result, nil
	}

	entries, err := os.ReadDir(downloadsPath)
	if err != nil {
		return result, nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".AppImage") && !strings.HasSuffix(name, ".appimage") {
			continue
		}
		fullPath := filepath.Join(downloadsPath, name)
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		size := info.Size()
		result.Items = append(result.Items, types.CleanableItem{
			Path:        fullPath,
			Size:        size,
			Type:        "file",
			Description: "AppImage: " + name,
		})
		result.TotalSize += size
	}

	return result, nil
}

func (m *AppImageModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %d AppImages", len(analysis.Items)))
		return result, nil
	}

	for _, item := range analysis.Items {
		err := os.Remove(item.Path)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Falha ao remover %s: %s", item.Path, err.Error()))
		} else {
			result.SpaceFreed += item.Size
			result.ItemsRemoved++
		}
	}

	if result.ItemsRemoved > 0 {
		utils.Item(m.Name(), fmt.Sprintf("%d AppImages removidas", result.ItemsRemoved))
	}

	return result, nil
}