package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/appins/csvcmd/pkg/csvfilter"
	"github.com/appins/csvcmd/pkg/csvnewcol"
	"github.com/appins/csvcmd/pkg/csvtrunc"
)

type options struct {
	humanReadable bool
	startLine     int
	endLine       int
	filtersString string
	orFilter      bool
	columns       string

	// Split string, numerator, and denominator
	split  string
	splitN int
	splitD int

	// Formulas for new columns
	newcols string
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
	flag.StringVar(&opts.split, "split", "", "Return a porton of the file without any overlaps")
	flag.StringVar(&opts.newcols, "newcols", "", "Create a new column from existing columns")

	flag.Parse()

	var writer lineWriter
	if opts.humanReadable {
		writer = &formattedWrite{}
	} else {
		writer = csv.NewWriter(os.Stdout)
	}

	// If we're splitting the output, define a numerator and denominator
	if opts.split != "" {
		// First we split the paramter by /, it should be in the format of '1/2'
		splitFrac := strings.Split(opts.split, "/")
		if len(splitFrac) != 2 {
			fmt.Fprintf(os.Stderr, "Error: Expected two numbers split by a '/', found %d\n", len(splitFrac))
			return
		}

		// Then convert numerator and denominator to integers
		var err error
		opts.splitN, err = strconv.Atoi(splitFrac[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Split numerator is not a number\n")
			return
		}

		opts.splitD, err = strconv.Atoi(splitFrac[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Split denominator is not a number\n")
			return
		}

		// Some other misc error checking
		if opts.splitD == 0 || opts.splitN > opts.splitD {
			fmt.Fprintf(os.Stderr, "Error: Illogical split parameter\n")
			return
		}
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

// showColumns takes a list of bools and strings and only returns the strings
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
	// If the file is split into segments, we should overwrite the start/stop lines
	if opts.split != "" {
		count, err := lineCounter(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %s while counting lines\n", fname, err)
			return
		}
		// start of the file
		opts.startLine = (opts.splitN-1)*count/opts.splitD + 1

		// end of the file
		opts.endLine = opts.splitN * count / opts.splitD

	}

	// Create a truncated csv reader, using the csvtrunc package
	csvReader, cols, err := csvtrunc.NewReader(fil, opts.startLine, opts.endLine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing %s: %s while searching for columns\n", fname, err)
		return
	}

	// Create a reader with new columns
	exprs, err := csvnewcol.CreateNewColumnExprs(opts.newcols, cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing %s: %v while creating custom columns", fname, err)
	}
	newColReader := csvnewcol.NewReader(csvReader, exprs)

	// Overwrite the header column with our new column data
	cols = append(cols, csvnewcol.GenColumns(exprs, cols)...)

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
	filteredReader := csvfilter.NewReader(newColReader, filters, !opts.orFilter)
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
