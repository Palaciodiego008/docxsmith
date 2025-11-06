package cli

import (
	"fmt"
	"os"
)

const Version = "1.0.0"

// Run is the main entry point for the CLI
func Run(args []string) {
	if len(args) < 1 {
		PrintUsage()
		os.Exit(1)
	}

	command := args[0]

	switch command {
	// DOCX commands
	case "create":
		HandleCreate(args[1:])
	case "add":
		HandleAdd(args[1:])
	case "delete":
		HandleDelete(args[1:])
	case "replace":
		HandleReplace(args[1:])
	case "find":
		HandleFind(args[1:])
	case "extract":
		HandleExtract(args[1:])
	case "table":
		HandleTable(args[1:])
	case "clear":
		HandleClear(args[1:])
	case "info":
		HandleInfo(args[1:])

	// PDF commands
	case "pdf-create":
		HandlePDFCreate(args[1:])
	case "pdf-add":
		HandlePDFAdd(args[1:])
	case "pdf-info":
		HandlePDFInfo(args[1:])
	case "pdf-extract":
		HandlePDFExtract(args[1:])

	// Conversion
	case "convert":
		HandleConvert(args[1:])

	// Template Engine
	case "template-render":
		HandleTemplateRender(args[1:])
	case "template-variables":
		HandleTemplateVariables(args[1:])
	case "template-example":
		HandleTemplateExample(args[1:])

	// Merge & Split
	case "merge":
		HandleMerge(args[1:])
	case "split":
		HandleSplit(args[1:])
	case "merge-info":
		HandleMergeInfo(args[1:])

	// Document Diff
	case "diff":
		HandleDiff(args[1:])

	// Utility
	case "version":
		fmt.Printf("DocxSmith v%s\n", Version)
	case "help":
		PrintUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		PrintUsage()
		os.Exit(1)
	}
}

// PrintUsage prints the usage information
func PrintUsage() {
	usage := `DocxSmith - The Document Forge
A powerful tool for manipulating .docx and .pdf files

Usage:
  docxsmith <command> [options]

DOCX Commands:
  create      Create a new DOCX document
  add         Add content to a DOCX document
  delete      Delete content from a DOCX document
  replace     Replace text in a DOCX document
  find        Find text in a DOCX document
  extract     Extract text from a DOCX document
  table       Manipulate tables in a DOCX document
  clear       Clear all content from a DOCX document
  info        Display DOCX document information

PDF Commands:
  pdf-create  Create a new PDF document
  pdf-add     Add content to a PDF document
  pdf-info    Display PDF document information
  pdf-extract Extract text from a PDF document

Conversion:
  convert     Convert between DOCX and PDF formats

Template Engine:
  template-render     Render a template with data (JSON/YAML)
  template-variables  List variables in a template
  template-example    Create example template and data files

Merge & Split:
  merge        Merge multiple documents into one
  split        Split a document into multiple files
  merge-info   Show information about merge operation

Comparison:
  diff         Compare two documents and show differences

Utility:
  version     Show version information
  help        Show this help message

Examples:
  # DOCX operations
  docxsmith create -output sample.docx -text "Hello World"
  docxsmith add -input doc.docx -output new.docx -text "New paragraph" -bold
  docxsmith table -input doc.docx -output new.docx -create -rows 3 -cols 4

  # PDF operations
  docxsmith pdf-create -output sample.pdf -text "Hello PDF"
  docxsmith pdf-add -input doc.pdf -output new.pdf -text "New text" -bold
  docxsmith pdf-info -input document.pdf

  # Conversion
  docxsmith convert -input document.docx -output document.pdf
  docxsmith convert -input document.pdf -output document.docx

  # Template Engine
  docxsmith template-example -template invoice.docx -data data.json
  docxsmith template-render -template invoice.docx -data data.json -output result.docx
  docxsmith template-variables -template invoice.docx

  # Merge & Split
  docxsmith merge -inputs doc1.docx,doc2.docx,doc3.docx -output combined.docx
  docxsmith split -input large.pdf -count 3 -pattern "chapter{n}.pdf"
  docxsmith split -input book.docx -by-heading -heading-level 1

  # Document Comparison
  docxsmith diff -old v1.docx -new v2.docx -output changes.html
  docxsmith diff -old v1.docx -new v2.docx -format markdown -ignore-whitespace

For more information on a command:
  docxsmith <command> -help
`
	fmt.Println(usage)
}
