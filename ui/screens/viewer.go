package screens

import "github.com/never00rei/a7/ui/layout"

func Viewer(layout layout.Layout, viewerTitle string, viewerContent string) string {
	title := viewerTitle
	if title == "" {
		title = "Journal"
	}
	pane := layout.TitledPaneWithWidth(title, viewerContent, layout.ContentWidth())
	return layout.CenterContent(pane)
}
