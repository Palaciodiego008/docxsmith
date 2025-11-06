package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

// HandleCreate handles the create command
func HandleCreate(args []string) {
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
