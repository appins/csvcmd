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
	orFilter      bool
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
	var opts options

	flag.BoolVar(&opts.humanReadable, "h", false, "Print in an easy to read format")

	flag.IntVar(&opts.startLine, "start", 1, "The first line, after the initial column line, that should be read (inclusive, 1-based index)")

	flag.IntVar(&opts.endLine, "end", -1, "The last line, after the initial column line, that should be read (inclusive, 1-based index)")

	flag.StringVar(&opts.filtersString, "filter", "", "Filters on columns, see GitHub for examples")

	flag.BoolVar(&opts.orFilter, "or", false, "Line will print if any single filter is matched")

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
		col_num := fmt.Sprintf("_%d", i+1)
		colsToInt[col_num] = i
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
	// Create a truncated csv reader, using the csvtrunc package
	csvReader, cols, err := csvtrunc.NewReader(fil, opts.startLine, opts.endLine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing %s: %s while searching for columns\n", fname, err)
		return
	}
	// Write the columns. This will automatically set the widths for the formatted writer
	output.Write(cols)

	// Create a filtered reader from the csv reader, using the csvfilter package
	filters, err := genFilters(opts.filtersString, cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", fname, err)
		return
	}

	// Create a filtered reader, which only reads out rows that meet the filter
	// criteria. Then we read all the rows from it into output
	filteredReader := csvfilter.NewReader(csvReader, filters, !opts.orFilter)
	for filteredReader.Scan() {
		err := output.Write(filteredReader.Row())
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
