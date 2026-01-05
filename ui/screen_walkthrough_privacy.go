package ui

func (m AppModel) viewWalkthroughPrivacy() string {
	bodyText := "Want encrypted notes? You can choose an SSH key now\n" +
		"or skip and set it up later in settings."

	formView := ""
	if m.privacyForm != nil {
		formView = m.privacyForm.View()
	}
	content := bodyText
	if formView != "" {
		content = bodyText + "\n\n" + formView
	}
	pane := m.titledPaneWithWidth("Privacy and Keys (Optional)", content, m.primaryPaneWidth())
	return m.centerContent(pane)
}
