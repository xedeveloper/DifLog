package component

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpinnerComponent struct {
	spinner spinner.Model
	message string
	style   lipgloss.Style
}

func NewSpinnerComponent(message string) SpinnerComponent {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	return SpinnerComponent{
		spinner: s,
		message: message,
		style:   lipgloss.NewStyle().Foreground(lipgloss.Color("86")),
	}
}

func (s SpinnerComponent) Init() tea.Cmd {
	return s.spinner.Tick
}

func (s SpinnerComponent) Update(msg tea.Msg) (SpinnerComponent, tea.Cmd) {
	var cmd tea.Cmd
	s.spinner, cmd = s.spinner.Update(msg)
	return s, cmd
}

func (s SpinnerComponent) View() string {
	return s.style.Render(s.spinner.View() + " " + s.message)
}
