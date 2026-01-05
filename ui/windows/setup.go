package windows

import (
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/never00rei/a7/config"
)

type SetupModel struct {
	Path       string
	SshKeyPath string
	Form       *huh.Form
}

func (s *SetupModel) NewSetup() (*config.Conf, error) {
	var path string
	var sshkeypath string
	var conf *config.Conf

	s.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(&path).
				Title("Where would you like to store your journal?").
				Placeholder(config.Home).
				Suggestions([]string{filepath.Join(config.Home, "Documents", "journal/")}).
				Description("This is the path on the filesystem where you'll store your journal."),
			huh.NewFilePicker().
				Title("SSH Key Path").
				CurrentDirectory(config.SshPath).
				ShowHidden(true).
				Picking(true).
				FileAllowed(true).
				Description("Path to your SSH Key, for encrypting your A7 journal files.").
				Value(&sshkeypath),
		),
	)

	err := s.Form.Run()
	if err != nil {
		return conf, err
	}

	s.Path = path
	s.SshKeyPath = sshkeypath

	conf = config.NewConf(path, sshkeypath)

	err = conf.SaveConfig()
	if err != nil {
		return conf, err
	}

	return conf, nil
}
