package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/component"
	"github.com/xedeveloper/DifLog/internal/usecase/checkout"
)

type checkoutDoneMsg struct{ err error }

type CheckoutModel struct {
	spinner    component.SpinnerComponent
	useCase    *checkout.CheckoutCommitUseCase
	commitHash string
	loading    bool
	done       bool
	err        error
}

func NewCheckoutModel(useCase *checkout.CheckoutCommitUseCase, commitHash string) CheckoutModel {
	return CheckoutModel{
		spinner:    component.NewSpinnerComponent("Restoring context..."),
		useCase:    useCase,
		commitHash: commitHash,
		loading:    true,
	}
}

func (m CheckoutModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Init(), m.runCheckout())
}

func (m CheckoutModel) runCheckout() tea.Cmd {
	return func() tea.Msg {
		err := m.useCase.Execute(m.commitHash)
		return checkoutDoneMsg{err: err}
	}
}

func (m CheckoutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" || m.done {
			return m, tea.Quit
		}
	case checkoutDoneMsg:
		m.loading = false
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m CheckoutModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).Render("⚡ DifLog — Checkout") + "\n\n"

	if m.loading {
		return header + m.spinner.View()
	}

	if m.err != nil {
		return header + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗ "+m.err.Error())
	}

	return header + lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render("✓ Context restored to commit: "+m.commitHash[:12]+"...")
}
