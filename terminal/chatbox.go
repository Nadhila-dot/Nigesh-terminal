package terminal

import (
	"fmt"
	"strings"
)

type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

type Message struct {
	Role    MessageRole
	Content string
}

type ChatBox struct {
	messages []Message
	x        int
	y        int
	width    int
	height   int
	scroll   int
}

func NewChatBox(width, height int) *ChatBox {
	return &ChatBox{
		messages: []Message{},
		x:        0,
		y:        0,
		width:    width,
		height:   height,
		scroll:   0,
	}
}

func (cb *ChatBox) SetPosition(x, y int) {
	cb.x = x
	cb.y = y
}

func (cb *ChatBox) AddMessage(role MessageRole, content string) {
	cb.messages = append(cb.messages, Message{Role: role, Content: content})
}

func (cb *ChatBox) StreamMessage(role MessageRole, content string) {
	if len(cb.messages) > 0 && cb.messages[len(cb.messages)-1].Role == role {
		cb.messages[len(cb.messages)-1].Content += content
	} else {
		cb.AddMessage(role, content)
	}
	cb.Render()
}

func (cb *ChatBox) ScrollUp() {
	if cb.scroll > 0 {
		cb.scroll--
	}
}

func (cb *ChatBox) ScrollDown() {
	cb.scroll++
}

func (cb *ChatBox) Render() {
	lines := cb.renderMessages()

	moveCursor(cb.x, cb.y)

	// Top border
	fmt.Print("╔")
	for i := 0; i < cb.width-2; i++ {
		fmt.Print("═")
	}
	fmt.Println("╗")

	// Content area
	start := cb.scroll
	end := start + cb.height - 2
	if end > len(lines) {
		end = len(lines)
	}
	if start > len(lines) {
		start = len(lines)
	}

	for i := 0; i < cb.height-2; i++ {
		idx := start + i
		if idx < len(lines) {
			fmt.Printf("║ %-*s ║\n", cb.width-4, truncate(lines[idx], cb.width-4))
		} else {
			fmt.Printf("║ %-*s ║\n", cb.width-4, "")
		}
	}

	// Bottom border
	fmt.Print("╚")
	for i := 0; i < cb.width-2; i++ {
		fmt.Print("═")
	}
	fmt.Println("╝")
}

func (cb *ChatBox) renderMessages() []string {
	var lines []string

	for _, msg := range cb.messages {
		// Add role header
		var header string
		switch msg.Role {
		case RoleUser:
			header = "┌─ You ─────────────────────────────────"
		case RoleAssistant:
			header = "┌─ Assistant ───────────────────────────"
		case RoleSystem:
			header = "┌─ System ──────────────────────────────"
		}
		lines = append(lines, header)

		// Wrap content
		wrapped := wrapText(msg.Content, cb.width-6)
		for _, line := range wrapped {
			lines = append(lines, "│ "+line)
		}

		lines = append(lines, "└───────────────────────────────────────")
		lines = append(lines, "")
	}

	return lines
}

func wrapTexter(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		if len(currentLine) == 0 {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
