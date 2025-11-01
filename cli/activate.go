package cli

import (
	"fmt"
	"os"
)

func Activate(file string, verbose bool) {
	if verbose {
		fmt.Printf("chmod +x %s\n", file)
	}

	if err := os.Chmod(file, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("done\n")
	}
}