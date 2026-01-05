package forms

import (
	"time"

	"github.com/charmbracelet/huh"
	"github.com/never00rei/a7/journal"
	"github.com/never00rei/a7/utils"
)

var RandomNotePlaceholders []string = []string{"Dear a7, today I have...", "It's been one of those days", "Dear a7, my garden is a mess!!", "I've been thinking about..."}

type NoteModel struct {
	Content              string
	Title                string
	CurrentDateFormatted string
	Sensitive            bool
	Form                 *huh.Form
}

func (n *NoteModel) TakeNote() error {
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
				Title("What's the name of your journal?").
				Value(&title).
				Placeholder("My Journal"),

			huh.NewConfirm().
				Title("Is this journal sensitive?").
				Value(&sensitive).
				Affirmative("Yes").
				Negative("No"),

			huh.NewText().
				Title("Your journal for today:").
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

	n.Title = title
	n.Content = content
	n.Sensitive = sensitive
	n.CurrentDateFormatted = currentDateFormatted

	return nil
}

func (n *NoteModel) SaveNote(journalPath string) error {
	created, err := time.Parse("2006-01-02_15-04", n.CurrentDateFormatted)
	if err != nil {
		created = time.Now()
	}

	service := journal.NewService(journalPath)
	_, err = service.SaveNote(n.Title, n.Content, created)
	return err
}
