package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func constructNotePath(dir string, title string) string {
	return fmt.Sprintf("%s/%s.md", dir, title)
}

func createNote(filepath string, content string) {
	f, _ := os.Create(filepath)
	defer f.Close()

	_, err := f.Write([]byte(content))
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create note: %s", err)
		log.Fatal(errMsg)
	}
}

func checkIfNoteExists(rootDir string, name string) (bool, string) {
	pathToNote := ""
	name = name + ".md"

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && name == d.Name() {
			pathToNote = path
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if pathToNote != "" {
		return true, pathToNote
	}
	return false, ""
}
