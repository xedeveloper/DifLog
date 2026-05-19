package model

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/component"
	"github.com/xedeveloper/DifLog/internal/usecase/add"
)

type addDoneMsg struct {
	result add.AddResult
	err    error
}

type AddModel struct {
	spinner   component.SpinnerComponent
	useCase   *add.AddContextUseCase
	filePaths []string
	aiTool    entity.AITool
	loading   bool
	done      bool
	result    add.AddResult
	err       error
}

func NewAddModel(useCase *add.AddContextUseCase, filePaths []string, aiTool entity.AITool) AddModel {
	return AddModel{
		spinner:   component.NewSpinnerComponent("Staging context files..."),
		useCase:   useCase,
		filePaths: filePaths,
		aiTool:    aiTool,
		loading:   true,
	}
}

func (m AddModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Init(), m.runAdd())
}

func (m AddModel) runAdd() tea.Cmd {
	return func() tea.Msg {
		result, err := m.useCase.Execute(m.filePaths, m.aiTool)
		return addDoneMsg{result: result, err: err}
	}
}

func (m AddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" || m.done {
			return m, tea.Quit
		}
	case addDoneMsg:
		m.loading = false
		m.done = true
		m.result = msg.result
		m.err = msg.err
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m AddModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).Render("⚡ DifLog — Add") + "\n\n"

	if m.loading {
		return header + m.spinner.View()
	}

	if m.err != nil {
		return header + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗ "+m.err.Error())
	}

	var sb strings.Builder
	sb.WriteString(header)
	sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render("✓ Staged files:") + "\n")
	for _, f := range m.result.StagedFiles {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Render("  + "+f) + "\n")
	}

	return sb.String()
}
