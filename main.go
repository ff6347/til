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
	fileName := flag.String("f", "output.csv", "Output file name (alias for --filename)")
	fileNameLong := flag.String("filename", "", "Output file name")
	listFlag := flag.Bool("l", false, "List all TILs (alias for --list)")
	listFlagLong := flag.Bool("list", false, "List all TILs")
	helpFlag := flag.Bool("h", false, "Show help (alias for --help)")
	helpFlagLong := flag.Bool("help", false, "Show help")
	printPathFlagLong := flag.Bool("print-path", false, "Print path to the default output file")
	printPathFlag := flag.Bool("p", false, "Print path to the default output file (alias for --print-path)")
	flag.Parse()

	if *fileNameLong != "" {
		*fileName = *fileNameLong
	}

	if *listFlagLong {
		*listFlag = true
	}

	if *helpFlagLong {
		*helpFlag = true
	}
	if *printPathFlagLong {
		*printPathFlag = true
	}

	dataPath, err := xdg.DataFile(filepath.Join("til", *fileName))

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if *helpFlag {
		showHelp()
		return
	}
	if *printPathFlag {
		printDefaultFilePath(*fileName)
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

func stamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func showHelp() {
	helpText := `
TIL is a simple command-line tool for storing and listing short learnings.

Usage:
  til [flags] [--] [arguments ...]

Flags:
  -f, --filename      Set output file name (default: "output.csv")
  -l, --list          List all TILs
  -h, --help          Show help
  -p, --print-path    Print path to the default output file

Examples:
  1. Read input from the user and store it:
     til -f output.csv

  2. Read input from pipe:
     echo "Piped input" | til

  3. List all the contents from the CSV file:
     til -l

  4. Print the path to the default output file:
     til --print-path  `
	fmt.Println(helpText)
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

func printDefaultFilePath(filename string) {
	dataPath, err := xdg.DataFile(filepath.Join("til", filename))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Default output file path:", dataPath)
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
