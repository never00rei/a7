package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/never00rei/a7/ui"
)

func main() {
	if _, err := tea.NewProgram(ui.NewAppModel(), tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}
}
