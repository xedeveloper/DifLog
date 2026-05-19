package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/component"
	"github.com/xedeveloper/DifLog/internal/usecase/log"
)

type logLoadedMsg struct {
	commits []entity.Commit
	err     error
}

type LogModel struct {
	spinner  component.SpinnerComponent
	list     component.ListComponent
	useCase  *log.ViewCommitLogUseCase
	loading  bool
	commits  []entity.Commit
	err      error
	width    int
	height   int
}

func NewLogModel(useCase *log.ViewCommitLogUseCase) LogModel {
	return LogModel{
		spinner: component.NewSpinnerComponent("Loading commit history..."),
		useCase: useCase,
		loading: true,
		width:   80,
		height:  24,
	}
}

func (m LogModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Init(), m.loadCommits())
}

func (m LogModel) loadCommits() tea.Cmd {
	return func() tea.Msg {
		commits, err := m.useCase.Execute()
		return logLoadedMsg{commits: commits, err: err}
	}
}

func (m LogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case logLoadedMsg:
		m.loading = false
		m.err = msg.err
		m.commits = msg.commits
		if msg.err == nil {
			items := commitsToListItems(msg.commits)
			m.list = component.NewListComponent("Commit History", items, m.width, m.height-4)
		}
		return m, nil
	}

	if !m.loading && m.err == nil {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m LogModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).Render("⚡ DifLog — Log") + "\n\n"

	if m.loading {
		return header + m.spinner.View()
	}

	if m.err != nil {
		return header + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗ "+m.err.Error())
	}

	if len(m.commits) == 0 {
		return header + lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("No commits yet.")
	}

	return header + m.list.View() + "\n" +
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press q to quit")
}

func commitsToListItems(commits []entity.Commit) []list.Item {
	items := make([]list.Item, len(commits))
	for i, c := range commits {
		shortHash := c.Hash
		if len(shortHash) > 8 {
			shortHash = shortHash[:8]
		}
		desc := fmt.Sprintf("%s | %s | %d file(s)", c.Timestamp.Format("2006-01-02 15:04"), c.Author, len(c.Contexts))
		items[i] = component.NewListItem(shortHash+" "+c.Message, desc)
	}
	return items
}
