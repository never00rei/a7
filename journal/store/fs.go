package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type NoteInfo struct {
	Filename string
	ModTime  time.Time
}

type FS struct {
	Root string
}

func NewFS(root string) *FS {
	return &FS{Root: root}
}

func (s *FS) ListMarkdown() ([]NoteInfo, error) {
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

func (s *FS) Read(filename string) (string, time.Time, error) {
	path := filepath.Join(s.Root, filename)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("load note: %w", err)
	}

	var modTime time.Time
	if info, err := os.Stat(path); err == nil {
		modTime = info.ModTime()
	}

	return string(content), modTime, nil
}

func (s *FS) Write(filename, content string) error {
	if err := os.MkdirAll(s.Root, 0755); err != nil {
		return fmt.Errorf("create journal dir: %w", err)
	}
	path := filepath.Join(s.Root, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("save note: %w", err)
	}
	return nil
}
