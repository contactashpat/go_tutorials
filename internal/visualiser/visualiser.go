// Package visualiser exposes helpers to inspect how text is represented in bytes.
package visualiser

import (
	"errors"
	"fmt"
)

// Result captures descriptive data for a single rune in a string.
type Result struct {
	Character         string   // original rune formatted via %q for readability
	CodePointHex      string   // Unicode code point in U+XXXX form
	CodePointDec      int      // Unicode code point in decimal
	UTF8BytesHex      []string // each UTF-8 byte formatted as 0xHH
	UTF8BytesDec      []string // decimal byte values
	UTF8BytesBinary   []string // binary byte values
	HTMLEntityDecimal string   // e.g., &#65;
	HTMLEntityHex     string   // e.g., &#x0041;
}

// AnalyseString walks each rune in the input and returns byte-level metadata.
func AnalyseString(input string) ([]Result, error) {
	if len(input) == 0 {
		return nil, errors.New("input string is empty")
	}

	results := make([]Result, 0, len(input))
	for _, r := range input {
		b := []byte(string(r))
		hexParts := make([]string, len(b))
		decParts := make([]string, len(b))
		binParts := make([]string, len(b))
		for i, by := range b {
			hexParts[i] = fmt.Sprintf("0x%02X", by)
			decParts[i] = fmt.Sprintf("%d", by)
			binParts[i] = fmt.Sprintf("%08b", by)
		}

		res := Result{
			Character:         fmt.Sprintf("%q", r),
			CodePointHex:      fmt.Sprintf("U+%04X", r),
			CodePointDec:      int(r),
			UTF8BytesHex:      hexParts,
			UTF8BytesDec:      decParts,
			UTF8BytesBinary:   binParts,
			HTMLEntityDecimal: fmt.Sprintf("&#%d;", r),
			HTMLEntityHex:     fmt.Sprintf("&#x%04X;", r),
		}
		results = append(results, res)
	}
	return results, nil
}
