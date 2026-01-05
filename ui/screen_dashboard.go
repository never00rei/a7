package ui

func (m AppModel) viewDashboard() string {
	if m.storagePath == "" {
		bodyText := "Set a journal folder to see recent entries.\n" +
			"Run setup to choose a storage location."
		pane := m.titledPaneWithWidth("Dashboard", bodyText, m.primaryPaneWidth())
		return m.centerContent(pane)
	}

	if m.dashboardErr != nil {
		bodyText := "Unable to load journals right now.\n" +
			"Check your journal folder and try again."
		pane := m.titledPaneWithWidth("Dashboard", bodyText, m.primaryPaneWidth())
		return m.centerContent(pane)
	}

	left := m.notesList.View()
	if len(m.notes) == 0 {
		left = "No journals yet.\nCreate your first entry."
	}
	right := formatSelectedMeta(m.notesList.SelectedItem(), len(m.notes), m.dashboardNote, m.dashboardNoteErr)
	body := m.twoPaneWithRatioAndTitlesAndWidth("Saved Journals", "Journal Metadata", left, right, dashboardLeftRatio, m.contentWidth())
	return m.centerContent(body)
}
