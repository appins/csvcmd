package main

import (
	"encoding/csv"
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
	orFilter      bool
	columns       string
}

// The built in CSV writer satisfies this interface, and so does the
// formattedWrite type below. One writes raw CSV, the other writes
// formatted CSV
type lineWriter interface {
	Write([]string) error
	Flush()
}

type formattedWrite struct {
	widths []int
}

func main() {
	// Create an instance of the options struct and populate it with our command line ags
	var opts options
	flag.BoolVar(&opts.humanReadable, "h", false, "Print in an easy to read format")
	flag.IntVar(&opts.startLine, "start", 1, "The first line, after the initial column line, that should be read (inclusive, 1-based index)")
	flag.IntVar(&opts.endLine, "end", -1, "The last line, after the initial column line, that should be read (inclusive, 1-based index)")
	flag.StringVar(&opts.filtersString, "filter", "", "Filters on columns, see GitHub for examples")
	flag.BoolVar(&opts.orFilter, "or", false, "Line will print if any single filter is matched")
	flag.StringVar(&opts.columns, "shown", "", "Which columns should be output")

	flag.Parse()

	var writer lineWriter
	if opts.humanReadable {
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

// shwoColumns takes a list of bools and strings and only returns the strings
// whose matching index in the bool array is true. Ex: ([false, true, false],
// ["a", "b", "c"]) => ["b"]
func showColumns(enabled []bool, row []string) []string {
	var result []string
	for i, j := range row {
		if enabled[i] {
			result = append(result, j)
		}
	}

	return result
}

// Process a file and write each line ([]string) with output.Write. opts contains
// command line flags and options
func processFile(fil io.Reader, fname string, opts options, output lineWriter) {
	// Create a truncated csv reader, using the csvtrunc package
	csvReader, cols, err := csvtrunc.NewReader(fil, opts.startLine, opts.endLine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing %s: %s while searching for columns\n", fname, err)
		return
	}

	// Create the filter functions, that is, functions that take a row and return bools
	filters, err := genFilters(opts.filtersString, cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", fname, err)
		return
	}

	// Create the enabled column list
	enabled, err := genColumns(opts.columns, cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", fname, err)
		return
	}

	// Write the columns. This will automatically set the widths for the formatted writer
	// Notice that this is run through the show columns function, which limits the columns
	// that come out if it
	output.Write(showColumns(enabled, cols))

	// Create a filtered reader, which only reads out rows that meet the filter
	// criteria. Then we read all the rows from it into output
	filteredReader := csvfilter.NewReader(csvReader, filters, !opts.orFilter)
	for filteredReader.Scan() {
		// We only show the columns that we processed with showColumns
		processed := showColumns(enabled, filteredReader.Row())
		err := output.Write(processed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", fname, err)
		}
	}

	// Finally, flush output. Since this is called exactly once per file, the
	// flush function for a formatted reader resets the line lengths so that
	// they can be unique for each file.
	output.Flush()

}

// Write is defined on our formattedWrite struct as displaying each column with
// a fixed amount of space around it
func (w *formattedWrite) Write(row []string) error {
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
