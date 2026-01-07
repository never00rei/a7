package app

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/never00rei/a7/config"
	"github.com/never00rei/a7/journal"
	"github.com/never00rei/a7/ui/components"
	"github.com/never00rei/a7/ui/layout"
)

type screenID int

const (
	screenWelcome screenID = iota
	screenWalkthroughStorage
	screenWalkthroughPrivacy
	screenSetup
	screenSettings
	screenDashboard
	screenViewer
	screenEditor
)

type AppModel struct {
	screen    screenID
	width     int
	height    int
	config    ConfigState
	welcome   WelcomeModel
	storage   StorageModel
	privacy   PrivacyModel
	setup     SetupModel
	settings  SettingsModel
	dashboard DashboardModel
	viewer    ViewerModel
	editor    EditorModel
	lastError error
}

func NewAppModel() AppModel {
	model := AppModel{
		screen: screenWelcome,
		config: ConfigState{
			SshKeyPath: config.SshPath,
		},
	}
	if conf, err := config.LoadConf(); err == nil && conf.JournalPath != "" {
		model.config.StoragePath = conf.JournalPath
		model.config.SshKeyPath = conf.SshKeyFile
		model.config.SshPubKeyPath = conf.SshPubKey
		model.config.Encrypt = conf.Encrypt
		model.screen = screenDashboard
	}
	model.storage.Form = components.NewStorageForm(&model.config.StoragePath, 0)
	model.privacy.Form = components.NewPrivacyForm(&model.config.Encrypt, &model.config.SshKeyPath, &model.config.SshPubKeyPath, 0)
	model.settings.Form = components.NewSettingsForm(&model.config.StoragePath, &model.config.Encrypt, &model.config.SshKeyPath, &model.config.SshPubKeyPath, 0)
	model.dashboard.List = components.NewNotesList(nil, 0, 0)
	model.viewer.Viewport = viewport.New(0, 0)
	model.editor.Title = textinput.New()
	model.editor.Title.Placeholder = "Journal title"
	model.editor.Body = textarea.New()
	model.editor.Body.Placeholder = "Start writing..."
	model.editor.Body.CharLimit = 0
	return model
}

