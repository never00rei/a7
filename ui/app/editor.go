package app

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/never00rei/a7/journal"
	"github.com/never00rei/a7/ui/components"
	"github.com/never00rei/a7/ui/layout"
)

func (m *AppModel) startEditorForNew() {
	m.editorFile = ""
	m.editorCreated = time.Now()
	m.editorErr = nil
	m.editorTitle.SetValue("")
	m.editorBody.SetValue("")
	m.editorTitle.Focus()
	m.editorBody.Blur()
	m.screen = screenEditor
	m.updateEditorSize()
}

func (m *AppModel) startEditorForSelected() {
	if m.storagePath == "" {
		return
	}
	item := m.notesList.SelectedItem()
	noteItem, ok := item.(components.NoteItem)
	if !ok {
		return
	}
	service := journal.NewService(m.storagePath, journal.WithEncryption(m.encrypt, m.sshKeyPath))
	note, err := service.LoadNote(noteItem.Info.Filename)
	if err != nil {
		m.editorErr = err
		return
	}

	m.editorFile = noteItem.Info.Filename
	m.editorCreated = note.Created
	if m.editorCreated.IsZero() {
		if created, ok := components.ParseFilenameTimestamp(noteItem.Info.Filename); ok {
			m.editorCreated = created
		}
	}
	m.editorErr = nil
	m.editorTitle.SetValue(note.Title)
	m.editorBody.SetValue(note.Content)
	m.editorTitle.Focus()
	m.editorBody.Blur()
	m.screen = screenEditor
	m.updateEditorSize()
}

func (m *AppModel) startEditorForViewer() {
	if m.viewerNote == nil {
		m.startEditorForSelected()
		return
	}
	note := m.viewerNote
	m.editorFile = note.Filename
	m.editorCreated = note.Created
	m.editorErr = nil
	m.editorTitle.SetValue(note.Title)
	m.editorBody.SetValue(note.Content)
	m.editorTitle.Focus()
	m.editorBody.Blur()
	m.screen = screenEditor
	m.updateEditorSize()
}

func (m *AppModel) updateEditorSize() *AppModel {
	layout := m.layout()
	width := layout.PaneContentWidth(layout.EditorPaneWidth())
	if width < 0 {
		width = 0
	}

	m.editorTitle.Width = width
	m.editorBody.SetWidth(width - 4)

	_, bodyPaneHeight := m.editorPaneHeights(layout)
	contentHeight := layout.PaneContentHeight(bodyPaneHeight)
	if contentHeight < 3 {
		contentHeight = 3
	}
	m.editorBody.SetHeight(contentHeight)
	return m
}

func (m AppModel) editorPaneHeights(layout layout.Layout) (int, int) {
	titlePane := layout.TitledPaneWithWidthAndHeight("Title", m.editorTitle.View(), layout.EditorPaneWidth(), 0)
	titleHeight := lipgloss.Height(titlePane)
	totalHeight := layout.BodyHeight()
	bodyPaneHeight := totalHeight - titleHeight
	if bodyPaneHeight < 3 {
		bodyPaneHeight = 3
	}
	return titleHeight, bodyPaneHeight
}

func (m AppModel) saveEditorNote() (AppModel, tea.Cmd) {
	if m.storagePath == "" {
		m.editorErr = fmt.Errorf("journal path is not set")
		return m, nil
	}

	title := strings.TrimSpace(m.editorTitle.Value())
	if title == "" {
		title = "Untitled"
	}
	body := m.editorBody.Value()

	service := journal.NewService(m.storagePath, journal.WithEncryption(m.encrypt, m.sshKeyPath))
	if m.editorFile == "" {
		_, err := service.SaveNote(title, body, m.editorCreated)
		if err != nil {
			m.editorErr = err
			return m, nil
		}
	} else {
		if err := service.UpdateNote(m.editorFile, title, body, m.editorCreated); err != nil {
			m.editorErr = err
			return m, nil
		}
	}

	m.editorErr = nil
	m = m.loadDashboardNotes()
	m.screen = screenDashboard
	return m, nil
}
