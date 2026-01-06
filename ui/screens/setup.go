package screens

import "github.com/never00rei/a7/ui/layout"

func Setup(layout layout.Layout, storagePath string, sshKeyPath string, encrypt bool) string {
	journalPath := storagePath
	if journalPath == "" {
		journalPath = "Not set yet"
	}
	keyPath := sshKeyPath
	if sshKeyPath == "" {
		keyPath = "Not set"
	}
	encryptStatus := "Disabled"
	if encrypt {
		encryptStatus = "Enabled"
	}

	bodyText := "Review your choices before finishing setup.\n\n" +
		"Journal folder:\n" + journalPath + "\n\n" +
		"Encryption:\n" + encryptStatus + "\n\n" +
		"SSH key:\n" + keyPath + "\n\n" +
		"Press enter to continue."

	pane := layout.TitledPaneWithWidth("Setup Review", bodyText, layout.PrimaryPaneWidth())
	return layout.CenterContent(pane)
}
