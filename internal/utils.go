package internal

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func TitleToKebabCase(title string) string {
	title = strings.ToLower(title)

	title = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(title, "-")
	title = regexp.MustCompile(`^-+|-+$`).ReplaceAllString(title, "")

	return title
}

func ConstructNotePath(dir string, title string) string {
	titleWithExtension := fmt.Sprint(title, ".md")
	return filepath.Join(dir, titleWithExtension)
}

func CreateNote(filepath string, content string) {
	f, _ := os.Create(filepath)
	defer f.Close()

	_, err := f.Write([]byte(content))
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create note: %s", err)
		log.Fatal(errMsg)
	}
}

func CheckIfNoteExists(rootDir string, name string) (bool, string) {
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
