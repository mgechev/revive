package formatter

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// formatTable formats a 2D string array into a table with aligned columns.
// Each row must have the same number of columns.
func formatTable(rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	colWidths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, col := range row {
			if w := utf8.RuneCountInString(col); w > colWidths[i] {
				colWidths[i] = w
			}
		}
	}

	var buf strings.Builder
	indent := "  "
	for _, row := range rows {
		buf.WriteString(indent)
		for i, col := range row {
			fmt.Fprintf(&buf, "%-*s", colWidths[i]+2, col)
		}
		buf.WriteByte('\n')
	}

	return buf.String()
}
