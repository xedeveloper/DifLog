package pathutil

import (
	"fmt"
	"path/filepath"
	"strings"
)

func ValidateBranchName(name string) error {
	if name == "" {
		return fmt.Errorf("branch name cannot be empty")
	}
	if filepath.IsAbs(name) {
		return fmt.Errorf("absolute branch names are not allowed: %s", name)
	}
	clean := filepath.Clean(name)
	if clean == "." || clean == ".." || strings.Contains(name, "..") || strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("invalid branch name: %s", name)
	}
	return nil
}

func SafeRelativePath(projectRoot, absPath string) (string, error) {
	relPath, err := filepath.Rel(projectRoot, absPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve relative path: %w", err)
	}
	if relPath == ".." || strings.HasPrefix(relPath, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes project root: %s", absPath)
	}
	return relPath, nil
}

func SafeJoin(projectRoot, relPath string) (string, error) {
	if filepath.IsAbs(relPath) {
		return "", fmt.Errorf("absolute paths are not allowed: %s", relPath)
	}
	clean := filepath.Clean(relPath)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes project root: %s", relPath)
	}
	return filepath.Join(projectRoot, clean), nil
}
