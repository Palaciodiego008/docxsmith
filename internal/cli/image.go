package cli

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

// ImageCommand handles image-related operations
func ImageCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("image command requires subcommand: add, insert, count")
	}

	switch args[0] {
	case "add":
		return imageAddCommand(args[1:])
	case "insert":
		return imageInsertCommand(args[1:])
	case "count":
		return imageCountCommand(args[1:])
	default:
		return fmt.Errorf("unknown image subcommand: %s", args[0])
	}
}

// imageAddCommand adds an image to the document
func imageAddCommand(args []string) error {
	fs := flag.NewFlagSet("image add", flag.ExitOnError)

	var (
		inputPath  = fs.String("input", "", "Input .docx file path (required)")
		outputPath = fs.String("output", "", "Output .docx file path (required)")
		imagePath  = fs.String("image", "", "Image file path (required)")
		width      = fs.Int("width", 200, "Image width in pixels (default: 200)")
		height     = fs.Int("height", 150, "Image height in pixels (default: 150)")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Validate required flags
	if *inputPath == "" {
		return fmt.Errorf("input file path is required")
	}
	if *outputPath == "" {
		return fmt.Errorf("output file path is required")
	}
	if *imagePath == "" {
		return fmt.Errorf("image file path is required")
	}

	// Open document
	doc, err := docx.Open(*inputPath)
	if err != nil {
		return fmt.Errorf("failed to open document: %v", err)
	}

	// Add image with options
	var opts []docx.ImageOption
	if *width != 200 {
		opts = append(opts, docx.WithImageWidth(*width))
	}
	if *height != 150 {
		opts = append(opts, docx.WithImageHeight(*height))
	}

	err = doc.AddImage(*imagePath, opts...)
	if err != nil {
		return fmt.Errorf("failed to add image: %v", err)
	}

	// Save document
	err = doc.Save(*outputPath)
	if err != nil {
		return fmt.Errorf("failed to save document: %v", err)
	}

	fmt.Printf("Image added successfully. Document saved as %s\n", *outputPath)
	return nil
}

// imageInsertCommand inserts an image at a specific position
func imageInsertCommand(args []string) error {
	fs := flag.NewFlagSet("image insert", flag.ExitOnError)

	var (
		inputPath  = fs.String("input", "", "Input .docx file path (required)")
		outputPath = fs.String("output", "", "Output .docx file path (required)")
		imagePath  = fs.String("image", "", "Image file path (required)")
		position   = fs.String("at", "", "Position to insert image (paragraph index, required)")
		width      = fs.Int("width", 200, "Image width in pixels (default: 200)")
		height     = fs.Int("height", 150, "Image height in pixels (default: 150)")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Validate required flags
	if *inputPath == "" {
		return fmt.Errorf("input file path is required")
	}
	if *outputPath == "" {
		return fmt.Errorf("output file path is required")
	}
	if *imagePath == "" {
		return fmt.Errorf("image file path is required")
	}
	if *position == "" {
		return fmt.Errorf("position is required")
	}

	// Parse position
	pos, err := strconv.Atoi(*position)
	if err != nil {
		return fmt.Errorf("invalid position: %v", err)
	}

	// Open document
	doc, err := docx.Open(*inputPath)
	if err != nil {
		return fmt.Errorf("failed to open document: %v", err)
	}

	// Add image with options
	var opts []docx.ImageOption
	if *width != 200 {
		opts = append(opts, docx.WithImageWidth(*width))
	}
	if *height != 150 {
		opts = append(opts, docx.WithImageHeight(*height))
	}

	err = doc.AddImageAt(pos, *imagePath, opts...)
	if err != nil {
		return fmt.Errorf("failed to insert image: %v", err)
	}

	// Save document
	err = doc.Save(*outputPath)
	if err != nil {
		return fmt.Errorf("failed to save document: %v", err)
	}

	fmt.Printf("Image inserted at position %d. Document saved as %s\n", pos, *outputPath)
	return nil
}

// imageCountCommand counts images in the document
func imageCountCommand(args []string) error {
	fs := flag.NewFlagSet("image count", flag.ExitOnError)

	var (
		inputPath = fs.String("input", "", "Input .docx file path (required)")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Validate required flags
	if *inputPath == "" {
		return fmt.Errorf("input file path is required")
	}

	// Open document
	doc, err := docx.Open(*inputPath)
	if err != nil {
		return fmt.Errorf("failed to open document: %v", err)
	}

	// Get image count
	count := doc.GetImageCount()
	fmt.Printf("Document contains %d image(s)\n", count)

	return nil
}
