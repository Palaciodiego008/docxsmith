package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

// HandleAdd handles the add command
func HandleAdd(args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	input := fs.String("input", "", "Input file path (required)")
	output := fs.String("output", "", "Output file path (required)")
	text := fs.String("text", "", "Text to add (required)")
	index := fs.Int("at", -1, "Insert at specific index (default: append)")
	bold := fs.Bool("bold", false, "Make text bold")
	italic := fs.Bool("italic", false, "Make text italic")
	size := fs.String("size", "", "Font size (e.g., '24' for 12pt)")
	color := fs.String("color", "", "Text color (hex without #, e.g., 'FF0000')")
	align := fs.String("align", "", "Alignment: left, center, right, both")
	fs.Parse(args)

	if *input == "" || *output == "" || *text == "" {
		fmt.Fprintln(os.Stderr, "Error: -input, -output, and -text are required")
		fs.Usage()
		os.Exit(1)
	}

	doc, err := docx.Open(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening document: %v\n", err)
		os.Exit(1)
	}

	// Build paragraph options
	var opts []docx.ParagraphOption
	if *bold {
		opts = append(opts, docx.WithBold())
	}
	if *italic {
		opts = append(opts, docx.WithItalic())
	}
	if *size != "" {
		opts = append(opts, docx.WithSize(*size))
	}
	if *color != "" {
		opts = append(opts, docx.WithColor(*color))
	}
	if *align != "" {
		opts = append(opts, docx.WithAlignment(*align))
	}

	if *index >= 0 {
		if err := doc.AddParagraphAt(*index, *text, opts...); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding paragraph: %v\n", err)
			os.Exit(1)
		}
	} else {
		doc.AddParagraph(*text, opts...)
	}

	if err := doc.Save(*output); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving document: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Content added successfully to: %s\n", *output)
}

// HandleDelete handles the delete command
func HandleDelete(args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	input := fs.String("input", "", "Input file path (required)")
	output := fs.String("output", "", "Output file path (required)")
	paragraph := fs.Int("paragraph", -1, "Paragraph index to delete")
	start := fs.Int("start", -1, "Start index for range deletion")
	end := fs.Int("end", -1, "End index for range deletion")
	table := fs.Int("table", -1, "Table index to delete")
	fs.Parse(args)

	if *input == "" || *output == "" {
		fmt.Fprintln(os.Stderr, "Error: -input and -output are required")
		fs.Usage()
		os.Exit(1)
	}

	doc, err := docx.Open(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening document: %v\n", err)
		os.Exit(1)
	}

	if *start >= 0 && *end >= 0 {
		if err := doc.DeleteParagraphsRange(*start, *end); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting paragraphs: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted paragraphs %d to %d\n", *start, *end)
	} else if *paragraph >= 0 {
		if err := doc.DeleteParagraph(*paragraph); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting paragraph: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted paragraph %d\n", *paragraph)
	} else if *table >= 0 {
		if err := doc.DeleteTable(*table); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting table: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted table %d\n", *table)
	} else {
		fmt.Fprintln(os.Stderr, "Error: specify -paragraph, -table, or -start/-end")
		fs.Usage()
		os.Exit(1)
	}

	if err := doc.Save(*output); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving document: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Document saved: %s\n", *output)
}

// HandleClear handles the clear command
func HandleClear(args []string) {
	fs := flag.NewFlagSet("clear", flag.ExitOnError)
	input := fs.String("input", "", "Input file path (required)")
	output := fs.String("output", "", "Output file path (required)")
	fs.Parse(args)

	if *input == "" || *output == "" {
		fmt.Fprintln(os.Stderr, "Error: -input and -output are required")
		fs.Usage()
		os.Exit(1)
	}

	doc, err := docx.Open(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening document: %v\n", err)
		os.Exit(1)
	}

	doc.Clear()

	if err := doc.Save(*output); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving document: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Document cleared and saved: %s\n", *output)
}
