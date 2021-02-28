package main

import (
	"errors"
	"fmt"
	"strings"
)

func genWidths(row []string) []int {
	var lengths []int
	for _, j := range row {
		lenJ := len(j) + 2

		// Rows should not be longer than 10, nor shorter than 5
		if lenJ > 10 {
			lenJ = 10
		} else if lenJ < 5 {
			lenJ = 5
		}

		lengths = append(lengths, lenJ)
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

func genColumns(columnFlag string, cols []string) ([]bool, error) {
	// Create a map of the columns, just like in genFilters
	colsToInt := make(map[string]int)
	for i, col := range cols {
		colsToInt[col] = i
		col_num := fmt.Sprintf("_%d", i+1)
		colsToInt[col_num] = i
	}

	// Each column has a bool assocaited with it, if it should show or not
	enabled := make([]bool, len(cols))

	// If nothing specified, we just show all the columns (each bool in the array = true)
	if len(columnFlag) == 0 {
		for i := range enabled {
			enabled[i] = true
		}
	} else {
		// Otherwise we split up the flag by ; and set the bools for the columns in that list
		shownCols := strings.Split(columnFlag, ";")
		for _, j := range shownCols {
			if col, ok := colsToInt[j]; ok {
				enabled[col] = true
			} else {
				// If we can't find the column we error out
				return enabled, errors.New("Couldn't find column " + j)
			}
		}
	}
	return enabled, nil
}
