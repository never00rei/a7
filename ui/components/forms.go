package components

import (
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/never00rei/a7/config"
)

const (
	StoragePathKey   = "journal_path"
	SshKeyPathKey    = "ssh_key_path"
	SshPubKeyPathKey = "ssh_pub_key_path"
	EncryptKey       = "encrypt"
)

func NewStorageForm(path *string, width int) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key(StoragePathKey).
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

func NewPrivacyForm(encrypt *bool, sshKeyPath *string, sshPubKeyPath *string, width int) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key(EncryptKey).
				Value(encrypt).
				Title("Encrypt sensitive journals?").
				Affirmative("Yes").
				Negative("No"),
		),
		huh.NewGroup(
			huh.NewFilePicker().
				Key(SshKeyPathKey).
				Value(sshKeyPath).
				Title("SSH private key (required)").
				CurrentDirectory(config.SshPath).
				ShowHidden(true).
				Picking(true).
				FileAllowed(true).
				Height(12).
				Description("Private key used to decrypt encrypted journals."),
			huh.NewFilePicker().
				Key(SshPubKeyPathKey).
				Value(sshPubKeyPath).
				Title("SSH public key (required)").
				CurrentDirectory(config.SshPath).
				ShowHidden(true).
				Picking(true).
				FileAllowed(true).
				Height(12).
				Description("Public key used to encrypt journal content."),
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

func NewSettingsForm(path *string, encrypt *bool, sshKeyPath *string, sshPubKeyPath *string, width int) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key(StoragePathKey).
				Value(path).
				Title("Journal folder").
				Placeholder(config.Home).
				Suggestions([]string{filepath.Join(config.Home, "Documents", "journal/")}).
				Description("Where A7 should save your journal files."),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Key(EncryptKey).
				Value(encrypt).
				Title("Encrypt sensitive journals?").
				Affirmative("Yes").
				Negative("No"),
		),
		huh.NewGroup(
			huh.NewFilePicker().
				Key(SshKeyPathKey).
				Value(sshKeyPath).
				Title("SSH private key (required)").
				CurrentDirectory(config.SshPath).
				ShowHidden(true).
				Picking(true).
				FileAllowed(true).
				Height(12).
				Description("Private key used to decrypt encrypted journals."),
			huh.NewFilePicker().
				Key(SshPubKeyPathKey).
				Value(sshPubKeyPath).
				Title("SSH public key (required)").
				CurrentDirectory(config.SshPath).
				ShowHidden(true).
				Picking(true).
				FileAllowed(true).
				Height(12).
				Description("Public key used to encrypt journal content."),
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
