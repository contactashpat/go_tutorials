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

	var err error
	switch os.Args[1] {
	case "see":
		err = NewSeeCommand().Run(os.Args[2:])
	case "decode":
		err = NewDecodeCommand().Run(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
		return
	default:
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
