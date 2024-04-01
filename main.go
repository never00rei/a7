package main

import (
	"log"

	"github.com/never00rei/a7/forms"
)

func main() {
	var setup forms.SetupModel
	var note forms.NoteModel
	var err error

	err = setup.NewSetup()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("The path you chose was: %s", setup.Path)

	err = note.TakeNote()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Your title is: %s", note.Title)
	log.Printf("Your note is: %s", note.Content)

}
