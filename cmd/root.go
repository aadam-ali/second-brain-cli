package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sb",
	Short: "sb is a note taking management tool",
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
