package main

import (
	"errors"
	"flag"
	"fmt"
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

func main() {
	// Manual test
	fileData, err := getFileData()

	fmt.Println(fileData, err)
}
