package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("235")).
		Padding(0, 2).
		Width(60)

	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func renderHeader(subtitle string) string {
	title := headerStyle.Render("⚡ DifLog — AI Context Version Control")
	if subtitle != "" {
		sub := dimStyle.Render("  " + subtitle)
		return title + "\n" + sub
	}
	return title
}

func renderSuccess(msg string) string {
	return successStyle.Render("✓ " + msg)
}

func renderError(msg string) string {
	return errorStyle.Render("✗ " + msg)
}

func renderInfo(msg string) string {
	return infoStyle.Render("ℹ " + msg)
}
