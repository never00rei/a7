package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/never00rei/a7/config"
)

type screenID int

const (
	screenWelcome screenID = iota
	screenWalkthroughStorage
	screenWalkthroughPrivacy
	screenSetup
	screenDashboard
)

type AppModel struct {
	screen      screenID
	width       int
	height      int
	storageForm *huh.Form
	privacyForm *huh.Form
	storagePath string
	sshKeyPath  string
	lastError   error
}

func NewAppModel() AppModel {
	model := AppModel{
		screen:     screenWelcome,
		sshKeyPath: config.SshPath,
	}
	model.storageForm = newStorageForm(&model.storagePath, 0)
	model.privacyForm = newPrivacyForm(&model.sshKeyPath, 0)
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
	case error:
		m.lastError = msg
		return m, nil
	case screenID:
		m.screen = msg
		return m, m.initActiveFormCmd()
	}

	if updated, cmd, handled := m.updateActiveForm(msg); handled {
		return updated, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
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
		case "shift+tab":
			m.screen = prevScreen(m.screen)
			return m, m.initActiveFormCmd()
		}
	}

	return m, nil
}

func (m AppModel) saveConfigCmd() tea.Cmd {
	journalPath := m.storagePath
	sshKeyPath := m.sshKeyPath
	return func() tea.Msg {
		conf := config.NewConf(journalPath, sshKeyPath)
		if err := conf.SaveConfig(); err != nil {
			return err
		}
		return nil
	}
}

func (m AppModel) updateActiveForm(msg tea.Msg) (AppModel, tea.Cmd, bool) {
	currentScreen := m.screen
	switch m.screen {
	case screenWalkthroughStorage:
		if m.storageForm == nil {
			return m, nil, false
		}
		model, cmd := m.storageForm.Update(msg)
		m.storageForm = model.(*huh.Form)
		if m.storageForm.State == huh.StateCompleted {
			m.storagePath = m.storageForm.GetString(storagePathKey)
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
			m.sshKeyPath = m.privacyForm.GetString(sshKeyPathKey)
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
	width := m.formWidth()
	if m.storageForm != nil {
		m.storageForm.WithWidth(width)
	}
	if m.privacyForm != nil {
		m.privacyForm.WithWidth(width)
	}
	return m
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

func (m AppModel) View() string {
	switch m.screen {
	case screenWelcome:
		return m.frame(m.viewWelcome(), m.helpText())
	case screenWalkthroughStorage:
		return m.frame(m.viewWalkthroughStorage(), m.helpText())
	case screenWalkthroughPrivacy:
		return m.frame(m.viewWalkthroughPrivacy(), m.helpText())
	case screenSetup:
		return m.frame(m.viewSetup(), m.helpText())
	case screenDashboard:
		return m.frame(m.viewDashboard(), m.helpText())
	default:
		return m.frame("unknown screen", m.helpText())
	}
}

func (m AppModel) helpText() string {
	switch m.screen {
	case screenWelcome:
		return "enter: begin  ctrl+c: quit"
	case screenWalkthroughStorage, screenWalkthroughPrivacy:
		return "enter/tab: next  shift+tab: back  ctrl+c: quit"
	case screenDashboard:
		return "ctrl+c: quit"
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
