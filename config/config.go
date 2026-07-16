package config

import (
	"fmt"
	"os"
	"time"
)

var version string = "development"
var Now = time.Now

// Configuration holds the configuration settings for the CLI
type Configuration struct {
	RootDir       string
	InboxDir      string
	JournalDir    string
	DailyNotePath string
	DayOfWeek     string
	Today         string
	Version       string
}

// GetConfig returns the Conifugration struct by reading environment
// variables and calculating values at runtime
func GetConfig() Configuration {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		userHomeDir = os.Getenv("HOME")
	}

	rootDir := getEnv("SB", fmt.Sprintf("%s/notes", userHomeDir))
	inboxDir := getEnv("SB_INBOX", fmt.Sprintf("%s/inbox", rootDir))

	today := Now().Format("2006-01-02")
	dayOfWeek := Now().Weekday().String()
	journalDir := fmt.Sprintf("%s/journal", rootDir)

	return Configuration{
		RootDir:    rootDir,
		InboxDir:   inboxDir,
		JournalDir: journalDir,
		DayOfWeek:  dayOfWeek,
		Today:      today,
		Version:    version,
	}
}

func getEnv(key string, defaultValue string) string {
	value, varExists := os.LookupEnv(key)

	if varExists {
		return value
	}
	return defaultValue
}
