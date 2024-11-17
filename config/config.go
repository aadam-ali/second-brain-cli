package config

import (
	"fmt"
	"os"
)

type Configuration struct {
	RootDir  string
	InboxDir string
}

func GetConfig() Configuration {
	userHomeDir, _ := os.UserHomeDir()

	rootDir := getEnv("SB", fmt.Sprintf("%s/SecondBrain", userHomeDir))
	inboxDir := getEnv("SB_INBOX", fmt.Sprintf("%s/inbox", rootDir))

	return Configuration{
		RootDir:  rootDir,
		InboxDir: inboxDir,
	}
}

func getEnv(key string, defaultValue string) string {
	value, varExists := os.LookupEnv(key)

	if varExists == true {
		return value
	}
	return defaultValue
}