func (m AppModel) Init() tea.Cmd {
	if m.screen == screenDashboard {
		return m.loadDashboardNotesCmd()
	}
	if model := m.activeScreenModel(); model != nil {
		return model.Init(&m)
	}
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
	case dashboardNotesMsg:
		return m.applyDashboardNotes(msg), nil
	case configSavedMsg:
		return m, nil
	case errMsg:
		m.lastError = msg.err
		return m, nil
	case error:
		m.lastError = msg
		return m, nil
	case screenID:
		m.screen = msg
		if m.screen == screenDashboard {
			m = m.resetDashboardNotes()
			return m, m.loadDashboardNotesCmd()
		}
		if m.screen == screenViewer {
			m = *m.updateViewerSize()
		}
		if m.screen == screenEditor {
			m = *m.updateEditorSize()
		}
		return m, m.initActiveFormCmd()
	}

	var cmds []tea.Cmd
	if model := m.activeScreenModel(); model != nil {
		cmd, handled := model.Update(&m, msg)
		if handled {
			return m, cmd
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.screen == screenWalkthroughPrivacy && msg.String() == "s" {
			m.config.Encrypt = false
			m.config.SshKeyPath = ""
			m.config.SshPubKeyPath = ""
			m.screen = screenSetup
			return m, m.initActiveFormCmd()
		}
		if m.screen == screenDashboard {
			switch msg.String() {
			case "enter":
				return m.openViewer()
			case "s":
				m.screen = screenSettings
				return m, m.initActiveFormCmd()
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
			if m.screen == screenViewer || m.screen == screenEditor || m.screen == screenSettings {
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
				if m.editor.Title.Focused() {
					m.editor.Title.Blur()
					m.editor.Body.Focus()
				} else {
					m.editor.Body.Blur()
					m.editor.Title.Focus()
				}
				return m, nil
			}
		case "shift+tab":
			if m.screen == screenEditor {
				if m.editor.Body.Focused() {
					m.editor.Body.Blur()
					m.editor.Title.Focus()
				} else {
					m.editor.Title.Blur()
					m.editor.Body.Focus()
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
	journalPath := m.config.StoragePath
	sshKeyPath := m.config.SshKeyPath
	sshPubKeyPath := m.config.SshPubKeyPath
	return func() tea.Msg {
		if !m.config.Encrypt {
			sshKeyPath = ""
			sshPubKeyPath = ""
		}
		conf := config.NewConf(journalPath, sshKeyPath, sshPubKeyPath, m.config.Encrypt)
		if err := conf.SaveConfig(); err != nil {
			return errMsg{err: err}
		}
		return configSavedMsg{}
	}
}

func (m AppModel) updateFormWidths() AppModel {
	layout := m.layout()
	width := layout.FormWidth()
	if m.storage.Form != nil {
		m.storage.Form.WithWidth(width)
	}
	if m.privacy.Form != nil {
		m.privacy.Form.WithWidth(width)
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
	m.dashboard.List.SetSize(width, height)
	return m
}

func (m AppModel) layout() layout.Layout {
	return layout.New(m.width, m.height)
}

func (m AppModel) initActiveFormCmd() tea.Cmd {
	if model := m.activeScreenModel(); model != nil {
		return model.Init(&m)
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

func (m AppModel) resetDashboardNotes() AppModel {
	m.dashboard.Err = nil
	m.dashboard.Notes = nil
	m.dashboard.SelectedNote = nil
	m.dashboard.SelectedErr = nil
	m.dashboard.SelectedFilename = ""
	m.dashboard.List.SetItems(nil)
	m.dashboard.List.Title = ""
	return m
}

func (m AppModel) loadDashboardNotesCmd() tea.Cmd {
	path := m.config.StoragePath
	return func() tea.Msg {
		if path == "" {
			return dashboardNotesMsg{path: path}
		}
		service := journal.NewService(path)
		notes, err := service.ListNotes()
		return dashboardNotesMsg{path: path, notes: notes, err: err}
	}
}

func (m AppModel) applyDashboardNotes(msg dashboardNotesMsg) AppModel {
	if msg.path != m.config.StoragePath {
		return m
	}
	m.dashboard.Err = msg.err
	if msg.err != nil {
		m.dashboard.Notes = nil
		m.dashboard.List.SetItems(nil)
		m.dashboard.List.Title = ""
		return m
	}

	m.dashboard.Notes = msg.notes
	m.dashboard.List.SetItems(components.BuildNoteItems(msg.notes))
	m.dashboard.List.Title = m.config.StoragePath
	if len(msg.notes) > 0 {
		m.dashboard.List.Select(0)
	}
	m = m.updateDashboardListSize()
	m = m.updateDashboardSelection()
	return m
}

func (m AppModel) updateDashboardSelection() AppModel {
	if m.config.StoragePath == "" {
		m.dashboard.SelectedNote = nil
		m.dashboard.SelectedErr = nil
		m.dashboard.SelectedFilename = ""
		return m
	}

	item := m.dashboard.List.SelectedItem()
	noteItem, ok := item.(components.NoteItem)
	if !ok {
		m.dashboard.SelectedNote = nil
		m.dashboard.SelectedErr = nil
		m.dashboard.SelectedFilename = ""
		return m
	}

	if noteItem.Info.Filename == m.dashboard.SelectedFilename && m.dashboard.SelectedErr == nil && m.dashboard.SelectedNote != nil {
		return m
	}

	service := journal.NewService(m.config.StoragePath, journal.WithEncryption(m.config.Encrypt, m.config.SshKeyPath))
	note, err := service.LoadNote(noteItem.Info.Filename)
	m.dashboard.SelectedFilename = noteItem.Info.Filename
	m.dashboard.SelectedNote = note
	m.dashboard.SelectedErr = err
	return m
}

func (m AppModel) View() string {
	layout := m.layout()
	model := m.activeScreenModel()
	if model == nil {
		return layout.Frame("unknown screen", m.helpText())
	}
	return layout.Frame(model.View(&m, layout), m.helpText())
}

func (m AppModel) helpText() string {
	switch m.screen {
	case screenWelcome:
		return "enter begin • ctrl+c quit"
	case screenWalkthroughStorage:
		return "⏎/enter/tab next • shift+tab back • ctrl+c quit"
	case screenWalkthroughPrivacy:
		return "⏎/enter/tab next • shift+tab back • s skip • ctrl+c quit"
	case screenDashboard:
		return "↑/k up • ↓/j down • / filter • ⏎/enter view • n new • e edit • s settings • ctrl+c quit"
	case screenViewer:
		return "esc back • e edit • ctrl+c quit"
	case screenEditor:
		return "tab switch • ctrl+s save • esc back • ctrl+c quit"
	case screenSettings:
		return "tab next • shift+tab back • esc back • ctrl+c quit"
	default:
		return "⏎/enter continue • shift+tab back • ctrl+c quit"
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
