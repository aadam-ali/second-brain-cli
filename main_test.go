package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

	createDirectoriesFlags := []bool{true, false}

	for _, tt := range testCases {
		for _, createDirectoriesFlag := range createDirectoriesFlags {
			Now = func() time.Time {
				return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
			}

			sb := prepareEnvironment(createDirectoriesFlag, []string{})
			defer os.RemoveAll(sb)

			var expectedError error
			expectedStdoutFilepath := filepath.Join(sb, "2025-07-13 "+tt.sanitisedTitle+".md")

			actualStdout, _, actualError := captureOutput(run, []string{"new", tt.inputTitle})
			captureOutput(run, []string{"new", tt.inputTitle})
			noteLines := readLines(expectedStdoutFilepath)

			assert.Equal(t, expectedError, actualError)
			assert.Contains(t, expectedStdoutFilepath+"\n", actualStdout)

			assert.FileExists(t, expectedStdoutFilepath)
			assert.Equal(t, len(noteLines), 5)
			assert.Equal(t, noteLines[0:5], []string{
				"---",
				fmt.Sprintf("title: '%s'", tt.inputTitle),
				"date: 2025-07-13",
				"tags: []",
				"---",
			})
		}
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
		sb := prepareEnvironment(true, []string{})
		defer os.RemoveAll(sb)

		src := filepath.Join(sb, "hello-world.md")
		dest := filepath.Join(sb, "todo", tt.destFilenameWithoutExtension+".md")

		os.Create(src)
		os.Create(dest)

		expectedOutput := fmt.Sprintf("[%s](todo/%s)", tt.destFilenameWithoutExtension, tt.urlEncodedDestFilename)
		actualOutput, _, acutalError := captureOutput(run, []string{"link", src, dest})

		assert.NoError(t, acutalError)
		assert.Equal(t, expectedOutput, actualOutput)
	}
}

func TestLinkCmdNested(t *testing.T) {
	sb := prepareEnvironment(true, []string{"1/2/3", "a/b"})
	defer os.RemoveAll(sb)

	src := filepath.Join(sb, "1/2/3", "hello-world.md")
	dest := filepath.Join(sb, "a/b", "file.md")

	os.Create(src)
	os.Create(dest)

	expectedOutput := fmt.Sprintf("[%s](../../../a/b/%s)", "file", "file.md")
	actualOutput, _, acutalError := captureOutput(run, []string{"link", src, dest})

	assert.NoError(t, acutalError)
	assert.Equal(t, expectedOutput, actualOutput)
}

// TEST HELPERS

func captureOutput(fn func(args []string) error, args []string) (string, string, error) {
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	var bufOut bytes.Buffer
	r1, w1, _ := os.Pipe()
	os.Stdout = w1

	var bufErr bytes.Buffer
	r2, w2, _ := os.Pipe()
	os.Stderr = w2

	err := fn(args)

	w1.Close()
	os.Stdout = originalStdout
	io.Copy(&bufOut, r1)

	w2.Close()
	os.Stderr = originalStderr
	io.Copy(&bufErr, r2)

	return bufOut.String(), bufErr.String(), err
}

func prepareEnvironment(createDirectories bool, additionalDirs []string) string {
	sb := fmt.Sprintf("/tmp/second-brain-cli-%d", time.Now().UnixNano())

	if createDirectories {
		os.Mkdir(sb, 0700)
		os.Mkdir(filepath.Join(sb, "todo"), 0700)

		for _, v := range additionalDirs {
			os.MkdirAll(filepath.Join(sb, v), 0700)
		}
	}

	os.Clearenv()
	os.Setenv("NOTES", sb)

	return sb
}

func readLines(path string) []string {
	contents, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimSuffix(string(contents), "\n"), "\n")

	return lines
}
