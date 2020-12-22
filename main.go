package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/appins/csvcmd/pkg/csvfilter"
)

func main() {
	var humanReadable bool
	flag.BoolVar(&humanReadable, "h", false, "Print in an easy to read format")

	var startLine int
	flag.IntVar(&startLine, "s", 0, "The first line that should be read (inclusive)")

	var endLine int
	flag.IntVar(&endLine, "e", -1, "The last line that should be read (inclusive)")

	flag.Parse()

	// The files that should be read
	if len(flag.Args()) != 0 {
		printNames := len(flag.Args()) > 1
		for _, fil := range flag.Args() {
			fmt.Printf("%s %s %t \n", fil, "hello", printNames)
		}

	} else {
		// In the case of no files being specified, read from stdin
		csvReader := csv.NewReader(os.Stdin)
		filters := []func([]string) bool{}
		filteredReader := csvfilter.NewReader(csvReader, filters, false)
		for filteredReader.Scan() {
			fmt.Println(filteredReader.Row())
		}

	}

}
