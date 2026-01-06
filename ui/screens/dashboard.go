package screens

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/never00rei/a7/journal"
	"github.com/never00rei/a7/ui/components"
	"github.com/never00rei/a7/ui/layout"
)

func Dashboard(layout layout.Layout, storagePath string, dashboardErr error, notes []journal.NoteInfo, notesList list.Model, dashboardNote *journal.Note, dashboardNoteErr error) string {
	if storagePath == "" {
		bodyText := "Set a journal folder to see recent entries.\n" +
			"Run setup to choose a storage location."
		pane := layout.TitledPaneWithWidth("Dashboard", bodyText, layout.PrimaryPaneWidth())
		return layout.CenterContent(pane)
	}

	if dashboardErr != nil {
		bodyText := "Unable to load journals right now.\n" +
			"Check your journal folder and try again."
		pane := layout.TitledPaneWithWidth("Dashboard", bodyText, layout.PrimaryPaneWidth())
		return layout.CenterContent(pane)
	}

	left := notesList.View()
	if len(notes) == 0 {
		left = "No journals yet.\nCreate your first entry."
	}
	right := components.FormatSelectedMeta(notesList.SelectedItem(), len(notes), dashboardNote, dashboardNoteErr)
	body := layout.TwoPaneWithRatioAndTitlesAndWidth("Saved Journals", "Journal Metadata", left, right, components.DashboardLeftRatio, layout.ContentWidth())
	return layout.CenterContent(body)
}
