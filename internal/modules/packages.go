package modules

import (
	"fmt"
	"strconv"
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

func getAptSizes(items []types.CleanableItem) {
	executor := utils.GetExecutor()
	for i, item := range items {
		result := executor.Exec("dpkg-query", "-W", "-f=${Installed-Size}", item.Path)
		if result.Success {
			sizeKB, err := strconv.ParseInt(strings.TrimSpace(result.Stdout), 10, 64)
			if err == nil {
				items[i].Size = sizeKB * 1024
			}
		}
	}
}

func getPacmanSizes(items []types.CleanableItem) {
	executor := utils.GetExecutor()
	for i, item := range items {
		result := executor.Exec("pacman", "-Si", item.Path)
		if !result.Success {
			continue
		}
		for _, line := range strings.Split(result.Stdout, "\n") {
			if strings.HasPrefix(line, "Installed Size") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					items[i].Size = parseSize(strings.TrimSpace(parts[1]))
				}
				break
			}
		}
	}
}

func getDnfSizes(items []types.CleanableItem) {
	executor := utils.GetExecutor()
	for i, item := range items {
		result := executor.Exec("dnf", "info", item.Path)
		if !result.Success {
			continue
		}
		for _, line := range strings.Split(result.Stdout, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Size") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					items[i].Size = parseSize(strings.TrimSpace(parts[1]))
				}
				break
			}
		}
	}
}

func (m *PackagesModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:    []types.CleanableItem{},
		TotalSize: 0,
	}

	var cmd string
	var args []string

	switch m.packageManager {
	case "apt":
		cmd = "apt-get"
		args = []string{"--just-print", "autoremove"}
	case "pacman":
		cmd = "pacman"
		args = []string{"-Qttd"}
	case "dnf":
		cmd = "dnf"
		args = []string{"autoremove", "--assumeno"}
	}

	executor := utils.GetExecutor()
	execResult := executor.Exec(cmd, args...)
	if !execResult.Success {
		if execResult.Code != 1 {
			return result, utils.NewAnalysisError(m.id, fmt.Sprintf("falha ao listar pacotes órfãos com %s", m.packageManager), fmt.Errorf("%s", execResult.Stderr))
		}
		stderr := strings.ToLower(execResult.Stderr)
		if strings.Contains(stderr, "e:") || strings.Contains(stderr, "error") {
			return result, utils.NewAnalysisError(m.id, fmt.Sprintf("falha ao listar pacotes órfãos com %s", m.packageManager), fmt.Errorf("%s", execResult.Stderr))
		}
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
			if !strings.HasPrefix(line, "Remv") {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				orphanCount++
				result.Items = append(result.Items, types.CleanableItem{
					Path:        parts[1],
					Size:        0,
					Type:        "package",
					Description: "Pacote órfão: " + parts[1],
				})
			}
		case "pacman":
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				orphanCount++
				result.Items = append(result.Items, types.CleanableItem{
					Path:        parts[0],
					Size:        0,
					Type:        "package",
					Description: "Pacote órfão: " + parts[0],
				})
			}
		case "dnf":
			if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
				fields := strings.Fields(line)
				if len(fields) >= 1 && fields[0] != "==" && fields[0] != "--" {
					orphanCount++
					result.Items = append(result.Items, types.CleanableItem{
						Path:        fields[0],
						Size:        0,
						Type:        "package",
						Description: "Pacote órfão: " + fields[0],
					})
				}
			}
		}
	}

	if orphanCount > 0 {
		switch m.packageManager {
		case "apt":
			getAptSizes(result.Items)
		case "pacman":
			getPacmanSizes(result.Items)
		case "dnf":
			getDnfSizes(result.Items)
		}

		var total int64
		for _, item := range result.Items {
			total += item.Size
		}
		if total == 0 {
			total = int64(orphanCount) * 10 * utils.MB
		}
		result.TotalSize = total
	}

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

	executor := utils.GetExecutor()
	var execResult types.CommandResult

	switch m.packageManager {
	case "apt":
		execResult = executor.Exec("sudo", "apt", "autoremove", "-y")
	case "pacman":
		var pkgNames []string
		for _, item := range analysis.Items {
			pkgNames = append(pkgNames, item.Path)
		}
		if len(pkgNames) == 0 {
			return result, nil
		}
		args := append([]string{"pacman", "-Rns", "--noconfirm"}, pkgNames...)
		execResult = executor.Exec("sudo", args...)
	case "dnf":
		execResult = executor.Exec("sudo", "dnf", "autoremove", "-y")
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