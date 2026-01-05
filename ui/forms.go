package ui

import (
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/never00rei/a7/config"
)

const (
	storagePathKey = "journal_path"
	sshKeyPathKey  = "ssh_key_path"
	encryptKey     = "encrypt"
)

func newStorageForm(path *string, width int) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key(storagePathKey).
				Value(path).
				Title("Journal folder").
				Placeholder(config.Home).
				Suggestions([]string{filepath.Join(config.Home, "Documents", "journal/")}).
				Description("Where A7 should save your journal files."),
		),
	).WithShowHelp(false)

	if width > 0 {
		form.WithWidth(width)
	}

	return form
}

func newPrivacyForm(encrypt *bool, sshKeyPath *string, width int) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key(encryptKey).
				Value(encrypt).
				Title("Encrypt sensitive notes?").
				Affirmative("Yes").
				Negative("No"),
		),
		huh.NewGroup(
			huh.NewFilePicker().
				Key(sshKeyPathKey).
				Value(sshKeyPath).
				Title("SSH key path (optional)").
				CurrentDirectory(config.SshPath).
				ShowHidden(true).
				Picking(true).
				FileAllowed(true).
				Height(12).
				Description("Choose an SSH key to encrypt sensitive notes."),
		).WithHideFunc(func() bool {
			if encrypt == nil {
				return true
			}
			return !*encrypt
		}),
	).WithShowHelp(false)

	if width > 0 {
		form.WithWidth(width)
	}

	return form
}
