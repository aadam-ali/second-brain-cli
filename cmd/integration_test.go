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
	var testCases = []struct {
		inputTitle     string
		sanitisedTitle string
	}{
		{"lower case only", "lower case only"},
		{"UPPER CASE ONLY", "UPPER CASE ONLY"},
		{"Mixed Case", "Mixed Case"},
		{"kebab-case", "kebab-case"},
		{"squash----hyphens", "squash----hyphens"},
		{" Leading space", "Leading space"},
		{"-Leading hyphen", "Leading hyphen"},
		{"Trailing hyphen-", "Trailing hyphen"},
		{"1 c4n c0unt 123456789", "1 c4n c0unt 123456789"},
		{"h3llo@world!", "h3llo world"},
		{"Keyboard special keys `~!@#$%^&*()-_=+[{]}\\|;:'\",<.>/?", "Keyboard special keys"},
	}

	for _, tt := range testCases {
		sb := prepareEnvironment()
		defer os.RemoveAll(sb)

		var wantError error
		wantStdoutFilepath := filepath.Join(sb, "inbox", tt.sanitisedTitle+".md")

		newCmd.Flags().Set("no-open", "true")
		gotStdout, _, gotError := captureOutput(newCmdFunction, newCmd, []string{tt.inputTitle})
		_, newNoteErr := os.Stat(wantStdoutFilepath)

		assert.Equal(t, wantError, gotError)
		assert.Contains(t, gotStdout, wantStdoutFilepath)
		assert.NoError(t, newNoteErr)
	}
}

func TestNewCmdExistingNote(t *testing.T) {
	var testCases = []struct {
		inputTitle     string
		sanitisedTitle string
	}{
		{"Hello World", "Hello World"},
		{"HELLO WORLD", "HELLO WORLD"},
		{"hello world", "hello world"},
		{"hello   world", "hello world"},
		{"   hello   world   ", "hello world"},
		{"-_- hello world _-_", "hello world"},
	}

	for _, tt := range testCases {
		sb := prepareEnvironment()
		defer os.RemoveAll(sb)

		wantStdoutFilepath := filepath.Join(sb, "inbox", "Hello World"+".md")
		wantStderr := fmt.Sprintf("Note with title %q already exists at %s", tt.sanitisedTitle, wantStdoutFilepath)

		os.Create(wantStdoutFilepath)

		newCmd.Flags().Set("no-open", "true")
		_, gotStderr, gotError := captureOutput(newCmdFunction, newCmd, []string{tt.inputTitle})
		_, newNoteErr := os.Stat(wantStdoutFilepath)

		assert.Contains(t, gotStderr, wantStderr)
		assert.NoError(t, newNoteErr)
		assert.ErrorContains(t, gotError, wantStderr)
	}
}

func TestNewCmdDateFlag(t *testing.T) {
	var testCases = []struct {
		inputTitle     string
		sanitisedTitle string
	}{
		{"Hello World", "20250713 Hello World"},
	}

	for _, tt := range testCases {
		config.Now = func() time.Time {
			return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
		}

		sb := prepareEnvironment()
		defer os.RemoveAll(sb)

		var wantError error
		wantStdoutFilepath := filepath.Join(sb, "inbox", tt.sanitisedTitle+".md")

		newCmd.Flags().Set("no-open", "true")
		newCmd.Flags().Set("date", "true")
		gotStdout, _, gotError := captureOutput(newCmdFunction, newCmd, []string{tt.inputTitle})
		_, newNoteErr := os.Stat(wantStdoutFilepath)

		assert.Equal(t, wantError, gotError)
		assert.Contains(t, gotStdout, wantStdoutFilepath)
		assert.NoError(t, newNoteErr)
	}
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
	var testCases = []struct {
		filepathOutput string
		filenameInput  string
	}{
		{"inbox/hello-world.md", "hello-world.md"},
		{"inbox/hello world.md", "hello%20world.md"},
	}

	for _, tt := range testCases {
		sb := prepareEnvironment()
		defer os.RemoveAll(sb)

		wantStdoutFilepath := filepath.Join(sb, tt.filepathOutput)
		os.Create(wantStdoutFilepath)

		gotStdout, _, gotError := captureOutput(pathCmdFunction, pathCmd, []string{tt.filenameInput})
		_, statErr := os.Stat(wantStdoutFilepath)

		assert.Contains(t, gotStdout, wantStdoutFilepath)
		assert.NoError(t, gotError)
		assert.NoError(t, statErr)

	}
}

