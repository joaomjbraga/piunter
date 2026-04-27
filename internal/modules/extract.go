package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golift.io/xtractr"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type ExtractModule struct {
	BaseModule
	outputDir string
	password  string
	listOnly  bool
	files     []string
}

func NewExtractModule() *ExtractModule {
	return &ExtractModule{
		BaseModule: BaseModule{
			id:          "extract",
			name:        "Extrair",
			description: "Extrai arquivos compactados (rar, 7z, zip, tar, gz, bz2, xz)",
		},
		files: []string{},
	}
}

func (m *ExtractModule) ID()          string { return m.id }
func (m *ExtractModule) Name()        string { return m.name }
func (m *ExtractModule) Description() string { return m.description }

func (m *ExtractModule) IsAvailable() bool { return true }
func (m *ExtractModule) AddFile(path string) { m.files = append(m.files, path) }
func (m *ExtractModule) SetOutputDir(dir string) { m.outputDir = dir }
func (m *ExtractModule) SetPassword(pwd string) { m.password = pwd }
func (m *ExtractModule) SetListOnly(list bool) { m.listOnly = list }

func (m *ExtractModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	if len(m.files) == 0 {
		return &types.AnalysisResult{Module: m.id, Items: []types.CleanableItem{}, TotalSize: 0}, nil
	}
	
	var items []types.CleanableItem
	var totalSize int64
	
	for _, file := range m.files {
		size := getArchiveSize(file)
		if size > 0 {
			totalSize += int64(size)
			items = append(items, types.CleanableItem{
				Path: file,
				Size: int64(size),
				Type: "archive",
				Description: filepath.Base(file),
			})
		}
	}
	
	return &types.AnalysisResult{Module: m.id, Items: items, TotalSize: totalSize}, nil
}

func (m *ExtractModule) Clean(dryRun bool) (*types.CleaningResult, error) {
	if len(m.files) == 0 {
		return &types.CleaningResult{Module: m.id, Success: true, SpaceFreed: 0, ItemsRemoved: 0, Errors: []string{}}, nil
	}

	result := &types.CleaningResult{Module: m.id, Success: true, SpaceFreed: 0, ItemsRemoved: 0, Errors: []string{}}

	for _, archivePath := range m.files {
		err := m.extractOne(archivePath, result, dryRun)
		if err != nil {
			result.Errors = append(result.Errors, err.Error())
			result.Success = false
		}
		result.ItemsRemoved++
	}
	
	return result, nil
}

func getArchiveSize(path string) uint64 {
	xf := &xtractr.XFile{FilePath: path, Password: ""}
	size, _, _, err := xtractr.ExtractFile(xf)
	if err != nil {
		return 0
	}
	return size
}

func (m *ExtractModule) extractOne(archivePath string, result *types.CleaningResult, dryRun bool) error {
	if !utils.FileExists(archivePath) {
		return fmt.Errorf("arquivo não encontrado: %s", archivePath)
	}

	outputDir := m.outputDir
	if outputDir == "" {
		nameWithoutExt := strings.TrimSuffix(archivePath, filepath.Ext(archivePath))
		outputDir = filepath.Dir(archivePath) + "/" + filepath.Base(nameWithoutExt)
	}

	if m.listOnly {
		return m.listArchive(archivePath)
	}

	if dryRun {
		size := getArchiveSize(archivePath)
		msg := "[DRY-RUN] Extrairia " + filepath.Base(archivePath) + " -> " + outputDir
		utils.Info(msg)
		result.SpaceFreed += int64(size)
		return nil
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	xf := &xtractr.XFile{
		FilePath:  archivePath,
		OutputDir: outputDir,
		Password: m.password,
		DirMode:  0755,
		FileMode: 0644,
	}

	size, filesList, archiveList, err := xtractr.ExtractFile(xf)
	if err != nil {
		return fmt.Errorf("erro ao extrair: %w", err)
	}

	result.SpaceFreed += int64(size)
	utils.Info("Extraído: " + filepath.Base(archivePath))
	utils.Item("Destino", outputDir)
	utils.Item("Arquivos", fmt.Sprintf("%d", len(filesList)))
	utils.Item("Subarquivos", fmt.Sprintf("%d", len(archiveList)))
	utils.Item("Tamanho", utils.FormatBytes(int64(size)))

	return nil
}

func (m *ExtractModule) listArchive(archivePath string) error {
	xf := &xtractr.XFile{FilePath: archivePath, Password: m.password}
	size, filesList, archiveList, err := xtractr.ExtractFile(xf)
	if err != nil {
		return fmt.Errorf("erro ao abrir: %w", err)
	}

	utils.Info("Conteúdo de: " + filepath.Base(archivePath))
	fmt.Println()

	for _, f := range filesList {
		fmt.Println("  " + f)
	}

	utils.Space()
	utils.Item("Total", fmt.Sprintf("%d arquivos", len(filesList)))
	utils.Item("Subarquivos", fmt.Sprintf("%d", len(archiveList)))
	utils.Item("Tamanho Total", utils.FormatBytes(int64(size)))

	return nil
}

func GetSupportedExtensions() []string {
	return []string{
		".zip", ".rar", ".7z",
		".tar", ".tar.gz", ".tgz", ".tar.bz2", ".tbz2", ".tar.xz", ".txz", ".tar.z",
		".gz", ".bz2", ".xz", ".z",
	}
}

func IsArchive(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, supported := range GetSupportedExtensions() {
		if ext == supported {
			return true
		}
	}
	return false
}