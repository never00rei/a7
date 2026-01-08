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

// Title() or Description() are used by bubbles/list "NewDefaultDelegate"
// to display filename and timestamp in the list viewport.
func (n NoteItem) Title() string {
	return n.Info.Filename
}

func (n NoteItem) Description() string {
	if !n.Info.Created.IsZero() {
		return n.Info.Created.Local().Format(time.RFC822)
	}
	return "Could not determine creation date"
}

func (n NoteItem) FilterValue() string {
	if strings.TrimSpace(n.Info.Title) != "" {
		return fmt.Sprintf("%s %s", n.Info.Title, n.Info.Filename)
	}
	return n.Info.Filename
}

func NewNotesList(items []list.Item, width, height int) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), width, height)
	l.SetShowHelp(false)
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
		boldLabel("Title"),
		noteTitle(noteItem),
		"",
		boldLabel("Last modified"),
		noteItem.Info.ModTime.Local().Format(time.RFC822),
	}

	created := time.Time{}
	if note != nil && !note.Created.IsZero() {
		created = note.Created
	} else if !noteItem.Info.Created.IsZero() {
		created = noteItem.Info.Created
	} else if parsed, ok := ParseFilenameTimestamp(noteItem.Info.Filename); ok {
		created = parsed
	}
	if !created.IsZero() {
		lines = append(lines, "", boldLabel("Created"), created.Local().Format(time.RFC822))
	} else {
		lines = append(lines, "", boldLabel("Created"), "Could not determine creation date")
	}

	encrypted := noteItem.Info.Encrypted
	if note != nil {
		encrypted = note.Encrypted
	}
	encryptedLabel := "No"
	if encrypted {
		encryptedLabel = "Yes"
	}
	lines = append(lines, "", boldLabel("Encrypted"), encryptedLabel)

	wordCount := noteItem.Info.WordCount
	if note != nil && note.WordCount >= 0 {
		wordCount = note.WordCount
	}
	if wordCount >= 0 {
		lines = append(lines, "", boldLabel("Word count"), fmt.Sprintf("%d", wordCount))
	} else if loadErr != nil {
		lines = append(lines, "", boldLabel("Word count"), "Unavailable")
	} else {
		lines = append(lines, "", boldLabel("Word count"), "Unavailable")
	}

	lines = append(lines, "", fmt.Sprintf("%s: %d", boldLabel("Total journals"), total))

	return strings.Join(lines, "\n")
}

func boldLabel(label string) string {
	return metadataLabelStyle.Render(label)
}

func noteTitle(item NoteItem) string {
	if strings.TrimSpace(item.Info.Title) != "" {
		return item.Info.Title
	}
	return item.Info.Filename
}
