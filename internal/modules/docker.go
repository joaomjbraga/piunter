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
			description: "Remove todos os recursos Docker",
		},
	}
}

func (m *DockerModule) IsAvailable() bool {
	return utils.IsCommandAvailable("docker")
}

func (m *DockerModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:    []types.CleanableItem{},
		TotalSize: 0,
	}

	executor := utils.GetExecutor()

	execResult := executor.Exec("docker", "system", "df", "--format", "{{.Type}}\t{{.Size}}")
	if !execResult.Success {
		return result, utils.NewAnalysisError(m.id, "falha ao analisar Docker", fmt.Errorf("%s: exit code %d", execResult.Stderr, execResult.Code))
	}

	var totalSize int64
	lines := utils.SplitLines(execResult.Stdout)
	for _, line := range lines {
		parts := utils.SplitColumns(line)
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

	containerResult := executor.Exec("docker", "ps", "-aq")
	containerCount := 0
	if containerResult.Success {
		containerCount = len(utils.SplitLines(containerResult.Stdout))
	}

	imageResult := executor.Exec("docker", "images", "-aq")
	imageCount := 0
	if imageResult.Success {
		imageCount = len(utils.SplitLines(imageResult.Stdout))
	}

	volumeResult := executor.Exec("docker", "volume", "ls", "-q")
	volumeCount := 0
	if volumeResult.Success {
		volumeCount = len(utils.SplitLines(volumeResult.Stdout))
	}

	networkResult := executor.Exec("docker", "network", "ls", "--filter", "type=custom", "-q")
	networkCount := 0
	if networkResult.Success {
		networkCount = len(utils.SplitLines(networkResult.Stdout))
	}

	if containerCount > 0 {
		result.Items = append(result.Items, types.CleanableItem{
			Path:        fmt.Sprintf("%d containers", containerCount),
			Size:        0,
			Type:        "docker",
			Description: fmt.Sprintf("%d containers para remover", containerCount),
		})
	}
	if imageCount > 0 {
		result.Items = append(result.Items, types.CleanableItem{
			Path:        fmt.Sprintf("%d imagens", imageCount),
			Size:        0,
			Type:        "docker",
			Description: fmt.Sprintf("%d imagens para remover", imageCount),
		})
	}
	if volumeCount > 0 {
		result.Items = append(result.Items, types.CleanableItem{
			Path:        fmt.Sprintf("%d volumes", volumeCount),
			Size:        0,
			Type:        "docker",
			Description: fmt.Sprintf("%d volumes para remover", volumeCount),
		})
	}
	if networkCount > 0 {
		result.Items = append(result.Items, types.CleanableItem{
			Path:        fmt.Sprintf("%d redes", networkCount),
			Size:        0,
			Type:        "docker",
			Description: fmt.Sprintf("%d redes para remover", networkCount),
		})
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

	if dryRun {
		result.SpaceFreed = analysis.TotalSize
		result.ItemsRemoved = len(analysis.Items)
		utils.Info(fmt.Sprintf("[DRY-RUN] Removeria todos os recursos Docker (%s)", utils.FormatBytes(analysis.TotalSize)))
		return result, nil
	}

	executor := utils.GetExecutor()

	runningResult := executor.Exec("docker", "ps", "-q")
	if !runningResult.Success {
		result.Errors = append(result.Errors, fmt.Sprintf("Falha ao listar containers: %s", runningResult.Stderr))
		return result, nil
	}
	containerIDs := strings.Fields(runningResult.Stdout)
	if len(containerIDs) > 0 {
		stopArgs := append([]string{"stop"}, containerIDs...)
		stopResult := executor.Exec("docker", stopArgs...)
		if !stopResult.Success {
			result.Errors = append(result.Errors, fmt.Sprintf("Falha ao parar containers: %s", stopResult.Stderr))
		}
	}

	pruneResult := executor.Exec("docker", "system", "prune", "-a", "--volumes", "-f")
	if pruneResult.Success {
		result.SpaceFreed = analysis.TotalSize
		result.ItemsRemoved = len(analysis.Items)
		utils.Item(m.Name(), "Docker totalmente limpo")
	} else {
		result.Success = false
		result.Errors = append(result.Errors, utils.NewCleaningError(m.Name(), "falha ao limpar Docker", fmt.Errorf("%s", pruneResult.Stderr)).Error())
	}

	return result, nil
}

var sizeUnits = []struct {
	suffix string
	mult   int64
}{
	{"TB", 1024 * 1024 * 1024 * 1024},
	{"T", 1024 * 1024 * 1024 * 1024},
	{"GB", 1024 * 1024 * 1024},
	{"G", 1024 * 1024 * 1024},
	{"MB", 1024 * 1024},
	{"M", 1024 * 1024},
	{"KB", 1024},
	{"K", 1024},
	{"B", 1},
}

func parseSize(s string) int64 {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	for _, u := range sizeUnits {
		if strings.HasSuffix(s, u.suffix) {
			var num float64
			fmt.Sscanf(strings.TrimSuffix(s, u.suffix), "%f", &num)
			return int64(num * float64(u.mult))
		}
	}

	var size int64
	fmt.Sscanf(s, "%d", &size)
	return size
}
