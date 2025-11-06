package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create":
		handleCreate(os.Args[2:])
	case "add":
		handleAdd(os.Args[2:])
	case "delete":
		handleDelete(os.Args[2:])
	case "replace":
		handleReplace(os.Args[2:])
	case "find":
		handleFind(os.Args[2:])
	case "extract":
		handleExtract(os.Args[2:])
	case "table":
		handleTable(os.Args[2:])
	case "clear":
		handleClear(os.Args[2:])
	case "info":
		handleInfo(os.Args[2:])
	case "version":
		fmt.Printf("DocxSmith v%s\n", version)
	case "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	usage := `DocxSmith - The Document Forge
A powerful tool for manipulating .docx files

Usage:
  docxsmith <command> [options]

Commands:
  create      Create a new document
  add         Add content to a document
  delete      Delete content from a document
  replace     Replace text in a document
  find        Find text in a document
  extract     Extract text from a document
  table       Manipulate tables in a document
  clear       Clear all content from a document
  info        Display document information
  version     Show version information
  help        Show this help message

Examples:
  docxsmith create -output sample.docx -text "Hello World"
  docxsmith add -input doc.docx -output new.docx -text "New paragraph"
  docxsmith replace -input doc.docx -output new.docx -old "foo" -new "bar"
  docxsmith find -input doc.docx -text "search term"
  docxsmith table -input doc.docx -output new.docx -create -rows 3 -cols 4

For more information on a command:
  docxsmith <command> -help
