package utils

import (
	"encoding/csv"
	"os"
)

func AppendToCsv(filename string, records [][]string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	return writer.WriteAll(records)
}
