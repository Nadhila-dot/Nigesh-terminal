package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func RunPython(args string) (*CommandResult, error) {
	start := time.Now()

	// Parse args: filename|code
	parts := strings.SplitN(args, "|", 2)
	if len(parts) != 2 {
		return &CommandResult{
			Error:    "invalid format: use filename.py|code",
			ExitCode: 1,
			Duration: time.Since(start),
		}, nil
	}

	scriptName := strings.TrimSpace(parts[0])
	code := strings.TrimSpace(parts[1])

	// Ensure .py extension
	if !strings.HasSuffix(scriptName, ".py") {
		scriptName += ".py"
	}

	// Create .nigesh/workspace directory
	cwd, _ := os.Getwd()
	workspaceDir := filepath.Join(cwd, ".nigesh", "workspace")
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		return &CommandResult{
			Error:    "failed to create workspace: " + err.Error(),
			ExitCode: 1,
			Duration: time.Since(start),
		}, nil
	}

	scriptPath := filepath.Join(workspaceDir, scriptName)

	// Write Python code to file
	if err := os.WriteFile(scriptPath, []byte(code), 0644); err != nil {
		return &CommandResult{
			Error:    "failed to write script: " + err.Error(),
			ExitCode: 1,
			Duration: time.Since(start),
		}, nil
	}

	// Execute the Python script
	cmd := fmt.Sprintf("cd %s && python3 %s", workspaceDir, scriptName)
	result, err := RunCommand(cmd)

	// Keep the script file for reuse and reference
	if result != nil {
		result.Output = fmt.Sprintf("Script saved as: %s\n%s", scriptName, result.Output)
	}

	return result, err
}

func InstallPythonPackage(packageName string) (*CommandResult, error) {
	cmd := fmt.Sprintf("pip3 install %s", packageName)
	return RunCommand(cmd)
}
