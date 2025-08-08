package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

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

func openFileInVim(rootDir string, filepath string) {
	cmd := exec.Command("nvim", filepath)
	cmd.Dir = rootDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
}
