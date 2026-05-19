package ai_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/infrastructure/ai"
)

func TestAIToolDetector_DetectClaude(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# Claude"), 0644))

	detector := ai.NewAIToolDetector()
	tool, err := detector.DetectAITool(dir)

	assert.NoError(t, err)
	assert.Equal(t, entity.AIToolClaude, tool)
}

func TestAIToolDetector_DetectCopilot(t *testing.T) {
	dir := t.TempDir()
	githubDir := filepath.Join(dir, ".github")
	require.NoError(t, os.MkdirAll(githubDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(githubDir, "copilot-instructions.md"), []byte("# Copilot"), 0644))

	detector := ai.NewAIToolDetector()
	tool, err := detector.DetectAITool(dir)

	assert.NoError(t, err)
	assert.Equal(t, entity.AIToolCopilot, tool)
}

func TestAIToolDetector_DetectUnknown(t *testing.T) {
	dir := t.TempDir()

	detector := ai.NewAIToolDetector()
	tool, err := detector.DetectAITool(dir)

	assert.NoError(t, err)
	assert.Equal(t, entity.AIToolUnknown, tool)
}

func TestAIToolDetector_GetClaudeContextFiles(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# Claude"), 0644))
	claudeDir := filepath.Join(dir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(claudeDir, "settings.json"), []byte("{}"), 0644))

	detector := ai.NewAIToolDetector()
	files := detector.GetContextFilePaths(dir, entity.AIToolClaude)

	assert.Contains(t, files, filepath.Join(dir, "CLAUDE.md"))
	assert.Contains(t, files, filepath.Join(claudeDir, "settings.json"))
}

func TestAIToolDetector_DetectOpenCode(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "opencode.json"), []byte("{}"), 0644))

	detector := ai.NewAIToolDetector()
	tool, err := detector.DetectAITool(dir)

	assert.NoError(t, err)
	assert.Equal(t, entity.AIToolOpenCode, tool)
}

func TestAIToolDetector_DetectOpenCodeJsonc(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "opencode.jsonc"), []byte("{}"), 0644))

	detector := ai.NewAIToolDetector()
	tool, err := detector.DetectAITool(dir)

	assert.NoError(t, err)
	assert.Equal(t, entity.AIToolOpenCode, tool)
}

func TestAIToolDetector_DetectAllTools_Conflict(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# Claude"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "opencode.json"), []byte("{}"), 0644))

	detector := ai.NewAIToolDetector()
	tools, err := detector.DetectAllAITools(dir)

	assert.NoError(t, err)
	assert.Contains(t, tools, entity.AIToolClaude)
	assert.Contains(t, tools, entity.AIToolOpenCode)
	assert.Len(t, tools, 2)
}

func TestAIToolDetector_GetOpenCodeContextFiles(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "opencode.json"), []byte("{}"), 0644))
	openCodeDir := filepath.Join(dir, ".opencode")
	agentsDir := filepath.Join(openCodeDir, "agents")
	require.NoError(t, os.MkdirAll(agentsDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(agentsDir, "my-agent.md"), []byte("# Agent"), 0644))

	detector := ai.NewAIToolDetector()
	files := detector.GetContextFilePaths(dir, entity.AIToolOpenCode)

	assert.Contains(t, files, filepath.Join(dir, "opencode.json"))
	assert.Contains(t, files, filepath.Join(agentsDir, "my-agent.md"))
}

func TestAIToolDetector_GetOpenCodeContextFiles_SkipsNodeModules(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "opencode.json"), []byte("{}"), 0644))

	openCodeDir := filepath.Join(dir, ".opencode")
	agentsDir := filepath.Join(openCodeDir, "agents")
	require.NoError(t, os.MkdirAll(agentsDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(agentsDir, "my-agent.md"), []byte("# Agent"), 0644))

	// Simulate a node_modules directory inside .opencode (as OpenCode places its runtime there)
	nodeModDir := filepath.Join(openCodeDir, "node_modules", "some-pkg")
	require.NoError(t, os.MkdirAll(nodeModDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(nodeModDir, "index.js"), []byte("module.exports={}"), 0644))

	detector := ai.NewAIToolDetector()
	files := detector.GetContextFilePaths(dir, entity.AIToolOpenCode)

	assert.Contains(t, files, filepath.Join(dir, "opencode.json"))
	assert.Contains(t, files, filepath.Join(agentsDir, "my-agent.md"))

	// node_modules files must NOT be staged
	for _, f := range files {
		if filepath.Base(filepath.Dir(f)) == "node_modules" || filepath.Base(filepath.Dir(filepath.Dir(f))) == "node_modules" {
			t.Errorf("node_modules file should have been excluded: %s", f)
		}
	}
}

func TestAIToolDetector_GetOpenCodeContextFiles_SkipsVendorAndGit(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "opencode.json"), []byte("{}"), 0644))

	openCodeDir := filepath.Join(dir, ".opencode")
	skillsDir := filepath.Join(openCodeDir, "skills")
	require.NoError(t, os.MkdirAll(skillsDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(skillsDir, "skill.md"), []byte("# Skill"), 0644))

	// vendor and .git inside .opencode should also be skipped
	vendorDir := filepath.Join(openCodeDir, "vendor", "lib")
	require.NoError(t, os.MkdirAll(vendorDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(vendorDir, "lib.go"), []byte("package lib"), 0644))

	gitDir := filepath.Join(openCodeDir, ".git")
	require.NoError(t, os.MkdirAll(gitDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main"), 0644))

	detector := ai.NewAIToolDetector()
	files := detector.GetContextFilePaths(dir, entity.AIToolOpenCode)

	assert.Contains(t, files, filepath.Join(skillsDir, "skill.md"))
	assert.Len(t, files, 2) // opencode.json + skill.md only
}

func TestAIToolDetector_ClaudeTakesPriorityOverOpenCode(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# Claude"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "opencode.json"), []byte("{}"), 0644))

	detector := ai.NewAIToolDetector()
	tool, err := detector.DetectAITool(dir)

	assert.NoError(t, err)
	assert.Equal(t, entity.AIToolClaude, tool)
}
