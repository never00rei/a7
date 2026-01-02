package main

import (
	"log"

	"github.com/never00rei/a7/config"
	"github.com/never00rei/a7/forms"
	"github.com/never00rei/a7/ui/windows"
	"github.com/never00rei/a7/utils"
)

func main() {
	var configuration *config.Conf
	var setup windows.SetupModel
	var note forms.NoteModel
	var err error

	configPath, err := config.BuildConfPath(config.Home, config.XdgConfigHome)
	if err != nil {
		log.Fatal(err)
	}

	configPathExists, err := utils.PathExists(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if configPathExists {
		configuration, err = config.LoadConf()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = setup.NewSetup()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Pass in configuration.JournalPath as the save point in "TakeNote".
	err = note.TakeNote()
	if err != nil {
		log.Fatal(err)
	}

	err = note.SaveNote(configuration.JournalPath)
	if err != nil {
		log.Fatal(err)
	}
}
