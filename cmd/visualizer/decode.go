package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// DecodeCommand handles the `decode` sub-command.
type DecodeCommand struct{}

// NewDecodeCommand creates a decode command handler.
func NewDecodeCommand() *DecodeCommand {
	return &DecodeCommand{}
}

// Run executes the decode command.
func (c *DecodeCommand) Run(args []string) error {
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
