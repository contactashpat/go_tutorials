// Fun Name Visualizer showcases how bytes represent human-friendly names and
// lets you decode hex or binary strings back into text.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	name := os.Args[1]
	if name == "help" || name == "-h" || name == "--help" {
		printUsage()
		return
	}

	if err := dispatchCommand(name, os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
