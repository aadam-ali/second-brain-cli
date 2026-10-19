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

func TestTodoCmd(t *testing.T) {
	Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}
	title := "Go to the shops"

	sb := prepareEnvironment(true, []string{})
	defer os.RemoveAll(sb)

	var expectedError error
	expectedStdoutFilepath := filepath.Join(sb, "todo", "2025-07-13 "+title+".md")

	actualStdout, _, actualError := captureOutput(run, []string{"todo", title})

	todoLines := readLines(expectedStdoutFilepath)
	todoListLines := readLines(filepath.Join(sb, "todo", "TODO.md"))

	assert.Equal(t, expectedError, actualError)
	assert.Contains(t, expectedStdoutFilepath+"\n", actualStdout)

	assert.FileExists(t, expectedStdoutFilepath)
	assert.Equal(t, len(todoLines), 5)
	assert.Equal(t, todoLines[0:5], []string{
		"---",
		"title: 'Go to the shops'",
		"date: 2025-07-13",
		"tags: []",
		"---",
	})

	assert.FileExists(t, filepath.Join(sb, "todo", "TODO.md"))
	assert.Equal(t, len(todoListLines), 3)
	assert.Equal(t, todoListLines[0:3], []string{
		"# TODO",
		"",
		"- [ ] [2025-07-13 Go to the shops](2025-07-13%20Go%20to%20the%20shops.md)",
	})
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

func TestVersionCmd(t *testing.T) {
	sb := prepareEnvironment(true, []string{})
	os.RemoveAll(sb)

	expected := "sb development\n"

	actualStdout, _, actualErr := captureOutput(run, []string{"version"})

	assert.Equal(t, expected, actualStdout)
	assert.NoError(t, actualErr)
}

func TestHelpCmd(t *testing.T) {
	sb := prepareEnvironment(true, []string{})
	os.RemoveAll(sb)

	expected := `Usage: sb COMMAND [ARGS]...

Commands:

  link [src] [dest]       Print a relative markdown link from src to dest
  new [title]             Create a new note
  path [dir] [relpath]    Print the absolute path to a note from dir (if one exists)
  todo [title]            Create a new todo note (also updates $NOTES_TODO/TODO.md)
  version                 Show the current version
  help                    Output this message
`

	actualStdout, _, actualErr := captureOutput(run, []string{"help"})

	assert.Equal(t, expected, actualStdout)
	assert.NoError(t, actualErr)
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

func TestPathCmdExists(t *testing.T) {
	var testCases = []struct {
		filepathOutput string
		filenameInput  string
	}{
		{"hello-world.md", "hello-world.md"},
		{"hello world.md", "hello%20world.md"},
		{"todo/hello world.md", "todo/hello%20world.md"},
	}

	for _, tt := range testCases {
		sb := prepareEnvironment(true, []string{})
		defer os.RemoveAll(sb)

		expectedStdoutFilepath := filepath.Join(sb, tt.filepathOutput)
		os.Create(expectedStdoutFilepath)

		actualStdout, _, actualError := captureOutput(run, []string{"path", sb, tt.filenameInput})

		assert.Equal(t, expectedStdoutFilepath, actualStdout)
		assert.NoError(t, actualError)
		assert.FileExists(t, expectedStdoutFilepath)
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
		sb := prepareEnvironment(true, []string{})
		defer os.RemoveAll(sb)

		expectedStdoutFilepath := filepath.Join(sb, tt.filepathOutput)

		actualStdout, _, actualError := captureOutput(run, []string{"path", sb, tt.filenameInput})

		assert.ErrorIs(t, NoteDoesntExistError, actualError)
		assert.Equal(t, "note doesn't exist\n", actualStdout)
		assert.NoFileExists(t, expectedStdoutFilepath)
	}
}
