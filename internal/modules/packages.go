package modules

import (
	"fmt"
	"strings"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type PackagesModule struct {
	BaseModule
	packageManager string
}

func NewPackagesModule() *PackagesModule {
	return &PackagesModule{
		BaseModule: BaseModule{
			id:          "packages",
			name:        "Pacotes Órfãos",
			description: "Remove pacotes órfãos do sistema",
		},
		packageManager: detectPackageManager(),
	}
}

func detectPackageManager() string {
	if utils.IsCommandAvailable("apt") {
		return "apt"
	}
	if utils.IsCommandAvailable("pacman") {
		return "pacman"
	}
	if utils.IsCommandAvailable("dnf") {
		return "dnf"
	}
	return ""
}

func (m *PackagesModule) IsAvailable() bool {
	return m.packageManager != "" && utils.IsCommandAvailable(m.packageManager)
}

func (m *PackagesModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:     []types.CleanableItem{},
		TotalSize: 0,
	}

	var cmd string
	var args []string

	switch m.packageManager {
	case "apt":
		cmd = "apt"
		args = []string{"list", "--manual-installed", "--installed"}
	case "pacman":
		cmd = "pacman"
		args = []string{"-Qttd"}
	case "dnf":
		cmd = "dnf"
		args = []string{"autoremove", "--assumeno", "--verbose"}
	}

	execResult := utils.Exec(cmd, args...)
	if !execResult.Success && execResult.Code != 1 {
		return result, nil
	}

	lines := strings.Split(execResult.Stdout, "\n")
	orphanCount := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch m.packageManager {
		case "apt":
			if strings.Contains(line, "[") {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				orphanCount++
				result.Items = append(result.Items, types.CleanableItem{
					Path:        parts[0],
					Size:        0,
					Type:        "package",
					Description: "Pacote manual: " + parts[0],
				})
			}
		case "pacman", "dnf":
			orphanCount++
			result.Items = append(result.Items, types.CleanableItem{
				Path:        line,
				Size:        0,
				Type:        "package",
				Description: "Pacote órfão: " + line,
			})
		}
	}

	result.TotalSize = int64(orphanCount) * 10 * 1024 * 1024

	return result, nil
}

func (m *PackagesModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %d pacotes órfãos", len(analysis.Items)))
		return result, nil
	}

	var execResult types.CommandResult

	switch m.packageManager {
	case "apt":
		execResult = utils.Exec("sudo", "apt", "autoremove", "-y")
	case "pacman":
		execResult = utils.Exec("sudo", "pacman", "-R", "-s", "--nosave")
	case "dnf":
		execResult = utils.Exec("sudo", "dnf", "autoremove", "-y")
	}

	if execResult.Success {
		result.SpaceFreed = analysis.TotalSize
		result.ItemsRemoved = len(analysis.Items)
		utils.Item(m.Name(), fmt.Sprintf("%d pacotes removidos", len(analysis.Items)))
	} else {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("Falha ao limpar pacotes: %s", execResult.Stderr))
	}

	return result, nil
}