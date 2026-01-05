package ui

import "strings"

func (m AppModel) viewEditor() string {
	bodyParts := []string{m.editorBody.View()}
	if m.editorErr != nil {
		bodyParts = append(bodyParts, "", "Error: "+m.editorErr.Error())
	}
	bodyContent := strings.Join(bodyParts, "\n")

	width := m.editorPaneWidth()
	titlePane := m.titledPaneWithWidthAndHeight("Title", m.editorTitle.View(), width, 0)
	_, bodyPaneHeight := m.editorPaneHeights()
	bodyPane := m.titledPaneWithWidthAndHeight("Journal", bodyContent, width, bodyPaneHeight)

	content := titlePane + "\n" + bodyPane
	return m.centerContent(content)
}
