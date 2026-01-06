package screens

import (
	"strings"

	"github.com/never00rei/a7/ui/layout"
)

func Editor(layout layout.Layout, titleView string, bodyView string, editorErr error, paneWidth int, bodyPaneHeight int) string {
	bodyParts := []string{bodyView}
	if editorErr != nil {
		bodyParts = append(bodyParts, "", "Error: "+editorErr.Error())
	}
	bodyContent := strings.Join(bodyParts, "\n")

	titlePane := layout.TitledPaneWithWidthAndHeight("Title", titleView, paneWidth, 0)
	bodyPane := layout.TitledPaneWithWidthAndHeight("Journal", bodyContent, paneWidth, bodyPaneHeight)

	content := titlePane + "\n" + bodyPane
	return layout.CenterContent(content)
}
