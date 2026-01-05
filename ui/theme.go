package ui

import "github.com/charmbracelet/lipgloss"

type theme struct {
	FrameBorder lipgloss.Color
	PaneBorder  lipgloss.Color
	Text        lipgloss.Color
	Help        lipgloss.Color
}

func currentTheme() theme {
	return theme{
		FrameBorder: lipgloss.Color("240"),
		PaneBorder:  lipgloss.Color("245"),
		Text:        lipgloss.Color("252"),
		Help:        lipgloss.Color("244"),
	}
}
