package modules

import (
	"fmt"
	"strings"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type PackageCacheModule struct {
	BaseModule
	packageManager string
}

func NewPackageCacheModule() *PackageCacheModule {
	return &PackageCacheModule{
		BaseModule: BaseModule{
			id:          "package-cache",
			name:        "Cache de Pacotes",
			description: "Limpa o cache de downloads dos gerenciadores de pacotes",
		},
		packageManager: detectPackageManager(),
	}
}

func (m *PackageCacheModule) IsAvailable() bool {
	return m.packageManager != ""
}

func (m *PackageCacheModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{Module: m.id, Items: []types.CleanableItem{}, TotalSize: 0}

	executor := utils.GetExecutor()
	var cmd string
	var args []string

	switch m.packageManager {
	case "apt":
		cmd = "apt-get"
		args = []string{"clean"}
	case "pacman":
		cmd = "pacman"
		args = []string{"-Sc"}
	case "dnf":
		cmd = "dnf"
		args = []string{"clean", "all"}
	default:
		return result, nil
	}

	execResult := executor.Exec(cmd, args...)
	if !execResult.Success {
		return result, utils.NewAnalysisError(m.id, fmt.Sprintf("falha ao analisar cache de %s", m.packageManager), fmt.Errorf("%s", execResult.Stderr))
	}

	result.Items = append(result.Items, types.CleanableItem{
		Path:        fmt.Sprintf("cache-%s", m.packageManager),
		Size:        0,
		Type:        "package-cache",
		Description: fmt.Sprintf("Cache de downloads do %s", m.packageManager),
	})
	result.TotalSize = 0
	return result, nil
}

func (m *PackageCacheModule) Clean(dryRun bool) (*types.CleaningResult, error) {
	result := &types.CleaningResult{Module: m.id, Success: true, SpaceFreed: 0, ItemsRemoved: 0, Errors: []string{}}
	if !m.IsAvailable() {
		return result, nil
	}

	executor := utils.GetExecutor()
	var cmd string
	var args []string

	switch m.packageManager {
	case "apt":
		cmd = "sudo"
		args = []string{"apt-get", "clean"}
	case "pacman":
		cmd = "sudo"
		args = []string{"pacman", "-Sc"}
	case "dnf":
		cmd = "sudo"
		args = []string{"dnf", "clean", "all"}
	default:
		return result, nil
	}

	if dryRun {
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia o cache do %s", m.packageManager))
		result.SpaceFreed = 0
		result.ItemsRemoved = 1
		return result, nil
	}

	execResult := executor.Exec(cmd, args...)
	if execResult.Success {
		utils.Item(m.Name(), "Cache de pacotes limpo")
		result.ItemsRemoved = 1
		return result, nil
	}

	result.Success = false
	result.Errors = append(result.Errors, fmt.Sprintf("Falha ao limpar cache de pacotes: %s", execResult.Stderr))
	return result, nil
}

func init() {
	// keep module registration in index.go; no-op here
	_ = strings.TrimSpace
}
