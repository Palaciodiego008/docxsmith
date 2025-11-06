package cli

import (
	"flag"
	"fmt"
	"os"
)

// Common error messages
const (
	ErrMissingInput  = "input file is required"
	ErrMissingOutput = "output file is required"
	ErrMissingText   = "text is required"
	ErrMissingData   = "data file is required"
)

// ValidateRequired checks if required parameters are provided
func ValidateRequired(params map[string]string) error {
	for name, value := range params {
		if value == "" {
			return fmt.Errorf("%s is required", name)
		}
	}
	return nil
}

// ValidateFileExists checks if a file exists
func ValidateFileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}
	return nil
}

// ExitWithError prints an error and exits
func ExitWithError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

// PrintSuccess prints a success message
func PrintSuccess(format string, args ...interface{}) {
	fmt.Printf("✓ "+format+"\n", args...)
}

// PrintInfo prints an informational message
func PrintInfo(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// FormatList formats a list of items for display
func FormatList(items []string, indent string) string {
	result := ""
	for _, item := range items {
		result += indent + "- " + item + "\n"
	}
	return result
}

// CommonFlags represents commonly used flags across commands
type CommonFlags struct {
	Input  string
	Output string
}

// DocumentFlags represents flags for document operations
type DocumentFlags struct {
	CommonFlags
	Text   string
	Bold   bool
	Italic bool
	Size   string
	Color  string
	Align  string
}

// AddCommonFlags adds common flags to a FlagSet
func AddCommonFlags(fs *flag.FlagSet) (*string, *string) {
	input := fs.String("input", "", "Input file path")
	output := fs.String("output", "", "Output file path")
	return input, output
}

// AddTextFormattingFlags adds text formatting flags to a FlagSet
func AddTextFormattingFlags(fs *flag.FlagSet) (*bool, *bool, *string, *string, *string) {
	bold := fs.Bool("bold", false, "Make text bold")
	italic := fs.Bool("italic", false, "Make text italic")
	size := fs.String("size", "", "Font size")
	color := fs.String("color", "", "Text color (hex without #)")
	align := fs.String("align", "", "Alignment (left, center, right, both)")
	return bold, italic, size, color, align
}
