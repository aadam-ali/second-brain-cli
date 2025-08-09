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

	gotOutput, gotError := captureStdout(versionCmdFunction, versionCmd, []string{})

	if gotError != wantError {
		t.Errorf("got %q, want %q", gotError, wantError)
	}

	if gotOutput != wantOutput {
		t.Errorf("got %q, want %q", gotOutput, wantOutput)
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
	gotOutput, gotError := captureStdout(newCmdFunction, newCmd, []string{"Hello World"})

	if gotError != wantError {
		t.Errorf("got %q, want %q", gotError, wantError)
	}

	if !strings.Contains(gotOutput, wantOutputFilepath) {
		t.Errorf("expected to find %q in %q", wantOutputFilepath, gotOutput)
	}

	if strings.Contains(gotOutput, dontWantOutputDailyNote) {
		t.Errorf("expected to not find %q in %q", dontWantOutputDailyNote, gotOutput)
	}

	_, newNoteErr := os.Stat(wantOutputFilepath)
	if newNoteErr != nil {
		t.Errorf("expected %q exist", wantOutputFilepath)
	}
}

func TestNewCmdExistingNote(t *testing.T) {
	sb := prepareEnvironment()
	defer os.RemoveAll(sb)

	var wantError error
	wantOutputFilepath := filepath.Join(sb, "inbox/hello-world.md")
	os.Create(wantOutputFilepath)
	wantOutput := fmt.Sprintf("Note already exists: %s", wantOutputFilepath)

	newCmd.Flags().Set("no-open", "true")
	gotOutput, gotError := captureStdout(newCmdFunction, newCmd, []string{"Hello World"})

	if gotError != wantError {
		t.Errorf("got %q, want %q", gotError, wantError)
	}

	if !strings.Contains(gotOutput, wantOutput) {
		t.Errorf("expected to find %q in %q", wantOutput, gotOutput)
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
	gotOutput, gotError := captureStdout(newCmdFunction, newCmd, []string{"Hello World"})

	if gotError != wantError {
		t.Errorf("got %q, want %q", gotError, wantError)
	}

	if !strings.Contains(gotOutput, wantOutputFilepath) {
		t.Errorf("expected to find %q in %q", wantOutputFilepath, gotOutput)
	}

	if !strings.Contains(gotOutput, wantOutputDailyNote) {
		t.Errorf("expected to find %q in %q", wantOutputDailyNote, gotOutput)
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
	_, gotError := captureStdout(dailyCmdFunction, dailyCmd, []string{})

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
	wantOutput := fmt.Sprintf("Note already exists: %s", wantOutputFilepath)
	os.Create(wantOutputFilepath)

	dailyCmd.Flags().Set("no-open", "true")
	gotOutput, gotError := captureStdout(dailyCmdFunction, dailyCmd, []string{})

	if gotError != wantError {
		t.Errorf("got %q, want %q", gotError, wantError)
	}

	if !strings.Contains(gotOutput, wantOutput) {
		t.Errorf("expected to find %q in %q", wantOutputFilepath, gotOutput)
	}

	_, dailyNoteErr := os.Stat(fmt.Sprintf("%s/journal/2025-07-13.md", sb))
	if dailyNoteErr != nil {
		t.Errorf("expected %q to exist", wantOutputFilepath)
	}
}
