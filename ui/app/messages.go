package app

import "github.com/never00rei/a7/journal"

type errMsg struct {
	err error
}

type configSavedMsg struct{}

type dashboardNotesMsg struct {
	path  string
	notes []journal.NoteInfo
	err   error
}
