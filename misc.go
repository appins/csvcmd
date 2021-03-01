package main

import (
	"encoding/csv"
	"os"
)

func lineCounter(fileName string) (int, error) {
	fil, err := os.Open(fileName)
	if err != nil {
		return -1, err
	}
	defer fil.Close()

	csvReader := csv.NewReader(fil)
	ctr := 0
	for {
		_, err := csvReader.Read()
		if err != nil {
			break
		}
		ctr++
	}

	return ctr, nil

}
