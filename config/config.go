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
	userHomeDir, _ := os.UserHomeDir()

	rootDir := getEnv("SB", fmt.Sprintf("%s/SecondBrain", userHomeDir))
	inboxDir := getEnv("SB_INBOX", fmt.Sprintf("%s/inbox", rootDir))
	journalDir := getEnv("SB_JOURNAL", fmt.Sprintf("%s/journal", rootDir))

	today := Now().Format("2006-01-02")
	dayOfWeek := Now().Weekday().String()
	dailyNotePath := fmt.Sprintf("%s/%s.md", journalDir, today)

	return Configuration{
		RootDir:       rootDir,
		InboxDir:      inboxDir,
		JournalDir:    journalDir,
		DailyNotePath: dailyNotePath,
		DayOfWeek:     dayOfWeek,
		Today:         today,
		Version:       version,
	}
}

func getEnv(key string, defaultValue string) string {
	value, varExists := os.LookupEnv(key)

	if varExists {
		return value
	}
	return defaultValue
}
