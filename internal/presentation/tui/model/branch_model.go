package model

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/component"
	branchUC "github.com/xedeveloper/DifLog/internal/usecase/branch"
)

type branchesLoadedMsg struct {
	branches []entity.Branch
	err      error
}

type branchActionDoneMsg struct{ err error }

type BranchModel struct {
	spinner    component.SpinnerComponent
	list       component.ListComponent
	listUC     *branchUC.ListBranchesUseCase
	loading    bool
	branches   []entity.Branch
	err        error
	width      int
	height     int
}

func NewBranchModel(listUC *branchUC.ListBranchesUseCase) BranchModel {
	return BranchModel{
		spinner: component.NewSpinnerComponent("Loading branches..."),
		listUC:  listUC,
		loading: true,
		width:   80,
		height:  24,
	}
}

func (m BranchModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Init(), m.loadBranches())
}

func (m BranchModel) loadBranches() tea.Cmd {
	return func() tea.Msg {
		branches, err := m.listUC.Execute()
		return branchesLoadedMsg{branches: branches, err: err}
	}
}

func (m BranchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case branchesLoadedMsg:
		m.loading = false
		m.err = msg.err
		m.branches = msg.branches
		if msg.err == nil {
			items := branchesToListItems(msg.branches)
			m.list = component.NewListComponent("Branches", items, m.width, m.height-4)
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

func (m BranchModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).Render("⚡ DifLog — Branches") + "\n\n"

	if m.loading {
		return header + m.spinner.View()
	}

	if m.err != nil {
		return header + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗ "+m.err.Error())
	}

	return header + m.list.View() + "\n" +
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press q to quit")
}

func branchesToListItems(branches []entity.Branch) []list.Item {
	items := make([]list.Item, len(branches))
	for i, b := range branches {
		prefix := "  "
		if b.IsActive {
			prefix = "* "
		}
		shortHash := b.CommitHash
		if len(shortHash) > 8 {
			shortHash = shortHash[:8]
		}
		desc := "latest: " + shortHash
		if shortHash == "" {
			desc = "no commits yet"
		}
		items[i] = component.NewListItem(prefix+b.Name, desc)
	}
	return items
}
