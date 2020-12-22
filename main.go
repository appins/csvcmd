package main

import (
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
			processFile(io.Reader(fil), flag.Args()[i], opts)
		}

	} else {
		// In the case of no files being specified, read from stdin
		processFile(io.Reader(os.Stdin), "STDIN", opts)
	}

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

func processFile(fil io.Reader, fname string, opts options) {
	csvReader, cols, err := csvtrunc.NewReader(fil, opts.startLine, opts.endLine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", fname, err)
		return
	}

	filters, err := genFilters(opts.filtersString, cols)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", fname, err)
		return
	}
	filteredReader := csvfilter.NewReader(csvReader, filters, true)
	for filteredReader.Scan() {
		fmt.Println(filteredReader.Row())
	}

}
