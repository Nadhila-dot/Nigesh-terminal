package tools

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type CommandResult struct {
	Output   string
	Error    string
	ExitCode int
	Duration time.Duration
}

func RunCommand(command string) (*CommandResult, error) {
	start := time.Now()

	if strings.TrimSpace(command) == "" {
		return &CommandResult{
			Error:    "empty command",
			ExitCode: 1,
			Duration: time.Since(start),
		}, nil
	}

	// Detect long-running commands that should run in background
	commandLower := strings.ToLower(command)
	isLongRunning := strings.Contains(commandLower, "http.server") ||
		strings.Contains(commandLower, "serve") ||
		strings.Contains(commandLower, "npm start") ||
		strings.Contains(commandLower, "yarn start") ||
		strings.Contains(commandLower, "python -m http") ||
		strings.Contains(commandLower, "python3 -m http") ||
		strings.Contains(commandLower, "flask run") ||
		strings.Contains(commandLower, "rails server") ||
		strings.Contains(commandLower, "ng serve")

	// If it's a long-running command, run it in background
	if isLongRunning {
		// Run in background with nohup
		bgCommand := fmt.Sprintf("nohup %s > /dev/null 2>&1 &", command)
		cmd := exec.Command("sh", "-c", bgCommand)

		if err := cmd.Start(); err != nil {
			return &CommandResult{
				Error:    "failed to start background process: " + err.Error(),
				ExitCode: 1,
				Duration: time.Since(start),
			}, nil
		}

		// Give it a moment to start
		time.Sleep(500 * time.Millisecond)

		return &CommandResult{
			Output:   fmt.Sprintf("Started in background (PID: %d)\nCommand is running in the background.", cmd.Process.Pid),
			ExitCode: 0,
			Duration: time.Since(start),
		}, nil
	}

	// Use shell to execute the command properly (handles quotes, redirection, etc.)
	cmd := exec.Command("sh", "-c", command)

	// Set up pipes for live output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return &CommandResult{
			Error:    "failed to create stdout pipe: " + err.Error(),
			ExitCode: 1,
			Duration: time.Since(start),
		}, nil
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return &CommandResult{
			Error:    "failed to create stderr pipe: " + err.Error(),
			ExitCode: 1,
			Duration: time.Since(start),
		}, nil
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return &CommandResult{
			Error:    "failed to start command: " + err.Error(),
			ExitCode: 1,
			Duration: time.Since(start),
		}, nil
	}

	var outputBuilder strings.Builder
	var errorBuilder strings.Builder

	// Read stdout and stderr concurrently with live display
	done := make(chan error, 1)
	go func() {
		defer func() {
			stdout.Close()
			stderr.Close()
		}()

		// Print command output header
		fmt.Printf("\033[90m📐 Command Output:\033[0m\n")

		// Read stdout
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Printf("%s\n", line)
				outputBuilder.WriteString(line + "\n")
			}
		}()

		// Read stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Printf("\033[91m%s\033[0m\n", line)
				errorBuilder.WriteString(line + "\n")
			}
		}()

		err := cmd.Wait()
		done <- err
	}()

	// Set timeout for long-running commands
	timeout := 5 * time.Minute // Increased for installations
	select {
	case err := <-done:
		duration := time.Since(start)
		result := &CommandResult{
			Output:   outputBuilder.String(),
			Duration: duration,
		}

		if errorBuilder.Len() > 0 {
			result.Error = errorBuilder.String()
		}

		if err != nil {
			if result.Error == "" {
				result.Error = err.Error()
			}
			if exitError, ok := err.(*exec.ExitError); ok {
				result.ExitCode = exitError.ExitCode()
			} else {
				result.ExitCode = 1
			}
		}

		return result, nil

	case <-time.After(timeout):
		cmd.Process.Kill()
		return &CommandResult{
			Output:   outputBuilder.String(),
			Error:    "command timed out after 5 minutes",
			ExitCode: 124,
			Duration: timeout,
		}, nil
	}
}
