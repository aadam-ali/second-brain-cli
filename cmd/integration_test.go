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

func TestDailyCmd(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	wantStdout := `### 2025-07-13 Sunday
@ <location>

* Todo:

* Notes:
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
		sb := prepareEnvironment(true)
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
		sb := prepareEnvironment(true)
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
		sb := prepareEnvironment(true)
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
		sb := prepareEnvironment(true)
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
