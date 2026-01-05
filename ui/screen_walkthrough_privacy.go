package ui

import "github.com/charmbracelet/lipgloss"

func (m AppModel) viewWalkthroughPrivacy() string {
	header := lipgloss.NewStyle().Bold(true).Padding(0, 2).Render("Privacy and Keys (Optional)")

	bodyText := "Want encrypted notes? You can choose an SSH key now\n" +
		"or skip and set it up later in settings."

	body := m.singlePane(bodyText)
	footer := m.footer("enter: choose key  s: skip  b: back  q: quit")

	return header + "\n" + body + footer
}
