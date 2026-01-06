package journal

import (
	"fmt"
	"time"

	"github.com/never00rei/a7/journal/codec"
	"github.com/never00rei/a7/journal/crypto"
	"github.com/never00rei/a7/journal/store"
)

type NoteInfo struct {
	Filename string
	ModTime  time.Time
}

type Note struct {
	Filename  string
	Title     string
	Content   string
	Created   time.Time
	ModTime   time.Time
	Updated   time.Time
	Encrypted bool
	WordCount int
}

type Service struct {
	Root       string
	Encrypt    bool
	SSHKeyPath string
	store      *store.FS
}

type Option func(*Service)

func NewService(root string, opts ...Option) *Service {
	svc := &Service{Root: root, store: store.NewFS(root)}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func WithEncryption(enabled bool, sshKeyPath string) Option {
	return func(s *Service) {
		s.Encrypt = enabled
		s.SSHKeyPath = sshKeyPath
	}
}

func (s *Service) ListNotes() ([]NoteInfo, error) {
	entries, err := s.store.ListMarkdown()
	if err != nil {
		return nil, err
	}

	notes := make([]NoteInfo, 0, len(entries))
	for _, entry := range entries {
		notes = append(notes, NoteInfo{
			Filename: entry.Filename,
			ModTime:  entry.ModTime,
		})
	}
	return notes, nil
}

func (s *Service) LoadNote(filename string) (*Note, error) {
	content, modTime, err := s.store.Read(filename)
	if err != nil {
		return nil, err
	}

	note := &Note{
		Filename:  filename,
		WordCount: -1,
	}
	note.ModTime = modTime

	matter, remaining := codec.ParseFrontMatter(content)
	if matter.Title != "" || !matter.Created.IsZero() || !matter.Updated.IsZero() {
		note.Title = matter.Title
		note.Created = matter.Created
		note.Updated = matter.Updated
		note.Encrypted = matter.Encrypted
		note.WordCount = matter.WordCount
		if matter.Encrypted {
			decrypted, err := crypto.DecryptBody(remaining, s.SSHKeyPath)
			if err != nil {
				return note, fmt.Errorf("decrypt note: %w", err)
			}
			note.Content = decrypted
		} else {
			note.Content = remaining
		}
		return note, nil
	}

	note.Title, note.Created, note.Content = codec.ParseHeader(content)
	return note, nil
}

func (s *Service) SaveNote(title, body string, created time.Time) (string, error) {
	if created.IsZero() {
		created = time.Now()
	}

	updated := time.Now()
	filename := codec.BuildFilename(title, created)
	contentBody, encrypted, err := crypto.MaybeEncryptBody(body, s.Encrypt, s.SSHKeyPath)
	if err != nil {
		return "", err
	}
	wordCount := codec.CountWords(body)
	if err := s.writeNoteFile(filename, title, contentBody, created, updated, encrypted, wordCount); err != nil {
		return "", err
	}

	return filename, nil
}

func (s *Service) UpdateNote(filename, title, body string, created time.Time) error {
	if created.IsZero() {
		created = time.Now()
	}
	updated := time.Now()
	contentBody, encrypted, err := crypto.MaybeEncryptBody(body, s.Encrypt, s.SSHKeyPath)
	if err != nil {
		return err
	}
	wordCount := codec.CountWords(body)
	return s.writeNoteFile(filename, title, contentBody, created, updated, encrypted, wordCount)
}

func (s *Service) writeNoteFile(filename, title, body string, created, updated time.Time, encrypted bool, wordCount int) error {
	content := codec.RenderContent(title, body, created, updated, encrypted, wordCount)
	return s.store.Write(filename, content)
}
