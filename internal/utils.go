package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func SanitiseTitle(title string) string {
	alphanumericTitle := regexp.MustCompile(`[^A-Za-z0-9-_]+`).ReplaceAllString(title, " ")
	squashedWhitespaceTitle := regexp.MustCompile(`\s+`).ReplaceAllString(alphanumericTitle, " ")
	sanitisedTitle := regexp.MustCompile(`^[-_\s]+|[-_\s]+$`).ReplaceAllString(squashedWhitespaceTitle, "")

	return sanitisedTitle
}

func ConstructNotePath(dir string, title string) string {
	titleWithExtension := fmt.Sprint(title, ".md")
	return filepath.Join(dir, titleWithExtension)
}

func CreateNote(pathToFile string, content string) {
	if err := os.MkdirAll(filepath.Dir(pathToFile), 0770); err != nil {
		return
	}

	f, _ := os.Create(pathToFile)
	defer f.Close()

	_, err := f.Write([]byte(content))
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create note: %s", err)
		log.Fatal(errMsg)
	}
}

func CheckIfNoteExists(rootDir string, name string) (bool, string) {
	pathToNote := ""

	if _, err := os.Stat(rootDir); err != nil {
		return false, ""
	}

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.EqualFold(name, d.Name()) {
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

func OpenFileInVim(rootDir string, filepath string) {
	cmd := exec.Command("nvim", filepath)
	cmd.Dir = rootDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
}

func GetError(template string, a ...any) error {
	error := fmt.Sprintf(template, a...)
	fmt.Fprintln(os.Stderr, error)
	return errors.New(error)
}
