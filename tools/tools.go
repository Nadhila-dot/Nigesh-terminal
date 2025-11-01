package tools

import (
	"fmt"
	"regexp"
	"strings"
)

type ToolCall struct {
	Name string
	Args string
}

func ExtractToolCalls(response string) []ToolCall {
	var tools []ToolCall

	// Look for <Tool>ToolName(args)</Tool> pattern - more flexible regex
	toolPattern := `<Tool>(\w+)\((.*?)\)</Tool>`
	re := regexp.MustCompile(`(?s)` + toolPattern) // (?s) makes . match newlines
	matches := re.FindAllStringSubmatch(response, -1)

	for _, match := range matches {
		if len(match) > 2 {
			tools = append(tools, ToolCall{
				Name: match[1],
				Args: strings.TrimSpace(match[2]),
			})
		}
	}

	// Also look for tool calls that might be at the end of responses
	if len(tools) == 0 {
		// Try a more lenient pattern
		lenientPattern := `<Tool>(\w+)\(([^)]*)\)</Tool>`
		re2 := regexp.MustCompile(`(?s)` + lenientPattern)
		matches2 := re2.FindAllStringSubmatch(response, -1)

		for _, match := range matches2 {
			if len(match) > 2 {
				tools = append(tools, ToolCall{
					Name: match[1],
					Args: strings.TrimSpace(match[2]),
				})
			}
		}
	}

	return tools
}

func RemoveToolCalls(response string) string {
	// Remove <Tool>...</Tool> patterns from response
	toolPattern := `<Tool>\w+\([^)]*\)</Tool>`
	re := regexp.MustCompile(`(?s)` + toolPattern)
	cleaned := re.ReplaceAllString(response, "")
	return strings.TrimSpace(cleaned)
}

func ExecuteTool(tool ToolCall, verbose bool) string {
	switch tool.Name {
	case "Search":
		if verbose {
			fmt.Printf("🌐 searching: %s\n", tool.Args)
		}

		result, err := Search(tool.Args)
		if err != nil {
			return fmt.Sprintf("Search failed for '%s': %v\n", tool.Args, err)
		}

		var output strings.Builder
		output.WriteString(fmt.Sprintf("Search results for '%s':\n", tool.Args))

		// Add image links first
		if len(result.ImageLinks) > 0 {
			output.WriteString("Image links found:\n")
			for i, img := range result.ImageLinks {
				if i >= 5 { // Show up to 5 images
					break
				}
				output.WriteString(fmt.Sprintf("- %s\n", img))
			}
			output.WriteString("\n")
		}

		// Add snippets
		for i, snippet := range result.Snippets {
			if i >= 3 {
				break
			}
			output.WriteString(fmt.Sprintf("- %s\n", snippet))
		}

		// Add related links
		if len(result.RelatedLinks) > 0 {
			output.WriteString("\nRelevant links:\n")
			for _, link := range result.RelatedLinks {
				output.WriteString(fmt.Sprintf("- %s\n", link))
			}
		}

		return output.String()

	case "Command":

		fmt.Printf("⚡ running: %s\n", tool.Args)

		result, err := RunCommand(tool.Args)
		if err != nil {
			return fmt.Sprintf("Command execution failed: %v\n", err)
		}

		var output strings.Builder
		output.WriteString(fmt.Sprintf("Command: %s\n", tool.Args))
		output.WriteString(fmt.Sprintf("Duration: %v\n", result.Duration))
		output.WriteString(fmt.Sprintf("Exit Code: %d\n", result.ExitCode))

		if result.Output != "" {
			output.WriteString("Output:\n")
			output.WriteString(result.Output)
			output.WriteString("\n")
		}

		if result.Error != "" {
			output.WriteString("Error:\n")
			output.WriteString(result.Error)
			output.WriteString("\n")
		}

		return output.String()

	case "Python":
		if verbose {
			fmt.Printf("🐍 executing python code\n")
		}

		result, err := RunPython(tool.Args)
		if err != nil {
			return fmt.Sprintf("Python execution failed: %v\n", err)
		}

		var output strings.Builder
		output.WriteString("Python Code Executed:\n")
		output.WriteString(fmt.Sprintf("Duration: %v\n", result.Duration))
		output.WriteString(fmt.Sprintf("Exit Code: %d\n", result.ExitCode))

		if result.Output != "" {
			output.WriteString("Output:\n")
			output.WriteString(result.Output)
			output.WriteString("\n")
		}

		if result.Error != "" {
			output.WriteString("Error:\n")
			output.WriteString(result.Error)
			output.WriteString("\n")
		}

		return output.String()

	case "PipInstall":
		if verbose {
			fmt.Printf("📦 installing python package: %s\n", tool.Args)
		}

		result, err := InstallPythonPackage(tool.Args)
		if err != nil {
			return fmt.Sprintf("Package installation failed: %v\n", err)
		}

		var output strings.Builder
		output.WriteString(fmt.Sprintf("Installing package: %s\n", tool.Args))
		output.WriteString(fmt.Sprintf("Duration: %v\n", result.Duration))
		output.WriteString(fmt.Sprintf("Exit Code: %d\n", result.ExitCode))

		if result.Output != "" {
			output.WriteString("Output:\n")
			output.WriteString(result.Output)
			output.WriteString("\n")
		}

		if result.Error != "" {
			output.WriteString("Error:\n")
			output.WriteString(result.Error)
			output.WriteString("\n")
		}

		return output.String()

	default:
		return fmt.Sprintf("Unknown tool: %s\n", tool.Name)
	}
}
