package cmd

import (
	"fmt"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().BoolP("no-open", "n", false, "prevents opening of file in editor")
}

var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "create a new note",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()

		no_open, _ := cmd.Flags().GetBool("no-open")
		var filepath string

		title := args[0]

		noteExists, existingNoteFilepath := checkIfNoteExists(cfg.RootDir, title)

		if !noteExists {
			filepath = constructNotePath(cfg.InboxDir, title)
			content := renderStdNoteContent(title)
			createNote(filepath, content)

			fmt.Println(filepath)
		} else {
			filepath = existingNoteFilepath
			fmt.Printf("Note already exists: %s\n", filepath)
		}

		if no_open == false {
			openFileInVim(cfg.RootDir, filepath)
		}
	},
}

func renderStdNoteContent(title string) string {
	return fmt.Sprintf("# %s\n\n", title)
}
