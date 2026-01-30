# Nigesh Terminal UI Components

A collection of terminal UI components for building interactive chat interfaces, similar to Claude or Gemini, directly in your terminal.

## Components

### ChatBox
A scrollable chat interface with support for multiple message roles (user, assistant, system).

```go
chatBox := terminal.NewChatBox(width, height)
chatBox.SetPosition(x, y)
chatBox.AddMessage(terminal.RoleUser, "Hello!")
chatBox.StreamMessage(terminal.RoleAssistant, "Hi there!")
chatBox.Render()
```

### Panel
Customizable bordered panels with titles and styling.

```go
panel := terminal.NewPanel("Title", width, height)
panel.SetPosition(x, y)
panel.SetStyle(terminal.UserStyle)
panel.AppendContent("Line 1")
panel.Render()
```

### InputBox
Input field with placeholder support.

```go
inputBox := terminal.NewInputBox("Type here...")
inputBox.SetPosition(x, y)
inputBox.SetWidth(50)
inputBox.Print(currentInput)
```

### ComponentManager
Manage multiple components dynamically.

```go
manager := terminal.NewComponentManager()
manager.Add("panel1", myPanel)
manager.MoveComponent("panel1", 10, 5)
manager.RenderAll()
```

### Layout
Get terminal dimensions and manage screen layout.

```go
layout := terminal.NewLayout()
width := layout.GetWidth()
height := layout.GetHeight()
```

## Demos

### AI Chat Demo
Full-featured chat interface with streaming responses:
```bash
go run tests/ai-chat-demo.go
```

### Dynamic Demo
Interactive component manipulation:
```bash
go run tests/dynamic-demo.go
```

### Panel Demo
Multiple styled panels:
```bash
go run tests/panel-demo.go
```

### Chat Demo
Basic chat with streaming:
```bash
go run tests/chat-demo.go
```

## Features

- ✨ Streaming text support (like AI responses)
- 🎨 Customizable colors and styles
- 📦 Multiple message roles (user/assistant/system)
- 📏 Automatic text wrapping
- 🔄 Dynamic component positioning
- 📊 Grid layout support
- ⌨️ Real-time input handling
- 🖥️ Responsive to terminal size

## Usage Example

```go
package main

import (
    "nadhi/do-t/terminal"
    "time"
)

func main() {
    terminal.ClearScreenFull()
    terminal.HideCursor()
    defer terminal.ShowCursor()
    
    // Create chat interface
    chat := terminal.NewChatBox(80, 20)
    chat.SetPosition(2, 2)
    
    // Add messages
    chat.AddMessage(terminal.RoleUser, "Hello!")
    chat.Render()
    
    // Stream response
    words := []string{"Hi", "there!", "How", "can", "I", "help?"}
    for _, word := range words {
        chat.StreamMessage(terminal.RoleAssistant, word+" ")
        time.Sleep(100 * time.Millisecond)
    }
}
```

## Styles

Pre-defined styles available:
- `DefaultStyle` - Cyan borders
- `UserStyle` - Blue borders
- `AssistantStyle` - Magenta borders

Create custom styles:
```go
customStyle := terminal.PanelStyle{
    BorderColor: "\033[32m", // Green
    TitleColor:  "\033[1;32m",
    TextColor:   "\033[0m",
}
panel.SetStyle(customStyle)
```
