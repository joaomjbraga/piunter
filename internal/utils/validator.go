package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PathValidator struct {
	allowedBasePaths []string
	excludedPaths  []string
}

func NewPathValidator() *PathValidator {
	return &PathValidator{
		allowedBasePaths: []string{"/home", "/var", "/tmp", "/opt", "/usr"},
		excludedPaths: []string{},
	}
}

func (v *PathValidator) WithAllowedBasePaths(paths []string) *PathValidator {
	v.allowedBasePaths = paths
	return v
}

func (v *PathValidator) WithExcludedPaths(paths []string) *PathValidator {
	v.excludedPaths = paths
	return v
}

func (v *PathValidator) Validate(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("caminho inválido: %w", err)
	}

	if v.isUnderSymlink(absPath) {
		return fmt.Errorf("symlink detectado: %s", path)
	}

	if v.isExcluded(absPath) {
		return fmt.Errorf("caminho excluded: %s", path)
	}

	return nil
}

func (v *PathValidator) ValidateForWrite(path string) error {
	if err := v.Validate(path); err != nil {
		return err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	parent := filepath.Dir(absPath)
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		return fmt.Errorf("diretório pai não existe: %s", parent)
	}

	return nil
}

func (v *PathValidator) isUnderSymlink(path string) bool {
	var current string
	if filepath.IsAbs(path) {
		current = string(filepath.Separator)
	} else {
		current, _ = os.Getwd()
	}

	parts := strings.Split(filepath.Clean(path), string(filepath.Separator))
	for i := range parts {
		current = current + string(filepath.Separator) + parts[i]
		linkTarget, err := os.Readlink(current)
		if err == nil {
			if !filepath.IsAbs(linkTarget) {
				linkTarget = filepath.Join(filepath.Dir(current), linkTarget)
			}
			linkAbs, _ := filepath.Abs(linkTarget)
			pathAbs, _ := filepath.Abs(path)
			if !strings.HasPrefix(linkAbs, pathAbs) {
				return true
			}
		}
	}

	return false
}

func (v *PathValidator) isExcluded(path string) bool {
	for _, excluded := range v.excludedPaths {
		if strings.HasPrefix(path, excluded) {
			return true
		}
	}
	return false
}

func ValidatePath(path string) error {
	validator := NewPathValidator()
	return validator.Validate(path)
}

func ValidatePathForWrite(path string) error {
	validator := NewPathValidator()
	return validator.ValidateForWrite(path)
}

type SafePathCleaner struct{}

func NewSafePathCleaner() *SafePathCleaner {
	return &SafePathCleaner{}
}

func (c *SafePathCleaner) Clean(path string) (string, error) {
	cleaned := filepath.Clean(path)
	
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("path traversal não permitido")
	}
	
	abs, err := filepath.Abs(cleaned)
	if err != nil {
		return "", err
	}
	
	return abs, nil
}

func (c *SafePathCleaner) IsSafe(path string) bool {
	_, err := c.Clean(path)
	return err == nil
}

func SafeCleanPath(path string) (string, error) {
	return NewSafePathCleaner().Clean(path)
}

func IsPathSafe(path string) bool {
	return NewSafePathCleaner().IsSafe(path)
}

type PathTraversalChecker struct {
	rootDirs []string
}

func NewPathTraversalChecker(rootDirs []string) *PathTraversalChecker {
	return &PathTraversalChecker{rootDirs: rootDirs}
}

func (c *PathTraversalChecker) Check(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("caminho inválido: %w", err)
	}

	for _, root := range c.rootDirs {
		absRoot, err := filepath.Abs(root)
		if err != nil {
			continue
		}
		if strings.HasPrefix(absPath, absRoot) {
			return nil
		}
	}

	return fmt.Errorf("caminho fora dos diretórios permitidos: %s", path)
}

func IsPathWithin(path string, roots []string) bool {
	checker := NewPathTraversalChecker(roots)
	return checker.Check(path) == nil
}