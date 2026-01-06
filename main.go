package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/never00rei/a7/ui/app"
)

func main() {
	if _, err := tea.NewProgram(app.NewAppModel(), tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}
}
