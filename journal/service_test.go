package journal

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/never00rei/a7/journal/codec"
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
	foundFirst := false
	foundSecond := false
	for _, note := range notes {
		switch note.Filename {
		case firstFile:
			foundFirst = true
		case secondFile:
			foundSecond = true
		}
	}
	if !foundFirst || !foundSecond {
		t.Fatalf("ListNotes missing files: first=%t second=%t", foundFirst, foundSecond)
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
	if loaded.Encrypted {
		t.Fatalf("LoadNote encrypted = true, want false")
	}
	if loaded.Content != "hello" {
		t.Fatalf("LoadNote content = %q, want %q", loaded.Content, "hello")
	}
	if loaded.Updated.IsZero() {
		t.Fatalf("LoadNote updated missing")
	}

	path := filepath.Join(root, firstFile)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("note file missing: %v", err)
	}
}

func TestUpdateNotePreservesFilename(t *testing.T) {
	root := t.TempDir()
	svc := NewService(root)

	created := time.Date(2024, 3, 4, 5, 6, 0, 0, time.UTC)
	filename, err := svc.SaveNote("Original", "first", created)
	if err != nil {
		t.Fatalf("SaveNote: %v", err)
	}

	if err := svc.UpdateNote(filename, "Updated Title", "changed", created); err != nil {
		t.Fatalf("UpdateNote: %v", err)
	}

	loaded, err := svc.LoadNote(filename)
	if err != nil {
		t.Fatalf("LoadNote: %v", err)
	}
	if loaded.Title != "Updated Title" {
		t.Fatalf("Title = %q, want %q", loaded.Title, "Updated Title")
	}
	if loaded.Encrypted {
		t.Fatalf("Encrypted = true, want false")
	}
	if loaded.Content != "changed" {
		t.Fatalf("Content = %q, want %q", loaded.Content, "changed")
	}
	if loaded.Created != created {
		t.Fatalf("Created = %v, want %v", loaded.Created, created)
	}
	if loaded.Updated.IsZero() {
		t.Fatalf("Updated missing")
	}
}

func TestListNotesLoadsMetadata(t *testing.T) {
	root := t.TempDir()
	svc := NewService(root)

	created := time.Date(2024, 4, 5, 6, 7, 0, 0, time.UTC)
	filename, err := svc.SaveNote("Meta Note", "hello world", created)
	if err != nil {
		t.Fatalf("SaveNote: %v", err)
	}

	notes, err := svc.ListNotes()
	if err != nil {
		t.Fatalf("ListNotes: %v", err)
	}

	var info *NoteInfo
	for i := range notes {
		if notes[i].Filename == filename {
			info = &notes[i]
			break
		}
	}
	if info == nil {
		t.Fatalf("ListNotes missing %q", filename)
	}
	if info.Title != "Meta Note" {
		t.Fatalf("Title = %q, want %q", info.Title, "Meta Note")
	}
	if !info.Created.Equal(created) {
		t.Fatalf("Created = %v, want %v", info.Created, created)
	}
	if info.Updated.IsZero() {
		t.Fatalf("Updated missing")
	}
	if info.Encrypted {
		t.Fatalf("Encrypted = true, want false")
	}
	if info.WordCount != codec.CountWords("hello world") {
		t.Fatalf("WordCount = %d, want %d", info.WordCount, codec.CountWords("hello world"))
	}
}
