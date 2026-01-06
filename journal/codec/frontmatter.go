package codec

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/never00rei/a7/utils"
)

const TimestampLayout = "2006-01-02_15-04"

type FrontMatter struct {
	Title     string
	Created   time.Time
	Updated   time.Time
	Encrypted bool
	WordCount int
}

func BuildFilename(title string, created time.Time) string {
	sanitizedTitle := utils.SanitizeSpecialChars(title)
	return fmt.Sprintf("%s_%s.md", created.Format(TimestampLayout), sanitizedTitle)
}

func RenderContent(title, body string, created, updated time.Time, encrypted bool, wordCount int) string {
	return fmt.Sprintf(
		"---\n"+
			"title: %s\n"+
			"created: %s\n"+
			"updated: %s\n"+
			"encrypted: %t\n"+
			"word_count: %d\n"+
			"---\n\n%s",
		title,
		created.Format(time.RFC3339),
		updated.Format(time.RFC3339),
		encrypted,
		wordCount,
		body,
	)
}

func ParseFrontMatter(content string) (FrontMatter, string) {
	matter := FrontMatter{WordCount: -1}
	lines := strings.Split(content, "\n")
	if len(lines) < 3 {
		return matter, content
	}
	if strings.TrimSpace(lines[0]) != "---" {
		return matter, content
	}

	end := -1
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "---" {
			end = i
			break
		}
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		switch key {
		case "title":
			matter.Title = value
		case "created":
			matter.Created = ParseTimestamp(value)
		case "updated":
			matter.Updated = ParseTimestamp(value)
		case "encrypted":
			matter.Encrypted = strings.EqualFold(value, "true")
		case "word_count":
			if parsed, err := strconv.Atoi(value); err == nil {
				matter.WordCount = parsed
			}
		}
	}
	if end == -1 {
		return FrontMatter{WordCount: -1}, content
	}

	body := strings.Join(lines[end+1:], "\n")
	body = strings.TrimPrefix(body, "\n")
	return matter, body
}

func ParseHeader(content string) (string, time.Time, string) {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return "", time.Time{}, content
	}

	line := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(line, "# ") {
		return "", time.Time{}, content
	}

	line = strings.TrimPrefix(line, "# ")
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return "", time.Time{}, content
	}

	created, err := time.Parse(TimestampLayout, parts[0])
	if err != nil {
		return strings.TrimSpace(parts[1]), time.Time{}, strings.Join(lines[1:], "\n")
	}

	body := strings.Join(lines[1:], "\n")
	body = strings.TrimPrefix(body, "\n")
	return strings.TrimSpace(parts[1]), created, body
}

func ParseTimestamp(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	if ts, err := time.Parse(time.RFC3339, value); err == nil {
		return ts
	}
	if ts, err := time.Parse(TimestampLayout, value); err == nil {
		return ts
	}
	return time.Time{}
}

func CountWords(content string) int {
	fields := strings.Fields(content)
	return len(fields)
}
