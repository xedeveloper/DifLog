package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/component"
	"github.com/xedeveloper/DifLog/internal/usecase/detect"
)

type detectDoneMsg struct {
	result detect.DetectionResult
	err    error
}

type DetectModel struct {
	spinner  component.SpinnerComponent
	useCase  *detect.DetectAIToolUseCase
	loading  bool
	result   detect.DetectionResult
	err      error
}

func NewDetectModel(useCase *detect.DetectAIToolUseCase) DetectModel {
	return DetectModel{
		spinner: component.NewSpinnerComponent("Detecting AI tool..."),
		useCase: useCase,
		loading: true,
	}
}

func (m DetectModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Init(), m.runDetect())
}

func (m DetectModel) runDetect() tea.Cmd {
	return func() tea.Msg {
		result, err := m.useCase.Execute()
		return detectDoneMsg{result: result, err: err}
	}
}

func (m DetectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" || !m.loading {
			return m, tea.Quit
		}
	case detectDoneMsg:
		m.loading = false
		m.result = msg.result
		m.err = msg.err
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m DetectModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	header := headerStyle.Render("⚡ DifLog — Detect AI Tool") + "\n\n"

	if m.loading {
		return header + m.spinner.View()
	}

	if m.err != nil {
		return header + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗ "+m.err.Error())
	}

	toolStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	fileStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	if m.result.AITool == "unknown" {
		return header + dimStyle.Render("No AI tool detected in this project.")
	}

	out := header
	out += toolStyle.Render("  Detected: "+aiToolLabel(m.result.AITool)) + "\n\n"

	if len(m.result.ContextFiles) > 0 {
		out += lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render("  Context files:") + "\n"
		for _, f := range m.result.ContextFiles {
			out += "    " + fileStyle.Render("+ "+f) + "\n"
		}
	} else {
		out += dimStyle.Render("  No context files found yet.")
	}

	return out
}
