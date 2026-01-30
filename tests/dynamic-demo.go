package main

import (
	"bufio"
	"fmt"
	"nadhi/do-t/terminal"
	"os"
	"strconv"
	"strings"
	"time"
)

func Dynamic_demo() {
	terminal.ClearScreenFull()
	terminal.HideCursor()
	defer terminal.ShowCursor()

	manager := terminal.NewComponentManager()
	layout := manager.GetLayout()

	// Create initial components
	helpPanel := terminal.NewPanel("Commands", 50, 10)
	helpPanel.SetPosition(2, 2)
	helpPanel.SetContent([]string{
		"add <id> <x> <y> - Add panel",
		"move <id> <x> <y> - Move panel",
		"remove <id> - Remove panel",
		"list - List all panels",
		"clear - Clear all panels",
		"grid <rows> <cols> - Create grid",
		"quit - Exit",
	})
	manager.Add("help", helpPanel)

	inputBox := terminal.NewInputBox("Enter command...")
	inputBox.SetPosition(2, layout.GetHeight()-5)
	inputBox.SetWidth(layout.GetWidth() - 4)

	statusPanel := terminal.NewPanel("Status", layout.GetWidth()-4, 3)
	statusPanel.SetPosition(2, layout.GetHeight()-8)
	statusPanel.SetContent([]string{"Ready. Type 'help' for commands."})
	manager.Add("status", statusPanel)

	manager.RenderAll()

	reader := bufio.NewReader(os.Stdin)
	var input string
	panelCounter := 0

	for {
		inputBox.Print(input)

		char, _, err := reader.ReadRune()
		if err != nil {
			break
		}

		if char == '\n' {
			if input != "" {
				result := handleCommand(manager, input, &panelCounter)

				if result == "quit" {
					break
				}

				// Update status
				if comp, exists := manager.Get("status"); exists {
					if status, ok := comp.(*terminal.Panel); ok {
						status.Clear()
						status.AppendContent(result)
					}
				}

				manager.RenderAll()
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

func handleCommand(manager *terminal.ComponentManager, cmd string, counter *int) string {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "Empty command"
	}

	switch parts[0] {
	case "add":
		if len(parts) < 4 {
			return "Usage: add <id> <x> <y>"
		}
		x, _ := strconv.Atoi(parts[2])
		y, _ := strconv.Atoi(parts[3])

		panel := terminal.NewPanel(parts[1], 30, 6)
		panel.SetPosition(x, y)
		panel.SetContent([]string{
			fmt.Sprintf("Panel: %s", parts[1]),
			fmt.Sprintf("Position: (%d, %d)", x, y),
			fmt.Sprintf("Created: %s", time.Now().Format("15:04:05")),
		})
		manager.Add(parts[1], panel)
		*counter++
		return fmt.Sprintf("Added panel '%s' at (%d, %d)", parts[1], x, y)

	case "move":
		if len(parts) < 4 {
			return "Usage: move <id> <x> <y>"
		}
		x, _ := strconv.Atoi(parts[2])
		y, _ := strconv.Atoi(parts[3])
		manager.MoveComponent(parts[1], x, y)
		return fmt.Sprintf("Moved '%s' to (%d, %d)", parts[1], x, y)

	case "remove":
		if len(parts) < 2 {
			return "Usage: remove <id>"
		}
		manager.Remove(parts[1])
		return fmt.Sprintf("Removed panel '%s'", parts[1])

	case "list":
		ids := manager.List()
		if len(ids) == 0 {
			return "No panels"
		}
		return fmt.Sprintf("Panels: %s", strings.Join(ids, ", "))

	case "clear":
		manager.Clear()
		return "Cleared all panels"

	case "grid":
		if len(parts) < 3 {
			return "Usage: grid <rows> <cols>"
		}
		rows, _ := strconv.Atoi(parts[1])
		cols, _ := strconv.Atoi(parts[2])

		manager.CreateGrid(rows, cols, func(r, c int) terminal.Component {
			panel := terminal.NewPanel(fmt.Sprintf("%d,%d", r, c), 20, 5)
			panel.SetContent([]string{
				fmt.Sprintf("Row: %d", r),
				fmt.Sprintf("Col: %d", c),
			})
			return panel
		})
		return fmt.Sprintf("Created %dx%d grid", rows, cols)

	case "quit":
		return "quit"

	default:
		return fmt.Sprintf("Unknown command: %s", parts[0])
	}
}
