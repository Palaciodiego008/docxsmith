//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

func main() {
	// Create a sample image for testing
	sampleImagePath := createSampleImage()
	defer os.Remove(sampleImagePath)

	// Create a new document
	doc := docx.New()

	// Add title
	doc.AddParagraph("Image Examples", docx.WithBold(), docx.WithSize("32"), docx.WithAlignment("center"))
	doc.AddParagraph("") // Empty line

	// Add image with default size
	fmt.Println("Adding image with default size...")
	err := doc.AddImage(sampleImagePath)
	if err != nil {
		log.Fatalf("Failed to add image: %v", err)
	}

	// Add some text
	doc.AddParagraph("Above is an image with default dimensions (200x150 pixels)")
	doc.AddParagraph("") // Empty line

	// Add image with custom size
	fmt.Println("Adding image with custom size...")
	err = doc.AddImage(sampleImagePath, docx.WithImageWidth(300), docx.WithImageHeight(200))
	if err != nil {
		log.Fatalf("Failed to add custom sized image: %v", err)
	}

	doc.AddParagraph("Above is the same image with custom dimensions (300x200 pixels)")
	doc.AddParagraph("") // Empty line

	// Add text and then insert image at specific position
	doc.AddParagraph("This paragraph will be followed by an inserted image.")
	doc.AddParagraph("This paragraph will come after the inserted image.")

	// Insert image between the two paragraphs
	fmt.Println("Inserting image at specific position...")
	err = doc.AddImageAt(doc.GetParagraphCount()-1, sampleImagePath, docx.WithImageWidth(150), docx.WithImageHeight(100))
	if err != nil {
		log.Fatalf("Failed to insert image at position: %v", err)
	}

	// Add summary
	doc.AddParagraph("") // Empty line
	doc.AddParagraph("Document Statistics:", docx.WithBold())
	doc.AddParagraph(fmt.Sprintf("- Total paragraphs: %d", doc.GetParagraphCount()))
	doc.AddParagraph(fmt.Sprintf("- Total images: %d", doc.GetImageCount()))

	// Save the document
	outputPath := "image_examples.docx"
	err = doc.Save(outputPath)
	if err != nil {
		log.Fatalf("Failed to save document: %v", err)
	}

	fmt.Printf("Document saved successfully as %s\n", outputPath)
	fmt.Printf("Total images in document: %d\n", doc.GetImageCount())

	// Demonstrate error handling
	fmt.Println("\nDemonstrating error handling...")

	// Try to add non-existent image
	err = doc.AddImage("nonexistent.jpg")
	if err != nil {
		fmt.Printf("Expected error for non-existent file: %v\n", err)
	}

	// Try to add unsupported format
	textFile := "test.txt"
	os.WriteFile(textFile, []byte("not an image"), 0644)
	defer os.Remove(textFile)

	err = doc.AddImage(textFile)
	if err != nil {
		fmt.Printf("Expected error for unsupported format: %v\n", err)
	}
}

// createSampleImage creates a minimal PNG image for testing
func createSampleImage() string {
	// Create a minimal 1x1 PNG image
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1 dimensions
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, // bit depth, color type, etc.
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, // IDAT chunk
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0x0F, 0x00, 0x00,
		0x01, 0x00, 0x01, 0x5C, 0xC2, 0x8A, 0x8E, 0x00, // IDAT data
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, // IEND chunk
		0x42, 0x60, 0x82,
	}

	imagePath := filepath.Join(os.TempDir(), "sample_image.png")
	err := os.WriteFile(imagePath, pngData, 0644)
	if err != nil {
		log.Fatalf("Failed to create sample image: %v", err)
	}

	return imagePath
}
