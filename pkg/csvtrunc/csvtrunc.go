package csvtrunc

import (
	"encoding/csv"
	"io"
)

// Reader is used for reading a csv file between the start and end parameters
type Reader struct {
	reader *csv.Reader
	// The line that we're going to be reading next, 1-based indexing
	line int
	// The last line we want to read, 1-based indexing
	stop int
}

// NewReader creates a new object that automatically discards lines until start
// and only reads up until the stop line (inclusive). If stop is <= 0, we don't
// use it to calculate a stopping point. Note that this is 1-indexed, but only
// if you consider the first row (the column names) the 0th index. Note that the
// first row is returned, so that column names can be known.
func NewReader(r io.Reader, start int, stop int) (*Reader, []string, error) {

	// Create a csv.Reader
	csvReader := csv.NewReader(r)
	// Read the first row, which are the column names
	cols, err := csvReader.Read()
	if err != nil {
		return nil, cols, err
	}

	// Discard the first start-1 lines (start including
	for i := 1; i < start; i++ {
		_, err := csvReader.Read()
		// there shouldn't be any errors reading the first start lines
		if err != nil {
			return nil, cols, err
		}

	}

	// Note that we read start-1 lines, so we're at line start
	return &Reader{csvReader, start, stop}, cols, nil
}

// Read reads a line from the csv. It checks to make sure that the line is before end
func (r *Reader) Read() ([]string, error) {
	// If we reach our stopping line, we stop
	if r.line > r.stop && r.stop > 0 {
		return []string{}, io.EOF
	}

	row, err := r.reader.Read()
	if err == nil {
		// If there wasn't an error we should keep track of what line we read
		r.line++
	}

	return row, err
}

// Line returns the current line number
func (r *Reader) Line() int {
	return r.line
}
