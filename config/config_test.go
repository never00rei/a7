package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/ini.v1"
)

func TestSaveConfigExistingDirNoFile(t *testing.T) {
	tempDir := t.TempDir()

	origHome := Home
	origXdg := XdgConfigHome
	t.Cleanup(func() {
		Home = origHome
		XdgConfigHome = origXdg
	})

	Home = tempDir
	XdgConfigHome = ""

	confPath, err := BuildConfPath(Home, XdgConfigHome)
	if err != nil {
		t.Fatalf("BuildConfPath error: %v", err)
	}
	if err := os.MkdirAll(confPath, 0755); err != nil {
		t.Fatalf("MkdirAll error: %v", err)
	}

	conf := NewConf(filepath.Join(tempDir, "journal"), filepath.Join(tempDir, "id_ed25519"))
	if err := conf.SaveConfig(); err != nil {
		t.Fatalf("SaveConfig error: %v", err)
	}

	confFilePath := filepath.Join(confPath, ConfFileName)
	if _, err := os.Stat(confFilePath); err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	loaded, err := ini.Load(confFilePath)
	if err != nil {
		t.Fatalf("ini.Load error: %v", err)
	}

	section, err := loaded.GetSection("Settings")
	if err != nil {
		t.Fatalf("missing Settings section: %v", err)
	}

	if got := section.Key("journal_path").String(); got != conf.JournalPath {
		t.Fatalf("journal_path = %q, want %q", got, conf.JournalPath)
	}
	if got := section.Key("ssh_key_file").String(); got != conf.SshKeyFile {
		t.Fatalf("ssh_key_file = %q, want %q", got, conf.SshKeyFile)
	}
}
