package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/aadam-ali/second-brain-cli/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dailyCmd)

	dailyCmd.Flags().BoolP("no-open", "n", false, "prevents opening of file in editor")
}

func dailyCmdFunction(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()

	noOpen, _ := cmd.Flags().GetBool("no-open")
	filepath := cfg.DailyNotePath

	dailyNoteExists := checkIfDailyNoteExists(filepath)

	if !dailyNoteExists {
		content := renderDailyNoteContent(cfg.Yesterday, cfg.Today, cfg.Tomorrow)
		internal.CreateNote(filepath, content)
	}

	fmt.Println(filepath)

	if !noOpen {
		internal.OpenFileInVim(cfg.RootDir, cfg.DailyNotePath)
	}

	return nil
}

var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "create a daily note",
	Args:  cobra.MatchAll(cobra.ExactArgs(0)),
	RunE:  dailyCmdFunction,
}

func checkIfDailyNoteExists(filepath string) bool {
	_, err := os.Stat(filepath)

	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}

	log.Fatal(err)

	return false
}

func renderDailyNoteContent(yesterday string, today string, tomorrow string) string {
	return fmt.Sprintf(`# %s

[[%s]] - [[%s]]

## Journal


## Notes Created Today
`, today, yesterday, tomorrow)
}

func appendToDailyNote(cfg config.Configuration, title string) {
	if !checkIfDailyNoteExists(cfg.DailyNotePath) {
		content := renderDailyNoteContent(cfg.Yesterday, cfg.Today, cfg.Tomorrow)
		internal.CreateNote(cfg.DailyNotePath, content)
		fmt.Printf("Daily note not found; creating a new one: %s\n", cfg.DailyNotePath)
	}

	f, err := os.OpenFile(cfg.DailyNotePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("\n[[%s]]", title)); err != nil {
		log.Fatal(err)
	}
}
