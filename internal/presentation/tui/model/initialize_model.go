package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/component"
	"github.com/xedeveloper/DifLog/internal/usecase/initialize"
)

type initDoneMsg struct{ err error }

type InitializeModel struct {
	spinner     component.SpinnerComponent
	useCase     *initialize.InitializeUseCase
	aiTool      entity.AITool
	loading     bool
	done        bool
	err         error
	successMsg  string
}

func NewInitializeModel(useCase *initialize.InitializeUseCase, aiTool entity.AITool) InitializeModel {
	return InitializeModel{
		spinner: component.NewSpinnerComponent("Initializing DifLog..."),
		useCase: useCase,
		aiTool:  aiTool,
		loading: true,
	}
}

func (m InitializeModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Init(), m.runInit())
}

func (m InitializeModel) runInit() tea.Cmd {
	return func() tea.Msg {
		err := m.useCase.Execute(m.aiTool)
		return initDoneMsg{err: err}
	}
}

func (m InitializeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.done {
			return m, tea.Quit
		}
	case initDoneMsg:
		m.loading = false
		m.done = true
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.successMsg = "DifLog initialized successfully! (.difLog directory created)"
		}
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m InitializeModel) View() string {
	bannerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	taglineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	separatorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	header := bannerStyle.Render(diflogBanner) + "\n" +
		taglineStyle.Render("  AI Context Version Control") + "\n" +
		separatorStyle.Render("  ─────────────────────────────────────────────────") + "\n\n"

	if m.loading {
		return header + m.spinner.View()
	}

	if m.err != nil {
		return header + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗ "+m.err.Error())
	}

	return header + lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render("✓ "+m.successMsg)
}
