package ui

func (m AppModel) viewViewer() string {
	title := m.viewerTitle
	if title == "" {
		title = "Journal"
	}
	content := m.viewer.View()
	pane := m.titledPaneWithWidth(title, content, m.contentWidth())
	return m.centerContent(pane)
}
