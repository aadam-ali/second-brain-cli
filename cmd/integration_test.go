package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aadam-ali/second-brain-cli/config"
	"github.com/stretchr/testify/assert"
)

func TestVersionCmd(t *testing.T) {
	wantStdout := "sb development\n"

	gotStdout, _, gotError := captureOutput(versionCmdFunction, versionCmd, []string{})

	assert.NoError(t, gotError)
	assert.Equal(t, wantStdout, gotStdout)
}

func TestNewCmd(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	os.Create(filepath.Join(sb, "journal/2025-07-13.md"))

	var wantError error
	wantStdoutFilepath := filepath.Join(sb, "inbox/hello-world.md")

	newCmd.Flags().Set("no-open", "true")
	gotStdout, _, gotError := captureOutput(newCmdFunction, newCmd, []string{"Hello World"})
	_, newNoteErr := os.Stat(wantStdoutFilepath)

	assert.Equal(t, wantError, gotError)
	assert.Contains(t, gotStdout, wantStdoutFilepath)
	assert.NoError(t, newNoteErr)
}

func TestNewCmdExistingNote(t *testing.T) {
	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	wantStdoutFilepath := filepath.Join(sb, "inbox/hello-world.md")
	wantStderr := fmt.Sprintf("Note with title \"hello-world\" already exists at %s", wantStdoutFilepath)

	os.Create(wantStdoutFilepath)

	newCmd.Flags().Set("no-open", "true")
	_, gotStderr, gotError := captureOutput(newCmdFunction, newCmd, []string{"Hello World"})
	_, newNoteErr := os.Stat(wantStdoutFilepath)

	assert.Contains(t, gotStderr, wantStderr)
	assert.NoError(t, newNoteErr)
	assert.ErrorContains(t, gotError, wantStderr)
}

func TestDailyCmd(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	wantStdout := `### 2025-07-13 Sunday
@ <location>

- todo:
---
`

	gotStdout, _, gotError := captureOutput(dailyCmdFunction, dailyCmd, []string{})

	assert.Equal(t, wantStdout, gotStdout)
	assert.NoError(t, gotError)
}

func TestPathCmdExists(t *testing.T) {
	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	wantStdoutFilepath := filepath.Join(sb, "inbox/hello-world.md")
	os.Create(wantStdoutFilepath)

	gotStdout, _, gotError := captureOutput(pathCmdFunction, pathCmd, []string{"hello-world"})
	_, statErr := os.Stat(wantStdoutFilepath)

	assert.Contains(t, gotStdout, wantStdoutFilepath)
	assert.NoError(t, gotError)
	assert.NoError(t, statErr)
}

func TestPathCmdDoesNotExist(t *testing.T) {
	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	wantStdoutFilepath := filepath.Join(sb, "somefolder/hello-world.md")
	wantStderr := "Note with title \"hello-world\" (hello-world) does not exist"

	_, gotStderr, gotError := captureOutput(pathCmdFunction, pathCmd, []string{"hello-world"})
	_, statErr := os.Stat(wantStdoutFilepath)

	assert.Contains(t, gotStderr, wantStderr)
	assert.Error(t, statErr)
	assert.ErrorContains(t, gotError, wantStderr)
}

func TestLinkCmd(t *testing.T) {
	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	src := filepath.Join(sb, "journal/2025-07-13.md")
	dest := filepath.Join(sb, "inbox", "hello-world.md")

	os.Create(src)
	os.Create(dest)

	wantOutput := "[hello-world](../inbox/hello-world.md)"

	gotOutput, _, gotError := captureOutput(linkCmdFunction, newCmd, []string{src, dest})

	assert.NoError(t, gotError)
	assert.Equal(t, wantOutput, gotOutput)
}

func TestLinkCmdWikiLink(t *testing.T) {
	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	src := filepath.Join(sb, "journal/2025-07-13.md")
	dest := filepath.Join(sb, "inbox", "hello-world.md")

	os.Create(src)
	os.Create(dest)

	wantOutput := "[[hello-world]]"

	linkCmd.Flags().Set("wiki", "true")
	gotOutput, _, gotError := captureOutput(linkCmdFunction, linkCmd, []string{src, dest})

	assert.NoError(t, gotError)
	assert.Equal(t, wantOutput, gotOutput)
}

func TestLinkCmdDestNotExist(t *testing.T) {
	sb := prepareEnvironment()
	os.RemoveAll(sb)

	src := filepath.Join(sb, "journal/2025-07-13.md")
	dest := filepath.Join(sb, "hello-world.md")

	os.Create(src)

	wantError := fmt.Sprintf("stat %s: no such file or directory", dest)

	_, gotStderr, gotError := captureOutput(linkCmdFunction, linkCmd, []string{src, dest})

	assert.Contains(t, gotStderr, wantError)
	assert.Error(t, gotError)
	assert.ErrorContains(t, gotError, wantError)
}
