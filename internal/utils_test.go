package internal

import (
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
