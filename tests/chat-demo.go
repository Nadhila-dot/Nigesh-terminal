package main

import (
	"bufio"
	"fmt"
	"nadhi/do-t/terminal"
	"os"
	"strings"
	"time"
)

func Chat_demo() {
	// Setup terminal
	terminal.ClearScreenFull()
	terminal.HideCursor()
	defer terminal.ShowCursor()

	layout := terminal.NewLayout()

	// Create chat interface
	chatBox := terminal.NewChatBox(layout.GetWidth()-4, layout.GetHeight()-8)
	chatBox.SetPosition(2, 1)

	inputBox := terminal.NewInputBox("Type your message...")
	inputBox.SetPosition(2, layout.GetHeight()-5)
	inputBox.SetWidth(layout.GetWidth() - 4)

	// Add welcome message
	chatBox.AddMessage(terminal.RoleSystem, "Welcome to Nigesh Terminal Chat! Type your messages below.")
	chatBox.Render()

	reader := bufio.NewReader(os.Stdin)
	var input string

	for {
		inputBox.Print(input)

		char, _, err := reader.ReadRune()
		if err != nil {
			break
		}

		if char == '\n' {
			if input != "" {
				// Add user message
				chatBox.AddMessage(terminal.RoleUser, input)
				chatBox.Render()

				// Simulate AI response with streaming
				response := generateResponse(input)
				simulateStreaming(chatBox, response)

				input = ""
			}
		} else if char == 127 || char == 8 { // Backspace
			if len(input) > 0 {
				input = input[:len(input)-1]
			}
		} else if char == 3 { // Ctrl+C
			break
		} else if char >= 32 {
			input += string(char)
		}
	}
}

func generateResponse(userInput string) string {
	lower := strings.ToLower(userInput)

	if strings.Contains(lower, "hello") || strings.Contains(lower, "hi") {
		return "Hello! How can I help you today?"
	} else if strings.Contains(lower, "how are you") {
		return "I'm doing great! Thanks for asking. I'm a terminal-based chat interface built with Go."
	} else if strings.Contains(lower, "bye") || strings.Contains(lower, "goodbye") {
		return "Goodbye! Have a great day!"
	} else {
		return fmt.Sprintf("You said: '%s'. This is a demo response showing streaming text capabilities!", userInput)
	}
}

func simulateStreaming(chatBox *terminal.ChatBox, text string) {
	words := strings.Fields(text)
	for i, word := range words {
		if i > 0 {
			chatBox.StreamMessage(terminal.RoleAssistant, " ")
		}
		chatBox.StreamMessage(terminal.RoleAssistant, word)
		time.Sleep(50 * time.Millisecond)
	}
}
