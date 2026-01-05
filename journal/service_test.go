package journal

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveLoadAndListNotes(t *testing.T) {
	root := t.TempDir()
	svc := NewService(root)

	firstTime := time.Date(2024, 1, 2, 3, 4, 0, 0, time.UTC)
	secondTime := time.Date(2024, 2, 3, 4, 5, 0, 0, time.UTC)

	firstFile, err := svc.SaveNote("First Note", "hello", firstTime)
	if err != nil {
		t.Fatalf("SaveNote first: %v", err)
	}
	secondFile, err := svc.SaveNote("Second Note", "world", secondTime)
	if err != nil {
		t.Fatalf("SaveNote second: %v", err)
	}

	notes, err := svc.ListNotes()
	if err != nil {
		t.Fatalf("ListNotes: %v", err)
	}
	if len(notes) != 2 {
		t.Fatalf("ListNotes count = %d, want 2", len(notes))
	}
	if notes[0].Filename != secondFile {
		t.Fatalf("ListNotes order = %q first, want %q", notes[0].Filename, secondFile)
	}

	loaded, err := svc.LoadNote(firstFile)
	if err != nil {
		t.Fatalf("LoadNote: %v", err)
	}
	if loaded.Title != "First Note" {
		t.Fatalf("LoadNote title = %q, want %q", loaded.Title, "First Note")
	}
	if loaded.Created != firstTime {
		t.Fatalf("LoadNote created = %v, want %v", loaded.Created, firstTime)
	}
	if loaded.Content == "" {
		t.Fatalf("LoadNote content empty")
	}

	path := filepath.Join(root, firstFile)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("note file missing: %v", err)
	}
}
