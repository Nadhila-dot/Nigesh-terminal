package main
import (
	"fmt"
	"nadhi/do-t/terminal"
	"time"
)

func Panel_test() {
	terminal.ClearScreenFull()
	terminal.HideCursor()
	defer terminal.ShowCursor()

	layout := terminal.NewLayout()

	// Create multiple panels
	userPanel := terminal.NewPanel("User", 40, 8)
	userPanel.SetPosition(2, 2)
	userPanel.SetStyle(terminal.UserStyle)

	assistantPanel := terminal.NewPanel("Assistant", 40, 8)
	assistantPanel.SetPosition(2, 11)
	assistantPanel.SetStyle(terminal.AssistantStyle)

	statusPanel := terminal.NewPanel("Status", layout.GetWidth()-4, 5)
	statusPanel.SetPosition(2, 20)

	// Simulate conversation
	userPanel.AppendContent("Hello! Can you help me")
	userPanel.AppendContent("with something?")
	userPanel.Render()

	time.Sleep(500 * time.Millisecond)

	// Simulate streaming response
	response := []string{
		"Of course! I'd be",
		"happy to help you.",
		"What do you need",
		"assistance with?",
	}

	for _, line := range response {
		assistantPanel.AppendContent(line)
		assistantPanel.Render()
		time.Sleep(300 * time.Millisecond)
	}

	statusPanel.AppendContent("Messages: 2")
	statusPanel.AppendContent(fmt.Sprintf("Terminal: %dx%d", layout.GetWidth(), layout.GetHeight()))
	statusPanel.AppendContent("Status: Active")
	statusPanel.Render()

	// Keep display for a moment
	time.Sleep(3 * time.Second)
}
