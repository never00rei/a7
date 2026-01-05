package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/never00rei/a7/journal"
)

const dashboardLeftRatio = 0.6

var metadataLabelStyle = lipgloss.NewStyle().Bold(true)

type noteItem struct {
	info journal.NoteInfo
}

func (n noteItem) Title() string {
	return n.info.Filename
}

func (n noteItem) Description() string {
	return n.info.ModTime.Format(time.RFC822)
}

func (n noteItem) FilterValue() string {
	return n.info.Filename
}

func newNotesList(items []list.Item, width, height int) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), width, height)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Title = ""
	return l
}

func buildNoteItems(notes []journal.NoteInfo) []list.Item {
	items := make([]list.Item, 0, len(notes))
	for _, note := range notes {
		items = append(items, noteItem{info: note})
	}
	return items
}

func (m AppModel) dashboardListSize() (int, int) {
	leftWidth, _ := m.splitPaneContentWidths(dashboardLeftRatio)
	height := m.paneContentHeight(m.bodyHeight())
	if height < 0 {
		height = 0
	}
	return leftWidth, height
}

func parseFilenameTimestamp(filename string) (time.Time, bool) {
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

func formatSelectedMeta(item list.Item, total int, note *journal.Note, loadErr error) string {
	if item == nil {
		return fmt.Sprintf("%s: %d\n\nSelect a journal to see details.", boldLabel("Total journals"), total)
	}

	noteItem, ok := item.(noteItem)
	if !ok {
		return fmt.Sprintf("%s: %d", boldLabel("Total journals"), total)
	}

	lines := []string{
		boldLabel("Selected journal"),
		noteItem.info.Filename,
		"",
		boldLabel("Last modified"),
		noteItem.info.ModTime.Format(time.RFC822),
	}

	if created, ok := parseFilenameTimestamp(noteItem.info.Filename); ok {
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
