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
	m.editor.File = ""
	m.editor.Created = time.Now()
	m.editor.Err = nil
	m.editor.Title.SetValue("")
	m.editor.Body.SetValue("")
	m.editor.Title.Focus()
	m.editor.Body.Blur()
	m.screen = screenEditor
	m.updateEditorSize()
}

func (m *AppModel) startEditorForSelected() {
	if m.config.StoragePath == "" {
		return
	}
	item := m.dashboard.List.SelectedItem()
	noteItem, ok := item.(components.NoteItem)
	if !ok {
		return
	}
	service := journal.NewService(m.config.StoragePath, journal.WithEncryption(m.config.Encrypt, m.config.SshKeyPath))
	note, err := service.LoadNote(noteItem.Info.Filename)
	if err != nil {
		m.editor.Err = err
		return
	}

	m.editor.File = noteItem.Info.Filename
	m.editor.Created = note.Created
	if m.editor.Created.IsZero() {
		if created, ok := components.ParseFilenameTimestamp(noteItem.Info.Filename); ok {
			m.editor.Created = created
		}
	}
	m.editor.Err = nil
	m.editor.Title.SetValue(note.Title)
	m.editor.Body.SetValue(strings.TrimSuffix(note.Content, "\n"))
	m.editor.Title.Focus()
	m.editor.Body.Blur()
	m.screen = screenEditor
	m.updateEditorSize()
}

func (m *AppModel) startEditorForViewer() {
	if m.viewer.Note == nil {
		m.startEditorForSelected()
		return
	}
	note := m.viewer.Note
	m.editor.File = note.Filename
	m.editor.Created = note.Created
	m.editor.Err = nil
	m.editor.Title.SetValue(note.Title)
	m.editor.Body.SetValue(strings.TrimSuffix(note.Content, "\n"))
	m.editor.Title.Focus()
	m.editor.Body.Blur()
	m.screen = screenEditor
	m.updateEditorSize()
}

func (m *AppModel) updateEditorSize() *AppModel {
	layout := m.layout()
	width := layout.PaneContentWidth(layout.EditorPaneWidth())
	if width < 0 {
		width = 0
	}

	m.editor.Title.Width = width
	bodyWidth := width - 4
	if bodyWidth < 0 {
		bodyWidth = 0
	}
	m.editor.Body.SetWidth(bodyWidth)
	_, _, bodyContentHeight := m.editorLayout(layout, m.editor.Title.View(), paneWidth)
	m.editor.Body.SetHeight(bodyContentHeight)
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
	if m.config.StoragePath == "" {
		m.editor.Err = fmt.Errorf("journal path is not set")
		return m, nil
	}

	title := strings.TrimSpace(m.editor.Title.Value())
	if title == "" {
		title = "Untitled"
	}
	body := m.editor.Body.Value()

	service := journal.NewService(m.config.StoragePath, journal.WithEncryption(m.config.Encrypt, m.config.SshKeyPath))
	if m.editor.File == "" {
		_, err := service.SaveNote(title, body, m.editor.Created)
		if err != nil {
			m.editor.Err = err
			return m, nil
		}
	} else {
		if err := service.UpdateNote(m.editor.File, title, body, m.editor.Created); err != nil {
			m.editor.Err = err
			return m, nil
		}
	}

	m.editor.Err = nil
	m.screen = screenDashboard
	m = m.resetDashboardNotes()
	return m, m.loadDashboardNotesCmd()
}