func TestPathCmdWikiLink(t *testing.T) {
	var testCases = []string{
		"hello-world",
		"hello world",
	}

	for _, tt := range testCases {
		sb := prepareEnvironment()
		defer os.RemoveAll(sb)

		wantStdoutFilepath := filepath.Join(sb, tt+".md")
		os.Create(wantStdoutFilepath)

		pathCmd.Flags().Set("wiki", "true")
		gotStdout, _, gotError := captureOutput(pathCmdFunction, pathCmd, []string{tt})
		_, statErr := os.Stat(wantStdoutFilepath)

		assert.Contains(t, gotStdout, wantStdoutFilepath)
		assert.NoError(t, gotError)
		assert.NoError(t, statErr)
	}
}

func TestPathCmdDoesNotExist(t *testing.T) {
	var testCases = []struct {
		filepathOutput string
		filenameInput  string
	}{
		{"inbox/hello-world.md", "hello-world.md"},
		{"inbox/hello world.md", "hello%20world.md"},
	}

	for _, tt := range testCases {
		sb := prepareEnvironment()
		defer os.RemoveAll(sb)

		wantStdoutFilepath := filepath.Join(sb, tt.filepathOutput)
		wantStderr := fmt.Sprintf("Note with title \"%[1]s\" (%[1]s) does not exist", tt.filenameInput)

		_, gotStderr, gotError := captureOutput(pathCmdFunction, pathCmd, []string{tt.filenameInput})
		_, statErr := os.Stat(wantStdoutFilepath)

		assert.Contains(t, gotStderr, wantStderr)
		assert.Error(t, statErr)
		assert.ErrorContains(t, gotError, wantStderr)
	}
}

func TestPathCmdDoesNotExistWikiLink(t *testing.T) {
	var testCases = []string{
		"hello-world",
		"hello world",
	}

	for _, tt := range testCases {
		sb := prepareEnvironment()
		defer os.RemoveAll(sb)

		wantStdoutFilepath := filepath.Join(sb, tt+".md")
		wantStderr := fmt.Sprintf("Note with title \"%[1]s\" (%[1]s) does not exist", tt)

		pathCmd.Flags().Set("wiki", "true")
		_, gotStderr, gotError := captureOutput(pathCmdFunction, pathCmd, []string{tt})
		_, statErr := os.Stat(wantStdoutFilepath)

		assert.Contains(t, gotStderr, wantStderr)
		assert.Error(t, statErr)
		assert.ErrorContains(t, gotError, wantStderr)
	}
}

func TestLinkCmd(t *testing.T) {
	var testCases = []struct {
		destFilenameWithoutExtension string
		urlEncodedDestFilename       string
	}{
		{"bonjour-world", "bonjour-world.md"},
		{"bonjour world", "bonjour%20world.md"},
	}

	for _, tt := range testCases {
		sb := prepareEnvironment()
		defer os.RemoveAll(sb)

		src := filepath.Join(sb, "inbox", "hello-world.md")
		dest := filepath.Join(sb, "journal", tt.destFilenameWithoutExtension+".md")

		os.Create(src)
		os.Create(dest)

		wantOutput := fmt.Sprintf("[%s](../journal/%s)", tt.destFilenameWithoutExtension, tt.urlEncodedDestFilename)
		gotOutput, _, gotError := captureOutput(linkCmdFunction, newCmd, []string{src, dest})

		assert.NoError(t, gotError)
		assert.Equal(t, wantOutput, gotOutput)
	}
}

func TestLinkCmdWikiLink(t *testing.T) {
	var testCases = []string{
		"bonjour-world",
		"bonjour world",
	}

	for _, tt := range testCases {
		sb := prepareEnvironment()
		defer os.RemoveAll(sb)

		src := filepath.Join(sb, "inbox", "hello-world.md")
		dest := filepath.Join(sb, "journal", tt+".md")

		os.Create(src)
		os.Create(dest)

		wantOutput := fmt.Sprintf("[[%s]]", tt)

		linkCmd.Flags().Set("wiki", "true")
		gotOutput, _, gotError := captureOutput(linkCmdFunction, linkCmd, []string{src, dest})

		assert.NoError(t, gotError)
		assert.Equal(t, wantOutput, gotOutput)
	}
}
