package forms

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/never00rei/a7/config"
	"github.com/never00rei/a7/utils"
)

var RandomNotePlaceholders []string = []string{"Dear a7, today I have...", "It's been one of those days", "Dear a7, my garden is a mess!!", "I've been thinking about..."}

type SetupModel struct {
	Path string
	Form *huh.Form
}

func (s *SetupModel) NewSetup() error {
	var path string

	s.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(&path).
				Title("Where would you like to store your journal?").
				Placeholder(config.Home).
				Description("This is the path on the filesystem where you'll store your journal."),
		),
	)

	err := s.Form.Run()
	if err != nil {
		return err
	}

	s.Path = path

	return nil
}

type NoteModel struct {
	Content   string
	Title     string
	Timestamp time.Time
	Sensitive bool
	Form      *huh.Form
}

func (n *NoteModel) TakeNote(journalPath string) error {
	var (
		currentDate          time.Time = time.Now()
		currentDateFormatted           = currentDate.Format("2006-01-02_15-04")
		content              string
		title                string
		sensitive            bool
		randomPlaceholder    string = utils.RandomStringFromSlice(RandomNotePlaceholders)
	)

	n.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What's the name of your note?").
				Value(&title).
				Placeholder("My Note"),

			huh.NewConfirm().
				Title("Is this note sensitive?").
				Value(&sensitive).
				Affirmative("Yes").
				Negative("No"),

			huh.NewText().
				Title("Your note for today:").
				Value(&content).
				Placeholder(randomPlaceholder).
				CharLimit(25000).
				Lines(25),
		),
	)

	err := n.Form.Run()
	if err != nil {
		return err
	}

	n.Title = fmt.Sprintf("%s_%s", currentDateFormatted, title)
	n.Content = content
	n.Sensitive = sensitive

	return nil
}
