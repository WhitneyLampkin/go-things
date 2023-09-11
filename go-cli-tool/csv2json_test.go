package main

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Csv2json", func() {
	var one int
	var defaultInputs *inputFile //, noInputs, semicolonInput, prettyInput, semicolonPrettyInput, noSeparatorInput *inputFile

	BeforeEach(func() {
		one = 1

		defaultInputs = &inputFile{
			filepath:  "test.csv",
			separator: "comma",
			pretty:    false,
		}
	})

	Describe("Get file data", func() {
		Context("when testing the value of one", func() {
			It("should return 2 after the count is incremented", func() {
				one++
				Expect(one).To(Equal(2))
			})
		})

		Context("when using default parameters", func() {
			It("should return the filepath provided, comma and false", func() {
				//actualOsArgs := os.Args
				os.Args = []string{"cmd", "test.csv"}
				result, err := getFileData()

				Expect(err).To(BeNil())
				Expect(result.filepath).To(Equal(defaultInputs.filepath))
			})
		})
	})
})

/* func Test_getFileData(t *testing.T) {
	tests := []struct {
		name    string    // The name of the test
		want    inputFile // What inputFile instance we want our function to return.
		wantErr bool      // whether or not we want an error.
		osArgs  []string  // The command arguments used for this test
	}{
		// Declaring each unit test input and output data as defined before
		{"Default parameters", inputFile{"test.csv", "comma", false}, false, []string{"cmd", "test.csv"}},
		{"No parameters", inputFile{}, true, []string{"cmd"}},
		{"Semicolon enabled", inputFile{"test.csv", "semicolon", false}, false, []string{"cmd", "--separator=semicolon", "test.csv"}},
		{"Pretty enabled", inputFile{"test.csv", "comma", true}, false, []string{"cmd", "--pretty", "test.csv"}},
		{"Pretty and semicolon enabled", inputFile{"test.csv", "semicolon", true}, false, []string{"cmd", "--pretty", "--separator=semicolon", "test.csv"}},
		{"Separator not identified", inputFile{}, true, []string{"cmd", "--separator=pipe", "test.csv"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Saving the original os.Args reference
			actualOsArgs := os.Args
			// This defer function will run after the test is done
			defer func() {
				os.Args = actualOsArgs                                           // Restoring the original os.Args reference
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // Reseting the Flag command line. So that we can parse flags again
			}()

			os.Args = tt.osArgs // Setting the specific command args for this test
			got, err := getFileData()
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFileData() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
