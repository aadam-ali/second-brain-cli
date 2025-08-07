package internal

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

func TitleToKebabCase(title string) string {
	title = strings.ToLower(title)

	title = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(title, "-")
	title = regexp.MustCompile(`^-+|-+$`).ReplaceAllString(title, "")

	return title
}

func ConstructNotePath(dir string, title string) string {
	titleWithExtension := fmt.Sprint(title, ".md")
	return filepath.Join(dir, titleWithExtension)
}
