package component

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListItem struct {
	title       string
	description string
}

func NewListItem(title, description string) ListItem {
	return ListItem{title: title, description: description}
}

func (i ListItem) Title() string       { return i.title }
func (i ListItem) Description() string { return i.description }
func (i ListItem) FilterValue() string { return i.title }

type ListComponent struct {
	list  list.Model
	style lipgloss.Style
}

func NewListComponent(title string, items []list.Item, width, height int) ListComponent {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("86"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("86"))

	l := list.New(items, delegate, width, height)
	l.Title = title
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	return ListComponent{
		list:  l,
		style: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")),
	}
}

func (l ListComponent) Update(msg tea.Msg) (ListComponent, tea.Cmd) {
	var cmd tea.Cmd
	l.list, cmd = l.list.Update(msg)
	return l, cmd
}

func (l ListComponent) View() string {
	return l.list.View()
}

func (l ListComponent) SelectedItem() list.Item {
	return l.list.SelectedItem()
}
