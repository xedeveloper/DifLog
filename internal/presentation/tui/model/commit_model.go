package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/component"
	"github.com/xedeveloper/DifLog/internal/usecase/commit"
)

type commitDoneMsg struct {
	commit entity.Commit
	err    error
}

type CommitModel struct {
	spinner  component.SpinnerComponent
	input    component.InputComponent
	useCase  *commit.CommitContextUseCase
	message  string
	loading  bool
	done     bool
	result   entity.Commit
	err      error
	inputMode bool
}

func NewCommitModel(useCase *commit.CommitContextUseCase, message string) CommitModel {
	m := CommitModel{
		spinner: component.NewSpinnerComponent("Committing contexts..."),
		input:   component.NewInputComponent("Commit Message:", "Enter commit message..."),
		useCase: useCase,
		message: message,
	}

	if message == "" {
		m.inputMode = true
	} else {
		m.loading = true
	}

	return m
}

func (m CommitModel) Init() tea.Cmd {
	if m.inputMode {
		return m.input.Init()
	}
	return tea.Batch(m.spinner.Init(), m.runCommit(m.message))
}

func (m CommitModel) runCommit(msg string) tea.Cmd {
	return func() tea.Msg {
		c, err := m.useCase.Execute(msg)
		return commitDoneMsg{commit: c, err: err}
	}
}

func (m CommitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.inputMode && msg.String() == "enter" {
			val := m.input.Value()
			if val != "" {
				m.inputMode = false
				m.loading = true
				m.message = val
				return m, tea.Batch(m.spinner.Init(), m.runCommit(val))
			}
		}
		if m.done {
			return m, tea.Quit
		}
	case commitDoneMsg:
		m.loading = false
		m.done = true
		m.result = msg.commit
		m.err = msg.err
		return m, tea.Quit
	}

	if m.inputMode {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m CommitModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).Render("⚡ DifLog — Commit") + "\n\n"

	if m.inputMode {
		return header + m.input.View() + "\n\n" +
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press Enter to commit, Ctrl+C to cancel")
	}

	if m.loading {
		return header + m.spinner.View()
	}

	if m.err != nil {
		return header + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗ "+m.err.Error())
	}

	hashStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	return header +
		lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render("✓ Committed successfully") + "\n" +
		hashStyle.Render("  Hash: "+m.result.Hash[:12]+"...") + "\n" +
		lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render("  Message: "+m.result.Message)
}
