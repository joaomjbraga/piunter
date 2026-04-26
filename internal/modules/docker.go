package modules

import (
	"fmt"
	"strings"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type DockerModule struct {
	BaseModule
}

func NewDockerModule() *DockerModule {
	return &DockerModule{
		BaseModule: BaseModule{
			id:          "docker",
			name:        "Docker",
			description: "Remove containers e imagens Docker não utilizados",
		},
	}
}

func (m *DockerModule) IsAvailable() bool {
	return utils.IsCommandAvailable("docker")
}

func (m *DockerModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:     []types.CleanableItem{},
		TotalSize: 0,
	}

	execResult := utils.Exec("docker", "system", "df", "--format", "{{.Type}}\t{{.Size}}")
	if !execResult.Success {
		return result, nil
	}

	var totalSize int64
	lines := splitLines(execResult.Stdout)
	for _, line := range lines {
		parts := splitColumns(line)
		if len(parts) >= 2 {
			size := parseSize(parts[1])
			totalSize += size
			result.Items = append(result.Items, types.CleanableItem{
				Path:        parts[0],
				Size:        size,
				Type:        "docker",
				Description: fmt.Sprintf("Docker %s", parts[0]),
			})
		}
	}

	result.TotalSize = totalSize
	return result, nil
}

func (m *DockerModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %s do Docker", utils.FormatBytes(analysis.TotalSize)))
		return result, nil
	}

	execResult := utils.Exec("docker", "system", "prune", "-a", "-f")
	if execResult.Success {
		result.SpaceFreed = analysis.TotalSize
		result.ItemsRemoved = len(analysis.Items)
		utils.Item(m.Name(), "Docker limpo")
	} else {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("Falha ao limpar Docker: %s", execResult.Stderr))
	}

	return result, nil
}

func splitLines(s string) []string {
	if s == "" {
		return []string{}
	}
	var lines []string
	for _, line := range splitSimple(s, '\n') {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func splitColumns(s string) []string {
	return splitSimple(s, '\t')
}

func splitSimple(s string, sep rune) []string {
	var result []string
	var current string
	for _, r := range s {
		if r == sep {
			result = append(result, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	result = append(result, current)
	return result
}

func parseSize(s string) int64 {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	multipliers := map[string]int64{
		"B":  1,
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}

	for unit, mult := range multipliers {
		if strings.HasSuffix(s, unit) {
			var num float64
			fmt.Sscanf(strings.TrimSuffix(s, unit), "%f", &num)
			return int64(num * float64(mult))
		}
	}

	var size int64
	fmt.Sscanf(s, "%d", &size)
	return size
}