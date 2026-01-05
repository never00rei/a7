package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/never00rei/a7/journal"
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
	noteItem, ok := item.(noteItem)
	if !ok {
		return
	}
	service := journal.NewService(m.storagePath, journal.WithEncryption(m.encrypt, m.sshKeyPath))
	note, err := service.LoadNote(noteItem.info.Filename)
	if err != nil {
		m.editorErr = err
		return
	}

	m.editorFile = noteItem.info.Filename
	m.editorCreated = note.Created
	if m.editorCreated.IsZero() {
		if created, ok := parseFilenameTimestamp(noteItem.info.Filename); ok {
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
	width := m.paneContentWidth(m.editorPaneWidth())
	if width < 0 {
		width = 0
	}

	m.editorTitle.Width = width
	m.editorBody.SetWidth(width - 4)

	_, bodyPaneHeight := m.editorPaneHeights()
	contentHeight := m.paneContentHeight(bodyPaneHeight)
	if contentHeight < 3 {
		contentHeight = 3
	}
	m.editorBody.SetHeight(contentHeight)
	return m
}

func (m AppModel) editorPaneHeights() (int, int) {
	titlePane := m.titledPaneWithWidthAndHeight("Title", m.editorTitle.View(), m.editorPaneWidth(), 0)
	titleHeight := lipgloss.Height(titlePane)
	totalHeight := m.bodyHeight()
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
