package screens

import (
	"github.com/charmbracelet/huh"
	"github.com/never00rei/a7/ui/layout"
)

func Settings(layout layout.Layout, form *huh.Form) string {
	bodyText := "Update your journal storage and encryption settings."
	formView := ""
	if form != nil {
		formView = form.View()
	}
	content := bodyText
	if formView != "" {
		content = bodyText + "\n\n" + formView
	}
	pane := layout.TitledPaneWithWidth("Settings", content, layout.PrimaryPaneWidth())
	return layout.CenterContent(pane)
}
