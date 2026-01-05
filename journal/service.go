package journal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/never00rei/a7/utils"
)

const timestampLayout = "2006-01-02_15-04"

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
}

type Option func(*Service)

func NewService(root string, opts ...Option) *Service {
	svc := &Service{Root: root}
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
	entries, err := os.ReadDir(s.Root)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}

	notes := make([]NoteInfo, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("stat note: %w", err)
		}
		notes = append(notes, NoteInfo{
			Filename: entry.Name(),
			ModTime:  info.ModTime(),
		})
	}

	sort.Slice(notes, func(i, j int) bool {
		return notes[i].ModTime.After(notes[j].ModTime)
	})

	return notes, nil
}

func (s *Service) LoadNote(filename string) (*Note, error) {
	content, err := os.ReadFile(filepath.Join(s.Root, filename))
	if err != nil {
		return nil, fmt.Errorf("load note: %w", err)
	}

	note := &Note{
		Filename:  filename,
		WordCount: -1,
	}

	body := string(content)
	if info, err := os.Stat(filepath.Join(s.Root, filename)); err == nil {
		note.ModTime = info.ModTime()
	}

	title, created, updated, encrypted, wordCount, remaining := parseFrontMatter(body)
	if title != "" || !created.IsZero() || !updated.IsZero() {
		note.Title = title
		note.Created = created
		note.Updated = updated
		note.Encrypted = encrypted
		note.WordCount = wordCount
		if encrypted {
			decrypted, err := decryptBody(remaining, s.SSHKeyPath)
			if err != nil {
				return note, fmt.Errorf("decrypt note: %w", err)
			}
			note.Content = decrypted
		} else {
			note.Content = remaining
		}
		return note, nil
	}

	note.Title, note.Created, note.Content = parseHeader(body)

	return note, nil
}

func (s *Service) SaveNote(title, body string, created time.Time) (string, error) {
	if created.IsZero() {
		created = time.Now()
	}

	updated := time.Now()
	filename := buildFilename(title, created)
	contentBody, encrypted, err := maybeEncryptBody(body, s.Encrypt, s.SSHKeyPath)
	if err != nil {
		return "", err
	}
	wordCount := countWords(body)
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
	contentBody, encrypted, err := maybeEncryptBody(body, s.Encrypt, s.SSHKeyPath)
	if err != nil {
		return err
	}
	wordCount := countWords(body)
	return s.writeNoteFile(filename, title, contentBody, created, updated, encrypted, wordCount)
}

func (s *Service) writeNoteFile(filename, title, body string, created, updated time.Time, encrypted bool, wordCount int) error {
	if err := os.MkdirAll(s.Root, 0755); err != nil {
		return fmt.Errorf("create journal dir: %w", err)
	}

	content := fmt.Sprintf(
		"---\n"+
			"title: %s\n"+
			"created: %s\n"+
			"updated: %s\n"+
			"encrypted: %t\n"+
			"word_count: %d\n"+
			"---\n\n%s",
		title,
		created.Format(time.RFC3339),
		updated.Format(time.RFC3339),
		encrypted,
		wordCount,
		body,
	)
	path := filepath.Join(s.Root, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("save note: %w", err)
	}
	return nil
}

func buildFilename(title string, created time.Time) string {
	sanitizedTitle := utils.SanitizeSpecialChars(title)
	return fmt.Sprintf("%s_%s.md", created.Format(timestampLayout), sanitizedTitle)
}

func parseHeader(content string) (string, time.Time, string) {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return "", time.Time{}, content
	}

	line := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(line, "# ") {
		return "", time.Time{}, content
	}

	line = strings.TrimPrefix(line, "# ")
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return "", time.Time{}, content
	}

	created, err := time.Parse(timestampLayout, parts[0])
	if err != nil {
		return strings.TrimSpace(parts[1]), time.Time{}, strings.Join(lines[1:], "\n")
	}

	body := strings.Join(lines[1:], "\n")
	body = strings.TrimPrefix(body, "\n")
	return strings.TrimSpace(parts[1]), created, body
}

func parseFrontMatter(content string) (string, time.Time, time.Time, bool, int, string) {
	lines := strings.Split(content, "\n")
	if len(lines) < 3 {
		return "", time.Time{}, time.Time{}, false, -1, content
	}
	if strings.TrimSpace(lines[0]) != "---" {
		return "", time.Time{}, time.Time{}, false, -1, content
	}

	title := ""
	var created time.Time
	var updated time.Time
	encrypted := false
	wordCount := -1
	end := -1
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "---" {
			end = i
			break
		}
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		switch key {
		case "title":
			title = value
		case "created":
			created = parseTimestamp(value)
		case "updated":
			updated = parseTimestamp(value)
		case "encrypted":
			encrypted = strings.EqualFold(value, "true")
		case "word_count":
			if parsed, err := strconv.Atoi(value); err == nil {
				wordCount = parsed
			}
		}
	}
	if end == -1 {
		return "", time.Time{}, time.Time{}, false, -1, content
	}

	body := strings.Join(lines[end+1:], "\n")
	body = strings.TrimPrefix(body, "\n")
	return title, created, updated, encrypted, wordCount, body
}

func parseTimestamp(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	if ts, err := time.Parse(time.RFC3339, value); err == nil {
		return ts
	}
	if ts, err := time.Parse(timestampLayout, value); err == nil {
		return ts
	}
	return time.Time{}
}

func countWords(content string) int {
	fields := strings.Fields(content)
	return len(fields)
}
