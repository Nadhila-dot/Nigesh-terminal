package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Search(term, fileType string, verbose bool) {
	//fmt.Println(term + " This is the term")
	ext := ""
	if fileType != "" {
		ext = strings.TrimPrefix(fileType, ".")
	}

	if verbose {
		fmt.Printf("searching '%s'\n", term)
	}

	found := false
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if ext != "" && !strings.HasSuffix(path, "."+ext) {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 1
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, term) {
				fmt.Printf("%s:%d:%s\n", path, lineNum, line)
				found = true
			}
			lineNum++
		}
		return nil
	})

	if !found {
		fmt.Printf("no matches\n")
	}
}
