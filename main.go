package main

import (
	"log"

	"github.com/never00rei/a7/forms"
)

func main() {
	var setup forms.SetupModel
	var note forms.NoteModel

	setup.NewSetup()

	log.Printf("The path you chose was: %s", setup.Path)

	note.TakeNote()

	log.Printf("Your title is: %s", note.Title)
	log.Printf("Your note is: %s", note.Content)

}
