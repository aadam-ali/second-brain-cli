package internal

import (
	"fmt"
	"io/fs"
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

func CreateNote(pathToFile string, content string) error {
	if err := os.MkdirAll(filepath.Dir(pathToFile), 0770); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	f, err := os.Create(pathToFile)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	_, err = f.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("write note: %w", err)
	}

	return nil
}

func CheckIfNoteExists(rootDir string, name string) (bool, string, error) {
	pathToNote := ""

	if _, err := os.Stat(rootDir); err != nil {
		return false, "", nil
	}

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.EqualFold(name, d.Name()) {
			pathToNote = path
		}
		return nil
	})
	if err != nil {
		return false, "", err
	}

	if pathToNote != "" {
		return true, pathToNote, nil
	}
	return false, "", nil
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
