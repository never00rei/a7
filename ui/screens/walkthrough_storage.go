package screens

import (
	"github.com/charmbracelet/huh"
	"github.com/never00rei/a7/ui/layout"
)

func WalkthroughStorage(layout layout.Layout, form *huh.Form) string {
	bodyText := "Journal entries are plain Markdown files you control.\n" +
		"Pick a folder on disk and A7 will write your journal entries there.\n\n" +
		"Example: ~/Documents/journal/"

	formView := ""
	if form != nil {
		formView = form.View()
	}
	content := bodyText
	if formView != "" {
		content = bodyText + "\n\n" + formView
	}
	pane := layout.TitledPaneWithWidth("How A7 Stores Journals", content, layout.PrimaryPaneWidth())
	return layout.CenterContent(pane)
}
