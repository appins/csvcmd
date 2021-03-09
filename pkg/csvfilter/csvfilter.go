package csvfilter

import "io"

type csvReader interface {
	Read() ([]string, error)
}

// Reader reads only rows that meet the filters from a csv filter
type Reader struct {
	reader     csvReader
	filters    []func([]string) bool
	filterType bool
	row        []string
	err        error
}

const (
	OR_FILTER  = false
	AND_FILTER = true
)

// NewReader creates a reader from a csvReader interface
func NewReader(r csvReader, filters []func([]string) bool, filterType bool) *Reader {
	return &Reader{r, filters, filterType, []string{}, nil}
}

// scanReader is used to scan the underlying reader, non-filtered rows
func (r *Reader) scanReader() bool {
	r.row, r.err = r.reader.Read()
	if r.err != nil {
		return false
	}
	return true
}

// Scan loads a row and err and returns if it could find a row that met
// the specified filters
func (r *Reader) Scan() bool {
	for r.scanReader() {
		// If zero filters, then we let the row through
		if len(r.filters) == 0 {
			return true
		}

		// if 1+ filters, we OR/AND the results together
		valid := r.filterType
		for _, j := range r.filters {
			result := j(r.row)
			if r.filterType == OR_FILTER {
				valid = valid || result
			} else {
				valid = valid && result
			}
		}

		if valid {
			return true
		}
	}
	// scanReader will already have put an error in err
	return false
}

// Returns an error, if there was one
func (r *Reader) Err() error {
	if r.err == io.EOF {
		return nil
	}
	return r.err
}

func (r *Reader) Row() []string {
	return r.row
}
