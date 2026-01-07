package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/never00rei/a7/ui/components"
	"github.com/never00rei/a7/ui/layout"
	"github.com/never00rei/a7/ui/screens"
)

type ScreenModel interface {
	Init(app *AppModel) tea.Cmd
	Update(app *AppModel, msg tea.Msg) (tea.Cmd, bool)
	View(app *AppModel, layout layout.Layout) string
}

func (m *AppModel) activeScreenModel() ScreenModel {
	switch m.screen {
	case screenWelcome:
		return &m.welcome
	case screenWalkthroughStorage:
		return &m.storage
	case screenWalkthroughPrivacy:
		return &m.privacy
	case screenSetup:
		return &m.setup
	case screenDashboard:
		return &m.dashboard
	case screenViewer:
		return &m.viewer
	case screenEditor:
		return &m.editor
	default:
		return nil
	}
}

func (m *WelcomeModel) Init(app *AppModel) tea.Cmd {
	return nil
}

func (m *WelcomeModel) Update(app *AppModel, msg tea.Msg) (tea.Cmd, bool) {
	return nil, false
}

func (m *WelcomeModel) View(app *AppModel, layout layout.Layout) string {
	return screens.Welcome(layout)
}

func (m *StorageModel) Init(app *AppModel) tea.Cmd {
	if m.Form != nil {
		return m.Form.Init()
	}
	return nil
}

func (m *StorageModel) Update(app *AppModel, msg tea.Msg) (tea.Cmd, bool) {
	currentScreen := app.screen
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "shift+tab" {
		app.screen = prevScreen(app.screen)
		return app.initActiveFormCmd(), true
	}
	if m.Form == nil {
		return nil, false
	}
	model, cmd := m.Form.Update(msg)
	m.Form = model.(*huh.Form)
	if m.Form.State == huh.StateCompleted {
		app.config.StoragePath = m.Form.GetString(components.StoragePathKey)
		app.screen = nextScreen(app.screen)
	}
	if m.Form.State == huh.StateAborted {
		return tea.Quit, true
	}
	return app.batchFormCmd(cmd, currentScreen), true
}

func (m *StorageModel) View(app *AppModel, layout layout.Layout) string {
	return screens.WalkthroughStorage(layout, m.Form)
}

func (m *PrivacyModel) Init(app *AppModel) tea.Cmd {
	if m.Form != nil {
		return m.Form.Init()
	}
	return nil
}

func (m *PrivacyModel) Update(app *AppModel, msg tea.Msg) (tea.Cmd, bool) {
	currentScreen := app.screen
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "shift+tab" {
		app.screen = prevScreen(app.screen)
		return app.initActiveFormCmd(), true
	}
	if m.Form == nil {
		return nil, false
	}
	model, cmd := m.Form.Update(msg)
	m.Form = model.(*huh.Form)
	if m.Form.State == huh.StateCompleted {
		app.config.Encrypt = m.Form.GetBool(components.EncryptKey)
		app.config.SshKeyPath = m.Form.GetString(components.SshKeyPathKey)
		app.config.SshPubKeyPath = m.Form.GetString(components.SshPubKeyPathKey)
		if !app.config.Encrypt {
			app.config.SshKeyPath = ""
			app.config.SshPubKeyPath = ""
		}
		app.screen = nextScreen(app.screen)
	}
	if m.Form.State == huh.StateAborted {
		return tea.Quit, true
	}
	return app.batchFormCmd(cmd, currentScreen), true
}

func (m *PrivacyModel) View(app *AppModel, layout layout.Layout) string {
	return screens.WalkthroughPrivacy(layout, m.Form)
}

func (m *SetupModel) Init(app *AppModel) tea.Cmd {
	return nil
}

func (m *SetupModel) Update(app *AppModel, msg tea.Msg) (tea.Cmd, bool) {
	return nil, false
}

func (m *SetupModel) View(app *AppModel, layout layout.Layout) string {
	return screens.Setup(layout, app.config.StoragePath, app.config.SshKeyPath, app.config.Encrypt)
}

func (m *DashboardModel) Init(app *AppModel) tea.Cmd {
	return nil
}

func (m *DashboardModel) Update(app *AppModel, msg tea.Msg) (tea.Cmd, bool) {
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	app.dashboard.List = m.List
	app.updateDashboardSelection()
	return cmd, false
}

func (m *DashboardModel) View(app *AppModel, layout layout.Layout) string {
	return screens.Dashboard(layout, app.config.StoragePath, m.Err, m.Notes, m.List, m.SelectedNote, m.SelectedErr)
}

func (m *ViewerModel) Init(app *AppModel) tea.Cmd {
	return nil
}

func (m *ViewerModel) Update(app *AppModel, msg tea.Msg) (tea.Cmd, bool) {
	var cmd tea.Cmd
	m.Viewport, cmd = m.Viewport.Update(msg)
	app.viewer.Viewport = m.Viewport
	return cmd, false
}

func (m *ViewerModel) View(app *AppModel, layout layout.Layout) string {
	return screens.Viewer(layout, m.Title, m.Viewport.View())
}

func (m *EditorModel) Init(app *AppModel) tea.Cmd {
	return nil
}

func (m *EditorModel) Update(app *AppModel, msg tea.Msg) (tea.Cmd, bool) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.Title, cmd = m.Title.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	m.Body, cmd = m.Body.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	app.editor.Title = m.Title
	app.editor.Body = m.Body
	if len(cmds) > 0 {
		return tea.Batch(cmds...), false
	}
	return nil, false
}

func (m *EditorModel) View(app *AppModel, layout layout.Layout) string {
	paneWidth := layout.EditorPaneWidth()
	_, bodyPaneHeight, _ := app.editorLayout(layout, m.Title.View(), paneWidth)
	return screens.Editor(layout, m.Title.View(), m.Body.View(), m.Err, paneWidth, bodyPaneHeight)
}
