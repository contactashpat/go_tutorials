# Fun Name Visualizer

Small Go CLI that shows how computers perceive names by printing the UTF-8 bytes for each letter in both hexadecimal and binary. The companion decoder reverses the process by turning machine-friendly hex or binary back into text, highlighting how encoding schemes bridge people and computers.

## Requirements

- Go 1.24+ (uses only the standard library)

## Getting Started

```bash
cd /path/to/go_tutorials
go run . help
```

## Commands

### `see`

Show the byte-level representation for each letter in a name.

```bash
go run . see --name "Ada Lovelace"
```

Sample output:

```
Name: Ada Lovelace
This is how a computer represents your name byte-by-byte:

Letter           UTF-8 Hex Bytes        Binary Bytes
--------------  --------------------  ------------------------------
'A'             0x41                  01000001
'd'             0x64                  01100100
'a'             0x61                  01100001
' '             0x20                  00100000
'L'             0x4C                  01001100
...
```

### `decode`

Convert a stream of bytes (hexadecimal or binary) back to UTF-8 text. The hex mode accepts prefixes (`0x`) and both spaced or continuous strings.

```bash
# Hex input (bytes can be spaced or continuous)
go run . decode --hex "41 73 68"
go run . decode --hex "417368"

# Binary input (space-separated 8-bit chunks)
go run . decode --bin "01000001 01110011 01101000"
```

Invalid UTF-8 sequences trigger a warning but are still displayed with Go's best-effort decoding.

## Tips

- Use `go build` to create a static binary if you want to share the program: `go build -o funname .`
- Pipe output into `less` when trying long names or multi-line hex input.
- Modify `printUsage` in `main.go` if you add new subcommands or flags.
