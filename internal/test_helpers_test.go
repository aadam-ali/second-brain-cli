package internal

import (
	"os"
	"path/filepath"
)

func createNoteInTempDir(title string, nested bool) (string, string, string) {
	var pathToNote string

	filename := title + ".md"

	rootDir, _ := os.MkdirTemp("", "second-brain-cli")

	if nested {
		dir, _ := os.MkdirTemp(rootDir, "nested")
		pathToNote = filepath.Join(dir, filename)
	} else {
		pathToNote = filepath.Join(rootDir, filename)
	}

	CreateNote(pathToNote, "content")

	return rootDir, filename, pathToNote
}
