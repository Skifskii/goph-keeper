package screens

import "github.com/charmbracelet/lipgloss"

var (
	LabelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	ItalicStyle  = lipgloss.NewStyle().Italic(true)
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	HelpStyle    = BlurredStyle
)
