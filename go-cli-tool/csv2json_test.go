package main_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Csv2json", func() {
	BeforeEach(func() {
		tests := []struct {
			name    string    // The name of the test
			want    inputFile // What inputFile instance we want our function to return.
			wantErr bool      // whether or not we want an error.
			osArgs  []string  // The command arguments used for this test
		}{
			// Here we're declaring each unit test input and output data as defined before
			{"Default parameters", inputFile{"test.csv", "comma", false}, false, []string{"cmd", "test.csv"}},
			{"No parameters", inputFile{}, true, []string{"cmd"}},
			{"Semicolon enabled", inputFile{"test.csv", "semicolon", false}, false, []string{"cmd", "--separator=semicolon", "test.csv"}},
			{"Pretty enabled", inputFile{"test.csv", "comma", true}, false, []string{"cmd", "--pretty", "test.csv"}},
			{"Pretty and semicolon enabled", inputFile{"test.csv", "semicolon", true}, false, []string{"cmd", "--pretty", "--separator=semicolon", "test.csv"}},
			{"Separator not identified", inputFile{}, true, []string{"cmd", "--separator=pipe", "test.csv"}},
		}
	})

	Describe("Get file data", func() {

	})
})
