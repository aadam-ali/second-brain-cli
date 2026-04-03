package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func captureOutput(fn func(cmd *cobra.Command, args []string) error, cmd *cobra.Command, args []string) (string, string, error) {
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	var bufOut bytes.Buffer
	r1, w1, _ := os.Pipe()
	os.Stdout = w1

	var bufErr bytes.Buffer
	r2, w2, _ := os.Pipe()
	os.Stderr = w2

	err := fn(cmd, args)

	w1.Close()
	os.Stdout = originalStdout
	io.Copy(&bufOut, r1)

	w2.Close()
	os.Stderr = originalStderr
	io.Copy(&bufErr, r2)

	return bufOut.String(), bufErr.String(), err
}

func prepareEnvironment(createDirectories bool) string {
	sb := fmt.Sprintf("/tmp/second-brain-cli-%d", time.Now().UnixNano())

	if createDirectories {
		os.Mkdir(sb, 0700)
		os.Mkdir(filepath.Join(sb, "inbox"), 0700)
		os.Mkdir(filepath.Join(sb, "journal"), 0700)
	}

	os.Clearenv()
	os.Setenv("SB", sb)

	return sb
}
