package screens

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/never00rei/a7/ui/layout"
)

func Editor(layout layout.Layout, titleView string, bodyView string, editorErr error) string {
	bodyParts := []string{bodyView}
	if editorErr != nil {
		bodyParts = append(bodyParts, "", "Error: "+editorErr.Error())
	}
	bodyContent := strings.Join(bodyParts, "\n")

	width := layout.EditorPaneWidth()
	titlePane := layout.TitledPaneWithWidthAndHeight("Title", titleView, width, 0)
	titleHeight := lipgloss.Height(titlePane)
	bodyPaneHeight := layout.BodyHeight() - titleHeight
	if bodyPaneHeight < 3 {
		bodyPaneHeight = 3
	}
	bodyPane := layout.TitledPaneWithWidthAndHeight("Journal", bodyContent, width, bodyPaneHeight)

	content := titlePane + "\n" + bodyPane
	return layout.CenterContent(content)
}
