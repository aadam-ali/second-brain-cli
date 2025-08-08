package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTitleToKebabCase(t *testing.T) {
	var testCases = []struct {
		input string
		want  string
	}{
		{"lower case only", "lower-case-only"},
		{"UPPER CASE ONLY", "upper-case-only"},
		{"Mixed Case", "mixed-case"},
		{"kebab-case", "kebab-case"},
		{"squash----hyphens", "squash-hyphens"},
		{" Leading space", "leading-space"},
		{"-Leading hyphen", "leading-hyphen"},
		{"Trailing hyphen-", "trailing-hyphen"},
		{"1 c4n c0unt 123456789", "1-c4n-c0unt-123456789"},
		{"h3llo@world!", "h3llo-world"},
		{"Keyboard special keys `~!@#$%^&*()-_=+[{]}\\|;:'\",<.>/?", "keyboard-special-keys"},
	}

	for _, tt := range testCases {
		got := TitleToKebabCase(tt.input)

		if got != tt.want {
			t.Errorf("got '%s', want '%s'", got, tt.want)
		}

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

		if got != tt.want {
			t.Errorf("got '%s', want '%s'", got, tt.want)
		}

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
		path, _ := os.MkdirTemp("", "second-brain-cli")
		filepath := filepath.Join(path, tt.filename)

		CreateNote(filepath, tt.content)
		got, _ := os.ReadFile(filepath)

		os.RemoveAll(path)

		if string(got) != tt.content {
			t.Errorf("got '%s', want '%s'", got, tt.content)
		}

	}
}
