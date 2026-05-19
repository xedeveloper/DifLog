package component

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProgressComponent struct {
	progress progress.Model
	label    string
	percent  float64
}

func NewProgressComponent(label string) ProgressComponent {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)
	return ProgressComponent{
		progress: p,
		label:    label,
		percent:  0,
	}
}

func (p ProgressComponent) SetPercent(pct float64) ProgressComponent {
	p.percent = pct
	return p
}

func (p ProgressComponent) Update(msg tea.Msg) (ProgressComponent, tea.Cmd) {
	pm, cmd := p.progress.Update(msg)
	if updated, ok := pm.(progress.Model); ok {
		p.progress = updated
	}
	return p, cmd
}

func (p ProgressComponent) View() string {
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	return labelStyle.Render(p.label) + "\n" + p.progress.ViewAs(p.percent)
}
