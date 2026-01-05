package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/never00rei/a7/journal"
)

func (m AppModel) openViewer() (AppModel, tea.Cmd) {
	if m.storagePath == "" {
		return m, nil
	}

	item := m.notesList.SelectedItem()
	noteItem, ok := item.(noteItem)
	if !ok {
		return m, nil
	}

	service := journal.NewService(m.storagePath, journal.WithEncryption(m.encrypt, m.sshKeyPath))
	note, err := service.LoadNote(noteItem.info.Filename)
	if err != nil {
		m.viewerTitle = "Unable to load journal"
		m.viewer.SetContent(fmt.Sprintf("Error: %v", err))
		m.viewer.YOffset = 0
		m.viewerNote = nil
		m.screen = screenViewer
		m.updateViewerSize()
		return m, nil
	}

	title := note.Title
	if title == "" {
		title = noteItem.info.Filename
	}

	m.viewerTitle = title
	m.viewerRaw = note.Content
	m.viewer.YOffset = 0
	m.viewerNote = note
	m.screen = screenViewer
	m.updateViewerSize()
	return m, nil
}

func (m *AppModel) updateViewerSize() *AppModel {
	width := m.paneContentWidth(m.contentWidth() - 2)
	height := m.paneContentHeight(m.bodyHeight())
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	if width == 0 {
		width = 1
	}
	if height == 0 {
		height = 1
	}
	m.viewer.Width = width
	m.viewer.Height = height
	m.renderViewerContent()
	return m
}

func (m *AppModel) renderViewerContent() {
	if m.viewerRaw == "" {
		m.viewer.SetContent("This journal is empty.")
		return
	}
	rendered, err := renderMarkdown(m.viewer.Width, m.viewerRaw)
	if err != nil || strings.TrimSpace(rendered) == "" {
		m.viewer.SetContent(m.viewerRaw)
		return
	}
	m.viewer.SetContent(rendered)
}

func renderMarkdown(width int, content string) (string, error) {
	if width <= 0 {
		width = 80
	}
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", err
	}
	defer renderer.Close()
	return renderer.Render(content)
}
