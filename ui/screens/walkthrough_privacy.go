package screens

import (
	"github.com/charmbracelet/huh"
	"github.com/never00rei/a7/ui/layout"
)

func WalkthroughPrivacy(layout layout.Layout, form *huh.Form) string {
	bodyText := "Want encrypted journals? You can choose an SSH key now\n" +
		"or skip and set it up later in settings."

	formView := ""
	if form != nil {
		formView = form.View()
	}
	content := bodyText
	if formView != "" {
		content = bodyText + "\n\n" + formView
	}
	pane := layout.TitledPaneWithWidth("Privacy and Keys (Optional)", content, layout.PrimaryPaneWidth())
	return layout.CenterContent(pane)
}
