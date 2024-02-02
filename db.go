package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func setupDatabase(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func saveToDb(filePath string, content [][]string) error {
	// Open the file in append mode, create it if it doesn't exist.
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, entry := range content {
		if err := writer.Write(entry); err != nil {
			return err
		}
	}

	return nil
}

func listDbContents(filePath string) error {
	// Open the file in read-only mode.
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

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
