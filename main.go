package main

import (
	"errors"
	"fmt"
	"github.com/aadam-ali/second-brain-cli/cmd"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// GLOBALS

var Now = time.Now

var NotImplementedError = errors.New("command isn't implemented")

// ENTRYPOINT

func main() {
	if err := run(os.Args[1:]); err != nil {
		if err == NotImplementedError {
			cmd.Execute()
		}

		os.Exit(1)
	}
}

func run(args []string) error {
	var err error
	var str string

	cfg := loadConfig()
	if err := os.MkdirAll(cfg.RootDir, 0700); err != nil {
		return err
	}

	switch args[0] {
	case "new", "n":
		if len(args) > 1 {
			str, err = new(cfg, args[1:])
		}
	case "link", "l":
		if len(args) == 3 {
			str, err = link(cfg, args[1:])
		}
	default:
		return nil
	}

	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	} else {
		fmt.Print(str)
		return nil
	}
}

// COMMANDS

func new(cfg Configuration, args []string) (string, error) {
	var title string

	if len(args) > 1 {
		title = strings.Join(args, " ")
	} else {
		title = args[0]
	}
	santisedTitle := sanitiseTitle(title)
	fileName := cfg.Date + " " + santisedTitle + ".md"
	filePath := filepath.Join(cfg.RootDir, fileName)

	createdFilePath, err := createNote(filePath, title, cfg.Date)
	if err != nil {
		return "", err
	}

	return createdFilePath + "\n", nil
}

func link(cfg Configuration, args []string) (string, error) {
	src, err := filepath.Abs(args[0])
	if err != nil {
		return "", err
	}

	dest, err := filepath.Abs(args[1])
	if err != nil {
		return "", err
	}

	destFileName := filepath.Base(dest)
	destTitle, _ := strings.CutSuffix(destFileName, ".md")

	relPath, err := filepath.Rel(filepath.Dir(src), dest)
	if err != nil {
		return "", err
	}

	urlEncodedDestFileName := url.PathEscape(destFileName)
	link := fmt.Sprintf("[%s](%s)", destTitle, filepath.Join(filepath.Dir(relPath), urlEncodedDestFileName))

	return link, nil
}

// HELPER

type Configuration struct {
	RootDir string
	TodoDir string
	Date    string
}

func loadConfig() Configuration {
	userHomeDir, _ := os.UserHomeDir()

	rootDir := getEnv("NOTES", filepath.Join(userHomeDir, "notes"))
	date := Now().Format("2006-01-02")

	return Configuration{
		RootDir: rootDir,
		Date:    date,
	}
}

func sanitiseTitle(title string) string {
	// Remove non alphanumeric and separators (_ / -)
	sanitisedTitle := regexp.MustCompile(`[^A-Za-z0-9-_]+`).ReplaceAllString(title, " ")

	// Squash consecutive whitespaces
	sanitisedTitle = regexp.MustCompile(`\s+`).ReplaceAllString(sanitisedTitle, " ")

	// Remove leading whitespaces and separators
	sanitisedTitle = regexp.MustCompile(`^[-_\s]+|[-_\s]+$`).ReplaceAllString(sanitisedTitle, "")

	return sanitisedTitle
}

func createFile(filePath string, content []byte) (string, error) {
	f, createErr := os.Create(filePath)
	os.Chmod(f.Name(), 0600)
	if createErr != nil {
		return "", createErr
	}

	_, writeErr := f.Write(content)
	if writeErr != nil {
		return "", writeErr
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	return absFilePath, nil
}

func createNote(filePath string, title string, date string) (string, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(filePath); err == nil {
		return absFilePath + "\n", nil
	}

	createdFilePath, err := createFile(absFilePath, fmt.Appendf(nil, `---
title: '%s'
date: %s
tags: []
---`, title, date))
	if err != nil {
		return "", err
	}

	return createdFilePath, nil
}

func getEnv(key string, defaultValue string) string {
	value, varExists := os.LookupEnv(key)

	if varExists {
		return value
	}
	return defaultValue
}