`
	fmt.Println(usage)
}

func handleCreate(args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	output := fs.String("output", "", "Output file path (required)")
	text := fs.String("text", "", "Initial text content")
	fs.Parse(args)

	if *output == "" {
		fmt.Fprintln(os.Stderr, "Error: -output is required")
		fs.Usage()
		os.Exit(1)
	}

	doc := docx.New()
	if *text != "" {
		doc.AddParagraph(*text)
	} else {
		doc.AddParagraph("Document created with DocxSmith")
	}

	if err := doc.Save(*output); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating document: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Document created successfully: %s\n", *output)
}

func handleAdd(args []string) {
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

func handleDelete(args []string) {
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

func handleReplace(args []string) {
	fs := flag.NewFlagSet("replace", flag.ExitOnError)
	input := fs.String("input", "", "Input file path (required)")
	output := fs.String("output", "", "Output file path (required)")
	oldText := fs.String("old", "", "Text to replace (required)")
	newText := fs.String("new", "", "Replacement text (required)")
	paragraph := fs.Int("paragraph", -1, "Only replace in specific paragraph")
	fs.Parse(args)

	if *input == "" || *output == "" || *oldText == "" || *newText == "" {
		fmt.Fprintln(os.Stderr, "Error: -input, -output, -old, and -new are required")
		fs.Usage()
		os.Exit(1)
	}

	doc, err := docx.Open(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening document: %v\n", err)
		os.Exit(1)
	}

	var count int
	if *paragraph >= 0 {
		if err := doc.ReplaceTextInParagraph(*paragraph, *oldText, *newText); err != nil {
			fmt.Fprintf(os.Stderr, "Error replacing text: %v\n", err)
			os.Exit(1)
		}
		count = 1
	} else {
		count = doc.ReplaceText(*oldText, *newText)
	}

	if err := doc.Save(*output); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving document: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Replaced %d occurrence(s) of '%s' with '%s'\n", count, *oldText, *newText)
	fmt.Printf("Document saved: %s\n", *output)
}

func handleFind(args []string) {
	fs := flag.NewFlagSet("find", flag.ExitOnError)
	input := fs.String("input", "", "Input file path (required)")
	text := fs.String("text", "", "Text to find (required)")
	fs.Parse(args)

	if *input == "" || *text == "" {
		fmt.Fprintln(os.Stderr, "Error: -input and -text are required")
		fs.Usage()
		os.Exit(1)
	}

	doc, err := docx.Open(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening document: %v\n", err)
		os.Exit(1)
	}

	indices := doc.FindText(*text)
	if len(indices) == 0 {
		fmt.Printf("Text '%s' not found in document\n", *text)
		return
	}

	fmt.Printf("Found '%s' in %d paragraph(s):\n", *text, len(indices))
	for _, idx := range indices {
		text, _ := doc.GetParagraphText(idx)
		preview := text
		if len(preview) > 80 {
			preview = preview[:77] + "..."
		}
		fmt.Printf("  Paragraph %d: %s\n", idx, preview)
	}
}

func handleExtract(args []string) {
	fs := flag.NewFlagSet("extract", flag.ExitOnError)
	input := fs.String("input", "", "Input file path (required)")
	output := fs.String("output", "", "Output text file (optional)")
	fs.Parse(args)

	if *input == "" {
		fmt.Fprintln(os.Stderr, "Error: -input is required")
		fs.Usage()
		os.Exit(1)
	}

	doc, err := docx.Open(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening document: %v\n", err)
		os.Exit(1)
	}

	text := doc.GetText()

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

func handleTable(args []string) {
	fs := flag.NewFlagSet("table", flag.ExitOnError)
	input := fs.String("input", "", "Input file path (required)")
	output := fs.String("output", "", "Output file path (required)")
	create := fs.Bool("create", false, "Create a new table")
	rows := fs.Int("rows", 2, "Number of rows")
	cols := fs.Int("cols", 2, "Number of columns")
	setCellText := fs.String("set", "", "Set cell text (format: 'tableIdx,row,col,text')")
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

	if *create {
		table := doc.AddTable(*rows, *cols)
		fmt.Printf("Created table with %d rows and %d columns\n", *rows, *cols)

		// Set header row as example
		if *rows > 0 && *cols > 0 {
			for i := 0; i < *cols; i++ {
				table.SetCellText(0, i, fmt.Sprintf("Header %d", i+1))
			}
		}
	}

	if *setCellText != "" {
		parts := strings.Split(*setCellText, ",")
		if len(parts) != 4 {
			fmt.Fprintln(os.Stderr, "Error: -set format must be 'tableIdx,row,col,text'")
			os.Exit(1)
		}

		var tableIdx, row, col int
		fmt.Sscanf(parts[0], "%d", &tableIdx)
		fmt.Sscanf(parts[1], "%d", &row)
		fmt.Sscanf(parts[2], "%d", &col)
		text := parts[3]

		if tableIdx < 0 || tableIdx >= len(doc.Body.Tables) {
			fmt.Fprintf(os.Stderr, "Error: table index %d out of range\n", tableIdx)
			os.Exit(1)
		}

		table := &doc.Body.Tables[tableIdx]
		if err := table.SetCellText(row, col, text); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting cell text: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Set cell [%d,%d] in table %d to: %s\n", row, col, tableIdx, text)
	}

	if err := doc.Save(*output); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving document: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Document saved: %s\n", *output)
}

func handleClear(args []string) {
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

func handleInfo(args []string) {
	fs := flag.NewFlagSet("info", flag.ExitOnError)
	input := fs.String("input", "", "Input file path (required)")
	fs.Parse(args)

	if *input == "" {
		fmt.Fprintln(os.Stderr, "Error: -input is required")
		fs.Usage()
		os.Exit(1)
	}

	doc, err := docx.Open(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening document: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Document Information: %s\n", *input)
	fmt.Printf("  Paragraphs: %d\n", doc.GetParagraphCount())
	fmt.Printf("  Tables: %d\n", doc.GetTableCount())

	wordCount := len(strings.Fields(doc.GetText()))
	charCount := len(doc.GetText())
	fmt.Printf("  Words: %d\n", wordCount)
	fmt.Printf("  Characters: %d\n", charCount)

	if doc.GetTableCount() > 0 {
		fmt.Println("\nTable Details:")
		for i, table := range doc.Body.Tables {
			fmt.Printf("  Table %d: %d rows × %d columns\n",
				i, table.GetRowCount(), table.GetColumnCount())
		}
	}
}
