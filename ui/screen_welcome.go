package ui

func (m AppModel) viewWelcome() string {
	bodyText := "A7 is a calm space for daily journaling.\n" +
		"Write in plain Markdown and keep journals on your machine.\n" +
		"Let's take 60 seconds to set things up.\n\n" +
		"No accounts.\nNo cloud.\nJust you and your journal."

	pane := m.titledPaneWithWidth("Welcome to your Terminal Journal", bodyText, m.primaryPaneWidth())
	return m.centerContent(pane)
}
