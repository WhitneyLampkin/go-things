package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
		return nil, errors.New("line doesn't match headers format (skipping)")
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

func writeJSONFile(csvPath string, writerChannel <-chan map[string]string, done chan<- bool, pretty bool) {
	writeString := createStringWriter(csvPath) // Instanciating a JSON writer function
	jsonFunc, breakLine := getJSONFunc(pretty) // Instanciating the JSON parse function and the breakline character
	// Log for informing
	fmt.Println("Writing JSON file...")
	// Writing the first character of our JSON file. We always start with a "[" since we always generate array of record
	writeString("["+breakLine, false)
	first := true
	for {
		// Waiting for pushed records into our writerChannel
		record, more := <-writerChannel
		if more {
			if !first { // If it's not the first record, we break the line
				writeString(","+breakLine, false)
			} else {
				first = false // If it's the first one, we don't break the line
			}

			jsonData := jsonFunc(record) // Parsing the record into JSON
			writeString(jsonData, false) // Writing the JSON string with our writer function
		} else { // If we get here, it means there aren't more record to parse. So we need to close the file
			writeString(breakLine+"]", true) // Writing the final character and closing the file
			fmt.Println("Completed!")        // Logging that we're done
			fmt.Println("The new json file has been added to the current directory with the same name as the one provided.")
			done <- true // Sending the signal to the main function so it can correctly exit out.
			break        // Stoping the for-loop
		}
	}
}

func createStringWriter(csvPath string) func(string, bool) {
	jsonDir := filepath.Dir(csvPath)                                                       // Getting the directory where the CSV file is
	jsonName := fmt.Sprintf("%s.json", strings.TrimSuffix(filepath.Base(csvPath), ".csv")) // Declaring the JSON filename, using the CSV file name as base
	finalLocation := filepath.Join(jsonDir, jsonName)                                      // Declaring the JSON file location, using the previous variables as base
	// Opening the JSON file that we want to start writing
	f, err := os.Create(finalLocation)
	check(err)
	// This is the function we want to return, we're going to use it to write the JSON file
	return func(data string, close bool) { // 2 arguments: The piece of text we want to write, and whether or not we should close the file
		_, err := f.WriteString(data) // Writing the data string into the file
		check(err)
		// If close is "true", it means there are no more data left to be written, so we close the file
		if close {
			f.Close()
		}
	}
}

func getJSONFunc(pretty bool) (func(map[string]string) string, string) {
	// Declaring the variables we're going to return at the end
	var jsonFunc func(map[string]string) string
	var breakLine string
	if pretty { //Pretty is enabled, so we should return a well-formatted JSON file (multi-line)
		breakLine = "\n"
		jsonFunc = func(record map[string]string) string {
			jsonData, _ := json.MarshalIndent(record, "   ", "   ") // By doing this we're ensuring the JSON generated is indented and multi-line
			return "   " + string(jsonData)                         // Transforming from binary data to string and adding the indent characets to the front
		}
	} else { // Now pretty is disabled so we should return a compact JSON file (one single line)
		breakLine = "" // It's an empty string because we never break lines when adding a new JSON object
		jsonFunc = func(record map[string]string) string {
			jsonData, _ := json.Marshal(record) // Now we're using the standard Marshal function, which generates JSON without formating
			return string(jsonData)             // Transforming from binary data to string
		}
	}

	return jsonFunc, breakLine // Returning everythinbg
}

func main() {
	// Provide help info when users enter the --help option using the anonymous function below
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <csvFile>\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Get file data entered by user
	fileData, err := getFileData()
	if err != nil {
		exitGracefully(err)
	}

	// Validate the file provided if no errors
	if _, err := checkIfValidFile(fileData.filepath); err != nil {
		exitGracefully(err)
	}

	// Declare channels for the go-routines to use
	writerChannel := make(chan map[string]string)
	done := make(chan bool)

	// Starting the go-routines used for reading and writing
	// These two will run asynchronously
	go processCsvFile(fileData, writerChannel)
	go writeJSONFile(fileData.filepath, writerChannel, done, fileData.pretty)

	// Terminate program execution when done
	<-done
}
