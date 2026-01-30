package main

import (
	"bufio"
	"fmt"
	"nadhi/do-t/terminal"
	"os"
	"strings"
	"time"
)

type ConversationUI struct {
	manager   *terminal.ComponentManager
	chatBox   *terminal.ChatBox
	inputBox  *terminal.InputBox
	statusBar *terminal.Panel
	layout    *terminal.Layout
	msgCount  int
}

func NewConversationUI() *ConversationUI {
	manager := terminal.NewComponentManager()
	layout := manager.GetLayout()

	// Chat area (main conversation)
	chatBox := terminal.NewChatBox(layout.GetWidth()-4, layout.GetHeight()-10)
	chatBox.SetPosition(2, 1)

	// Input box at bottom
	inputBox := terminal.NewInputBox("Ask me anything...")
	inputBox.SetPosition(2, layout.GetHeight()-6)
	inputBox.SetWidth(layout.GetWidth() - 4)

	// Status bar
	statusBar := terminal.NewPanel("", layout.GetWidth()-4, 3)
	statusBar.SetPosition(2, layout.GetHeight()-9)

	return &ConversationUI{
		manager:   manager,
		chatBox:   chatBox,
		inputBox:  inputBox,
		statusBar: statusBar,
		layout:    layout,
		msgCount:  0,
	}
}

func (ui *ConversationUI) Start() {
	terminal.ClearScreenFull()
	terminal.HideCursor()
	defer terminal.ShowCursor()

	// Welcome message
	ui.chatBox.AddMessage(terminal.RoleSystem, "Welcome to Nigesh AI Chat! I'm your terminal assistant.")
	ui.chatBox.AddMessage(terminal.RoleAssistant, "Hi! I'm here to help. Ask me anything or type 'help' for commands.")
	ui.updateStatus("Ready")
	ui.render()

	reader := bufio.NewReader(os.Stdin)
	var input string

	for {
		ui.inputBox.Print(input)

		char, _, err := reader.ReadRune()
		if err != nil {
			break
		}

		if char == '\n' {
			if input != "" {
				if ui.handleInput(input) {
					break
				}
				input = ""
			}
		} else if char == 127 || char == 8 {
			if len(input) > 0 {
				input = input[:len(input)-1]
			}
		} else if char == 3 {
			break
		} else if char >= 32 {
			input += string(char)
		}
	}
}

func (ui *ConversationUI) handleInput(input string) bool {
	lower := strings.ToLower(strings.TrimSpace(input))

	// Check for exit commands
	if lower == "exit" || lower == "quit" || lower == "bye" {
		ui.chatBox.AddMessage(terminal.RoleUser, input)
		ui.chatBox.Render()
		ui.streamResponse("Goodbye! Thanks for chatting with me. Have a great day!")
		time.Sleep(2 * time.Second)
		return true
	}

	// Add user message
	ui.chatBox.AddMessage(terminal.RoleUser, input)
	ui.msgCount++
	ui.updateStatus("Thinking...")
	ui.render()

	// Simulate thinking delay
	time.Sleep(300 * time.Millisecond)

	// Generate and stream response
	response := ui.generateResponse(input)
	ui.updateStatus("Responding...")
	ui.streamResponse(response)

	ui.msgCount++
	ui.updateStatus(fmt.Sprintf("%d messages", ui.msgCount))

	return false
}

func (ui *ConversationUI) generateResponse(input string) string {
	lower := strings.ToLower(input)

	if strings.Contains(lower, "help") {
		return "I can help you with various tasks! Try asking me about:\n• Programming questions\n• General knowledge\n• Terminal commands\n• Or just chat with me!\n\nType 'exit' or 'quit' to leave."
	}

	if strings.Contains(lower, "hello") || strings.Contains(lower, "hi") {
		return "Hello! Great to see you. How can I assist you today?"
	}

	if strings.Contains(lower, "how are you") {
		return "I'm functioning perfectly! I'm a terminal-based AI assistant built with Go. Thanks for asking! How are you doing?"
	}

	if strings.Contains(lower, "code") || strings.Contains(lower, "program") {
		return "I'd love to help with coding! I can assist with Go, Python, JavaScript, and many other languages. What would you like to build?"
	}

	if strings.Contains(lower, "terminal") || strings.Contains(lower, "ui") {
		return "This terminal UI is built using Go with custom rendering components. It features dynamic panels, chat boxes, and streaming text - similar to modern AI chat interfaces!"
	}

	if strings.Contains(lower, "thank") {
		return "You're very welcome! I'm always here to help. Is there anything else you'd like to know?"
	}

	// Default response
	return fmt.Sprintf("Interesting! You mentioned: '%s'. This is a demo showing streaming responses in a terminal UI. The system supports multi-line messages, word wrapping, and smooth text streaming!", input)
}

func (ui *ConversationUI) streamResponse(text string) {
	words := strings.Fields(text)
	for i, word := range words {
		if i > 0 {
			ui.chatBox.StreamMessage(terminal.RoleAssistant, " ")
		}
		ui.chatBox.StreamMessage(terminal.RoleAssistant, word)
		time.Sleep(40 * time.Millisecond)
	}
}

func (ui *ConversationUI) updateStatus(status string) {
	ui.statusBar.Clear()
	ui.statusBar.AppendContent(fmt.Sprintf("Status: %s | Terminal: %dx%d",
		status, ui.layout.GetWidth(), ui.layout.GetHeight()))
}

func (ui *ConversationUI) render() {
	ui.chatBox.Render()
	ui.statusBar.Render()
}


