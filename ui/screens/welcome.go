package screens

import "github.com/never00rei/a7/ui/layout"

func Welcome(layout layout.Layout) string {
	bodyText := "A7 is a calm space for daily journaling.\n" +
		"Write in plain Markdown and keep journals on your machine.\n\n" +
		"So, let's take 60 seconds to set things up.\n\n" +
		"No accounts.\nNo cloud.\nJust you and your journal."

	pane := layout.TitledPaneWithWidth("Welcome to your Terminal Journal", bodyText, layout.PrimaryPaneWidth())
	return layout.CenterContent(pane)
}
