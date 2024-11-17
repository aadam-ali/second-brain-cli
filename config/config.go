package config

import (
	"fmt"
	"os"
	"time"
)

var version string = "development"

type Configuration struct {
	RootDir       string
	InboxDir      string
	JournalDir    string
	DailyNotePath string
	Yesterday     string
	Today         string
	Tomorrow      string
	Version       string
}

func GetConfig() Configuration {
	userHomeDir, _ := os.UserHomeDir()

	rootDir := getEnv("SB", fmt.Sprintf("%s/SecondBrain", userHomeDir))
	inboxDir := getEnv("SB_INBOX", fmt.Sprintf("%s/inbox", rootDir))
	journalDir := getEnv("SB_JOURNAL", fmt.Sprintf("%s/journal", rootDir))

	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	dailyNotePath := fmt.Sprintf("%s/%s.md", journalDir, today)

	return Configuration{
		RootDir:       rootDir,
		InboxDir:      inboxDir,
		JournalDir:    journalDir,
		DailyNotePath: dailyNotePath,
		Yesterday:     yesterday,
		Today:         today,
		Tomorrow:      tomorrow,
		Version:       version,
	}
}

func getEnv(key string, defaultValue string) string {
	value, varExists := os.LookupEnv(key)

	if varExists == true {
		return value
	}
	return defaultValue
}
