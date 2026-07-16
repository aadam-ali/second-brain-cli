package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func createNoteInTempDir(t *testing.T, title string, nested bool) (string, string, string) {
	t.Helper()

	var pathToNote string

	filename := title + ".md"

	rootDir, err := os.MkdirTemp("", "second-brain-cli")
	if err != nil {
		t.Fatal(err)
	}

	if nested {
		dir, err := os.MkdirTemp(rootDir, "nested")
		if err != nil {
			t.Fatal(err)
		}
		pathToNote = filepath.Join(dir, filename)
	} else {
		pathToNote = filepath.Join(rootDir, filename)
	}

	if err := CreateNote(pathToNote, "content"); err != nil {
		t.Fatal(err)
	}

	return rootDir, filename, pathToNote
}
