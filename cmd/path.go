package cmd

import (
	"fmt"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pathCmd)
}

var pathCmd = &cobra.Command{
	Use:   "path [title]",
	Short: "outputs path of note if it exists",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		title := args[0]

		title = titleToKebabCase(title)

		noteExists, filepath := checkIfNoteExists(cfg.RootDir, title)

		if noteExists {
			fmt.Println(filepath)
		}
	},
}
