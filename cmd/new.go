package cmd

import (
	"fmt"
	"os"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/aadam-ali/second-brain-cli/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().BoolP("no-open", "n", false, "prevents opening of file in editor")
	newCmd.Flags().BoolP("no-date", "d", false, "prepends the date to the filename")
}

func newCmdFunction(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()

	noOpen, err := cmd.Flags().GetBool("no-open")
	if err != nil {
		return fmt.Errorf("flag no-open: %w", err)
	}
	var notePath string
	title := args[0]

	sanitisedTitle := internal.SanitiseTitle(title)

	dateFlag, err := cmd.Flags().GetBool("no-date")
	if err != nil {
		return fmt.Errorf("flag no-date: %w", err)
	}
	if !dateFlag {
		sanitisedTitle = cfg.Today + " " + sanitisedTitle
	}

	noteExists, existingNoteFilepath, err := internal.CheckIfNoteExists(cfg.RootDir, sanitisedTitle+".md")
	if err != nil {
		return fmt.Errorf("check note: %w", err)
	}

	if !noteExists {
		notePath = internal.ConstructNotePath(cfg.InboxDir, sanitisedTitle)
		content := renderStdNoteContent(title, cfg.Today)
		if err := internal.CreateNote(notePath, content); err != nil {
			return fmt.Errorf("create note: %w", err)
		}

		fmt.Println(notePath)
	} else {
		return fmt.Errorf("%w: %q at %s", internal.ErrNoteExists, sanitisedTitle, existingNoteFilepath)
	}

	if !noOpen {
		if err := internal.OpenFileInEditor(cfg.Editor, cfg.RootDir, notePath); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	return nil
}

var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "create a new note",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	RunE:  newCmdFunction}

func renderStdNoteContent(title string, date string) string {
	return fmt.Sprintf("---\ntitle: %s\ndate: %s\n---\n", title, date)
}
