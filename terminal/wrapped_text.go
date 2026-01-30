package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"
)

type WrappedText struct {
	content string
	width   int
	label   string
}

func getTerminalWidth() int {
	width := 80
	if w, ok := os.LookupEnv("COLUMNS"); ok {
		fmt.Sscanf(w, "%d", &width)
	} else {
		out, err := execCommand("tput", "cols")
		if err == nil {
			fmt.Sscanf(strings.TrimSpace(out), "%d", &width)
		}
	}
	return width
}

func execCommand(name string, arg string) (string, error) {
	cmd := exec.Command(name, arg)
	out, err := cmd.Output()
	return string(out), err
}

func CreateWrappedText(label string) *WrappedText {
	width := getTerminalWidth() - 8
	if width < 20 {
		width = 20
	}
	return &WrappedText{width: width, label: label}
}

func (w *WrappedText) Set(text string) {
	w.content = text
	w.printBox()
}

func (w *WrappedText) printBox() {
	lines := wrapMultilineText(w.content, w.width-4)
	longest := maxLineLength(lines)
	label := w.label
	if label == "" {
		label = "You"
	}
	labelPart := fmt.Sprintf("─ %s ", label)
	borderLen := max(1, longest-len(labelPart)+2)
	// Top border
	fmt.Printf("╭%s%s╮\n", labelPart, strings.Repeat("─", borderLen)+ "──")
	// Content
	for _, line := range lines {
		fmt.Printf("│ %-*s │\n", longest, line)
	}
	// Bottom border
	fmt.Printf("╰%s╯\n", strings.Repeat("─", longest+2))
}

// New function to handle multi-line and wrapping
func wrapMultilineText(text string, width int) []string {
	var result []string
	for _, rawLine := range strings.Split(text, "\n") {
		rawLine = strings.ReplaceAll(rawLine, "\r", "")
		wrapped := wrapText(rawLine, width)
		result = append(result, wrapped...)
	}
	return result
}

func (w *WrappedText) Exit() {
	// Optional: could clear state or print a closing message
}

func wrapText(text string, width int) []string {
	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}
	line := words[0]
	for _, word := range words[1:] {
		if utf8.RuneCountInString(line)+1+utf8.RuneCountInString(word) > width {
			lines = append(lines, line)
			line = word
		} else {
			line += " " + word
		}
	}
	lines = append(lines, line)
	return lines
}

func maxLineLength(lines []string) int {
	max := 0
	for _, l := range lines {
		if utf8.RuneCountInString(l) > max {
			max = utf8.RuneCountInString(l)
		}
	}
	return max
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
