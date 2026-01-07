package app

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/huh"
	"github.com/never00rei/a7/journal"
)

type ConfigState struct {
	StoragePath   string
	SshKeyPath    string
	SshPubKeyPath string
	Encrypt       bool
}

type WelcomeModel struct{}

type StorageModel struct {
	Form *huh.Form
}

type PrivacyModel struct {
	Form *huh.Form
}

type SetupModel struct{}

type DashboardModel struct {
	List             list.Model
	Notes            []journal.NoteInfo
	Err              error
	SelectedNote     *journal.Note
	SelectedErr      error
	SelectedFilename string
}

type ViewerModel struct {
	Viewport viewport.Model
	Title    string
	Note     *journal.Note
	Raw      string
}

type EditorModel struct {
	Title   textinput.Model
	Body    textarea.Model
	Created time.Time
	File    string
	Err     error
}
