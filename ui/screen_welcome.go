package ui

func (m AppModel) viewWelcome() string {
	left := "A7 is a calm space for daily journaling.\n" +
		"Write in plain Markdown and keep notes on your machine.\n" +
		"Let's take 60 seconds to set things up."

	right := "No accounts.\nNo cloud.\nJust you and your journal."

	body := m.twoPane(left, right)
	footer := m.footer("enter: begin  q: quit")

	return body + footer
}
