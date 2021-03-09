package csvnewcol

import (
	"errors"
	"strings"
)

type csvReader interface {
	Read() ([]string, error)
}

// Expression allows generating a new column given a row
type Expression interface {
	evaluate([]string) string
}

type stringExpr struct {
	str string
}

// Evaluate on StringExpr just returns its contained string
func (se *stringExpr) evaluate(row []string) string {
	return se.str
}

type lookupExpr struct {
	index int
}

// Evaluate on LookupExpr returns the cooresponding cell in a row
func (le *lookupExpr) evaluate(row []string) string {
	if le.index >= len(row) {
		return "missing"
	}
	return row[le.index]
}

// CreateExpressions creates a list of expressions based on a string and
// the header row (column names)
func CreateExpressions(exprString string, cols []string) ([]Expression, error) {
	var exprs []Expression

	openCurlyError := errors.New("`{` found before first one closed")
	closeCurlyError := errors.New("`}` found before opening curly brace")

	currentExpr := ""
	currentIsString := true
	for _, j := range exprString {
		// If we're looking at a {, we've entered a LookupExpr or
		// done something illegal
		if j == '{' {
			if currentIsString {
				newStrExpr := &stringExpr{currentExpr}
				exprs = append(exprs, newStrExpr)
				currentExpr = ""
				currentIsString = false
				continue
			} else {
				return exprs, openCurlyError
			}
		}
		if j == '}' {
			if currentIsString {
				return exprs, closeCurlyError
			}
			index := -1
			for ii, jj := range cols {
				if jj == currentExpr {
					index = ii
					break
				}
			}
			if index == -1 {
				return exprs, errors.New(
					"Couldn't find column " + currentExpr)
			}
			newLookupExpr := &lookupExpr{index}
			exprs = append(exprs, newLookupExpr)
			currentExpr = ""
			currentIsString = true
			continue

		}
		currentExpr += string(j)
	}

	if !currentIsString {
		return exprs, errors.New("lookup expression opened but never closed")
	}

	exprs = append(exprs, &stringExpr{currentExpr})

	return exprs, nil
}

// CreateNewColumnExprs creates an a new cell (slice of expressions) for
// each semicolon separated item
func CreateNewColumnExprs(newColString string, cols []string) ([][]Expression, error) {
	var exprs [][]Expression

	if newColString == "" {
		return exprs, nil
	}

	newcells := strings.Split(newColString, ";")
	for _, j := range newcells {
		set, err := CreateExpressions(j, cols)
		if err != nil {
			return exprs, err
		}
		exprs = append(exprs, set)
	}

	return exprs, nil
}

// GenColumns will return newly generated columns from a 2d expression slice and a row
func GenColumns(newcols [][]Expression, row []string) []string {
	var newcells []string
	for _, col := range newcols {
		var cell string
		for _, expr := range col {
			cell += expr.evaluate(row)
		}
		newcells = append(newcells, cell)
	}

	return newcells
}

// Reader reads a CSV row that may have new columns
type Reader struct {
	reader csvReader
	exprs  [][]Expression
}

// NewReader creates a reader from a csvReader interface and expression list
func NewReader(r csvReader, exprs [][]Expression) *Reader {
	return &Reader{r, exprs}
}

// Read reads a line with the added columns
func (r *Reader) Read() ([]string, error) {
	row, err := r.reader.Read()
	if err != nil {
		return []string{}, err
	}

	newCols := GenColumns(r.exprs, row)

	return append(row, newCols...), nil
}
