package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/appins/csvcmd/pkg/csvfilter"
	"github.com/appins/csvcmd/pkg/csvtrunc"
)

type options struct {
	humanReadable bool
	startLine     int
	endLine       int
	filtersString string
}

// The built in CSV writer satisfies this interface, and so does the
// type for formatted CSV below
type lineWriter interface {
	Write([]string) error
	Flush()
}

type formattedWrite struct {
	Widths []int
}

// Write is defined on our formattedWrite struct as displaying each column with
// a fixed amount of space around it
func (w *formattedWrite) Write(row []string) error {
	// We overwrite lengths if they're empty
	if len(w.Widths) == 0 {
		w.Widths = genWidths(row)
	}
	for i, j := range row {
		// Limit the width to either Widths[column] or the cells lenth itself
		// Whichever is less
		width := w.Widths[i]
		if width >= len(j) {
			width = len(j)
		} else {
			j = j[:width-3] + "..."
		}
		// Store the original width, so we can equally size columns
		ow := w.Widths[i]
		fmt.Print(j[:width] + strings.Repeat(" ", ow-width) + " ")
	}
	fmt.Print("\n")
	return nil
}

// As flush is called exactly once per file, we can use it to reset the widths
// Flush is called on the csv.Writer so that the text actually displays in the console
func (w *formattedWrite) Flush() {
	w.Widths = []int{}
}

func main() {
	var humanReadable bool
	flag.BoolVar(&humanReadable, "h", false, "Print in an easy to read format")

	var startLine int
	flag.IntVar(&startLine, "start", 0, "The first line, after the initial column line, that should be read (inclusive, 1-based index)")

	var endLine int
	flag.IntVar(&endLine, "end", -1, "The last line, after the initial column line, that should be read (inclusive, 1-based index)")

	var filtersString string
	flag.StringVar(&filtersString, "filter", "", "Filters on columns, see GitHub for examples")

	flag.Parse()

	opts := options{humanReadable, startLine, endLine, filtersString}
	var writer lineWriter
	if humanReadable {
		writer = &formattedWrite{}
	} else {
		writer = csv.NewWriter(os.Stdout)
	}

	// The files that should be read
	if len(flag.Args()) != 0 {
		files := []io.Reader{}
		// Open the files and add them to the files slice, to ensure that
		// we can open all of them
		for _, fil := range flag.Args() {
			fp, err := os.Open(fil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				return
			}
			files = append(files, io.Reader(fp))
			defer fp.Close()
		}
		// Then we call process file with each file
		for i, fil := range files {
			processFile(io.Reader(fil), flag.Args()[i], opts, writer)
		}

	} else {
		// In the case of no files being specified, read from stdin
		processFile(io.Reader(os.Stdin), "STDIN", opts, writer)
	}

}

func genWidths(row []string) []int {
	var lengths []int
	for _, j := range row {
		if len(j) < 10 {
			lengths = append(lengths, len(j))
		} else {
			lengths = append(lengths, 10)
		}
	}

	return lengths
}

func genFilters(filterString string, cols []string) ([]func([]string) bool, error) {
	var filters []func([]string) bool
	// Create a map of all columns, so that we can search for them
	colsToInt := make(map[string]int)
	for i, col := range cols {
		colsToInt[col] = i
	}

	// Break the filterString into indivisual filters
	filterStrings := strings.Split(filterString, ";")
	for _, filter := range filterStrings {
		// Skip empty filters
		if len(filter) == 0 {
			continue
		}

		// Break apart the filter by =. If there is an equal sign, it's an = filter
		if parts := strings.Split(filter, "="); len(parts) == 2 {
			// We gotta check if the column exists. If it does, we can refer to it as col
			if col, ok := colsToInt[parts[0]]; ok {
				filters = append(filters, func(row []string) bool {
					return row[col] == parts[1]
				})
			} else {
				// No col -> fail to parse all of the filters
				// IDEA: have some character that can ignore faulty columns
				return filters, errors.New("Couldn't find column " + parts[0])
			}

		} else {
			return filters, errors.New("Couldn't parse filter " + filter)
		}
	}

	return filters, nil

}

func processFile(fil io.Reader, fname string, opts options, output lineWriter) {
	csvReader, cols, err := csvtrunc.NewReader(fil, opts.startLine, opts.endLine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", fname, err)
		return
	}
	output.Write(cols)

	filters, err := genFilters(opts.filtersString, cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", fname, err)
		return
	}
	filteredReader := csvfilter.NewReader(csvReader, filters, true)
	for filteredReader.Scan() {
		err := output.Write(filteredReader.Row())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", fname, err)
		}
	}

	output.Flush()

}
