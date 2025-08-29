package cmd

import (
	"fmt"
	"net/url"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/aadam-ali/second-brain-cli/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pathCmd)
}

func pathCmdFunction(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()
	title := args[0]

	urlEscapedTitle, _ := url.PathUnescape(title)

	noteExists, filepath := internal.CheckIfNoteExists(cfg.RootDir, urlEscapedTitle)

	if !noteExists {
		return internal.GetError("Note with title %q (%s) does not exist", args[0], title)
	}

	fmt.Println(filepath)
	return nil
}

var pathCmd = &cobra.Command{
	Use:   "path [title]",
	Short: "Output path of note if it exists",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	RunE:  pathCmdFunction,
}
