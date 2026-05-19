package model

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/component"
	"github.com/xedeveloper/DifLog/internal/usecase/diff"
)

type diffLoadedMsg struct {
	result diff.CommitDiffResult
	err    error
}

type DiffModel struct {
	spinner  component.SpinnerComponent
	viewer   component.DiffViewerComponent
	useCase  *diff.DiffCommitsUseCase
	oldHash  string
	newHash  string
	loading  bool
	err      error
	width    int
	height   int
}

func NewDiffModel(useCase *diff.DiffCommitsUseCase, oldHash, newHash string) DiffModel {
	return DiffModel{
		spinner: component.NewSpinnerComponent("Computing diff..."),
		useCase: useCase,
		oldHash: oldHash,
		newHash: newHash,
		loading: true,
		width:   80,
		height:  24,
	}
}

func (m DiffModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Init(), m.loadDiff())
}

func (m DiffModel) loadDiff() tea.Cmd {
	return func() tea.Msg {
		result, err := m.useCase.Execute(m.oldHash, m.newHash)
		return diffLoadedMsg{result: result, err: err}
	}
}

func (m DiffModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case diffLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			content := buildDiffContent(msg.result)
			m.viewer = component.NewDiffViewerComponent(m.width-4, m.height-6)
			m.viewer = m.viewer.SetContent(content)
		}
		return m, nil
	}

	if !m.loading && m.err == nil {
		var cmd tea.Cmd
		m.viewer, cmd = m.viewer.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m DiffModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).Render("⚡ DifLog — Diff") + "\n\n"

	if m.loading {
		return header + m.spinner.View()
	}

	if m.err != nil {
		return header + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗ "+m.err.Error())
	}

	return header + m.viewer.View() + "\n" +
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("↑/↓ scroll • q quit")
}

func buildDiffContent(result diff.CommitDiffResult) string {
	var sb strings.Builder
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	removeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	fileStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)

	for _, fd := range result.FileDiffs {
		sb.WriteString(fileStyle.Render("=== "+fd.FilePath+" ===") + "\n")
		for _, line := range strings.Split(fd.DiffText, "\n") {
			if strings.HasPrefix(line, "+") {
				sb.WriteString(addStyle.Render(line) + "\n")
			} else if strings.HasPrefix(line, "-") {
				sb.WriteString(removeStyle.Render(line) + "\n")
			} else {
				sb.WriteString(line + "\n")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
