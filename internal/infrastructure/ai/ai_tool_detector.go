package ai

import (
	"os"
	"path/filepath"

	"github.com/xedeveloper/DifLog/internal/domain/entity"
)

type AIToolDetector struct{}

func NewAIToolDetector() *AIToolDetector {
	return &AIToolDetector{}
}

func (d *AIToolDetector) DetectAllAITools(projectRoot string) ([]entity.AITool, error) {
	var detected []entity.AITool

	claudeIndicators := []string{"CLAUDE.md", ".claude"}
	for _, indicator := range claudeIndicators {
		if _, err := os.Stat(filepath.Join(projectRoot, indicator)); err == nil {
			detected = append(detected, entity.AIToolClaude)
			break
		}
	}

	openCodeIndicators := []string{"opencode.json", "opencode.jsonc", ".opencode"}
	for _, indicator := range openCodeIndicators {
		if _, err := os.Stat(filepath.Join(projectRoot, indicator)); err == nil {
			detected = append(detected, entity.AIToolOpenCode)
			break
		}
	}

	copilotIndicators := []string{
		filepath.Join(".github", "copilot-instructions.md"),
		".github",
	}
	for _, indicator := range copilotIndicators {
		if _, err := os.Stat(filepath.Join(projectRoot, indicator)); err == nil {
			detected = append(detected, entity.AIToolCopilot)
			break
		}
	}

	return detected, nil
}

func (d *AIToolDetector) DetectAITool(projectRoot string) (entity.AITool, error) {
	all, err := d.DetectAllAITools(projectRoot)
	if err != nil {
		return entity.AIToolUnknown, err
	}
	if len(all) == 0 {
		return entity.AIToolUnknown, nil
	}
	return all[0], nil
}

func (d *AIToolDetector) GetContextFilePaths(projectRoot string, tool entity.AITool) []string {
	switch tool {
	case entity.AIToolClaude:
		return d.claudeContextFilePaths(projectRoot)
	case entity.AIToolOpenCode:
		return d.openCodeContextFilePaths(projectRoot)
	case entity.AIToolCopilot:
		return d.copilotContextFilePaths(projectRoot)
	default:
		return []string{}
	}
}

func (d *AIToolDetector) claudeContextFilePaths(projectRoot string) []string {
	candidates := []string{
		filepath.Join(projectRoot, "CLAUDE.md"),
		filepath.Join(projectRoot, ".claude", "settings.json"),
		filepath.Join(projectRoot, ".claude", "settings.local.json"),
	}
	var existing []string
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			existing = append(existing, p)
		}
	}
	claudeDir := filepath.Join(projectRoot, ".claude")
	if entries, err := os.ReadDir(claudeDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
				existing = append(existing, filepath.Join(claudeDir, entry.Name()))
			}
		}
	}
	return existing
}

func (d *AIToolDetector) openCodeContextFilePaths(projectRoot string) []string {
	candidates := []string{
		filepath.Join(projectRoot, "opencode.json"),
		filepath.Join(projectRoot, "opencode.jsonc"),
	}
	var existing []string
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			existing = append(existing, p)
		}
	}
	// Directories that are never AI context files and should be skipped entirely.
	skipDirs := map[string]bool{
		"node_modules": true,
		"vendor":       true,
		".git":         true,
	}

	openCodeDir := filepath.Join(projectRoot, ".opencode")
	if err := filepath.WalkDir(openCodeDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		existing = append(existing, path)
		return nil
	}); err != nil {
		// directory may not exist — that is fine
	}
	return existing
}

func (d *AIToolDetector) copilotContextFilePaths(projectRoot string) []string {
	candidates := []string{
		filepath.Join(projectRoot, ".github", "copilot-instructions.md"),
	}
	var existing []string
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			existing = append(existing, p)
		}
	}
	githubDir := filepath.Join(projectRoot, ".github")
	if entries, err := os.ReadDir(githubDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
				full := filepath.Join(githubDir, entry.Name())
				if full != filepath.Join(projectRoot, ".github", "copilot-instructions.md") {
					existing = append(existing, full)
				}
			}
		}
	}
	return existing
}
