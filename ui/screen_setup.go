package ui

import "github.com/charmbracelet/lipgloss"

func (m AppModel) viewSetup() string {
	header := lipgloss.NewStyle().Bold(true).Padding(0, 2).Render("Setup")

	bodyText := "Setup form goes here.\n" +
		"Choose your journal folder and optional SSH key.\n\n" +
		"Press enter to simulate setup complete."

	body := m.singlePane(bodyText)
	footer := m.footer("")

	return header + "\n" + body + footer
}
