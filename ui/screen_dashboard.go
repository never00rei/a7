package ui

import "github.com/charmbracelet/lipgloss"

func (m AppModel) viewDashboard() string {
	header := lipgloss.NewStyle().Bold(true).Padding(0, 2).Render("Dashboard")

	bodyText := "This is a placeholder dashboard.\n" +
		"Recent notes and quick actions will live here."

	body := m.singlePane(bodyText)
	footer := m.footer("q: quit")

	return header + "\n" + body + footer
}
