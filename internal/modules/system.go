package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type LogsModule struct {
	BaseModule
}

func NewLogsModule() *LogsModule {
	return &LogsModule{
		BaseModule: BaseModule{
			id:          "logs",
			name:        "Logs do Sistema",
			description: "Limpa logs antigos do sistema",
		},
	}
}

func (m *LogsModule) IsAvailable() bool {
	return utils.IsRoot() || utils.HasSudoPassword()
}

func getJournalSize() int64 {
	executor := utils.GetExecutor()
	var execResult types.CommandResult
	if utils.IsRoot() {
		execResult = executor.Exec("journalctl", "--disk-usage")
	} else {
		execResult = executor.Exec("sudo", "journalctl", "--disk-usage")
	}
	if !execResult.Success {
		return 0
	}
	for _, line := range strings.Split(execResult.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		for i := len(fields) - 1; i >= 0; i-- {
			size := parseSize(fields[i])
			if size > 0 {
				return size
			}
		}
	}
	return 0
}

func isLogGzFile(name string) bool {
	base := strings.TrimSuffix(name, ".gz")
	if strings.HasSuffix(base, ".log") {
		return true
	}
	if base != "" && base[len(base)-1] >= '0' && base[len(base)-1] <= '9' {
		return true
	}
	return false
}

func getOldGzSize(root string) int64 {
	var total int64
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if isLogGzFile(info.Name()) && time.Since(info.ModTime()).Hours() > 30*24 {
			total += info.Size()
		}
		return nil
	})
	return total
}

func (m *LogsModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:     []types.CleanableItem{},
		TotalSize: 0,
	}

	journalSize := getJournalSize()
	if journalSize > 0 {
		result.Items = append(result.Items, types.CleanableItem{
			Path:        "/var/log/journal",
			Size:        journalSize,
			Type:        "logs",
			Description: "Logs do systemd-journald",
		})
		result.TotalSize += journalSize
	}

	gzSize := getOldGzSize("/var/log")
	if gzSize > 0 {
		result.Items = append(result.Items, types.CleanableItem{
			Path:        "/var/log/*.gz",
			Size:        gzSize,
			Type:        "logs",
			Description: "Ficheiros .gz antigos (>30 dias) em /var/log",
		})
		result.TotalSize += gzSize
	}

	return result, nil
}

func (m *LogsModule) Clean(dryRun bool) (*types.CleaningResult, error) {
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
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %s de logs", utils.FormatBytes(analysis.TotalSize)))
		return result, nil
	}

	executor := utils.GetExecutor()

	var journalResult, findResult types.CommandResult
	if utils.IsRoot() {
		journalResult = executor.Exec("journalctl", "--vacuum-time=7d")
		findResult = executor.Exec("find", "/var/log", "-type", "f", "-name", "*.gz", "-mtime", "+30", "-delete")
	} else {
		journalResult = executor.Exec("sudo", "journalctl", "--vacuum-time=7d")
		findResult = executor.Exec("sudo", "find", "/var/log", "-type", "f", "-name", "*.gz", "-mtime", "+30", "-delete")
	}

	var journalSize, gzSize int64
	for _, item := range analysis.Items {
		switch item.Path {
		case "/var/log/journal":
			journalSize = item.Size
		case "/var/log/*.gz":
			gzSize = item.Size
		}
	}

	var totalFreed int64
	if journalResult.Success {
		totalFreed += journalSize
		result.ItemsRemoved++
	} else {
		result.Errors = append(result.Errors, fmt.Sprintf("Falha ao limpar journal: %s", journalResult.Stderr))
	}

	if findResult.Success {
		totalFreed += gzSize
		result.ItemsRemoved++
	} else {
		result.Errors = append(result.Errors, fmt.Sprintf("Falha ao limpar logs antigos: %s", findResult.Stderr))
	}

	if totalFreed > 0 {
		result.SpaceFreed = totalFreed
		utils.Item(m.Name(), "Logs limpos")
	} else if len(result.Errors) > 0 {
		result.Success = false
	}

	return result, nil
}

