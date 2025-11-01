package ai

import (
	"fmt"
	"regexp"
	"strings"
)

func FormatTerminal(text string) string {
	// Bold **text**
	boldRe := regexp.MustCompile(`\*\*(.*?)\*\*`)
	text = boldRe.ReplaceAllString(text, "\033[1m$1\033[0m")

	// Underline __text__
	underlineRe := regexp.MustCompile(`__(.*?)__`)
	text = underlineRe.ReplaceAllString(text, "\033[4m$1\033[0m")

	// Headers ###
	headerRe := regexp.MustCompile(`(?m)^### (.+)$`)
	text = headerRe.ReplaceAllStringFunc(text, func(match string) string {
		header := strings.TrimPrefix(match, "### ")
		return fmt.Sprintf("\033[1;36m%s\033[0m", header)
	})

	// Code blocks ```
	codeRe := regexp.MustCompile("```([\\s\\S]*?)```")
	text = codeRe.ReplaceAllString(text, "\033[90m$1\033[0m")

	// Inline code `text`
	inlineCodeRe := regexp.MustCompile("`([^`]+)`")
	text = inlineCodeRe.ReplaceAllString(text, "\033[93m$1\033[0m")

	return text
}
