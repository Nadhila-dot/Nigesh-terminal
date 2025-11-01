package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	dataStore = make(map[string]interface{}) // Use `interface{}` to support any value type
	mu        sync.Mutex
)

// Store saves the key-value pair to the specified file.
func Store(key, value string, verbose bool) {
	mu.Lock()
	defer mu.Unlock()

	// Parse the --file flag to get the file name
	fileFlag := flag.Lookup("file")
	var fileName string
	if fileFlag != nil {
		fileName = fileFlag.Value.String()
	}
	if fileName == "" {
		fileName = "data.json" // Default file name
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get current directory: %v\n", err)
		os.Exit(1)
	}

	// Construct the file path
	dir := filepath.Join(cwd, ".nigesh")
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create data directory %s: %v\n", dir, err)
		os.Exit(1)
	}
	filePath := filepath.Join(dir, fileName)

	// Load existing data from the file if it exists
	if _, err := os.Stat(filePath); err == nil {
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to read file %s: %v\n", filePath, err)
			os.Exit(1)
		}
		if len(fileData) > 0 {
			if err := json.Unmarshal(fileData, &dataStore); err != nil {
				fmt.Fprintf(os.Stderr, "error: failed to unmarshal JSON data: %v\n", err)
				os.Exit(1)
			}
		}
	}

	// Try to parse the value as JSON
	var parsedValue interface{}
	if err := json.Unmarshal([]byte(value), &parsedValue); err != nil {
		// If parsing fails, treat it as a plain string
		parsedValue = value
	}

	// Check if the key already exists and is a JSON object
	if existingValue, exists := dataStore[key]; exists {
		if existingMap, ok := existingValue.(map[string]interface{}); ok {
			// If the existing value is a map, merge the new data into it
			if newMap, ok := parsedValue.(map[string]interface{}); ok {
				for k, v := range newMap {
					existingMap[k] = v
				}
				parsedValue = existingMap
			}
		}
	}

	dataStore[key] = parsedValue

	if verbose {
		fmt.Printf("storing key: %s, value: %v\n", key, parsedValue)
	}

	// Save the updated data back to the file
	fileData, err := json.MarshalIndent(dataStore, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to marshal data: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to write to file: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("data saved to %s\n", filePath)
	}
}