type FlatpakModule struct {
	BaseModule
}

func NewFlatpakModule() *FlatpakModule {
	return &FlatpakModule{
		BaseModule: BaseModule{
			id:          "flatpak",
			name:        "Flatpak",
			description: "Remove dados órfãos do Flatpak",
		},
	}
}

func (m *FlatpakModule) IsAvailable() bool {
	return utils.IsCommandAvailable("flatpak")
}

func getFlatpakDirSize(dir string) int64 {
	if !utils.FileExists(dir) {
		return 0
	}
	size, err := utils.GetDirSize(dir)
	if err != nil {
		return 0
	}
	return size
}

func getFlatpakDataDirs() []string {
	dirs := []string{"/var/lib/flatpak"}
	if home, err := os.UserHomeDir(); err == nil {
		userDir := filepath.Join(home, ".local", "share", "flatpak")
		dirs = append(dirs, userDir)
	}
	return dirs
}

func (m *FlatpakModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:     []types.CleanableItem{},
		TotalSize: 0,
	}

	executor := utils.GetExecutor()

	// List installed apps (for info, not for size estimation)
	execResult := executor.Exec("flatpak", "list", "--app", "--columns=ref")
	if execResult.Success {
		lines := strings.Split(execResult.Stdout, "\n")
		appCount := 0
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				appCount++
			}
		}
		if appCount > 0 {
			result.Items = append(result.Items, types.CleanableItem{
				Path:        "/var/lib/flatpak/app",
				Size:        0,
				Type:        "flatpak-apps",
				Description: fmt.Sprintf("%d aplicações Flatpak instaladas", appCount),
			})
		}
	}

	// Measure actual size of runtimes and .removed across all install locations
	flatpakDirs := getFlatpakDataDirs()
	var totalCleanable int64
	for _, base := range flatpakDirs {
		runtimeSize := getFlatpakDirSize(base + "/runtime")
		removedSize := getFlatpakDirSize(base + "/.removed")
		if runtimeSize > 0 || removedSize > 0 {
			totalCleanable += runtimeSize + removedSize
		}
	}

	if totalCleanable > 0 {
		result.Items = append(result.Items, types.CleanableItem{
			Path:        "flatpak-data",
			Size:        totalCleanable,
			Type:        "flatpak-runtimes",
			Description: "Dados Flatpak removíveis (runtimes + .removed)",
		})
		result.TotalSize += totalCleanable
	}

	return result, nil
}

func (m *FlatpakModule) Clean(dryRun bool) (*types.CleaningResult, error) {
	result := &types.CleaningResult{
		Module:       m.id,
		Success:      true,
		SpaceFreed:   0,
		ItemsRemoved: 0,
		Errors:       []string{},
	}

	flatpakDirs := getFlatpakDataDirs()
	var totalCleanable int64
	for _, base := range flatpakDirs {
		totalCleanable += getFlatpakDirSize(base + "/runtime")
		totalCleanable += getFlatpakDirSize(base + "/.removed")
	}

	if totalCleanable == 0 {
		return result, nil
	}

	if dryRun {
		result.SpaceFreed = totalCleanable
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %s de dados Flatpak", utils.FormatBytes(totalCleanable)))
		return result, nil
	}

	executor := utils.GetExecutor()
	execResult := executor.Exec("flatpak", "remove", "--unused", "-y")
	if execResult.Success {
		result.SpaceFreed = totalCleanable
		result.ItemsRemoved = 1
		utils.Item(m.Name(), "Flatpak limpo")
	} else {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("Falha ao limpar Flatpak: %s", execResult.Stderr))
	}

	return result, nil
}

type snapRevision struct {
	name string
	rev  int
}

