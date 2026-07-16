package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvExists(t *testing.T) {
	want := "CustomValue"

	os.Setenv("SB_TEST_VAR", want)

	got := getEnv("SB_TEST_VAR", "default")

	assert.Equal(t, want, got)
}

func TestGetEnvDoesNotExist(t *testing.T) {
	want := "DefaultValue"

	os.Unsetenv("SB_TEST_VAR")

	got := getEnv("SB_TEST_VAR", want)

	assert.Equal(t, want, got)
}

func TestGetConfigDefaultValues(t *testing.T) {
	Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	os.Clearenv()
	os.Setenv("HOME", "/home/test")
	rootDir := "/home/test/notes"

	want := Configuration{
		RootDir:    rootDir,
		InboxDir:   rootDir + "/inbox",
		JournalDir: rootDir + "/journal",
		DayOfWeek:  "Sunday",
		Today:      "2025-07-13",
		Version:    "development",
		Editor:     "vim",
	}
	got := GetConfig()

	assert.Equal(t, want, got)
}

func TestResolveEditorDefaults(t *testing.T) {
	os.Unsetenv("VISUAL")
	os.Unsetenv("EDITOR")

	got := resolveEditor()
	assert.Equal(t, "vim", got)
}

func TestResolveEditorUsesVISUAL(t *testing.T) {
	os.Setenv("VISUAL", "nano")
	defer os.Unsetenv("VISUAL")

	got := resolveEditor()
	assert.Equal(t, "nano", got)
}

func TestResolveEditorFallsBackToEDITOR(t *testing.T) {
	os.Unsetenv("VISUAL")
	os.Setenv("EDITOR", "emacs")
	defer os.Unsetenv("EDITOR")

	got := resolveEditor()
	assert.Equal(t, "emacs", got)
}

func TestResolveEditorVISUALOverridesEDITOR(t *testing.T) {
	os.Setenv("VISUAL", "code --wait")
	os.Setenv("EDITOR", "vim")
	defer os.Unsetenv("VISUAL")
	defer os.Unsetenv("EDITOR")

	got := resolveEditor()
	assert.Equal(t, "code --wait", got)
}

func TestGetConfigOverriddenValues(t *testing.T) {
	Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	sb := "/home/test/Documents/Notes"
	sbInbox := "/home/test/Documents/Notes/Entrypoint"

	os.Clearenv()
	os.Setenv("HOME", "/home/test")
	os.Setenv("SB", sb)
	os.Setenv("SB_INBOX", sbInbox)

	want := Configuration{
		RootDir:    sb,
		InboxDir:   sbInbox,
		JournalDir: sb + "/journal",
		DayOfWeek:  "Sunday",
		Today:      "2025-07-13",
		Version:    "development",
		Editor:     "vim",
	}

	got := GetConfig()

	assert.Equal(t, want, got)
}
