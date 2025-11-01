package ai

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func StreamResponse(text string) {
	// Apply terminal formatting first
	formatted := FormatTerminal(text)

	// Check if this contains tool calls - don't show them, they'll be executed
	if strings.Contains(text, "<Tool>") {
		// Just show a simple message instead of the ugly tool syntax
		fmt.Printf("\033[90m⚙️  nigesh is working...\033[0m\n")
	} else {
		// Stream the response with typewriter effect, no box
		for _, char := range formatted {
			fmt.Print(string(char))
			time.Sleep(15 * time.Millisecond) // Typewriter speed
		}
		fmt.Println() // Add newline at the end
	}
}

func StreamResponseInBox(text string) {
	maxWidth := 60 // Fixed reasonable width

	// Wrap long lines
	lines := strings.Split(text, "\n")
	var wrappedLines []string

	for _, line := range lines {
		cleanLine := stripAnsiCodes(line)
		if len(cleanLine) <= maxWidth {
			wrappedLines = append(wrappedLines, line)
		} else {
			// Simple word wrapping
			words := strings.Fields(line)
			currentLine := ""
			for _, word := range words {
				if len(stripAnsiCodes(currentLine+" "+word)) <= maxWidth {
					if currentLine == "" {
						currentLine = word
					} else {
						currentLine += " " + word
					}
				} else {
					if currentLine != "" {
						wrappedLines = append(wrappedLines, currentLine)
					}
					currentLine = word
				}
			}
			if currentLine != "" {
				wrappedLines = append(wrappedLines, currentLine)
			}
		}
	}

	boxWidth := maxWidth + 4

	// Top border
	fmt.Printf("\033[36m┌%s┐\033[0m\n", strings.Repeat("─", boxWidth-2))

	// Stream each line with typewriter effect
	for _, line := range wrappedLines {
		fmt.Print("\033[36m│\033[0m ")

		// Stream characters one by one
		for _, char := range line {
			fmt.Print(string(char))
			time.Sleep(15 * time.Millisecond) // Typewriter speed
		}

		// Pad the line to box width (using clean length)
		cleanLine := stripAnsiCodes(line)
		padding := maxWidth - len(cleanLine)
		if padding > 0 {
			fmt.Print(strings.Repeat(" ", padding))
		}

		fmt.Printf(" \033[36m│\033[0m\n")
	}

	// Bottom border
	fmt.Printf("\033[36m└%s┘\033[0m\n", strings.Repeat("─", boxWidth-2))
}

func stripAnsiCodes(s string) string {
	// Remove ANSI escape codes for length calculation
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(s, "")
}

func StreamGeminiResponse(prompt string) (string, error) {
	// Show thinking indicator
	fmt.Printf("\033[90m💭 nigesh is thinking...\033[0m\n")

	// Call Gemini API
	response, err := callGemini(prompt)
	if err != nil {
		return "", err
	}

	// Clear the thinking line
	fmt.Print("\033[A\033[K")

	// Stream the response with typewriter effect
	StreamResponse(response)

	return response, nil
}
