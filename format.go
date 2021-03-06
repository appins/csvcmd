package main

import (
	"fmt"
	"strings"
)

type formattedWrite struct {
	widths []int
}

// The built in CSV writer satisfies this interface, and so does the
// formattedWrite type below. One writes raw CSV, the other writes
// formatted CSV
type lineWriter interface {
	Write([]string) error
	Flush()
}

// Write is defined on our formattedWrite struct as displaying each column with
// a fixed amount of space around it
func (w *formattedWrite) Write(row []string) error {
	// We must trim off white space and unicode identifiers
	for i, j := range row {
		j = strings.Trim(j, string([]rune{65279}))
		j = strings.TrimSpace(j)
		row[i] = j
	}
	// We overwrite lengths if they're empty
	if len(w.widths) == 0 {
		w.widths = genWidths(row)
	}
	for i, j := range row {
		// Limit the width to either widths[column] or the cells lenth itself
		// Whichever is less
		width := w.widths[i]
		if width >= len(j) {
			width = len(j)
		} else {
			j = j[:width-3] + "..."
		}
		// Store the original width, so we can equally size columns
		ow := w.widths[i]
		fmt.Print(j[:width] + strings.Repeat(" ", ow-width) + " ")
	}
	fmt.Print("\n")
	return nil
}

// As flush is called exactly once per file, we can use it to reset the widths
// Flush is called on the csv.Writer so that the text actually displays in the console
func (w *formattedWrite) Flush() {
	w.widths = []int{}
}
