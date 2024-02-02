package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/adrg/xdg"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	//Parse CLI arguments
	fileName := flag.String("f", "til.csv", "Output database name (alias for --filename)")
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
	db, err := setupDatabase(dataPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer db.Close()

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
		err = listDbContents(db)
		if err != nil {
			fmt.Println("Error:", err)
		}
		return
	}

	inputArgs := flag.Args()
	if len(inputArgs) > 0 {
		entry := strings.Join(inputArgs, " ")
		data := [][]string{[]string{"", entry}}
		err = saveToDb(db, data)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Text saved to %s\n", dataPath)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	var data [][]string

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
			lines = append(lines, strings.TrimSpace(line))

		}
		completeText := strings.Join(lines, "\\n")
		data = append(data, []string{stamp(), completeText})
	}
	//
	err = saveToDb(db, data)
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
  -f, --filename      Set output file name (default: "til.csv")
  -l, --list          List all TILs
  -h, --help          Show help
  -p, --print-path    Print path to the default output file

Notes:
  If additional input (other than the flags) is provided, it will be treated as an entry and saved to the database.

Examples:
  1. Save a text entry directly:
     til Learned about the Go language today

  2. Read input from the user and store it:
     til -f output.csv

  3. Read input from pipe:
     echo "Piped input" | til

  4. List all the contents from the database:
     til -l

  5. Print the path to the default output file:
     til --print-path`
	fmt.Println(helpText)
}

func printDefaultFilePath(filename string) {
	dataPath, err := xdg.DataFile(filepath.Join("til", filename))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Default output file path:", dataPath)
}
