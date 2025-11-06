package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

// HandleTable handles the table command
func HandleTable(args []string) {
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
		if n, err := fmt.Sscanf(parts[0], "%d", &tableIdx); err != nil || n != 1 {
			fmt.Fprintf(os.Stderr, "Error: invalid tableIdx value '%s'\n", parts[0])
			os.Exit(1)
		}
		if n, err := fmt.Sscanf(parts[1], "%d", &row); err != nil || n != 1 {
			fmt.Fprintf(os.Stderr, "Error: invalid row value '%s'\n", parts[1])
			os.Exit(1)
		}
		if n, err := fmt.Sscanf(parts[2], "%d", &col); err != nil || n != 1 {
			fmt.Fprintf(os.Stderr, "Error: invalid col value '%s'\n", parts[2])
			os.Exit(1)
		}
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
