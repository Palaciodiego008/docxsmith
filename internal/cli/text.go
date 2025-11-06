package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

// HandleReplace handles the replace command
func HandleReplace(args []string) {
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
		count, err := doc.ReplaceTextInParagraph(*paragraph, *oldText, *newText)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error replacing text: %v\n", err)
			os.Exit(1)
		}
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

// HandleFind handles the find command
func HandleFind(args []string) {
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

// HandleExtract handles the extract command
func HandleExtract(args []string) {
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
