package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Palaciodiego008/docxsmith/pkg/diff"
)

// HandleDiff handles the diff command using improved architecture
func HandleDiff(args []string) {
	fs := flag.NewFlagSet("diff", flag.ExitOnError)

	// Use common flag helpers
	oldFile := fs.String("old", "", "Old/original document (required)")
	newFile := fs.String("new", "", "New/modified document (required)")
	output := fs.String("output", "", "Output file (default: stdout)")
	format := fs.String("format", "html", "Output format: html, markdown, text")
	ignoreWhitespace := fs.Bool("ignore-whitespace", false, "Ignore whitespace differences")
	ignoreCase := fs.Bool("ignore-case", false, "Ignore case differences")
	showStats := fs.Bool("stats", true, "Show statistics")

	if err := fs.Parse(args); err != nil {
		ExitWithError("Failed to parse flags: %v", err)
	}

	// Validate required parameters using common utility
	if err := ValidateRequired(map[string]string{
		"old": *oldFile,
		"new": *newFile,
	}); err != nil {
		ExitWithError("%v", err)
	}

	// Validate files exist
	if err := ValidateFileExists(*oldFile); err != nil {
		ExitWithError("%v", err)
	}
	if err := ValidateFileExists(*newFile); err != nil {
		ExitWithError("%v", err)
	}

	// Configure diff options
	opts := diff.DiffOptions{
		IgnoreWhitespace: *ignoreWhitespace,
		IgnoreCase:       *ignoreCase,
		ContextLines:     3,
		MinChangeLength:  1,
	}

	// Compare documents
	PrintInfo("Comparing documents...")
	result, err := diff.CompareDOCX(*oldFile, *newFile, opts)
	if err != nil {
		ExitWithError("Failed to compare documents: %v", err)
	}

	// Choose renderer based on format
	var renderer diff.Renderer
	switch *format {
	case "html":
		renderer = diff.NewHTMLRenderer(*showStats)
	case "markdown", "md":
		renderer = diff.NewMarkdownRenderer(*showStats)
	case "text", "txt":
		renderer = diff.NewPlainTextRenderer(*showStats, true)
	default:
		ExitWithError("Unknown format: %s (use: html, markdown, text)", *format)
	}

	// Render diff
	outputContent, err := renderer.Render(result)
	if err != nil {
		ExitWithError("Failed to render diff: %v", err)
	}

	// Output result
	if *output != "" {
		if err := os.WriteFile(*output, []byte(outputContent), 0644); err != nil {
			ExitWithError("Failed to write output file: %v", err)
		}
		PrintSuccess("Diff saved to: %s", *output)
	} else {
		fmt.Println(outputContent)
	}

	// Print summary
	if result.Stats.TotalChanges == 0 {
		PrintSuccess("Documents are identical - no changes detected")
	} else {
		PrintInfo("\nSummary:")
		PrintInfo("  Total changes: %d", result.Stats.TotalChanges)
		PrintInfo("  Added lines:   %d", result.Stats.AddedLines)
		PrintInfo("  Deleted lines: %d", result.Stats.DeletedLines)
	}
}
