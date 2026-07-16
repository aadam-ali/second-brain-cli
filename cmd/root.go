package cmd

import (
	"log"
	"os"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sb",
	Short: "sb is a note taking management tool",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true

		config := config.GetConfig()
		os.MkdirAll(config.InboxDir, 0755)
		os.MkdirAll(config.JournalDir, 0755)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()

		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute handles the execution of the provided command
// this may be the root command or any of it's children
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
