package ui

func (m AppModel) viewWalkthroughStorage() string {
	bodyText := "Notes are plain Markdown files you control.\n" +
		"Pick a folder on disk and A7 will write your journal there.\n\n" +
		"Example: ~/Documents/journal/"

	formView := ""
	if m.storageForm != nil {
		formView = m.storageForm.View()
	}
	content := bodyText
	if formView != "" {
		content = bodyText + "\n\n" + formView
	}
	pane := m.titledPaneWithWidth("How A7 Stores Notes", content, m.primaryPaneWidth())
	return m.centerContent(pane)
}
