package main

import "fmt"

func printUsage() {
	fmt.Print(`Fun Name Visualizer

Usage:
  go run ./cmd/visualizer see --name "Ada Lovelace"
  go run ./cmd/visualizer see --reverse=codepoints --name "U+0041 0x0042 67"
  go run ./cmd/visualizer see --reverse=bytes --name "0xF0 0x9F 0x98 0x8A"
  go run ./cmd/visualizer decode --hex "41 64 61"
  go run ./cmd/visualizer decode --bin "01000001 01100100 01100001"

Commands:
  see      Show the hex and binary representation of every letter in a name.
  decode   Convert hex or binary bytes back into UTF-8 text.
`)
}
