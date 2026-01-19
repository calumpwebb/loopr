package ui

import "github.com/charmbracelet/lipgloss"

var (
	SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	WarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	HeaderStyle  = lipgloss.NewStyle().Bold(true).Border(lipgloss.RoundedBorder())
)
