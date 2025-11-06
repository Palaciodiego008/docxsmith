package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

// HandleInfo handles the info command
func HandleInfo(args []string) {
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
