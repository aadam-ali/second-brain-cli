package cmd

import (
	"fmt"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dailyCmd)
}

var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "create a daily note",
	Args:  cobra.MatchAll(cobra.ExactArgs(0)),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()

		filepath := cfg.DailyNotePath

		content := renderDailyNoteContent(cfg.Yesterday, cfg.Today, cfg.Tomorrow)

		createNote(filepath, content)

		fmt.Println(filepath)
	},
}

func renderDailyNoteContent(yesterday string, today string, tomorrow string) string {
	return fmt.Sprintf(`# %s

[[%s]] - [[%s]]

## Journal


## Notes Created Today
`, today, yesterday, tomorrow)
}
