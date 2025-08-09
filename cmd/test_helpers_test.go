package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func captureStdout(fn func(cmd *cobra.Command, args []string) error, cmd *cobra.Command, args []string) (string, error) {
	originalStdout := os.Stdout

	var buf bytes.Buffer
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := fn(cmd, args)

	w.Close()
	os.Stdout = originalStdout
	io.Copy(&buf, r)

	return buf.String(), err
}

func prepareEnvironment() string {
	sb, _ := os.MkdirTemp("", "second-brain-cli-")
	os.Mkdir(filepath.Join(sb, "inbox"), 0700)
	os.Mkdir(filepath.Join(sb, "journal"), 0700)

	os.Clearenv()
	os.Setenv("SB", sb)

	return sb
}
