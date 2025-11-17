package reverseinput

import "testing"

func TestTokenize(t *testing.T) {
	tokens := Tokenize("U+0041, U+0042;0x0043 68")
	if len(tokens) != 4 {
		t.Fatalf("expected 4 tokens, got %d", len(tokens))
	}
	want := []string{"U+0041", "U+0042", "0x0043", "68"}
	for i, tok := range tokens {
		if tok != want[i] {
			t.Fatalf("token %d: want %q, got %q", i, want[i], tok)
		}
	}
}

func TestBuildStringFromCodePoints(t *testing.T) {
	input := []string{"U+0041", "0x0042", "67"}
	got, err := BuildStringFromCodePoints(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ABC" {
		t.Fatalf("expected ABC, got %q", got)
	}
}

func TestBuildStringFromBytes(t *testing.T) {
	input := []string{"0x41", "0b01100010", "99"}
	got, err := BuildStringFromBytes(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Abc" {
		t.Fatalf("expected Abc, got %q", got)
	}
}

func TestBuildStringInvalidTokens(t *testing.T) {
	if _, err := BuildStringFromCodePoints([]string{"not-a-number"}); err == nil {
		t.Fatalf("expected error for invalid code point token")
	}
	if _, err := BuildStringFromBytes([]string{"999"}); err == nil {
		t.Fatalf("expected error for out-of-range byte")
	}
}
