package config

import (
	"os"
)

const ApplicationName string = "a7-journal"

var (
	HomeDir   string = os.Getenv("HOME")
	ConfigDir string = ".a7-journal"
)

type Conf struct {
	ConfigPath  string
	JournalPath string
}
