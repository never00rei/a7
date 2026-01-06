package app

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/never00rei/a7/config"
	"github.com/never00rei/a7/journal"
	"github.com/never00rei/a7/ui/components"
	"github.com/never00rei/a7/ui/layout"
	"github.com/never00rei/a7/ui/screens"
)

type screenID int

const (
	screenWelcome screenID = iota
	screenWalkthroughStorage
	screenWalkthroughPrivacy
	screenSetup
	screenDashboard
	screenViewer
	screenEditor
)

type AppModel struct {
	screen                screenID
	width                 int
	height                int
	storageForm           *huh.Form
	privacyForm           *huh.Form
	notesList             list.Model
	notes                 []journal.NoteInfo
	dashboardErr          error
	dashboardNote         *journal.Note
	dashboardNoteErr      error
	dashboardNoteFilename string
	viewer                viewport.Model
	viewerTitle           string
	viewerNote            *journal.Note
	viewerRaw             string
	editorTitle           textinput.Model
	editorBody            textarea.Model
	editorCreated         time.Time
	editorFile            string
	editorErr             error
	storagePath           string
	sshKeyPath            string
	sshPubKeyPath         string
	encrypt               bool
	lastError             error
}

func NewAppModel() AppModel {
	model := AppModel{
		screen:     screenWelcome,
		sshKeyPath: config.SshPath,
	}
	if conf, err := config.LoadConf(); err == nil && conf.JournalPath != "" {
		model.storagePath = conf.JournalPath
		model.sshKeyPath = conf.SshKeyFile
		model.sshPubKeyPath = conf.SshPubKey
		model.encrypt = conf.Encrypt
		model.screen = screenDashboard
	}
	model.storageForm = components.NewStorageForm(&model.storagePath, 0)
	model.privacyForm = components.NewPrivacyForm(&model.encrypt, &model.sshKeyPath, &model.sshPubKeyPath, 0)
	model.notesList = components.NewNotesList(nil, 0, 0)
	model.viewer = viewport.New(0, 0)
	model.editorTitle = textinput.New()
	model.editorTitle.Placeholder = "Journal title"
	model.editorBody = textarea.New()
	model.editorBody.Placeholder = "Start writing..."
	model.editorBody.CharLimit = 0
	if model.screen == screenDashboard {
		model = model.loadDashboardNotes()
	}
	return model
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.updateFormWidths()
		m = m.updateDashboardListSize()
		m = *m.updateViewerSize()
		m = *m.updateEditorSize()
	case error:
		m.lastError = msg
		return m, nil
	case screenID:
		m.screen = msg
		if m.screen == screenDashboard {
			m = m.loadDashboardNotes()
		}
		if m.screen == screenViewer {
			m = *m.updateViewerSize()
		}
		if m.screen == screenEditor {
			m = *m.updateEditorSize()
		}
		return m, m.initActiveFormCmd()
	}

	if updated, cmd, handled := m.updateActiveForm(msg); handled {
		return updated, cmd
	}

	var cmds []tea.Cmd
	if m.screen == screenDashboard {
		var cmd tea.Cmd
		m.notesList, cmd = m.notesList.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		m = m.updateDashboardSelection()
	}
	if m.screen == screenViewer {
		var cmd tea.Cmd
		m.viewer, cmd = m.viewer.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	if m.screen == screenEditor {
		var cmd tea.Cmd
		m.editorTitle, cmd = m.editorTitle.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		m.editorBody, cmd = m.editorBody.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.screen == screenWalkthroughPrivacy && msg.String() == "s" {
			m.encrypt = false
			m.sshKeyPath = ""
			m.sshPubKeyPath = ""
			m.screen = screenSetup
			return m, m.initActiveFormCmd()
		}
		if m.screen == screenDashboard {
			switch msg.String() {
			case "enter":
				return m.openViewer()
			case "n":
				m.startEditorForNew()
				return m, nil
			case "e":
				m.startEditorForSelected()
				return m, nil
			}
		}
		switch msg.String() {
		case "esc":
			if m.screen == screenViewer || m.screen == screenEditor {
				m.screen = screenDashboard
				return m, nil
			}
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.screen == screenSetup {
				cmd := tea.Sequence(m.saveConfigCmd(), func() tea.Msg {
					return screenDashboard
				})
				return m, cmd
			}
			m.screen = nextScreen(m.screen)
			return m, m.initActiveFormCmd()
		case "e":
			if m.screen == screenViewer {
				m.startEditorForViewer()
				return m, nil
			}
		case "tab":
			if m.screen == screenEditor {
				if m.editorTitle.Focused() {
					m.editorTitle.Blur()
					m.editorBody.Focus()
				} else {
					m.editorBody.Blur()
					m.editorTitle.Focus()
				}
				return m, nil
			}
		case "shift+tab":
			if m.screen == screenEditor {
				if m.editorBody.Focused() {
					m.editorBody.Blur()
					m.editorTitle.Focus()
				} else {
					m.editorTitle.Blur()
					m.editorBody.Focus()
				}
				return m, nil
			}
			if m.screen == screenViewer {
				m.screen = screenDashboard
				return m, nil
			}
			m.screen = prevScreen(m.screen)
			return m, m.initActiveFormCmd()
		case "ctrl+s":
			if m.screen == screenEditor {
				return m.saveEditorNote()
			}
		}
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

func (m AppModel) saveConfigCmd() tea.Cmd {
	journalPath := m.storagePath
	sshKeyPath := m.sshKeyPath
	sshPubKeyPath := m.sshPubKeyPath
	return func() tea.Msg {
		if !m.encrypt {
			sshKeyPath = ""
			sshPubKeyPath = ""
		}
		conf := config.NewConf(journalPath, sshKeyPath, sshPubKeyPath, m.encrypt)
		if err := conf.SaveConfig(); err != nil {
			return err
		}
		return nil
	}
}

func (m AppModel) updateActiveForm(msg tea.Msg) (AppModel, tea.Cmd, bool) {
	currentScreen := m.screen
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "shift+tab" {
		switch m.screen {
		case screenWalkthroughStorage, screenWalkthroughPrivacy:
			m.screen = prevScreen(m.screen)
			return m, m.initActiveFormCmd(), true
		}
	}
	switch m.screen {
	case screenWalkthroughStorage:
		if m.storageForm == nil {
			return m, nil, false
		}
		model, cmd := m.storageForm.Update(msg)
		m.storageForm = model.(*huh.Form)
		if m.storageForm.State == huh.StateCompleted {
			m.storagePath = m.storageForm.GetString(components.StoragePathKey)
			m.screen = nextScreen(m.screen)
		}
		if m.storageForm.State == huh.StateAborted {
			return m, tea.Quit, true
		}
		return m, m.batchFormCmd(cmd, currentScreen), true
	case screenWalkthroughPrivacy:
		if m.privacyForm == nil {
			return m, nil, false
		}
		model, cmd := m.privacyForm.Update(msg)
		m.privacyForm = model.(*huh.Form)
		if m.privacyForm.State == huh.StateCompleted {
			m.encrypt = m.privacyForm.GetBool(components.EncryptKey)
			m.sshKeyPath = m.privacyForm.GetString(components.SshKeyPathKey)
			m.sshPubKeyPath = m.privacyForm.GetString(components.SshPubKeyPathKey)
			if !m.encrypt {
				m.sshKeyPath = ""
				m.sshPubKeyPath = ""
			}
			m.screen = nextScreen(m.screen)
		}
		if m.privacyForm.State == huh.StateAborted {
			return m, tea.Quit, true
		}
		return m, m.batchFormCmd(cmd, currentScreen), true
	default:
		return m, nil, false
	}
}

func (m AppModel) updateFormWidths() AppModel {
	layout := m.layout()
	width := layout.FormWidth()
	if m.storageForm != nil {
		m.storageForm.WithWidth(width)
	}
	if m.privacyForm != nil {
		m.privacyForm.WithWidth(width)
	}
	return m
}

func (m AppModel) updateDashboardListSize() AppModel {
	layout := m.layout()
	leftWidth, _ := layout.SplitPaneContentWidths(components.DashboardLeftRatio)
	height := layout.PaneContentHeight(layout.BodyHeight())
	if height < 0 {
		height = 0
	}
	width := leftWidth
	m.notesList.SetSize(width, height)
	return m
}

func (m AppModel) layout() layout.Layout {
	return layout.New(m.width, m.height)
}

func (m AppModel) initActiveFormCmd() tea.Cmd {
	switch m.screen {
	case screenWalkthroughStorage:
		if m.storageForm != nil {
			return m.storageForm.Init()
		}
	case screenWalkthroughPrivacy:
		if m.privacyForm != nil {
			return m.privacyForm.Init()
		}
	}
	return nil
}

func (m AppModel) batchFormCmd(cmd tea.Cmd, previous screenID) tea.Cmd {
	if previous == m.screen {
		return cmd
	}
	nextCmd := m.initActiveFormCmd()
	if cmd == nil {
		return nextCmd
	}
	if nextCmd == nil {
		return cmd
	}
	return tea.Batch(cmd, nextCmd)
}

func (m AppModel) loadDashboardNotes() AppModel {
	m.dashboardErr = nil
	m.notes = nil
	m.dashboardNote = nil
	m.dashboardNoteErr = nil
	m.dashboardNoteFilename = ""

	if m.storagePath == "" {
		m.notesList.SetItems(nil)
		return m
	}

	service := journal.NewService(m.storagePath)
	notes, err := service.ListNotes()
	if err != nil {
		m.dashboardErr = err
		m.notesList.SetItems(nil)
		return m
	}

	m.notes = notes
	m.notesList.SetItems(components.BuildNoteItems(notes))
	if len(notes) > 0 {
		m.notesList.Select(0)
	}
	m = m.updateDashboardListSize()
	m = m.updateDashboardSelection()
	return m
}

func (m AppModel) updateDashboardSelection() AppModel {
	if m.storagePath == "" {
		m.dashboardNote = nil
		m.dashboardNoteErr = nil
		m.dashboardNoteFilename = ""
		return m
	}

	item := m.notesList.SelectedItem()
	noteItem, ok := item.(components.NoteItem)
	if !ok {
		m.dashboardNote = nil
		m.dashboardNoteErr = nil
		m.dashboardNoteFilename = ""
		return m
	}

	if noteItem.Info.Filename == m.dashboardNoteFilename && m.dashboardNoteErr == nil && m.dashboardNote != nil {
		return m
	}

	service := journal.NewService(m.storagePath, journal.WithEncryption(m.encrypt, m.sshKeyPath))
	note, err := service.LoadNote(noteItem.Info.Filename)
	m.dashboardNoteFilename = noteItem.Info.Filename
	m.dashboardNote = note
	m.dashboardNoteErr = err
	return m
}

func (m AppModel) View() string {
	layout := m.layout()
	switch m.screen {
	case screenWelcome:
		return layout.Frame(screens.Welcome(layout), m.helpText())
	case screenWalkthroughStorage:
		return layout.Frame(screens.WalkthroughStorage(layout, m.storageForm), m.helpText())
	case screenWalkthroughPrivacy:
		return layout.Frame(screens.WalkthroughPrivacy(layout, m.privacyForm), m.helpText())
	case screenSetup:
		return layout.Frame(screens.Setup(layout, m.storagePath, m.sshKeyPath, m.encrypt), m.helpText())
	case screenDashboard:
		return layout.Frame(screens.Dashboard(layout, m.storagePath, m.dashboardErr, m.notes, m.notesList, m.dashboardNote, m.dashboardNoteErr), m.helpText())
	case screenViewer:
		return layout.Frame(screens.Viewer(layout, m.viewerTitle, m.viewer.View()), m.helpText())
	case screenEditor:
		return layout.Frame(screens.Editor(layout, m.editorTitle.View(), m.editorBody.View(), m.editorErr), m.helpText())
	default:
		return layout.Frame("unknown screen", m.helpText())
	}
}

func (m AppModel) helpText() string {
	switch m.screen {
	case screenWelcome:
		return "enter: begin  ctrl+c: quit"
	case screenWalkthroughStorage:
		return "enter/tab: next  shift+tab: back  ctrl+c: quit"
	case screenWalkthroughPrivacy:
		return "enter/tab: next  shift+tab: back  s: skip  ctrl+c: quit"
	case screenDashboard:
		return "up/down: select  enter: view  n: new  e: edit  ctrl+c: quit"
	case screenViewer:
		return "esc: back  e: edit  ctrl+c: quit"
	case screenEditor:
		return "tab: switch  ctrl+s: save  esc: back  ctrl+c: quit"
	default:
		return "enter: continue  shift+tab: back  ctrl+c: quit"
	}
}

func nextScreen(current screenID) screenID {
	switch current {
	case screenWelcome:
		return screenWalkthroughStorage
	case screenWalkthroughStorage:
		return screenWalkthroughPrivacy
	case screenWalkthroughPrivacy:
		return screenSetup
	case screenSetup:
		return screenDashboard
	default:
		return current
	}
}

func prevScreen(current screenID) screenID {
	switch current {
	case screenWalkthroughStorage:
		return screenWelcome
	case screenWalkthroughPrivacy:
		return screenWalkthroughStorage
	case screenSetup:
		return screenWalkthroughPrivacy
	case screenDashboard:
		return screenSetup
	default:
		return current
	}
}
