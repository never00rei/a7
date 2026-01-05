package ui

import "github.com/charmbracelet/lipgloss"

func (m AppModel) viewWalkthroughStorage() string {
	header := lipgloss.NewStyle().Bold(true).Padding(0, 2).Render("How A7 Stores Notes")

	bodyText := "Notes are plain Markdown files you control.\n" +
		"Pick a folder on disk and A7 will write your journal there.\n\n" +
		"Example: ~/Documents/journal/"

	body := m.singlePane(bodyText)
	footer := m.footer("")

	return header + "\n" + body + footer
}
