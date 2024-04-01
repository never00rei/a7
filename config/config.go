package config

import (
	"errors"
	"log"
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

	HomeConfigEnvVarNotSetError error = errors.New("home and xdg_config_home are not set.")
)

type Conf struct {
	ConfPath    string
	JournalPath string
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

func NewConf(journalPath string) Conf {
	confPath, err := BuildConfPath(Home, XdgConfigHome)
	if err != nil {
		log.Fatal(err)
	}

	return Conf{
		ConfPath:    confPath,
		JournalPath: journalPath,
	}
}

func (c *Conf) ConfPathExists() (bool, error) {
	exists, err := utils.PathExists(c.ConfPath)
	return exists, err
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
			log.Fatal(err)
		}
	}

	return nil
}

func (c *Conf) SaveConfig() error {
	var conf *ini.File
	var err error

	confPathExists, err := c.ConfPathExists()
	if err != nil {
		return err
	}

	if !confPathExists {
		err := os.MkdirAll(c.ConfPath, 0755)
		if err != nil {
			log.Fatal(err)
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

	if err = conf.SaveTo(filepath.Join(AppConfDir, ConfFileName)); err != nil {
		return err
	}

	return nil
}
