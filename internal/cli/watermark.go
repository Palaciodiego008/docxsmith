package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

// HandleWatermark handles watermark operations
func HandleWatermark(args []string) {
	if len(args) == 0 {
		printWatermarkUsage()
		return
	}

	switch args[0] {
	case "add":
		if err := watermarkAddCommand(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "remove":
		if err := watermarkRemoveCommand(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "help":
		printWatermarkUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown watermark subcommand: %s\n", args[0])
		printWatermarkUsage()
		os.Exit(1)
	}
}

// watermarkAddCommand adds a watermark to a document
func watermarkAddCommand(args []string) error {
	fs := flag.NewFlagSet("watermark add", flag.ContinueOnError)

	var (
		inputPath  = fs.String("input", "", "Input .docx file path (required)")
		outputPath = fs.String("output", "", "Output .docx file path (required)")
		text       = fs.String("text", "", "Watermark text (required)")
		opacity    = fs.Float64("opacity", 0.5, "Watermark opacity (0.0-1.0, default: 0.5)")
		angle      = fs.Float64("angle", -45, "Watermark rotation angle in degrees (default: -45)")
		size       = fs.Int("size", 96, "Watermark font size in points (default: 96)")
		color      = fs.String("color", "CCCCCC", "Watermark color in hex (default: CCCCCC)")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Validate required flags
	if *inputPath == "" {
		return fmt.Errorf("input file path is required (-input)")
	}
	if *outputPath == "" {
		return fmt.Errorf("output file path is required (-output)")
	}
	if *text == "" {
		return fmt.Errorf("watermark text is required (-text)")
	}

	// Validate opacity range
	if *opacity < 0 || *opacity > 1 {
		return fmt.Errorf("opacity must be between 0 and 1, got %.2f", *opacity)
	}

	// Validate size
	if *size <= 0 {
		return fmt.Errorf("size must be greater than 0, got %d", *size)
	}

	// Open document
	doc, err := docx.Open(*inputPath)
	if err != nil {
		return fmt.Errorf("failed to open document: %v", err)
	}

	// Build watermark options
	opts := []docx.WatermarkOption{
		docx.WithWatermarkOpacity(*opacity),
		docx.WithWatermarkAngle(*angle),
		docx.WithWatermarkSize(*size),
		docx.WithWatermarkColor(*color),
	}

	// Add watermark
	if err := doc.AddWatermark(*text, opts...); err != nil {
		return fmt.Errorf("failed to add watermark: %v", err)
	}

	// Save document
	if err := doc.Save(*outputPath); err != nil {
		return fmt.Errorf("failed to save document: %v", err)
	}

	fmt.Printf("Watermark added successfully to %s\n", *outputPath)
	return nil
}

// watermarkRemoveCommand removes the watermark from a document
func watermarkRemoveCommand(args []string) error {
	fs := flag.NewFlagSet("watermark remove", flag.ContinueOnError)

	var (
		inputPath  = fs.String("input", "", "Input .docx file path (required)")
		outputPath = fs.String("output", "", "Output .docx file path (required)")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Validate required flags
	if *inputPath == "" {
		return fmt.Errorf("input file path is required (-input)")
	}
	if *outputPath == "" {
		return fmt.Errorf("output file path is required (-output)")
	}

	// Open document
	doc, err := docx.Open(*inputPath)
	if err != nil {
		return fmt.Errorf("failed to open document: %v", err)
	}

	// Remove watermark
	if err := doc.RemoveWatermark(); err != nil {
		return fmt.Errorf("failed to remove watermark: %v", err)
	}

	// Save document
	if err := doc.Save(*outputPath); err != nil {
		return fmt.Errorf("failed to save document: %v", err)
	}

	fmt.Printf("Watermark removed successfully from %s\n", *outputPath)
	return nil
}

// printWatermarkUsage prints the watermark command usage
func printWatermarkUsage() {
	usage := `Watermark - Add or remove watermarks from DOCX documents

Usage:
  docxsmith watermark <subcommand> [options]

Subcommands:
  add       Add a watermark to a document
  remove    Remove the watermark from a document
  help      Show this help message

Add Watermark Options:
  -input string        Input .docx file path (required)
  -output string       Output .docx file path (required)
  -text string         Watermark text (required)
  -opacity float       Watermark opacity between 0 and 1 (default: 0.5)
  -angle float         Watermark rotation angle in degrees (default: -45)
  -size int            Watermark font size in points (default: 96)
  -color string        Watermark color in hex format (default: CCCCCC)

Remove Watermark Options:
  -input string        Input .docx file path (required)
  -output string       Output .docx file path (required)

Examples:
  # Add a watermark with default settings
  docxsmith watermark add -input document.docx -output watermarked.docx -text "DRAFT"

  # Add a watermark with custom settings
  docxsmith watermark add -input document.docx -output watermarked.docx \
    -text "CONFIDENTIAL" -opacity 0.3 -size 120 -color "FF0000"

  # Remove watermark
  docxsmith watermark remove -input watermarked.docx -output clean.docx
`
	fmt.Println(usage)
}
