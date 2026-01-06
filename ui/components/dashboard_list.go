package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/never00rei/a7/journal"
)

const DashboardLeftRatio = 0.6

var metadataLabelStyle = lipgloss.NewStyle().Bold(true)

type NoteItem struct {
	Info journal.NoteInfo
}

func (n NoteItem) Title() string {
	return n.Info.Filename
}

func (n NoteItem) Description() string {
	return n.Info.ModTime.Format(time.RFC822)
}

func (n NoteItem) FilterValue() string {
	return n.Info.Filename
}

func NewNotesList(items []list.Item, width, height int) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), width, height)
	l.SetShowHelp(true)
	l.SetShowPagination(true)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.DisableQuitKeybindings()
	l.Title = ""
	return l
}

func BuildNoteItems(notes []journal.NoteInfo) []list.Item {
	items := make([]list.Item, 0, len(notes))
	for _, note := range notes {
		items = append(items, NoteItem{Info: note})
	}
	return items
}

func ParseFilenameTimestamp(filename string) (time.Time, bool) {
	if len(filename) < len("2006-01-02_15-04") {
		return time.Time{}, false
	}
	ts := filename[:len("2006-01-02_15-04")]
	parsed, err := time.Parse("2006-01-02_15-04", ts)
	if err != nil {
		return time.Time{}, false
	}
	return parsed, true
}

func FormatSelectedMeta(item list.Item, total int, note *journal.Note, loadErr error) string {
	if item == nil {
		return fmt.Sprintf("%s: %d\n\nSelect a journal to see details.", boldLabel("Total journals"), total)
	}

	noteItem, ok := item.(NoteItem)
	if !ok {
		return fmt.Sprintf("%s: %d", boldLabel("Total journals"), total)
	}

	lines := []string{
		boldLabel("Selected journal"),
		noteItem.Info.Filename,
		"",
		boldLabel("Last modified"),
		noteItem.Info.ModTime.Format(time.RFC822),
	}

	if created, ok := ParseFilenameTimestamp(noteItem.Info.Filename); ok {
		lines = append(lines, "", boldLabel("Created"), created.Format(time.RFC822))
	}

	encryptedLabel := "No"
	if note != nil && note.Encrypted {
		encryptedLabel = "Yes"
	}
	lines = append(lines, "", boldLabel("Encrypted"), encryptedLabel)

	if loadErr != nil {
		lines = append(lines, "", boldLabel("Word count"), "Unavailable")
	} else if note != nil && note.WordCount >= 0 {
		lines = append(lines, "", boldLabel("Word count"), fmt.Sprintf("%d", note.WordCount))
	} else if note != nil {
		lines = append(lines, "", boldLabel("Word count"), "Unavailable")
	}

	lines = append(lines, "", fmt.Sprintf("%s: %d", boldLabel("Total journals"), total))

	return strings.Join(lines, "\n")
}

func boldLabel(label string) string {
	return metadataLabelStyle.Render(label)
}
