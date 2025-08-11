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
	dontWantOutputDailyNote := fmt.Sprintf("Daily note not found; creating a new one: %s/journal/2025-07-13.md", sb)

	newCmd.Flags().Set("no-open", "true")
	gotStdout, _, gotError := captureOutput(newCmdFunction, newCmd, []string{"Hello World"})
	_, newNoteErr := os.Stat(wantStdoutFilepath)

	assert.Equal(t, wantError, gotError)
	assert.Contains(t, gotStdout, wantStdoutFilepath)
	assert.NotContains(t, gotStdout, dontWantOutputDailyNote)
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

func TestNewCmdCreateDailyNote(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	wantStdoutFilepath := filepath.Join(sb, "inbox/hello-world.md")
	wantStdoutDailyNote := fmt.Sprintf("Daily note not found; creating a new one: %s/journal/2025-07-13.md", sb)

	newCmd.Flags().Set("no-open", "true")
	gotStdout, _, gotError := captureOutput(newCmdFunction, newCmd, []string{"Hello World"})
	_, newNoteErr := os.Stat(wantStdoutFilepath)
	_, dailyNoteErr := os.Stat(fmt.Sprintf("%s/journal/2025-07-13.md", sb))

	assert.Contains(t, gotStdout, wantStdoutFilepath)
	assert.Contains(t, gotStdout, wantStdoutDailyNote)
	assert.NoError(t, gotError)
	assert.NoError(t, newNoteErr)
	assert.NoError(t, dailyNoteErr)
}

func TestDailyCmd(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	wantStdoutFilepath := filepath.Join(sb, "journal/2025-07-13.md")
	dailyCmd.Flags().Set("no-open", "true")
	gotStdout, _, gotError := captureOutput(dailyCmdFunction, dailyCmd, []string{})
	_, dailyNoteErr := os.Stat(fmt.Sprintf("%s/journal/2025-07-13.md", sb))

	assert.Contains(t, gotStdout, wantStdoutFilepath)
	assert.NoError(t, gotError)
	assert.NoError(t, dailyNoteErr)
}

func TestDailyCmdAlreadyExists(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	wantStdoutFilepath := filepath.Join(sb, "journal/2025-07-13.md")
	os.Create(wantStdoutFilepath)

	dailyCmd.Flags().Set("no-open", "true")
	gotStdout, _, gotError := captureOutput(dailyCmdFunction, dailyCmd, []string{})
	_, dailyNoteErr := os.Stat(fmt.Sprintf("%s/journal/2025-07-13.md", sb))

	assert.Contains(t, gotStdout, wantStdoutFilepath)
	assert.NoError(t, gotError)
	assert.NoError(t, dailyNoteErr)
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
