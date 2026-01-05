package ui

func (m AppModel) viewSetup() string {
	journalPath := m.storagePath
	if journalPath == "" {
		journalPath = "Not set yet"
	}
	sshKeyPath := m.sshKeyPath
	if sshKeyPath == "" {
		sshKeyPath = "Not set"
	}
	encryptStatus := "Disabled"
	if m.encrypt {
		encryptStatus = "Enabled"
	}

	bodyText := "Review your choices before finishing setup.\n\n" +
		"Journal folder:\n" + journalPath + "\n\n" +
		"Encryption:\n" + encryptStatus + "\n\n" +
		"SSH key:\n" + sshKeyPath + "\n\n" +
		"Press enter to continue."

	pane := m.titledPaneWithWidth("Setup Review", bodyText, m.primaryPaneWidth())
	return m.centerContent(pane)
}
