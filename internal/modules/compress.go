package modules

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
)

type CompressModule struct {
	BaseModule
	files     []string
	outputDir string
	format    string
	level     int
	exclude   []string
}

func NewCompressModule() *CompressModule {
	return &CompressModule{
		BaseModule: BaseModule{
			id:          "compress",
			name:        "Compactar",
			description: "Compacta arquivos (zip, tar.gz, tar.bz2, tar.xz)",
		},
		files:   []string{},
		format:  "zip",
		level:   6,
		exclude: []string{".git", "node_modules", "__pycache__", ".cache", ".DS_Store"},
	}
}

func (m *CompressModule) ID()          string { return m.id }
func (m *CompressModule) Name()        string { return m.name }
func (m *CompressModule) Description() string { return m.description }

func (m *CompressModule) IsAvailable() bool { return true }
func (m *CompressModule) AddFile(path string) { m.files = append(m.files, path) }
func (m *CompressModule) SetFormat(format string) { m.format = format }
func (m *CompressModule) SetOutputDir(dir string) { m.outputDir = dir }
func (m *CompressModule) SetLevel(level int) {
	if level < 1 { level = 1 }
	if level > 9 { level = 9 }
	m.level = level
}

func (m *CompressModule) Analyze(threshold int) (*types.AnalysisResult, error) {
	if len(m.files) == 0 {
		return &types.AnalysisResult{Module: m.id, Items: []types.CleanableItem{}, TotalSize: 0}, nil
	}
	
	var items []types.CleanableItem
	var totalSize int64
	
	for _, file := range m.files {
		size, _ := utils.GetDirSize(file)
		totalSize += size
		desc := filepath.Base(file) + " -> " + m.format
		items = append(items, types.CleanableItem{Path: file, Size: size, Type: "directory", Description: desc})
	}
	
	return &types.AnalysisResult{Module: m.id, Items: items, TotalSize: totalSize}, nil
}

func (m *CompressModule) Clean(dryRun bool) (*types.CleaningResult, error) {
	if len(m.files) == 0 {
		return &types.CleaningResult{Module: m.id, Success: true, SpaceFreed: 0, ItemsRemoved: 0, Errors: []string{}}, nil
	}

	result := &types.CleaningResult{Module: m.id, Success: true, SpaceFreed: 0, ItemsRemoved: 0, Errors: []string{}}

	if dryRun {
		var totalSize int64
		for _, file := range m.files {
			size, _ := utils.GetDirSize(file)
			totalSize += size
		}
		utils.Info(fmt.Sprintf("[DRY-RUN] Compactaria %d arquivos", len(m.files)))
		utils.Item("Tamanho Original", utils.FormatBytes(totalSize))
		utils.Item("Formato", m.format)
		return result, nil
	}

	for _, file := range m.files {
		outputFile, err := m.compressOne(file)
		if err != nil {
			result.Errors = append(result.Errors, err.Error())
			result.Success = false
			continue
		}
		result.ItemsRemoved++
		stat, _ := os.Stat(outputFile)
		result.SpaceFreed += stat.Size()
		utils.Info("Compactado: " + filepath.Base(file))
		utils.Item("Arquivo", outputFile)
		if stat != nil {
			utils.Item("Tamanho", utils.FormatBytes(stat.Size()))
		}
	}
	
	return result, nil
}

func (m *CompressModule) compressOne(file string) (string, error) {
	baseName := filepath.Base(file)
	var outputFile string
	
	switch m.format {
	case "zip":
		outputFile = baseName + ".zip"
		return outputFile, m.compressZip(file, outputFile)
	case "tar.gz", "tgz":
		outputFile = baseName + ".tar.gz"
		return outputFile, m.compressTarGz(file, outputFile)
	case "tar":
		outputFile = baseName + ".tar"
		return outputFile, m.compressTar(file, outputFile)
	default:
		outputFile = baseName + ".zip"
		return outputFile, m.compressZip(file, outputFile)
	}
}

func (m *CompressModule) compressZip(file, output string) error {
	zipFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	return m.walkFiles(file, func(path string, info os.FileInfo) error {
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name, _ = filepath.Rel(filepath.Dir(file), path)
		writer, err := w.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(writer, f)
		return err
	})
}

func (m *CompressModule) compressTarGz(file, output string) error {
	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()

	gzw, err := gzip.NewWriterLevel(outFile, m.level)
	if err != nil {
		return err
	}
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	return m.walkFiles(file, func(path string, info os.FileInfo) error {
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		name, _ := filepath.Rel(filepath.Dir(file), path)
		header.Name = name
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(tw, f)
		return err
	})
}

func (m *CompressModule) compressTar(file, output string) error {
	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()

	tw := tar.NewWriter(outFile)
	defer tw.Close()

	return m.walkFiles(file, func(path string, info os.FileInfo) error {
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		name, _ := filepath.Rel(filepath.Dir(file), path)
		header.Name = name
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(tw, f)
		return err
	})
}

type walkFunc func(string, os.FileInfo) error

func (m *CompressModule) walkFiles(basePath string, fn walkFunc) error {
	return filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if m.shouldExclude(info.Name()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		return fn(path, info)
	})
}

func (m *CompressModule) shouldExclude(name string) bool {
	for _, p := range m.exclude {
		if strings.Contains(name, p) {
			return true
		}
	}
	return false
}

func GetSupportedCompressionFormats() []string {
	return []string{"zip", "tar.gz", "tgz", "tar"}
}