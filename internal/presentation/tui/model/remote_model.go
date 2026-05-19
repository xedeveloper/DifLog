package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/component"
	"github.com/xedeveloper/DifLog/internal/usecase/push"
	"github.com/xedeveloper/DifLog/internal/usecase/pull"
)

type remoteDoneMsg struct{ err error }

type RemoteModel struct {
	spinner   component.SpinnerComponent
	progress  component.ProgressComponent
	pushUC    *push.PushToRemoteUseCase
	pullUC    *pull.PullFromRemoteUseCase
	remoteURL string
	operation string
	loading   bool
	done      bool
	err       error
}

func NewPushModel(pushUC *push.PushToRemoteUseCase, remoteURL string) RemoteModel {
	return RemoteModel{
		spinner:   component.NewSpinnerComponent("Pushing to remote..."),
		progress:  component.NewProgressComponent("Push progress"),
		pushUC:    pushUC,
		remoteURL: remoteURL,
		operation: "push",
		loading:   true,
	}
}

func NewPullModel(pullUC *pull.PullFromRemoteUseCase, remoteURL string) RemoteModel {
	return RemoteModel{
		spinner:   component.NewSpinnerComponent("Pulling from remote..."),
		progress:  component.NewProgressComponent("Pull progress"),
		pullUC:    pullUC,
		remoteURL: remoteURL,
		operation: "pull",
		loading:   true,
	}
}

func (m RemoteModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Init(), m.runOperation())
}

func (m RemoteModel) runOperation() tea.Cmd {
	return func() tea.Msg {
		var err error
		if m.operation == "push" {
			err = m.pushUC.Execute(m.remoteURL)
		} else {
			err = m.pullUC.Execute(m.remoteURL)
		}
		return remoteDoneMsg{err: err}
	}
}

func (m RemoteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" || m.done {
			return m, tea.Quit
		}
	case remoteDoneMsg:
		m.loading = false
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m RemoteModel) View() string {
	title := "Push"
	if m.operation == "pull" {
		title = "Pull"
	}
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).Render("⚡ DifLog — "+title) + "\n\n"

	if m.loading {
		return header + m.spinner.View()
	}

	if m.err != nil {
		return header + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗ "+m.err.Error())
	}

	return header + lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render("✓ "+title+" completed successfully")
}
