package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func setupDatabase(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func saveToDb(file *os.File, content [][]string) error {
	writer := csv.NewWriter(file)
	for _, entry := range content {
		err := writer.Write(entry)
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

func listDbContents(file *os.File) error {
	_, _ = file.Seek(0, 0)
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	fmt.Println("Timestamp\t\t\tContent")
	for _, record := range records {
		fmt.Println(record[0] + "\t\t" + record[1])
	}
	return nil
}
