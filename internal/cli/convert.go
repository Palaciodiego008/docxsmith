package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Palaciodiego008/docxsmith/pkg/converter"
)

// HandleConvert handles the convert command
func HandleConvert(args []string) {
	fs := flag.NewFlagSet("convert", flag.ExitOnError)
	input := fs.String("input", "", "Input file path (required)")
	output := fs.String("output", "", "Output file path (required)")
	pageSize := fs.String("page-size", "A4", "Page size (A4, Letter, Legal)")
	fontSize := fs.Float64("font-size", 12, "Default font size")
	fontFamily := fs.String("font-family", "Arial", "Default font family")
	fs.Parse(args)

	if *input == "" || *output == "" {
		fmt.Fprintln(os.Stderr, "Error: -input and -output are required")
		fs.Usage()
		os.Exit(1)
	}

	// Determine conversion direction based on file extensions
	inputExt := strings.ToLower(filepath.Ext(*input))
	outputExt := strings.ToLower(filepath.Ext(*output))

	opts := converter.ConvertOptions{
		PageSize:    *pageSize,
		Orientation: "Portrait",
		FontSize:    *fontSize,
		FontFamily:  *fontFamily,
		Margins:     [4]float64{20, 20, 20, 20},
	}

	var err error

	switch {
	case inputExt == ".docx" && outputExt == ".pdf":
		fmt.Println("Converting DOCX to PDF...")
		err = converter.ConvertDocxToPDF(*input, *output, opts)

	case inputExt == ".pdf" && outputExt == ".docx":
		fmt.Println("Converting PDF to DOCX...")
		err = converter.ConvertPDFToDocx(*input, *output, opts)

	default:
		fmt.Fprintf(os.Stderr, "Error: Unsupported conversion from %s to %s\n", inputExt, outputExt)
		fmt.Fprintln(os.Stderr, "Supported conversions:")
		fmt.Fprintln(os.Stderr, "  - .docx to .pdf")
		fmt.Fprintln(os.Stderr, "  - .pdf to .docx")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting document: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Conversion successful: %s -> %s\n", *input, *output)
}
