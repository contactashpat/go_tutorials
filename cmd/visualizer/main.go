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
	"unicode"
	"unicode/utf8"

	"go_tutorials/internal/visualiser"
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
	nameFlag := fs.String("name", "", "Name or tokens to visualize")
	reverseFlag := fs.String("reverse", "", "Reverse input: 'codepoints' or 'bytes' to build text from numbers")
	if err := fs.Parse(args); err != nil {
		return err
	}

	input := *nameFlag
	if input == "" {
		input = strings.Join(fs.Args(), " ")
	}
	if input == "" {
		return errors.New("no name provided; use --name or add it after the command")
	}

	name := input
	switch strings.ToLower(*reverseFlag) {
	case "":
	case "codepoint", "codepoints", "cp":
		tokens := tokenizeReverseInput(input)
		if len(tokens) == 0 {
			return errors.New("reverse=codepoints requires numeric code point tokens")
		}
		built, err := buildStringFromCodePoints(tokens)
		if err != nil {
			return err
		}
		name = built
	case "byte", "bytes":
		tokens := tokenizeReverseInput(input)
		if len(tokens) == 0 {
			return errors.New("reverse=bytes requires byte tokens")
		}
		built, err := buildStringFromBytes(tokens)
		if err != nil {
			return err
		}
		name = built
	default:
		return fmt.Errorf("unknown reverse mode %q (use 'codepoints' or 'bytes')", *reverseFlag)
	}

	fmt.Printf("Name: %s\n", name)
	if *reverseFlag != "" {
		fmt.Printf("  (built from %s: %s)\n", *reverseFlag, input)
	}
	fmt.Println("This is how a computer represents your name byte-by-byte:")
	fmt.Println()
	results, err := visualiser.AnalyseString(name)
	if err != nil {
		return err
	}

	fmt.Println("Letter           Code Point (dec)  Code Point (hex)  HTML Entity (dec)  HTML Entity (hex)  UTF-8 Hex Bytes        UTF-8 Dec Bytes        Binary Bytes")
	fmt.Println("--------------  -----------------  ----------------  ------------------  ------------------  --------------------  ---------------------  ------------------------------")

	for _, res := range results {
		fmt.Printf("%-14s  %-17d  %-16s  %-18s  %-18s  %-20s  %-21s  %s\n",
			res.Character,
			res.CodePointDec,
			res.CodePointHex,
			res.HTMLEntityDecimal,
			res.HTMLEntityHex,
			strings.Join(res.UTF8BytesHex, " "),
			strings.Join(res.UTF8BytesDec, " "),
			strings.Join(res.UTF8BytesBinary, " "),
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

func tokenizeReverseInput(input string) []string {
	return strings.FieldsFunc(input, func(r rune) bool {
		return unicode.IsSpace(r) || r == ',' || r == ';'
	})
}

func buildStringFromCodePoints(tokens []string) (string, error) {
	var runes []rune
	for _, tok := range tokens {
		val, err := parseCodePointToken(tok)
		if err != nil {
			return "", err
		}
		if val < 0 || val > utf8.MaxRune {
			return "", fmt.Errorf("code point %d out of range for token %q", val, tok)
		}
		runes = append(runes, rune(val))
	}
	return string(runes), nil
}

func buildStringFromBytes(tokens []string) (string, error) {
	bytes := make([]byte, len(tokens))
	for i, tok := range tokens {
		val, err := parseByteToken(tok)
		if err != nil {
			return "", err
		}
		if val < 0 || val > 255 {
			return "", fmt.Errorf("byte %d out of range for token %q", val, tok)
		}
		bytes[i] = byte(val)
	}
	return string(bytes), nil
}

func parseCodePointToken(token string) (int64, error) {
	val, err := parseNumericToken(token, true)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func parseByteToken(token string) (int64, error) {
	val, err := parseNumericToken(token, false)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func parseNumericToken(token string, allowCodePointPrefix bool) (int64, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return 0, errors.New("empty numeric token")
	}

	upper := strings.ToUpper(token)
	switch {
	case allowCodePointPrefix && strings.HasPrefix(upper, "U+"):
		return strconv.ParseInt(token[2:], 16, 32)
	case strings.HasPrefix(token, "0x") || strings.HasPrefix(token, "0X"):
		return strconv.ParseInt(token[2:], 16, 64)
	case strings.HasPrefix(token, "0b") || strings.HasPrefix(token, "0B"):
		return strconv.ParseInt(token[2:], 2, 64)
	default:
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid numeric token %q: %w", token, err)
		}
		return val, nil
	}
}
