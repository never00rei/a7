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
	Path       string
	SshKeyPath string
	Form       *huh.Form
}

func (s *SetupModel) NewSetup() error {
	var path string
	var sshkeypath string

	s.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(&path).
				Title("Where would you like to store your journal?").
				Placeholder(config.Home).
				Description("This is the path on the filesystem where you'll store your journal."),
			huh.NewFilePicker().
				Title("SSH Key Path").
				CurrentDirectory(config.SshPath).
				ShowHidden(true).
				Picking(true).
				FileAllowed(true).
				Description("Path to SSH key, for encrypting your A7 journal files.").
				Value(&sshkeypath),
		),
	)

	err := s.Form.Run()
	if err != nil {
		return err
	}

	s.Path = path
	s.SshKeyPath = sshkeypath

	return nil
}

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

	n.Title = title
	n.Content = content
	n.Sensitive = sensitive
	n.CurrentDateFormatted = currentDateFormatted

	return nil
}

func (n *NoteModel) SaveNote(journalPath string) error {
	sanitizedTitle := utils.SanitizeSpecialChars(n.Title)
	filename := fmt.Sprintf("%s_%s.md", n.CurrentDateFormatted, sanitizedTitle)

	markdown := fmt.Sprintf("# %s %s\n\n%s", n.CurrentDateFormatted, n.Title, n.Content)

	utils.SaveFile(journalPath, filename, markdown)

	return nil
}
