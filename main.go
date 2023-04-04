package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/xdg"
)

func main() {

	//Parse CLI arguments
	fileName := flag.String("filename", "output.csv", "Output file name")
	listFlag := flag.Bool("list", false, "List all TILs")

	flag.Parse()

	dataPath, err := xdg.DataFile(filepath.Join("til", *fileName))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// list all the TILs

	if *listFlag {
		err = listCsvContents(dataPath)
		if err != nil {
			fmt.Println("Error:", err)
		}
		return
	}
	reader := bufio.NewReader(os.Stdin)
	var data [][]string

	// scanner := bufio.NewScanner(reader)
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) != os.ModeCharDevice {
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error:", err)
				return
			}
			// line = strings.TrimSpace(line)
			// if line == "" {
			// 	break
			// }
			// timestamp := stamp() 
			data = append(data, []string{stamp(), strings.TrimSpace(line)})
		}

	} else {

		fmt.Println("Enter text to save (Shift + Enter for newline, Enter to save):")
		lines := []string{}
		for {
			line, err := reader.ReadString('\n')
			if errors.Is(err, io.EOF) || (err == nil && strings.TrimSpace(line) == "") {
				break
			} else if err != nil {
				fmt.Println("Error:", err)
				return
			}
			// line = strings.TrimSpace(line)
			// if line == "" {
			// 	break
			// }
			lines = append(lines, strings.TrimSpace(line))

		}
		// timestamp := stamp() 
		completeText := strings.Join(lines, "\\n")
		data = append(data, []string{stamp(), completeText})

	}
	// for scanner.Scan(){
	//   line := scanner.Text()
	//   timestamp := time.Now().Format("2006-01-02 15:04:05")
	//   data = append(data, []string{timestamp, line})
	// }
	// if err := scanner.Err(); err != nil{
	//   fmt.Println("Error:", err)
	//   return
	// }

	err = saveToCsvFile(dataPath, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Text saved to %s\n", dataPath)
}

func stamp () string {
  return time.Now().Format("2006-01-02 15:04:05") 
}

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
