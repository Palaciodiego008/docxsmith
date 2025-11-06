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
