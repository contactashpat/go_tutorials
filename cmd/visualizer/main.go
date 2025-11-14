// Fun Name Visualizer showcases how bytes represent human-friendly names and
// lets you decode hex or binary strings back into text.
package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "see":
		err = runSee(os.Args[2:])
	case "decode":
		err = runDecode(os.Args[2:])
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

// runSee prints a table showing the UTF-8 bytes in hex and binary for a name.
func runSee(args []string) error {
	fs := flag.NewFlagSet("see", flag.ContinueOnError)
	nameFlag := fs.String("name", "", "Name to visualize")
	if err := fs.Parse(args); err != nil {
		return err
	}

	name := *nameFlag
	if name == "" {
		name = strings.Join(fs.Args(), " ")
	}
	if name == "" {
		return errors.New("no name provided; use --name or add it after the command")
	}

	fmt.Printf("Name: %s\n", name)
	fmt.Println("This is how a computer represents your name byte-by-byte:")
	fmt.Println()
	fmt.Println("Letter           UTF-8 Hex Bytes        Binary Bytes")
	fmt.Println("--------------  --------------------  ------------------------------")

	for _, r := range name {
		bytes := []byte(string(r))
		hexParts := make([]string, len(bytes))
		binParts := make([]string, len(bytes))

		for i, b := range bytes {
			hexParts[i] = fmt.Sprintf("0x%02X", b)
			binParts[i] = fmt.Sprintf("%08b", b)
		}

		fmt.Printf("%-14q  %-20s  %s\n",
			r,
			strings.Join(hexParts, " "),
			strings.Join(binParts, " "),
		)
	}
	return nil
}

// runDecode reads hex or binary byte strings and prints the UTF-8 result.
func runDecode(args []string) error {
	fs := flag.NewFlagSet("decode", flag.ContinueOnError)
	hexInput := fs.String("hex", "", "Hex bytes (space separated or continuous, e.g. '41 73' or '4173')")
	binInput := fs.String("bin", "", "Binary bytes (space separated, e.g. '01000001 01110011')")
	if err := fs.Parse(args); err != nil {
		return err
	}

	switch {
	case *hexInput != "" && *binInput != "":
		return errors.New("please provide either --hex or --bin, not both")
	case *hexInput != "":
		bytes, err := parseHexInput(*hexInput)
		if err != nil {
			return err
		}
		printDecoded(bytes)
	case *binInput != "":
		bytes, err := parseBinaryInput(*binInput)
		if err != nil {
			return err
		}
		printDecoded(bytes)
	default:
		return errors.New("no input provided; use --hex or --bin")
	}
	return nil
}

// parseHexInput accepts "0x", spaces, or continuous hex bytes and returns raw bytes.
func parseHexInput(input string) ([]byte, error) {
	clean := strings.ReplaceAll(input, "0x", "")
	clean = strings.ReplaceAll(clean, "0X", "")
	clean = strings.ReplaceAll(clean, " ", "")
	if len(clean)%2 != 0 {
		clean = "0" + clean
	}
	bytes, err := hex.DecodeString(clean)
	if err != nil {
		return nil, fmt.Errorf("invalid hex input: %w", err)
	}
	return bytes, nil
}

// parseBinaryInput consumes space-separated binary chunks (up to 8 bits each).
func parseBinaryInput(input string) ([]byte, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return nil, errors.New("binary input is empty")
	}

	bytes := make([]byte, len(fields))
	for i, field := range fields {
		if len(field) > 8 {
			return nil, fmt.Errorf("binary chunk %q is longer than 8 bits", field)
		}
		val, err := strconv.ParseUint(field, 2, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid binary chunk %q: %w", field, err)
		}
		bytes[i] = byte(val)
	}
	return bytes, nil
}

// printDecoded displays the resulting UTF-8 text and warns on invalid sequences.
func printDecoded(bytes []byte) {
	if !utf8.Valid(bytes) {
		fmt.Println("Warning: byte sequence is not valid UTF-8, showing best effort:")
	}
	fmt.Printf("Decoded UTF-8: %s\n", string(bytes))
	fmt.Printf("Byte count: %d\n", len(bytes))
}

// printUsage documents available commands and sample invocations.
func printUsage() {
	fmt.Println(`Fun Name Visualizer

Usage:
  go run ./cmd/visualizer see --name "Ada Lovelace"
  go run ./cmd/visualizer decode --hex "41 64 61"
  go run ./cmd/visualizer decode --bin "01000001 01100100 01100001"

Commands:
  see      Show the hex and binary representation of every letter in a name.
  decode   Convert hex or binary bytes back into UTF-8 text.
`)
}
