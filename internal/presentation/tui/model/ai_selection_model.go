package model

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
)

type AISelectionModel struct {
	detectedTools []entity.AITool
	selectedIdx   int
	confirmed     bool
	cancelled     bool
}

const diflogBanner = `
 ██████╗ ██╗███████╗██╗      ██████╗  ██████╗ 
 ██╔══██╗██║██╔════╝██║     ██╔═══██╗██╔════╝ 
 ██║  ██║██║█████╗  ██║     ██║   ██║██║  ███╗
 ██║  ██║██║██╔══╝  ██║     ██║   ██║██║   ██║
 ██████╔╝██║██║     ███████╗╚██████╔╝╚██████╔╝
 ╚═════╝ ╚═╝╚═╝     ╚══════╝ ╚═════╝  ╚═════╝ `

var aiToolChoices = []struct {
	tool  entity.AITool
	label string
	desc  string
}{
	{entity.AIToolClaude, "ClaudeCode", "Tracks CLAUDE.md, .claude/ settings"},
	{entity.AIToolOpenCode, "OpenCode", "Tracks opencode.json, .opencode/ config"},
	{entity.AIToolCopilot, "GitHub Copilot CLI", "Tracks .github/copilot-instructions.md"},
}

func NewAISelectionModel(detectedTools []entity.AITool) AISelectionModel {
	selectedIdx := 0
	if len(detectedTools) >= 1 {
		for i, choice := range aiToolChoices {
			if choice.tool == detectedTools[0] {
				selectedIdx = i
				break
			}
		}
	}
	return AISelectionModel{
		detectedTools: detectedTools,
		selectedIdx:   selectedIdx,
	}
}

func (m AISelectionModel) Init() tea.Cmd {
	return nil
}

func (m AISelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.cancelled = true
			return m, tea.Quit
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down", "j":
			if m.selectedIdx < len(aiToolChoices)-1 {
				m.selectedIdx++
			}
		case "1":
			m.selectedIdx = 0
		case "2":
			m.selectedIdx = 1
		case "3":
			m.selectedIdx = 2
		case "enter", " ":
			m.confirmed = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m AISelectionModel) View() string {
	bannerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	taglineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	separatorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	detectedBadge := lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	conflictStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	out := bannerStyle.Render(diflogBanner) + "\n"
	out += taglineStyle.Render("  AI Context Version Control") + "\n"
	out += separatorStyle.Render("  ─────────────────────────────────────────────────") + "\n\n"

	switch len(m.detectedTools) {
	case 0:
		out += dimStyle.Render("  No AI tool detected in this project.") + "\n\n"
	case 1:
		out += detectedBadge.Render(fmt.Sprintf("  Auto-detected: %s", aiToolLabel(m.detectedTools[0]))) + "\n\n"
	default:
		names := aiToolLabel(m.detectedTools[0])
		for _, t := range m.detectedTools[1:] {
			names += " and " + aiToolLabel(t)
		}
		out += conflictStyle.Render(fmt.Sprintf("  ⚠  Both %s detected — please select one:", names)) + "\n\n"
	}

	out += lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render("  Select AI tool to track:") + "\n\n"

	for i, choice := range aiToolChoices {
		cursor := "  "
		var label string
		if i == m.selectedIdx {
			cursor = "▶ "
			label = selectedStyle.Render(fmt.Sprintf("[%d] %s", i+1, choice.label))
		} else {
			label = normalStyle.Render(fmt.Sprintf("[%d] %s", i+1, choice.label))
		}
		out += cursor + label + "\n"
		out += "    " + descStyle.Render(choice.desc) + "\n"
	}

	out += "\n" + dimStyle.Render("  ↑/↓ or 1/2/3 to select • Enter to confirm • q to cancel")
	return out
}

func (m AISelectionModel) SelectedTool() entity.AITool {
	if m.cancelled || m.selectedIdx >= len(aiToolChoices) {
		return entity.AIToolUnknown
	}
	return aiToolChoices[m.selectedIdx].tool
}

func (m AISelectionModel) WasCancelled() bool {
	return m.cancelled
}

func aiToolLabel(tool entity.AITool) string {
	for _, c := range aiToolChoices {
		if c.tool == tool {
			return c.label
		}
	}
	return string(tool)
}
