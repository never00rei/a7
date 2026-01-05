package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/never00rei/a7/journal"
)

func (m AppModel) viewDashboard() string {
	if m.storagePath == "" {
		bodyText := "Set a journal folder to see recent notes.\n" +
			"Run setup to choose a storage location."
		pane := m.titledPaneWithWidth("Dashboard", bodyText, m.primaryPaneWidth())
		return m.centerContent(pane)
	}

	service := journal.NewService(m.storagePath)
	notes, err := service.ListNotes()
	if err != nil {
		bodyText := "Unable to load notes right now.\n" +
			"Check your journal folder and try again."
		pane := m.titledPaneWithWidth("Dashboard", bodyText, m.primaryPaneWidth())
		return m.centerContent(pane)
	}

	left := buildNotesList(notes, 8)
	right := buildNotesMeta(notes)
	body := m.twoPaneWithRatioAndTitles("Your Journals", "Journal Metadata", left, right, 0.6)
	return m.centerContent(body)
}

func buildNotesList(notes []journal.NoteInfo, limit int) string {
	if len(notes) == 0 {
		return "No notes yet.\nCreate your first entry."
	}
	if limit <= 0 || limit > len(notes) {
		limit = len(notes)
	}

	lines := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		lines = append(lines, fmt.Sprintf("%d. %s", i+1, notes[i].Filename))
	}

	return strings.Join(lines, "\n")
}

func buildNotesMeta(notes []journal.NoteInfo) string {
	if len(notes) == 0 {
		return "No metadata yet.\nNotes will appear here."
	}

	latest := notes[0]
	meta := []string{
		"Latest note",
		latest.Filename,
		"",
		"Last modified",
		latest.ModTime.Format(time.RFC822),
		"",
		fmt.Sprintf("Total notes: %d", len(notes)),
	}

	return strings.Join(meta, "\n")
}
