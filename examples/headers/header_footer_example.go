//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

func main() {
	// Create a new document
	doc := docx.New()

	// Add document title
	doc.AddParagraph("Professional Document with Headers and Footers", 
		docx.WithBold(), 
		docx.WithSize("32"), 
		docx.WithAlignment("center"))
	doc.AddParagraph("") // Empty line

	// Set default header
	err := doc.SetHeader(docx.HeaderTypeDefault, "Company Name - Confidential", 
		docx.WithHFBold(), 
		docx.WithHFAlignment("center"),
		docx.WithHFTextColor("0066CC"))
	if err != nil {
		log.Fatalf("Failed to set default header: %v", err)
	}

	// Set first page header (different from default)
	err = doc.SetHeader(docx.HeaderTypeFirst, "DRAFT - Internal Use Only", 
		docx.WithHFItalic(), 
		docx.WithHFAlignment("right"),
		docx.WithHFTextColor("FF0000"))
	if err != nil {
		log.Fatalf("Failed to set first page header: %v", err)
	}

	// Set default footer with page numbering
	err = doc.SetFooter(docx.FooterTypeDefault, "Page {PAGE} of {NUMPAGES}", 
		docx.WithHFAlignment("center"),
		docx.WithHFFontSize("20"))
	if err != nil {
		log.Fatalf("Failed to set default footer: %v", err)
	}

	// Set first page footer
	err = doc.SetFooter(docx.FooterTypeFirst, "© 2024 Company Name. All rights reserved.", 
		docx.WithHFAlignment("center"),
		docx.WithHFFontSize("18"),
		docx.WithHFItalic())
	if err != nil {
		log.Fatalf("Failed to set first page footer: %v", err)
	}

	// Add some content
	doc.AddParagraph("Introduction", docx.WithBold(), docx.WithSize("28"))
	doc.AddParagraph("This document demonstrates the use of headers and footers in DocxSmith.")
	doc.AddParagraph("")

	doc.AddParagraph("Features Demonstrated:", docx.WithBold())
	doc.AddParagraph("• Default headers and footers for all pages")
	doc.AddParagraph("• Special first page headers and footers")
	doc.AddParagraph("• Text formatting in headers and footers")
	doc.AddParagraph("• Different alignments and colors")
	doc.AddParagraph("")

	// Add more content to show multiple pages
	for i := 1; i <= 5; i++ {
		doc.AddParagraph(fmt.Sprintf("Section %d", i), docx.WithBold(), docx.WithSize("24"))
		doc.AddParagraph(fmt.Sprintf("This is the content for section %d. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.", i))
		doc.AddParagraph("")
	}

	// Demonstrate header/footer management
	fmt.Println("Document Header/Footer Status:")
	fmt.Printf("Has default header: %v\n", doc.HasHeader(docx.HeaderTypeDefault))
	fmt.Printf("Has first page header: %v\n", doc.HasHeader(docx.HeaderTypeFirst))
	fmt.Printf("Has even page header: %v\n", doc.HasHeader(docx.HeaderTypeEven))
	fmt.Printf("Has default footer: %v\n", doc.HasFooter(docx.FooterTypeDefault))
	fmt.Printf("Has first page footer: %v\n", doc.HasFooter(docx.FooterTypeFirst))

	// Retrieve and display header content
	if header, err := doc.GetHeader(docx.HeaderTypeDefault); err == nil {
		fmt.Printf("Default header content: %s\n", getHeaderFooterText(header))
	}

	// Save the document
	outputPath := "professional_document.docx"
	err = doc.Save(outputPath)
	if err != nil {
		log.Fatalf("Failed to save document: %v", err)
	}

	fmt.Printf("Document saved successfully as %s\n", outputPath)

	// Demonstrate removing headers/footers
	fmt.Println("\nRemoving first page header...")
	err = doc.RemoveHeader(docx.HeaderTypeFirst)
	if err != nil {
		log.Printf("Failed to remove first page header: %v", err)
	} else {
		fmt.Printf("First page header removed. Has first page header: %v\n", 
			doc.HasHeader(docx.HeaderTypeFirst))
	}

	// Save modified document
	modifiedPath := "professional_document_modified.docx"
	err = doc.Save(modifiedPath)
	if err != nil {
		log.Fatalf("Failed to save modified document: %v", err)
	}

	fmt.Printf("Modified document saved as %s\n", modifiedPath)
}

// Helper function to extract text from header/footer
func getHeaderFooterText(hf *docx.HeaderFooter) string {
	if len(hf.Paragraphs) == 0 {
		return ""
	}
	
	var text string
	for _, run := range hf.Paragraphs[0].Runs {
		for _, t := range run.Text {
			text += t.Content
		}
	}
	return text
}
