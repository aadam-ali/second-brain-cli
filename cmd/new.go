package cmd

import (
	"fmt"

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
		content := renderStdNoteContent(title)

		createNote(filepath, content)

		fmt.Println(filepath)
	},
}

func renderStdNoteContent(title string) string {
	return fmt.Sprintf("# %s\n\n", title)
}
