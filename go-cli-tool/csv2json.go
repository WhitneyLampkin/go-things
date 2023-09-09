package main

import (
	"errors"
	"flag"
	"os"
)

// Struct that stores the user inputted file information needed
type inputFile struct {
	filepath  string
	separator string
	pretty    bool
}

// <TODO: Should an interface be used here?>
func getFileData() (inputFile, error) {
	if len(os.Args) < 2 {
		return inputFile{}, errors.New("A filepath argument is required.")
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

	// Validation for separator value
	if !(*separator == "comma" || *separator == "semicolon") {
		return inputFile{}, errors.New("Invalid separator. Use comma or semicolon.")
	}

	// After validation, return struct with arguments
	return inputFile{fileLocation, *separator, *pretty}, nil
}

func main() {

}
