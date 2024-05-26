package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/never00rei/a7/utils"
	"gopkg.in/ini.v1"
)

const ApplicationName string = "a7-journal"

var (
	Home          string = os.Getenv("HOME")
	XdgConfigHome string = os.Getenv("XDG_CONFIG_HOME")
	AppConfDir    string = ".a7-journal"
	ConfFileName  string = "conf.ini"
	SshPath       string = filepath.Join(Home, ".ssh")

	HomeConfigEnvVarNotSetError error = errors.New("home and xdg_config_home are not set.")
)

type Conf struct {
	JournalPath string
	SshKeyFile  string
	FirstSetup  bool
}

func BuildConfPath(homeDir, xdgConfigHomeDir string) (string, error) {
	var path string

	if xdgConfigHomeDir != "" {
		path = filepath.Join(xdgConfigHomeDir, AppConfDir)
		return path, nil
	} else if homeDir != "" {
		path = filepath.Join(homeDir, ".config", AppConfDir)
		return path, nil
	}
	// This is intentional, we don't want to return anything
	// if none of the environment variables are set.
	return path, HomeConfigEnvVarNotSetError
}

func NewConf(journalPath, sshKeyPath string) *Conf {
	return &Conf{
		JournalPath: journalPath,
		SshKeyFile:  sshKeyPath,
	}
}

func (c *Conf) JournalPathExists() (bool, error) {
	exists, err := utils.PathExists(c.JournalPath)
	return exists, err
}

func (c *Conf) CreateJournalPath() error {
	journalPathExists, err := c.JournalPathExists()
	if err != nil {
		return err
	}

	if !journalPathExists {
		err := os.MkdirAll(c.JournalPath, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Conf) SaveConfig() error {
	var conf *ini.File
	var err error

	configPath, err := BuildConfPath(Home, XdgConfigHome)
	if err != nil {
		return err
	}

	confPathExists, err := utils.PathExists(configPath)
	if err != nil {
		return err
	}

	if !confPathExists {
		err := os.MkdirAll(configPath, 0755)
		if err != nil {
			return err
		}
		conf = ini.Empty()
	}

	section, err := conf.NewSection("Settings")
	if err != nil {
		return err
	}

	if _, err = section.NewKey("journal_path", c.JournalPath); err != nil {
		return err
	}

	if _, err = section.NewKey("ssh_key_file", c.SshKeyFile); err != nil {
		return err
	}

	if err = conf.SaveTo(filepath.Join(configPath, ConfFileName)); err != nil {
		return err
	}

	err = c.CreateJournalPath()
	if err != nil {
		return err
	}

	return nil
}

func LoadConf() (*Conf, error) {
	confPath, err := BuildConfPath(Home, XdgConfigHome)
	if err != nil {
		return nil, err
	}

	confFilePath := filepath.Join(confPath, ConfFileName)

	confFile, err := ini.Load(confFilePath)
	if err != nil {
		return nil, err
	}

	section, err := confFile.GetSection("Settings")
	if err != nil {
		return nil, err
	}

	journalPath := section.Key("journal_path").String()
	sshKeyPath := section.Key("ssh_key_file").String()

	conf := NewConf(journalPath, sshKeyPath)

	return conf, nil
}
