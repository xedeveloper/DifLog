package component

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DiffViewerComponent struct {
	viewport viewport.Model
	style    lipgloss.Style
}

func NewDiffViewerComponent(width, height int) DiffViewerComponent {
	vp := viewport.New(width, height)
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	return DiffViewerComponent{
		viewport: vp,
		style:    lipgloss.NewStyle(),
	}
}

func (d DiffViewerComponent) SetContent(content string) DiffViewerComponent {
	d.viewport.SetContent(content)
	return d
}

func (d DiffViewerComponent) Update(msg tea.Msg) (DiffViewerComponent, tea.Cmd) {
	var cmd tea.Cmd
	d.viewport, cmd = d.viewport.Update(msg)
	return d, cmd
}

func (d DiffViewerComponent) View() string {
	return d.viewport.View()
}
