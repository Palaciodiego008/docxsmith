package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Palaciodiego008/docxsmith/pkg/pdf"
)

// HandlePDFCreate handles the PDF create command
func HandlePDFCreate(args []string) {
	fs := flag.NewFlagSet("pdf-create", flag.ExitOnError)
	output := fs.String("output", "", "Output PDF file path (required)")
	text := fs.String("text", "", "Initial text content")
	title := fs.String("title", "", "Document title")
	author := fs.String("author", "", "Document author")
	fs.Parse(args)

	if *output == "" {
		fmt.Fprintln(os.Stderr, "Error: -output is required")
		fs.Usage()
		os.Exit(1)
	}

	doc := pdf.New()

	// Set metadata
	if *title != "" || *author != "" {
		doc.SetMetadata(*title, *author, "")
	}

	// Add first page
	page := doc.AddPage()

	// Add content if provided
	if *text != "" {
		page.AddText(*text, 20, 30, 12)
	} else {
		page.AddText("PDF document created with DocxSmith", 20, 30, 12)
	}

	if err := doc.Save(*output); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("PDF created successfully: %s\n", *output)
}

// HandlePDFAdd handles adding content to PDF
func HandlePDFAdd(args []string) {
	fs := flag.NewFlagSet("pdf-add", flag.ExitOnError)
	input := fs.String("input", "", "Input PDF file path (required)")
	output := fs.String("output", "", "Output PDF file path (required)")
	text := fs.String("text", "", "Text to add (required)")
	bold := fs.Bool("bold", false, "Make text bold")
	italic := fs.Bool("italic", false, "Make text italic")
	size := fs.Float64("size", 12, "Font size")
	color := fs.String("color", "000000", "Text color (hex without #)")
	fs.Parse(args)

	if *input == "" || *output == "" || *text == "" {
		fmt.Fprintln(os.Stderr, "Error: -input, -output, and -text are required")
		fs.Usage()
		os.Exit(1)
	}

	doc, err := pdf.Open(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening PDF: %v\n", err)
		os.Exit(1)
	}

	// Get the last page or create a new one
	var page *pdf.Page
	if len(doc.Pages) > 0 {
		page = doc.Pages[len(doc.Pages)-1]
	} else {
		page = doc.AddPage()
	}

	// Add text with styling
	style := pdf.TextStyle{
		FontSize:   *size,
		FontFamily: "Arial",
		Bold:       *bold,
		Italic:     *italic,
		Color:      *color,
	}

	// Calculate Y position (simple placement)
	y := 30.0 + float64(len(page.Content))*10
	page.AddTextStyled(*text, 20, y, style)

	if err := doc.Save(*output); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Content added successfully to: %s\n", *output)
}

// HandlePDFInfo handles displaying PDF info
func HandlePDFInfo(args []string) {
	fs := flag.NewFlagSet("pdf-info", flag.ExitOnError)
	input := fs.String("input", "", "Input PDF file path (required)")
	fs.Parse(args)

	if *input == "" {
		fmt.Fprintln(os.Stderr, "Error: -input is required")
		fs.Usage()
		os.Exit(1)
	}

	doc, err := pdf.Open(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("PDF Document Information: %s\n", *input)
	fmt.Printf("  Pages: %d\n", doc.GetPageCount())

	if doc.Metadata != nil {
		if doc.Metadata.Title != "" {
			fmt.Printf("  Title: %s\n", doc.Metadata.Title)
		}
		if doc.Metadata.Author != "" {
			fmt.Printf("  Author: %s\n", doc.Metadata.Author)
		}
	}

	// Count total text content
	text := doc.GetAllText()
	wordCount := len(strings.Fields(text))
	charCount := len(text)

	fmt.Printf("  Words: %d\n", wordCount)
	fmt.Printf("  Characters: %d\n", charCount)
}

// HandlePDFExtract handles extracting text from PDF
func HandlePDFExtract(args []string) {
	fs := flag.NewFlagSet("pdf-extract", flag.ExitOnError)
	input := fs.String("input", "", "Input PDF file path (required)")
	output := fs.String("output", "", "Output text file (optional)")
	fs.Parse(args)

	if *input == "" {
		fmt.Fprintln(os.Stderr, "Error: -input is required")
		fs.Usage()
		os.Exit(1)
	}

	doc, err := pdf.Open(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening PDF: %v\n", err)
		os.Exit(1)
	}

	text := doc.GetAllText()

	if *output != "" {
		if err := os.WriteFile(*output, []byte(text), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Text extracted to: %s\n", *output)
	} else {
		fmt.Println(text)
	}
}
