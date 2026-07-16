package cmd

import (
	"fmt"
	"os"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sb",
	Short: "sb is a note taking management tool",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		cfg := config.GetConfig()
		if err := os.MkdirAll(cfg.InboxDir, 0755); err != nil {
			return fmt.Errorf("create inbox dir: %w", err)
		}
		if err := os.MkdirAll(cfg.JournalDir, 0755); err != nil {
			return fmt.Errorf("create journal dir: %w", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute handles the execution of the provided command
// this may be the root command or any of it's children
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