func getDisabledSnapRevisions() ([]snapRevision, error) {
	executor := utils.GetExecutor()
	execResult := executor.Exec("snap", "list", "--all")
	if !execResult.Success {
		return nil, fmt.Errorf("falha ao listar snaps: %s", execResult.Stderr)
	}

	lines := strings.Split(execResult.Stdout, "\n")
	var disabled []snapRevision
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		if fields[0] == "Name" {
			continue
		}
		notes := fields[len(fields)-1]
		if notes != "disabled" && !strings.Contains(notes, "disabled") {
			continue
		}
		rev, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}
		disabled = append(disabled, snapRevision{name: fields[0], rev: rev})
	}
	return disabled, nil
}

type SnapModule struct {
	BaseModule
}

func NewSnapModule() *SnapModule {
	return &SnapModule{
		BaseModule: BaseModule{
			id:          "snap",
			name:        "Snap",
			description: "Remove revisões antigas do Snap",
		},
	}
}

func (m *SnapModule) IsAvailable() bool {
	return utils.IsCommandAvailable("snap")
}

func (m *SnapModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	result := &types.AnalysisResult{
		Module:    m.id,
		Items:     []types.CleanableItem{},
		TotalSize: 0,
	}

	disabled, err := getDisabledSnapRevisions()
	if err != nil {
		return result, utils.NewAnalysisError(m.id, "falha ao analisar snaps", err)
	}

	var totalSize int64
	for _, s := range disabled {
		snapPath := fmt.Sprintf("/var/lib/snapd/snaps/%s_%d.snap", s.name, s.rev)
		size := int64(0)
		if fi, err := os.Stat(snapPath); err == nil {
			size = fi.Size()
		}
		result.Items = append(result.Items, types.CleanableItem{
			Path:        snapPath,
			Size:        size,
			Type:        "snap-revision",
			Description: fmt.Sprintf("Revisão %d de %s", s.rev, s.name),
		})
		totalSize += size
	}

	result.TotalSize = totalSize
	return result, nil
}

func (m *SnapModule) Clean(dryRun bool) (*types.CleaningResult, error) {
	result := &types.CleaningResult{
		Module:       m.id,
		Success:      true,
		SpaceFreed:   0,
		ItemsRemoved: 0,
		Errors:       []string{},
	}

	disabled, err := getDisabledSnapRevisions()
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, err.Error())
		return result, nil
	}

	if len(disabled) == 0 {
		return result, nil
	}

	if dryRun {
		var totalSize int64
		for _, s := range disabled {
			snapPath := fmt.Sprintf("/var/lib/snapd/snaps/%s_%d.snap", s.name, s.rev)
			if fi, err := os.Stat(snapPath); err == nil {
				totalSize += fi.Size()
			}
		}
		result.SpaceFreed = totalSize
		utils.Info(fmt.Sprintf("[DRY-RUN] Limparia %s de revisões Snap", utils.FormatBytes(totalSize)))
		return result, nil
	}

	var removed int
	var freed int64
	executor := utils.GetExecutor()
	for _, s := range disabled {
		snapPath := fmt.Sprintf("/var/lib/snapd/snaps/%s_%d.snap", s.name, s.rev)
		var execResult types.CommandResult
		if utils.IsRoot() {
			execResult = executor.Exec("snap", "remove", s.name, "--revision", strconv.Itoa(s.rev))
		} else {
			execResult = executor.Exec("sudo", "snap", "remove", s.name, "--revision", strconv.Itoa(s.rev))
		}
		if execResult.Success {
			removed++
			if fi, err := os.Stat(snapPath); err == nil {
				freed += fi.Size()
			}
		} else {
			result.Errors = append(result.Errors, fmt.Sprintf("Falha ao remover revisão %d de %s: %s", s.rev, s.name, execResult.Stderr))
		}
	}

	if removed > 0 {
		result.SpaceFreed = freed
		result.ItemsRemoved = removed
		utils.Item(m.Name(), "Snap limpo")
	} else {
		result.Success = false
	}

	return result, nil
}