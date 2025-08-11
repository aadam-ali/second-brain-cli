package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aadam-ali/second-brain-cli/config"
)

func TestVersionCmd(t *testing.T) {
	var wantError error
	wantOutput := "sb development\n"

	gotStdout, _, gotError := captureOutput(versionCmdFunction, versionCmd, []string{})

	if gotError != wantError {
		t.Errorf("got %q, want %q", gotError, wantError)
	}

	if gotStdout != wantOutput {
		t.Errorf("got %q, want %q", gotStdout, wantOutput)
	}
}

func TestNewCmd(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	os.Create(filepath.Join(sb, "journal/2025-07-13.md"))

	var wantError error
	wantOutputFilepath := filepath.Join(sb, "inbox/hello-world.md")
	dontWantOutputDailyNote := fmt.Sprintf("Daily note not found; creating a new one: %s/journal/2025-07-13.md", sb)

	newCmd.Flags().Set("no-open", "true")
	gotStdout, _, gotError := captureOutput(newCmdFunction, newCmd, []string{"Hello World"})

	if gotError != wantError {
		t.Errorf("got %q, want %q", gotError, wantError)
	}

	if !strings.Contains(gotStdout, wantOutputFilepath) {
		t.Errorf("expected to find %q in %q", wantOutputFilepath, gotStdout)
	}

	if strings.Contains(gotStdout, dontWantOutputDailyNote) {
		t.Errorf("expected to not find %q in %q", dontWantOutputDailyNote, gotStdout)
	}

	_, newNoteErr := os.Stat(wantOutputFilepath)
	if newNoteErr != nil {
		t.Errorf("expected %q exist", wantOutputFilepath)
	}
}

func TestNewCmdExistingNote(t *testing.T) {
	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	wantOutputFilepath := filepath.Join(sb, "inbox/hello-world.md")
	wantOutput := fmt.Sprintf("Note with title \"hello-world\" already exists at %s", wantOutputFilepath)

	os.Create(wantOutputFilepath)

	newCmd.Flags().Set("no-open", "true")
	_, gotStderr, gotError := captureOutput(newCmdFunction, newCmd, []string{"Hello World"})

	if gotError.Error() != wantOutput {
		t.Errorf("got %q, want %q", gotError, wantOutput)
	}

	if !strings.Contains(gotStderr, wantOutput) {
		t.Errorf("expected to find %q in %q", wantOutput, gotStderr)
	}

	_, newNoteErr := os.Stat(wantOutputFilepath)
	if newNoteErr != nil {
		t.Errorf("expected %q exist", wantOutputFilepath)
	}
}

func TestNewCmdCreateDailyNote(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	var wantError error
	wantOutputFilepath := filepath.Join(sb, "inbox/hello-world.md")
	wantOutputDailyNote := fmt.Sprintf("Daily note not found; creating a new one: %s/journal/2025-07-13.md", sb)

	newCmd.Flags().Set("no-open", "true")
	gotStdout, _, gotError := captureOutput(newCmdFunction, newCmd, []string{"Hello World"})

	if gotError != wantError {
		t.Errorf("got %q, want %q", gotError, wantError)
	}

	if !strings.Contains(gotStdout, wantOutputFilepath) {
		t.Errorf("expected to find %q in %q", wantOutputFilepath, gotStdout)
	}

	if !strings.Contains(gotStdout, wantOutputDailyNote) {
		t.Errorf("expected to find %q in %q", wantOutputDailyNote, gotStdout)
	}

	_, newNoteErr := os.Stat(wantOutputFilepath)
	if newNoteErr != nil {
		t.Errorf("expected %q exist", wantOutputFilepath)
	}

	_, dailyNoteErr := os.Stat(fmt.Sprintf("%s/journal/2025-07-13.md", sb))
	if dailyNoteErr != nil {
		t.Errorf("expected %q to exist", wantOutputFilepath)
	}
}

func TestDailyCmd(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	var wantError error
	wantOutputFilepath := filepath.Join(sb, "journal/2025-07-13.md")
	dailyCmd.Flags().Set("no-open", "true")
	_, _, gotError := captureOutput(dailyCmdFunction, dailyCmd, []string{})

	if gotError != wantError {
		t.Errorf("got %q, want %q", gotError, wantError)
	}

	_, dailyNoteErr := os.Stat(fmt.Sprintf("%s/journal/2025-07-13.md", sb))
	if dailyNoteErr != nil {
		t.Errorf("expected %q to exist", wantOutputFilepath)
	}
}

func TestDailyCmdAlreadyExists(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	var wantError error
	wantOutputFilepath := filepath.Join(sb, "journal/2025-07-13.md")
	os.Create(wantOutputFilepath)

	dailyCmd.Flags().Set("no-open", "true")
	gotStdout, _, gotError := captureOutput(dailyCmdFunction, dailyCmd, []string{})

	if gotError != wantError {
		t.Errorf("got %q, want %q", gotError, wantError)
	}

	if !strings.Contains(gotStdout, wantOutputFilepath) {
		t.Errorf("expected to find %q in %q", wantOutputFilepath, gotStdout)
	}

	_, dailyNoteErr := os.Stat(fmt.Sprintf("%s/journal/2025-07-13.md", sb))
	if dailyNoteErr != nil {
		t.Errorf("expected %q to exist", wantOutputFilepath)
	}
}

func TestPathCmdExists(t *testing.T) {
	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	var wantError error
	wantOutputFilepath := filepath.Join(sb, "inbox/hello-world.md")
	os.Create(wantOutputFilepath)

	gotStdout, _, gotError := captureOutput(pathCmdFunction, pathCmd, []string{"hello-world"})

	if gotError != wantError {
		t.Errorf("got %q, want %q", gotError, wantError)
	}

	if !strings.Contains(gotStdout, wantOutputFilepath) {
		t.Errorf("expected to find %q in %q", wantOutputFilepath, gotStdout)
	}

	_, statErr := os.Stat(wantOutputFilepath)
	if statErr != nil {
		t.Errorf("expected %q to exist", wantOutputFilepath)
	}
}

func TestPathCmdDoesNotExist(t *testing.T) {
	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	wantOutputFilepath := filepath.Join(sb, "somefolder/hello-world.md")
	wantStderr := "Note with title \"hello-world\" (hello-world) does not exist"

	_, gotStderr, gotError := captureOutput(pathCmdFunction, pathCmd, []string{"hello-world"})

	if gotError.Error() != wantStderr {
		t.Errorf("got %q, want %q", gotError, wantStderr)
	}

	if !strings.Contains(gotStderr, wantStderr) {
		t.Errorf("expected to not find %q in %q", wantStderr, gotStderr)
	}

	_, statErr := os.Stat(wantOutputFilepath)
	if statErr == nil {
		t.Errorf("expected %q not to exist", wantOutputFilepath)
	}
}
