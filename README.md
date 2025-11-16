# Go Tutorials Playground

This repository hosts small Go experiments. At the moment there are two CLI tools:

1. **Fun Name Visualizer** (`cmd/visualizer`): shows each character’s code point (decimal and hex), HTML entities, and UTF-8 bytes (hex/dec/binary) and decodes byte streams back to UTF-8.
2. **Simple TCP Echo Server** (`cmd/echo`): minimalist server that echoes whatever clients send, handy for networking demos.

Both binaries use only the Go standard library (tested with Go 1.24+).

## Running the Visualizer

```bash
go run ./cmd/visualizer see --name "Ada Lovelace"
```

Sample output:

```
Name: Ada Lovelace
This is how a computer represents your name byte-by-byte:

Letter           Code Point (dec)  Code Point (hex)  HTML Entity (dec)  HTML Entity (hex)  UTF-8 Hex Bytes        UTF-8 Dec Bytes        Binary Bytes
--------------  -----------------  ----------------  ------------------  ------------------  --------------------  ---------------------  ------------------------------
'A'             65                 U+0041            &#65;               &#x0041;            0x41                  65                     01000001
'd'             100                U+0064            &#100;             &#x0064;            0x64                  100                    01100100
'a'             97                 U+0061            &#97;              &#x0061;            0x61                  97                     01100001
' '             32                 U+0020            &#32;              &#x0020;            0x20                  32                     00100000
'L'             76                 U+004C            &#76;              &#x004C;            0x4C                  76                     01001100
...
```

### Reverse Input Mode

Feed numeric code points or raw bytes and let the tool reconstruct the text before visualising it:

```bash
# Code points in mixed formats
go run ./cmd/visualizer see --reverse=codepoints --name "U+0041 0x0042 67"

# Raw bytes (hex / binary / decimal)
go run ./cmd/visualizer see --reverse=bytes --name "0xF0 0x9F 0x98 0x8A"
```

### Understanding the Columns

- **Code Point (dec)**: Unicode scalar value in base 10 (what `rune` represents).
- **Code Point (hex)**: Same value in the canonical `U+XXXX` notation.
- **HTML Entity (dec/hex)**: Ready-to-use HTML entity escape sequences.
- **UTF-8 Hex Bytes / UTF-8 Dec Bytes / Binary Bytes**: How UTF-8 encodes that rune at the byte level.

Decode hex or binary back to text:

```bash
# Hex input (bytes can be spaced or continuous)
go run ./cmd/visualizer decode --hex "41 73 68"
go run ./cmd/visualizer decode --hex "417368"
# Binary input (space-separated 8-bit chunks)
go run ./cmd/visualizer decode --bin "01000001 01110011 01101000"
```

Invalid UTF-8 sequences trigger a warning but still print with Go's best-effort decoding.

## Running the Echo Server

```bash
go run ./cmd/echo --addr :9000
```

Then connect from another terminal (or `nc`):

```bash
nc localhost 9000
hello world
hello world
```

Every message you type is sent straight back, making it easy to inspect TCP traffic.

## Building Binaries

```bash
go build -o bin/visualizer ./cmd/visualizer
go build -o bin/echo ./cmd/echo
```

To sanity-check the CLI, run the helper script with a few sample commands:

```bash
./scripts/run_visualizer_examples.sh
```

Feel free to add more tools under `cmd/<your-tool>` and keep shared helpers under `internal/` if needed.
