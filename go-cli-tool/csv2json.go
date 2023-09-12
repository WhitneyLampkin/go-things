package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Struct that stores the user inputted file information needed
type inputFile struct {
	filepath  string
	separator string
	pretty    bool
}

// Gets user input from the terminal, validates it and returns the inputFile data and an error
// <TODO: Should an interface be used here?>
func getFileData() (inputFile, error) {
	// Validate the number of arguments
	if len(os.Args) < 2 {
		return inputFile{}, errors.New("a filepath argument is required")
	}

	// Define flags for the separator and pretty arguments
	// The flag package is used for command-line flag parsing
	// Flags take 3 arguments: name, default value and help description
	separator := flag.String("separator", "comma", "Column separator used in the csv file")
	pretty := flag.Bool("pretty", false, "Generate pretty JSON")

	// Called after all flag definitions
	// Parses arguments from the terminal
	flag.Parse()

	// Non-flag argument
	fileLocation := flag.Arg(0)

	// Validate separator value
	if !(*separator == "comma" || *separator == "semicolon") {
		return inputFile{}, errors.New("invalid separator. Use comma or semicolon")
	}

	fmt.Println(*pretty)

	// Return inputFile with validated arguments
	return inputFile{fileLocation, *separator, *pretty}, nil
}

func checkIfValidFile(fileLocation string) (bool, error) {
	// Validate file extension (CSV)
	if fileExtension := filepath.Ext(fileLocation); fileExtension != ".csv" {
		return false, fmt.Errorf("file %s is not CSV", fileLocation)
	}

	// Validate file exists
	if _, err := os.Stat(fileLocation); err != nil && os.IsNotExist(err) {
		return false, fmt.Errorf("file %s does not exist", fileLocation)
	}

	// If no errors, the file is valid.
	return true, nil
}

// Though go-channels and go-routines are unnecessary, they tutorial uses them for learning purposes
// This is the first go-routine that both reads and processes each line of the csv file
func processCsvFile(fileData inputFile, writerChannel chan<- map[string]string) {
	// Open the file for reading
	file, err := os.Open(fileData.filepath)
	// Check that no errors were received
	check(err)
	// Ensure that the file is closed after processing
	defer file.Close()

	// Defining the "headers" and "line" slices to represent the contents of the file
	var headers, line []string
	// Initialize reader
	reader := csv.NewReader(file)
	// Conditionally update the separator value if semicolon was entered
	if fileData.separator == "semicolon" {
		reader.Comma = ';'
	}

	// Start reading the file
	// First line will be the headers
	headers, err = reader.Read()
	// Check that no errors were received
	check(err)
	// Iterate over each line of the csv file
	for {
		// Read 1 row (line) of the csv file at a time
		// Each element is a column
		line, err = reader.Read()
		// Check to see if the end of the file was reached
		if err != nil {
			// Prefered way by GoLang doc
			if errors.Is(err, io.EOF) {
				fmt.Println("Reading file finished...")
				close(writerChannel)
				break
			} else {
				// Handle unexpected errors
				exitGracefully(err)
			}
		}
		// Process the current CSV line
		// Each line is returned as a map[string]string with the key being the column name
		record, err := processLine(headers, line)
		if err != nil {
			// The error returned from processLine would be for an incorrect number of items
			// This record should be skipped
			fmt.Printf("Line: %sError: %s\n", line, err)
			continue
		}

		// If the record is valid, write it to the channel
		writerChannel <- record
	}

}

func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func check(e error) {
	if e != nil {
		exitGracefully(e)
	}
}

func processLine(headers []string, dataList []string) (map[string]string, error) {
	// Validating if we're getting the same number of headers and columns. Otherwise, we return an error
	if len(dataList) != len(headers) {
		return nil, errors.New("Line doesn't match headers format. Skipping")
	}
	// Creating the map we're going to populate
	recordMap := make(map[string]string)
	// For each header we're going to set a new map key with the corresponding column value
	for i, name := range headers {
		recordMap[name] = dataList[i]
	}
	// Returning our generated map
	return recordMap, nil
}

func main() {
	// Manual test
	fileData, err := getFileData()

	fmt.Println(fileData, err)
}
