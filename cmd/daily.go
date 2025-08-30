package cmd

import (
	"fmt"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dailyCmd)
}

func dailyCmdFunction(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()

	fmt.Printf(`### %s %s
@ <location>

- todo:
---
`, cfg.Today, cfg.DayOfWeek)

	return nil
}

var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Generate a daily note entry",
	Args:  cobra.MatchAll(cobra.ExactArgs(0)),
	RunE:  dailyCmdFunction,
}
