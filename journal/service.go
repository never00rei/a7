package journal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	Filename string
	Title    string
	Content  string
	Created  time.Time
	ModTime  time.Time
}

type Service struct {
	Root string
}

func NewService(root string) *Service {
	return &Service{Root: root}
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
		Filename: filename,
		Content:  string(content),
	}

	if info, err := os.Stat(filepath.Join(s.Root, filename)); err == nil {
		note.ModTime = info.ModTime()
	}

	note.Title, note.Created = parseHeader(note.Content)

	return note, nil
}

func (s *Service) SaveNote(title, body string, created time.Time) (string, error) {
	if created.IsZero() {
		created = time.Now()
	}

	filename := buildFilename(title, created)
	if err := os.MkdirAll(s.Root, 0755); err != nil {
		return "", fmt.Errorf("create journal dir: %w", err)
	}

	content := fmt.Sprintf("# %s %s\n\n%s", created.Format(timestampLayout), title, body)
	path := filepath.Join(s.Root, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("save note: %w", err)
	}

	return filename, nil
}

func buildFilename(title string, created time.Time) string {
	sanitizedTitle := utils.SanitizeSpecialChars(title)
	return fmt.Sprintf("%s_%s.md", created.Format(timestampLayout), sanitizedTitle)
}

func parseHeader(content string) (string, time.Time) {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return "", time.Time{}
	}

	line := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(line, "# ") {
		return "", time.Time{}
	}

	line = strings.TrimPrefix(line, "# ")
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return "", time.Time{}
	}

	created, err := time.Parse(timestampLayout, parts[0])
	if err != nil {
		return strings.TrimSpace(parts[1]), time.Time{}
	}

	return strings.TrimSpace(parts[1]), created
}
