package ui

func (m AppModel) viewDashboard() string {
	bodyText := "This is a placeholder dashboard.\n" +
		"Recent notes and quick actions will live here."

	pane := m.titledPaneWithWidth("Dashboard", bodyText, m.primaryPaneWidth())
	return m.centerContent(pane)
}
