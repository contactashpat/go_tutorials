package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = orig
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestParseHexInput(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect string
	}{
		{"spaced", "41 42", "AB"},
		{"prefixed", "0x41 0x42", "AB"},
		{"continuous", "4142", "AB"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseHexInput(tc.input)
			if err != nil {
				t.Fatalf("parseHexInput error: %v", err)
			}
			if string(got) != tc.expect {
				t.Fatalf("expected %q, got %q", tc.expect, string(got))
			}
		})
	}
	if _, err := parseHexInput("zz"); err == nil {
		t.Fatalf("expected error for invalid hex")
	}
}

func TestParseBinaryInput(t *testing.T) {
	bytes, err := parseBinaryInput("01000001 01000010")
	if err != nil {
		t.Fatalf("parseBinaryInput error: %v", err)
	}
	if string(bytes) != "AB" {
		t.Fatalf("expected AB, got %q", string(bytes))
	}
	if _, err := parseBinaryInput("2"); err == nil {
		t.Fatalf("expected error for invalid binary chunk")
	}
	if _, err := parseBinaryInput(""); err == nil {
		t.Fatalf("expected error for empty binary input")
	}
}

func TestDecodeCommandHex(t *testing.T) {
	cmd := NewDecodeCommand()
	out := captureOutput(t, func() {
		if err := cmd.Run([]string{"--hex", "41 42"}); err != nil {
			t.Fatalf("run error: %v", err)
		}
	})
	if !strings.Contains(out, "Decoded UTF-8: AB") {
		t.Fatalf("expected decoded output to mention AB, got %q", out)
	}
}

func TestDecodeCommandBinary(t *testing.T) {
	cmd := NewDecodeCommand()
	out := captureOutput(t, func() {
		if err := cmd.Run([]string{"--bin", "01000001 01000010"}); err != nil {
			t.Fatalf("run error: %v", err)
		}
	})
	if !strings.Contains(out, "Decoded UTF-8: AB") {
		t.Fatalf("expected decoded output to mention AB, got %q", out)
	}
}

func TestDecodeCommandErrors(t *testing.T) {
	cmd := NewDecodeCommand()
	if err := cmd.Run([]string{"--hex", "41", "--bin", "01000001"}); err == nil {
		t.Fatalf("expected error when both hex and bin provided")
	}
	if err := cmd.Run([]string{}); err == nil {
		t.Fatalf("expected error when no input provided")
	}
}

func TestDispatchUnknownCommand(t *testing.T) {
	if err := dispatchCommand("does-not-exist", nil); err == nil {
		t.Fatalf("expected error for unknown command")
	}
}
