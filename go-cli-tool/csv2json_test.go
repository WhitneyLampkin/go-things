package main

import (
	"errors"
	"flag"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Csv2json", func() {
	var actualOsArgs []string
	var defaultInputs, noInputs, semicolonInput, prettyInput, semicolonPrettyInput, unknownSeparatorInput *inputFile

	BeforeEach(func() {
		actualOsArgs = os.Args

		defaultInputs = &inputFile{
			filepath:  "test.csv",
			separator: "comma",
			pretty:    false,
		}

		noInputs = &inputFile{}

		semicolonInput = &inputFile{
			filepath:  "test.csv",
			separator: "semicolon",
			pretty:    false,
		}

		prettyInput = &inputFile{
			filepath:  "test.csv",
			separator: "comma",
			pretty:    true,
		}

		semicolonPrettyInput = &inputFile{
			filepath:  "test.csv",
			separator: "semicolon",
			pretty:    true,
		}

		unknownSeparatorInput = &inputFile{}
	})

	AfterEach(func() {
		os.Args = actualOsArgs                                           // Restoring the original os.Args reference
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // Reseting the Flag command line. So that we can parse flags again
	})

	Describe("Get file data", func() {
		Context("when using default parameters", func() {
			It("should return the filepath provided, comma and false", func() {
				actualOsArgs = os.Args
				os.Args = []string{"cmd", "test.csv"}
				result, err := getFileData()

				Expect(err).To(BeNil())
				Expect(result.filepath).To(Equal(defaultInputs.filepath))
			})
		})

		Context("when using no parameters", func() {
			It("should return an error", func() {
				actualOsArgs = os.Args
				os.Args = []string{"cmd"}
				expectedErr := errors.New("a filepath argument is required")
				result, err := getFileData()

				Expect(err).To(Equal(expectedErr))
				Expect(result).To(Equal(*noInputs))
				Expect(result.filepath).To(Equal(noInputs.filepath))
			})
		})

		Context("when semicolon is enabled", func() {
			It("should return semicolon as the separator", func() {
				actualOsArgs = os.Args
				os.Args = []string{"cmd", "--separator=semicolon", "test.csv"}
				result, err := getFileData()

				Expect(err).To(BeNil())
				Expect(result.separator).To(Equal(semicolonInput.separator))
			})
		})

		Context("when pretty is enabled", func() {
			It("should return true", func() {
				actualOsArgs = os.Args
				os.Args = []string{"cmd", "--pretty", "test.csv"}
				result, err := getFileData()

				Expect(err).To(BeNil())
				Expect(result.pretty).To(Equal(prettyInput.pretty))
			})
		})

		Context("when semicolon and pretty are enabled", func() {
			It("should return semicolon and true", func() {
				actualOsArgs = os.Args
				os.Args = []string{"cmd", "--pretty", "--separator=semicolon", "test.csv"}
				result, err := getFileData()

				Expect(err).To(BeNil())
				Expect(result.separator).To(Equal(semicolonPrettyInput.separator))
				Expect(result.pretty).To(Equal(semicolonPrettyInput.pretty))
			})
		})

		Context("when an invalid separator is provided", func() {
			It("should return an error and empty inputFile type", func() {
				actualOsArgs = os.Args
				os.Args = []string{"cmd", "--separator=pipe", "test.csv"}
				expectedErr := errors.New("invalid separator. Use comma or semicolon")
				result, err := getFileData()

				Expect(err).To(Equal(expectedErr))
				Expect(result).To(Equal(*unknownSeparatorInput))
			})
		})
	})
})

/*

//////////// Auto-generated unit test written by VSCode & later modified ////////////

func Test_getFileData(t *testing.T) {
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
