package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aadam-ali/second-brain-cli/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(linkCmd)

	linkCmd.Flags().BoolP("wiki", "w", false, "returns a wikilink when set to true")
}

func linkCmdFunction(cmd *cobra.Command, args []string) error {
	src := args[0]
	dest := args[1]

	destTitle, _ := strings.CutSuffix(filepath.Base(dest), ".md")

	if _, err := os.Stat(dest); err != nil {
		return internal.GetError(err.Error())
	}

	if _, err := os.Stat(src); err != nil {
		return internal.GetError(err.Error())
	}

	relpath, err := filepath.Rel(filepath.Dir(src), dest)

	if err != nil {
		return internal.GetError(err.Error())
	}

	if useWikiLink, _ := cmd.Flags().GetBool("wiki"); useWikiLink {
		fmt.Printf("[[%s]]", destTitle)
	} else {
		fmt.Printf("[%s](%s)", destTitle, relpath)
	}

	return nil
}

var linkCmd = &cobra.Command{
	Use:   "link [src] [dest]",
	Short: "Display a link from src to dest",
	Args:  cobra.MatchAll(cobra.ExactArgs(2)),
	RunE:  linkCmdFunction,
}
