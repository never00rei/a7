package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	FrameBorder lipgloss.Color
	PaneBorder  lipgloss.Color
	Text        lipgloss.Color
	Help        lipgloss.Color
}

func CurrentTheme() Theme {
	return Theme{
		FrameBorder: lipgloss.Color("240"),
		PaneBorder:  lipgloss.Color("245"),
		Text:        lipgloss.Color("252"),
		Help:        lipgloss.Color("244"),
	}
}
