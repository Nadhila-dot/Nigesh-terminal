package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Layout struct {
	width  int
	height int
}

func NewLayout() *Layout {
	w, h := getTerminalSize()
	return &Layout{
		width:  w,
		height: h,
	}
}

func (l *Layout) GetWidth() int {
	return l.width
}

func (l *Layout) GetHeight() int {
	return l.height
}

func (l *Layout) Refresh() {
	l.width, l.height = getTerminalSize()
}

func getTerminalSize() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 80, 24 // Default fallback
	}

	parts := strings.Fields(string(out))
	if len(parts) != 2 {
		return 80, 24
	}

	height, _ := strconv.Atoi(parts[0])
	width, _ := strconv.Atoi(parts[1])

	return width, height
}

func HideCursor() {
	fmt.Print("\033[?25l")
}

func ShowCursor() {
	fmt.Print("\033[?25h")
}

func SaveCursor() {
	fmt.Print("\033[s")
}

func RestoreCursor() {
	fmt.Print("\033[u")
}

func ClearLine() {
	fmt.Print("\033[2K")
}

func MoveCursorUp(n int) {
	fmt.Printf("\033[%dA", n)
}

func MoveCursorDown(n int) {
	fmt.Printf("\033[%dB", n)
}
