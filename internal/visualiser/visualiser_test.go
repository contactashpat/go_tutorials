package visualiser

import (
	"reflect"
	"testing"
)

func TestAnalyseStringASCII(t *testing.T) {
	results, err := AnalyseString("A")
	if err != nil {
		t.Fatalf("AnalyseString returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	res := results[0]
	if res.Character != "'A'" {
		t.Errorf("unexpected character: %s", res.Character)
	}
	if res.CodePointDec != 65 {
		t.Errorf("expected CodePointDec 65, got %d", res.CodePointDec)
	}
	if res.CodePointHex != "U+0041" {
		t.Errorf("expected CodePointHex U+0041, got %s", res.CodePointHex)
	}
	if res.HTMLEntityDecimal != "&#65;" {
		t.Errorf("expected decimal entity &#65;, got %s", res.HTMLEntityDecimal)
	}
	if res.HTMLEntityHex != "&#x0041;" {
		t.Errorf("expected hex entity &#x0041;, got %s", res.HTMLEntityHex)
	}
	expectHex := []string{"0x41"}
	expectDec := []string{"65"}
	expectBin := []string{"01000001"}
	if !reflect.DeepEqual(res.UTF8BytesHex, expectHex) {
		t.Errorf("hex bytes mismatch: %v", res.UTF8BytesHex)
	}
	if !reflect.DeepEqual(res.UTF8BytesDec, expectDec) {
		t.Errorf("dec bytes mismatch: %v", res.UTF8BytesDec)
	}
	if !reflect.DeepEqual(res.UTF8BytesBinary, expectBin) {
		t.Errorf("bin bytes mismatch: %v", res.UTF8BytesBinary)
	}
}

func TestAnalyseStringMultiByte(t *testing.T) {
	input := "ई🙂"
	results, err := AnalyseString(input)
	if err != nil {
		t.Fatalf("AnalyseString returned error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	hindi := results[0]
	if hindi.CodePointHex != "U+0908" {
		t.Errorf("expected U+0908, got %s", hindi.CodePointHex)
	}
	expectedHindiHex := []string{"0xE0", "0xA4", "0x88"}
	expectedHindiDec := []string{"224", "164", "136"}
	expectedHindiBin := []string{"11100000", "10100100", "10001000"}
	if !reflect.DeepEqual(hindi.UTF8BytesHex, expectedHindiHex) {
		t.Errorf("hindi hex mismatch: %v", hindi.UTF8BytesHex)
	}
	if !reflect.DeepEqual(hindi.UTF8BytesDec, expectedHindiDec) {
		t.Errorf("hindi dec mismatch: %v", hindi.UTF8BytesDec)
	}
	if !reflect.DeepEqual(hindi.UTF8BytesBinary, expectedHindiBin) {
		t.Errorf("hindi bin mismatch: %v", hindi.UTF8BytesBinary)
	}
	if hindi.HTMLEntityDecimal != "&#2312;" {
		t.Errorf("unexpected decimal entity: %s", hindi.HTMLEntityDecimal)
	}
	if hindi.HTMLEntityHex != "&#x0908;" {
		t.Errorf("unexpected hex entity: %s", hindi.HTMLEntityHex)
	}

	emoji := results[1]
	if emoji.CodePointHex != "U+1F642" {
		t.Errorf("expected U+1F642, got %s", emoji.CodePointHex)
	}
	expectedEmojiHex := []string{"0xF0", "0x9F", "0x99", "0x82"}
	expectedEmojiDec := []string{"240", "159", "153", "130"}
	expectedEmojiBin := []string{"11110000", "10011111", "10011001", "10000010"}
	if !reflect.DeepEqual(emoji.UTF8BytesHex, expectedEmojiHex) {
		t.Errorf("emoji hex mismatch: %v", emoji.UTF8BytesHex)
	}
	if !reflect.DeepEqual(emoji.UTF8BytesDec, expectedEmojiDec) {
		t.Errorf("emoji dec mismatch: %v", emoji.UTF8BytesDec)
	}
	if !reflect.DeepEqual(emoji.UTF8BytesBinary, expectedEmojiBin) {
		t.Errorf("emoji bin mismatch: %v", emoji.UTF8BytesBinary)
	}
}

func TestAnalyseStringEmpty(t *testing.T) {
	if _, err := AnalyseString(""); err == nil {
		t.Fatalf("expected error for empty input")
	}
}
