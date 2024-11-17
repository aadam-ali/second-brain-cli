package cmd

import (
	"fmt"
	"log"
	"os"

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

		filepath := constructNotePath(cfg.InboxDir, title)
		content := renderNoteContent(title)

		createNote(filepath, content)

		fmt.Println(filepath)
	},
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

func renderNoteContent(title string) string {
	return fmt.Sprintf("# %s\n\n", title)
}

func constructNotePath(dir string, title string) string {
	return fmt.Sprintf("%s/%s.md", dir, title)
}
