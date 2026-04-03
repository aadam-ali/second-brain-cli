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
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true
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
	cfg := config.GetConfig()

	os.MkdirAll(cfg.RootDir, os.ModePerm)
	os.MkdirAll(cfg.InboxDir, os.ModePerm)
	os.MkdirAll(cfg.JournalDir, os.ModePerm)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
