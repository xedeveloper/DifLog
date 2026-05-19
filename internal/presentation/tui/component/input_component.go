package component

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputComponent struct {
	input textinput.Model
	label string
	style lipgloss.Style
}

func NewInputComponent(label, placeholder string) InputComponent {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	return InputComponent{
		input: ti,
		label: label,
		style: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
	}
}

func (i InputComponent) Update(msg tea.Msg) (InputComponent, tea.Cmd) {
	var cmd tea.Cmd
	i.input, cmd = i.input.Update(msg)
	return i, cmd
}

func (i InputComponent) View() string {
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	return labelStyle.Render(i.label) + "\n" + i.input.View()
}

func (i InputComponent) Value() string {
	return i.input.Value()
}

func (i *InputComponent) SetValue(v string) {
	i.input.SetValue(v)
}

func (i InputComponent) Init() tea.Cmd {
	return textinput.Blink
}
