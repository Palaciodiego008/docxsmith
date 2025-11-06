package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Palaciodiego008/docxsmith/pkg/operations"
	"github.com/Palaciodiego008/docxsmith/pkg/pdf"
)

// HandleMerge handles the merge command
func HandleMerge(args []string) {
	fs := flag.NewFlagSet("merge", flag.ExitOnError)
	inputs := fs.String("inputs", "", "Comma-separated list of input files (required)")
	output := fs.String("output", "", "Output file path (required)")
	pageBreaks := fs.Bool("page-breaks", true, "Add page breaks between documents")
	separator := fs.Bool("separator", false, "Add separator between documents")
	separatorText := fs.String("separator-text", "---", "Separator text")
	fs.Parse(args)

	if *inputs == "" || *output == "" {
		fmt.Fprintln(os.Stderr, "Error: -inputs and -output are required")
		fs.Usage()
		os.Exit(1)
	}

	// Parse input files
	inputFiles := strings.Split(*inputs, ",")
	for i := range inputFiles {
		inputFiles[i] = strings.TrimSpace(inputFiles[i])
	}

	if len(inputFiles) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No input files provided")
		os.Exit(1)
	}

	fmt.Printf("Merging %d documents...\n", len(inputFiles))

	// Configure options
	opts := operations.MergeOptions{
		AddPageBreaks:      *pageBreaks,
		AddSeparator:       *separator,
		SeparatorText:      *separatorText,
		PreserveFormatting: true,
	}

	// Merge documents
	err := operations.MergeDocuments(inputFiles, *output, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error merging documents: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully merged %d documents into: %s\n", len(inputFiles), *output)
}

// HandleSplit handles the split command
func HandleSplit(args []string) {
	fs := flag.NewFlagSet("split", flag.ExitOnError)
	input := fs.String("input", "", "Input file path (required)")
	outputPattern := fs.String("pattern", "{base}_part{n}", "Output filename pattern")
	outputDir := fs.String("dir", ".", "Output directory")
	pages := fs.String("pages", "", "Page ranges (e.g., '1-5,7,9-12')")
	count := fs.Int("count", 0, "Split into N equal parts")
	byHeading := fs.Bool("by-heading", false, "Split by heading levels")
	headingLevel := fs.Int("heading-level", 1, "Heading level to split by (1-6)")
	fs.Parse(args)

	if *input == "" {
		fmt.Fprintln(os.Stderr, "Error: -input is required")
		fs.Usage()
		os.Exit(1)
	}

	opts := operations.SplitOptions{
		OutputPattern: *outputPattern,
		OutputDir:     *outputDir,
	}

	var outputFiles []string
	var err error

	// Determine split method
	if *byHeading {
		// Split by headings (DOCX only)
		fmt.Printf("Splitting by heading level %d...\n", *headingLevel)
		outputFiles, err = operations.SplitDOCXByHeadings(*input, *headingLevel, opts)

	} else if *count > 0 {
		// Split into N parts
		fmt.Printf("Splitting into %d parts...\n", *count)

		// Detect file type
		if strings.HasSuffix(*input, ".pdf") {
			outputFiles, err = operations.SplitPDFByCount(*input, *count, opts)
		} else {
			outputFiles, err = operations.SplitDOCXByCount(*input, *count, opts)
		}

	} else if *pages != "" {
		// Split by page ranges (PDF only)
		fmt.Printf("Splitting by page ranges: %s\n", *pages)

		// First, get page count
		doc, openErr := pdf.Open(*input)
		if openErr != nil {
			fmt.Fprintf(os.Stderr, "Error opening PDF: %v\n", openErr)
			os.Exit(1)
		}

		ranges, parseErr := operations.ParsePageRanges(*pages, doc.GetPageCount())
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "Error parsing page ranges: %v\n", parseErr)
			os.Exit(1)
		}

		outputFiles, err = operations.SplitPDFByPages(*input, ranges, opts)

	} else {
		fmt.Fprintln(os.Stderr, "Error: Must specify one of: -pages, -count, or -by-heading")
		fs.Usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error splitting document: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully split into %d files:\n", len(outputFiles))
	for _, file := range outputFiles {
		fmt.Printf("  - %s\n", file)
	}
}

// HandleMergeInfo handles the merge-info command
func HandleMergeInfo(args []string) {
	fs := flag.NewFlagSet("merge-info", flag.ExitOnError)
	inputs := fs.String("inputs", "", "Comma-separated list of input files (required)")
	fs.Parse(args)

	if *inputs == "" {
		fmt.Fprintln(os.Stderr, "Error: -inputs is required")
		fs.Usage()
		os.Exit(1)
	}

	// Parse input files
	inputFiles := strings.Split(*inputs, ",")
	for i := range inputFiles {
		inputFiles[i] = strings.TrimSpace(inputFiles[i])
	}

	// Get info based on file type
	if strings.HasSuffix(inputFiles[0], ".pdf") {
		info, err := operations.GetMergePDFInfo(inputFiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting PDF info: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Merge Information (PDF):\n")
		fmt.Printf("  Documents: %d\n", info.TotalDocuments)
		fmt.Printf("  Total Pages: %d\n", info.TotalPages)

	} else {
		info, err := operations.GetMergeDOCXInfo(inputFiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting DOCX info: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Merge Information (DOCX):\n")
		fmt.Printf("  Documents: %d\n", info.TotalDocuments)
		fmt.Printf("  Total Paragraphs: %d\n", info.TotalParagraphs)
		fmt.Printf("  Total Tables: %d\n", info.TotalTables)
	}
}
