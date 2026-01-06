package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/never00rei/a7/journal"
	"github.com/never00rei/a7/ui/components"
)

func (m AppModel) openViewer() (AppModel, tea.Cmd) {
	if m.config.StoragePath == "" {
		return m, nil
	}

	item := m.dashboard.List.SelectedItem()
	noteItem, ok := item.(components.NoteItem)
	if !ok {
		return m, nil
	}

	service := journal.NewService(m.config.StoragePath, journal.WithEncryption(m.config.Encrypt, m.config.SshKeyPath))
	note, err := service.LoadNote(noteItem.Info.Filename)
	if err != nil {
		m.viewer.Title = "Unable to load journal"
		m.viewer.Viewport.SetContent(fmt.Sprintf("Error: %v", err))
		m.viewer.Viewport.YOffset = 0
		m.viewer.Note = nil
		m.screen = screenViewer
		m.updateViewerSize()
		return m, nil
	}

	title := note.Title
	if title == "" {
		title = noteItem.Info.Filename
	}

	m.viewer.Title = title
	m.viewer.Raw = note.Content
	m.viewer.Viewport.YOffset = 0
	m.viewer.Note = note
	m.screen = screenViewer
	m.updateViewerSize()
	return m, nil
}

func (m *AppModel) updateViewerSize() *AppModel {
	layout := m.layout()
	width := layout.PaneContentWidth(layout.ContentWidth() - 2)
	height := layout.PaneContentHeight(layout.BodyHeight())
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
	m.viewer.Viewport.Width = width
	m.viewer.Viewport.Height = height
	m.renderViewerContent()
	return m
}

func (m *AppModel) renderViewerContent() {
	if m.viewer.Raw == "" {
		m.viewer.Viewport.SetContent("This journal is empty.")
		return
	}
	rendered, err := renderMarkdown(m.viewer.Viewport.Width, m.viewer.Raw)
	if err != nil || strings.TrimSpace(rendered) == "" {
		m.viewer.Viewport.SetContent(m.viewer.Raw)
		return
	}
	m.viewer.Viewport.SetContent(rendered)
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
