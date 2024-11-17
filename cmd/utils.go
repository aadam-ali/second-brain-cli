package cmd

import (
	"fmt"
	"log"
	"os"
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
