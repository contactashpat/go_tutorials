package main

import (
	"fmt"
	"strings"

	"go_tutorials/internal/visualiser"
)

const (
	tableHeader  = "Letter           Code Point (dec)  Code Point (hex)  HTML Entity (dec)  HTML Entity (hex)  UTF-8 Hex Bytes        UTF-8 Dec Bytes        Binary Bytes"
	tableDivider = "--------------  -----------------  ----------------  ------------------  ------------------  --------------------  ---------------------  ------------------------------"
)

// renderTable prints high-level info followed by the per-character table.
func renderTable(resolvedText, note string, results []visualiser.Result) {
	fmt.Printf("Name: %s\n", resolvedText)
	if note != "" {
		fmt.Printf("  (%s)\n", note)
	}
	fmt.Println("This is how a computer represents your name byte-by-byte:")
	fmt.Println()
	fmt.Println(tableHeader)
	fmt.Println(tableDivider)

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
}
