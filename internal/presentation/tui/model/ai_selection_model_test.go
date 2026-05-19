package model_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/model"
)

func TestAISelectionModel_DefaultsToDetectedTool(t *testing.T) {
	m := model.NewAISelectionModel([]entity.AITool{entity.AIToolCopilot})
	assert_selectedTool(t, m, entity.AIToolCopilot)
}

func TestAISelectionModel_DefaultsToClaudeWhenUnknown(t *testing.T) {
	m := model.NewAISelectionModel([]entity.AITool{})
	assert_selectedTool(t, m, entity.AIToolClaude)
}

func TestAISelectionModel_ConfirmWithEnter(t *testing.T) {
	m := model.NewAISelectionModel([]entity.AITool{entity.AIToolClaude})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	result := updated.(model.AISelectionModel)

	if result.WasCancelled() {
		t.Fatal("expected not cancelled after Enter")
	}
	if result.SelectedTool() != entity.AIToolClaude {
		t.Fatalf("expected claude, got %s", result.SelectedTool())
	}
}

func TestAISelectionModel_CancelWithQ(t *testing.T) {
	m := model.NewAISelectionModel([]entity.AITool{entity.AIToolClaude})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	result := updated.(model.AISelectionModel)

	if !result.WasCancelled() {
		t.Fatal("expected cancelled after q")
	}
}

func TestAISelectionModel_NavigateDown(t *testing.T) {
	m := model.NewAISelectionModel([]entity.AITool{entity.AIToolClaude})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	result := updated.(model.AISelectionModel)

	if result.SelectedTool() != entity.AIToolOpenCode {
		t.Fatalf("expected opencode after down, got %s", result.SelectedTool())
	}
}

func TestAISelectionModel_SelectByNumber(t *testing.T) {
	m := model.NewAISelectionModel([]entity.AITool{entity.AIToolClaude})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2")})
	result := updated.(model.AISelectionModel)

	if result.SelectedTool() != entity.AIToolOpenCode {
		t.Fatalf("expected opencode after pressing 2, got %s", result.SelectedTool())
	}
}

func TestAISelectionModel_ViewRendersWithoutPanic(t *testing.T) {
	for _, tools := range [][]entity.AITool{
		{entity.AIToolClaude},
		{entity.AIToolOpenCode},
		{entity.AIToolCopilot},
		{},
	} {
		m := model.NewAISelectionModel(tools)
		view := m.View()
		if view == "" {
			t.Fatalf("View() returned empty string for tools %v", tools)
		}
	}
}

func TestAISelectionModel_OpenCodeOption(t *testing.T) {
	m := model.NewAISelectionModel([]entity.AITool{entity.AIToolOpenCode})
	assert_selectedTool(t, m, entity.AIToolOpenCode)
}

func TestAISelectionModel_ConflictHeader(t *testing.T) {
	m := model.NewAISelectionModel([]entity.AITool{entity.AIToolClaude, entity.AIToolOpenCode})
	view := m.View()
	if !strings.Contains(view, "⚠") {
		t.Fatal("expected conflict warning in view")
	}
}

func TestAISelectionModel_SelectOpenCodeByNumber(t *testing.T) {
	m := model.NewAISelectionModel([]entity.AITool{})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2")})
	result := updated.(model.AISelectionModel)

	if result.SelectedTool() != entity.AIToolOpenCode {
		t.Fatalf("expected opencode after pressing 2, got %s", result.SelectedTool())
	}
}

func TestAISelectionModel_SelectCopilotByNumber(t *testing.T) {
	m := model.NewAISelectionModel([]entity.AITool{})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("3")})
	result := updated.(model.AISelectionModel)

	if result.SelectedTool() != entity.AIToolCopilot {
		t.Fatalf("expected copilot after pressing 3, got %s", result.SelectedTool())
	}
}

func assert_selectedTool(t *testing.T, m model.AISelectionModel, expected entity.AITool) {
	t.Helper()
	if got := m.SelectedTool(); got != expected {
		t.Fatalf("expected selected tool %s, got %s", expected, got)
	}
}
