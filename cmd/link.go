package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(linkCmd)

	linkCmd.Flags().BoolP("wiki", "w", false, "returns a wikilink when set to true")
}

func linkCmdFunction(cmd *cobra.Command, args []string) error {
	src := args[0]
	dest := args[1]

	destFilename := filepath.Base(dest)
	destTitle, _ := strings.CutSuffix(destFilename, ".md")

	if _, err := os.Stat(dest); err != nil {
		return fmt.Errorf("access dest: %w", err)
	}

	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("access src: %w", err)
	}

	if useWikiLink, _ := cmd.Flags().GetBool("wiki"); useWikiLink {
		fmt.Printf("[[%s]]", destTitle)

		return nil
	}

	relpath, err := filepath.Rel(filepath.Dir(src), dest)

	if err != nil {
		return fmt.Errorf("resolve relpath: %w", err)
	}

	urlEncodedFilename := url.PathEscape(destFilename)
	fmt.Printf("[%s](%s)", destTitle, filepath.Join(filepath.Dir(relpath), urlEncodedFilename))

	return nil
}

var linkCmd = &cobra.Command{
	Use:   "link [src] [dest]",
	Short: "Output a link to another note",
	Args:  cobra.MatchAll(cobra.ExactArgs(2)),
	RunE:  linkCmdFunction,
}
