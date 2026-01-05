package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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
	screen screenID
	width  int
	height int
}

func NewAppModel() AppModel {
	return AppModel{screen: screenWelcome}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.screen = nextScreen(m.screen)
			return m, nil
		case "s":
			if m.screen == screenWalkthroughPrivacy {
				m.screen = screenSetup
				return m, nil
			}
		case "b":
			m.screen = prevScreen(m.screen)
			return m, nil
		}
	}

	return m, nil
}

func (m AppModel) View() string {
	content := fmt.Sprintf("size: %dx%d", m.width, m.height)
	return m.frame(content)
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
