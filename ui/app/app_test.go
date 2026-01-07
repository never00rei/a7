package app

import (
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/never00rei/a7/config"
	"github.com/never00rei/a7/journal"
	"github.com/never00rei/a7/ui/components"
)

func setupTestConfig(t *testing.T) {
	t.Helper()
	origHome := config.Home
	origXdg := config.XdgConfigHome
	origSsh := config.SshPath
	temp := t.TempDir()
	config.Home = temp
	config.XdgConfigHome = temp
	config.SshPath = filepath.Join(temp, ".ssh")
	t.Cleanup(func() {
		config.Home = origHome
		config.XdgConfigHome = origXdg
		config.SshPath = origSsh
	})
}

func applyCmd(model AppModel, cmd tea.Cmd) AppModel {
	for i := 0; i < 4 && cmd != nil; i++ {
		msg := cmd()
		if msg == nil {
			break
		}
		updated, nextCmd := model.Update(msg)
		model = updated.(AppModel)
		cmd = nextCmd
	}
	return model
}

func createTestJournal(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	svc := journal.NewService(root)
	filename, err := svc.SaveNote("Test Journal", "hello world", time.Now())
	if err != nil {
		t.Fatalf("SaveNote: %v", err)
	}
	return root, filename
}

func TestWelcomeEnterAdvances(t *testing.T) {
	setupTestConfig(t)
	model := NewAppModel()
	if model.screen != screenWelcome {
		t.Fatalf("start screen = %v, want %v", model.screen, screenWelcome)
	}
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	next := updated.(AppModel)
	if next.screen != screenWalkthroughStorage {
		t.Fatalf("after enter screen = %v, want %v", next.screen, screenWalkthroughStorage)
	}
}

func TestPrivacyShiftTabBack(t *testing.T) {
	setupTestConfig(t)
	model := NewAppModel()
	model.screen = screenWalkthroughPrivacy
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	next := updated.(AppModel)
	if next.screen != screenWalkthroughStorage {
		t.Fatalf("after shift+tab screen = %v, want %v", next.screen, screenWalkthroughStorage)
	}
}

func TestEscReturnsToDashboard(t *testing.T) {
	setupTestConfig(t)
	model := NewAppModel()
	model.screen = screenViewer
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	next := updated.(AppModel)
	if next.screen != screenDashboard {
		t.Fatalf("viewer esc screen = %v, want %v", next.screen, screenDashboard)
	}

	model.screen = screenEditor
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	next = updated.(AppModel)
	if next.screen != screenDashboard {
		t.Fatalf("editor esc screen = %v, want %v", next.screen, screenDashboard)
	}
}

func TestEditorTabSwitchesFocus(t *testing.T) {
	setupTestConfig(t)
	model := NewAppModel()
	model.startEditorForNew()
	if !model.editor.Title.Focused() {
		t.Fatalf("title should be focused on start")
	}
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	next := updated.(AppModel)
	if !next.editor.Body.Focused() {
		t.Fatalf("body should be focused after tab")
	}
	updated, _ = next.Update(tea.KeyMsg{Type: tea.KeyTab})
	next = updated.(AppModel)
	if !next.editor.Title.Focused() {
		t.Fatalf("title should be focused after second tab")
	}
}

func TestPrivacySkipClearsKeys(t *testing.T) {
	setupTestConfig(t)
	model := NewAppModel()
	model.screen = screenWalkthroughPrivacy
	model.config.Encrypt = true
	model.config.SshKeyPath = "/tmp/key"
	model.config.SshPubKeyPath = "/tmp/key.pub"

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	next := updated.(AppModel)
	if next.screen != screenSetup {
		t.Fatalf("privacy skip screen = %v, want %v", next.screen, screenSetup)
	}
	if next.config.Encrypt {
		t.Fatalf("encrypt = true, want false")
	}
	if next.config.SshKeyPath != "" || next.config.SshPubKeyPath != "" {
		t.Fatalf("ssh keys not cleared")
	}
}

func TestSetupEnterSavesAndAdvances(t *testing.T) {
	setupTestConfig(t)
	model := NewAppModel()
	model.screen = screenSetup
	model.config.StoragePath = t.TempDir()

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	next := updated.(AppModel)
	if cmd == nil {
		t.Fatalf("expected save config cmd")
	}
	next = applyCmd(next, cmd)
	updated, _ = next.Update(screenDashboard)
	next = updated.(AppModel)
	if next.screen != screenDashboard {
		t.Fatalf("after setup enter screen = %v, want %v", next.screen, screenDashboard)
	}
}

func TestDashboardSelectionLoadsMetadata(t *testing.T) {
	setupTestConfig(t)
	root, _ := createTestJournal(t)
	model := NewAppModel()
	model.config.StoragePath = root
	model.dashboard.List = components.NewNotesList(nil, 0, 0)

	svc := journal.NewService(root)
	notes, err := svc.ListNotes()
	if err != nil {
		t.Fatalf("ListNotes: %v", err)
	}
	model = model.applyDashboardNotes(dashboardNotesMsg{path: root, notes: notes})
	if model.dashboard.SelectedNote == nil {
		t.Fatalf("selected note is nil")
	}
	if model.dashboard.SelectedNote.Title == "" {
		t.Fatalf("selected note title missing")
	}
}

func TestDashboardEnterOpensViewer(t *testing.T) {
	setupTestConfig(t)
	root, _ := createTestJournal(t)
	model := NewAppModel()
	model.screen = screenDashboard
	model.config.StoragePath = root
	model.dashboard.List = components.NewNotesList(nil, 0, 0)

	svc := journal.NewService(root)
	notes, err := svc.ListNotes()
	if err != nil {
		t.Fatalf("ListNotes: %v", err)
	}
	model = model.applyDashboardNotes(dashboardNotesMsg{path: root, notes: notes})

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	next := updated.(AppModel)
	if next.screen != screenViewer {
		t.Fatalf("after enter screen = %v, want %v", next.screen, screenViewer)
	}
	if next.viewer.Title == "" {
		t.Fatalf("viewer title missing")
	}
	if next.viewer.Note == nil {
		t.Fatalf("viewer note missing")
	}
}
