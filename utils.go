package main



import (
  "os"
  "encoding/csv"
  "fmt"
)

func listCsvContents(filePath string) error {

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer f.Close()
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	fmt.Printf("Timestamp\t\t\tTIL\n")
	for _, row := range data {
		fmt.Printf("%s\t%s\n", row[0], row[1])
	}
	return nil
}



func saveToCsvFile(filePath string, content [][]string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	err = csvWriter.WriteAll(content)
	return err
}

func saveToFile(filePath, content string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

