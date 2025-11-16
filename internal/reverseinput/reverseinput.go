package reverseinput

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Tokenize splits user-provided tokens for reverse parsing.
func Tokenize(input string) []string {
	return strings.FieldsFunc(input, func(r rune) bool {
		return unicode.IsSpace(r) || r == ',' || r == ';'
	})
}

// BuildStringFromCodePoints converts code point tokens into a UTF-8 string.
func BuildStringFromCodePoints(tokens []string) (string, error) {
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

// BuildStringFromBytes converts byte tokens into a UTF-8 string.
func BuildStringFromBytes(tokens []string) (string, error) {
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
	return parseNumericToken(token, true)
}

func parseByteToken(token string) (int64, error) {
	return parseNumericToken(token, false)
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
