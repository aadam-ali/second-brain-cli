package cmd

import (
	"fmt"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionCmdFunction(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()

	fmt.Printf("sb %s\n", cfg.Version)

	return nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output version of sb",
	Args:  cobra.MatchAll(cobra.ExactArgs(0)),
	RunE:  versionCmdFunction,
}
