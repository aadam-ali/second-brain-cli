package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "create a new note",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()

		title := args[0]

		noteExists, existingNoteFilepath := checkIfNoteExists(cfg.RootDir, title)

		if !noteExists {
			filepath := constructNotePath(cfg.InboxDir, title)
			content := renderStdNoteContent(title)
			createNote(filepath, content)

			fmt.Println(filepath)
		} else {
			fmt.Printf("Note already exists: %s\n", existingNoteFilepath)
		}
	},
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

func renderStdNoteContent(title string) string {
	return fmt.Sprintf("# %s\n\n", title)
}
