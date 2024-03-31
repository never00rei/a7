package config

import (
	"os"
	"path/filepath"
)

const ApplicationName string = "a7-journal"

var (
	Home          string = os.Getenv("HOME")
	XdgConfigHome string = os.Getenv("XDG_CONFIG_HOME")
	AppConfDir    string = ".a7-journal"
)

type Conf struct {
	ConfigPath  string
	JournalPath string
}

func BuildConfPath(homeDir, xdgConfigHomeDir string) string {
	if xdgConfigHomeDir != "" {
		return filepath.Join(xdgConfigHomeDir, AppConfDir)
	} else if homeDir != "" {
		return filepath.Join(homeDir, ".config", AppConfDir)
	}
	// This is intentional, we don't want to return anything
	// if none of the environment variables are set.
	return ""
}

func NewConf() Conf {
	return Conf{}
}
