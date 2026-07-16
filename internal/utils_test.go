package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSanitiseTitle(t *testing.T) {
	var testCases = []struct {
		input string
		want  string
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
		got := SanitiseTitle(tt.input)

		assert.Equal(t, tt.want, got)
	}
}

func TestConstructNotePath(t *testing.T) {
	var testCases = []struct {
		path  string
		title string
		want  string
	}{
		{"/home/test", "note", "/home/test/note.md"},
		{"/home/test/", "note", "/home/test/note.md"},
	}

	for _, tt := range testCases {
		got := ConstructNotePath(tt.path, tt.title)

		assert.Equal(t, tt.want, got)
	}
}

func TestCreateNote(t *testing.T) {
	var testCases = []struct {
		filename string
		content  string
	}{
		{"single-line.md", "note"},
		{"newline-characters.md", "# Title\n\nThis is a title"},
		{"single-line-raw-string-literal.md", `note`},
		{"multiline-raw-string-literal.md", `# Another header

Some content

## Another header`},
	}

	for _, tt := range testCases {
		path, err := os.MkdirTemp("", "second-brain-cli")
		if err != nil {
			t.Fatal(err)
		}
		filepath := filepath.Join(path, tt.filename)

		err = CreateNote(filepath, tt.content)
		assert.NoError(t, err)
		got, err := os.ReadFile(filepath)
		assert.NoError(t, err)

		os.RemoveAll(path)

		assert.Equal(t, tt.content, string(got))
	}
}

func TestCheckIfNoteExistsReturnPathWhenExists(t *testing.T) {
	var testCases = []bool{true, false}

	for _, tt := range testCases {
		title := "matching-title"

		rootDir, _, want := createNoteInTempDir(t, title, tt)
		defer os.RemoveAll(rootDir)

		_, got, err := CheckIfNoteExists(rootDir, title+".md")

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	}
}

func TestCheckIfNoteExistsReturnsEmptyStringWhenNotExists(t *testing.T) {
	rootDir, _, _ := createNoteInTempDir(t, "this-one-exists", false)
	defer os.RemoveAll(rootDir)

	_, got, err := CheckIfNoteExists(rootDir, "but-this-one-does-not"+".md")

	assert.NoError(t, err)
	assert.Empty(t, got)
}

func TestCheckIfNoteExistsReturnsEmptyStringWhenDirNotExists(t *testing.T) {
	rootDir := fmt.Sprintf("/tmp/this-dir-does-not-exist-%d", time.Now().UnixNano())

	_, got, err := CheckIfNoteExists(rootDir, "but-this-one-does-not"+".md")

	assert.NoError(t, err)
	assert.Empty(t, got)
}

func TestOpenFileInEditorWithMissingBinary(t *testing.T) {
	err := OpenFileInEditor("nonexistent-binary-abc123", "/tmp", "/tmp/test.md")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrEditorLaunch)
}

func TestOpenFileInEditorWithFailingBinary(t *testing.T) {
	err := OpenFileInEditor("false", "/tmp", "/tmp/test.md")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrEditorLaunch)
}

func TestCreateNoteMkdirAllError(t *testing.T) {
	err := CreateNote("/dev/null/test.md", "content")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "create dir")
}

func TestCreateNoteWriteError(t *testing.T) {
	err := CreateNote("/dev/full", "content")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "write note")
}

func TestCheckIfNoteExistsReturnsBool(t *testing.T) {

	var testCases = []struct {
		createTitle   string
		expectedTitle string
		nested        bool
		want          bool
	}{
		{"this-exists", "this-exists", false, true},
		{"this-exists", "this-exists", true, true},
		{"this-one-exists", "but-this-one-does-not", false, false},
		{"this-nested-one-exists", "but-this-one-does-not", true, false},
	}

	for _, tt := range testCases {

		rootDir, _, _ := createNoteInTempDir(t, tt.createTitle, tt.nested)
		defer os.RemoveAll(rootDir)

		got, _, err := CheckIfNoteExists(rootDir, tt.expectedTitle+".md")

		assert.NoError(t, err)
		assert.Equal(t, tt.want, got)
	}
}
