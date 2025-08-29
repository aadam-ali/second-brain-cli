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
	rootDir := "/home/test/SecondBrain"

	want := Configuration{
		RootDir:       rootDir,
		InboxDir:      rootDir + "/inbox",
		JournalDir:    rootDir + "/journal",
		DailyNotePath: rootDir + "/journal/2025-07-13.md",
		DayOfWeek:     "Sunday",
		Today:         "2025-07-13",
		Version:       "development",
	}
	got := GetConfig()

	assert.Equal(t, want, got)
}

func TestGetConfigOverriddenValues(t *testing.T) {
	Now = func() time.Time {
		return time.Date(2025, 7, 13, 20, 0, 0, 0, time.UTC)
	}

	sb := "/home/test/Documents/Notes"
	sbJournal := "/home/test/Documents/Notes/Log"
	sbInbox := "/home/test/Documents/Notes/Entrypoint"

	os.Clearenv()
	os.Setenv("HOME", "/home/test")
	os.Setenv("SB", sb)
	os.Setenv("SB_JOURNAL", sbJournal)
	os.Setenv("SB_INBOX", sbInbox)

	want := Configuration{
		RootDir:       sb,
		InboxDir:      sbInbox,
		JournalDir:    sbJournal,
		DailyNotePath: sbJournal + "/2025-07-13.md",
		DayOfWeek:     "Sunday",
		Today:         "2025-07-13",
		Version:       "development",
	}

	got := GetConfig()

	assert.Equal(t, want, got)
}
