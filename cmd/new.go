package cmd

import (
	"fmt"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/aadam-ali/second-brain-cli/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().BoolP("no-open", "n", false, "prevents opening of file in editor")
	newCmd.Flags().BoolP("date", "d", false, "prepends the date to the filename")
}

func newCmdFunction(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()

	noOpen, _ := cmd.Flags().GetBool("no-open")
	var filepath string
	title := args[0]

	sanitisedTitle := internal.SanitiseTitle(title)

	if dateFlag, _ := cmd.Flags().GetBool("date"); dateFlag {
		sanitisedTitle = config.Now().Format("20060102 ") + sanitisedTitle
	}

	noteExists, existingNoteFilepath := internal.CheckIfNoteExists(cfg.RootDir, sanitisedTitle+".md")

	if !noteExists {
		filepath = internal.ConstructNotePath(cfg.InboxDir, sanitisedTitle)
		content := renderStdNoteContent(title)
		internal.CreateNote(filepath, content)

		fmt.Println(filepath)
	} else {
		return internal.GetError("Note with title %q already exists at %s", sanitisedTitle, existingNoteFilepath)
	}

	if !noOpen {
		internal.OpenFileInVim(cfg.RootDir, filepath)
	}

	return nil
}

var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "create a new note",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	RunE:  newCmdFunction}

func renderStdNoteContent(title string) string {
	return fmt.Sprintf("# %s\n\n", title)
}
