package main

import (
	"fmt"
	"log"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

func main() {
	// Example 1: Create a new document
	fmt.Println("=== Example 1: Creating a new document ===")
	doc := docx.New()
	doc.AddParagraph("Welcome to DocxSmith!")
	doc.AddParagraph("This is a simple document manipulation library.", docx.WithBold())
	doc.AddParagraph("You can format text easily.", docx.WithItalic(), docx.WithColor("0000FF"))

	if err := doc.Save("example_new.docx"); err != nil {
		log.Fatalf("Error saving document: %v", err)
	}
	fmt.Println("Created: example_new.docx")

	// Example 2: Open and modify an existing document
	fmt.Println("\n=== Example 2: Modifying an existing document ===")
	doc2, err := docx.Open("example_new.docx")
	if err != nil {
		log.Fatalf("Error opening document: %v", err)
	}

	// Add more content
	doc2.AddParagraph("")
	doc2.AddParagraph("Additional Content", docx.WithBold(), docx.WithSize("28"))
	doc2.AddParagraph("This paragraph was added after creation.")

	if err := doc2.Save("example_modified.docx"); err != nil {
		log.Fatalf("Error saving modified document: %v", err)
	}
	fmt.Println("Created: example_modified.docx")

	// Example 3: Search and replace
	fmt.Println("\n=== Example 3: Search and Replace ===")
	doc3, err := docx.Open("example_modified.docx")
	if err != nil {
		log.Fatalf("Error opening document: %v", err)
	}

	count := doc3.ReplaceText("DocxSmith", "DocxSmith Pro")
	fmt.Printf("Replaced %d occurrence(s)\n", count)

	if err := doc3.Save("example_replaced.docx"); err != nil {
		log.Fatalf("Error saving document: %v", err)
	}
	fmt.Println("Created: example_replaced.docx")

	// Example 4: Find text
	fmt.Println("\n=== Example 4: Finding Text ===")
	indices := doc3.FindText("paragraph")
	fmt.Printf("Found 'paragraph' in %d paragraph(s): %v\n", len(indices), indices)

	// Example 5: Working with tables
	fmt.Println("\n=== Example 5: Creating Tables ===")
	doc4 := docx.New()
	doc4.AddParagraph("Employee List", docx.WithBold(), docx.WithSize("32"))
	doc4.AddParagraph("")

	table := doc4.AddTable(4, 3)

	// Header row
	table.SetCellText(0, 0, "Name")
	table.SetCellText(0, 1, "Position")
	table.SetCellText(0, 2, "Department")

	// Data rows
	table.SetCellText(1, 0, "John Doe")
	table.SetCellText(1, 1, "Developer")
	table.SetCellText(1, 2, "Engineering")

	table.SetCellText(2, 0, "Jane Smith")
	table.SetCellText(2, 1, "Designer")
	table.SetCellText(2, 2, "UX/UI")

	table.SetCellText(3, 0, "Bob Johnson")
	table.SetCellText(3, 1, "Manager")
	table.SetCellText(3, 2, "Operations")

	if err := doc4.Save("example_table.docx"); err != nil {
		log.Fatalf("Error saving document: %v", err)
	}
	fmt.Println("Created: example_table.docx")

	// Example 6: Document info
	fmt.Println("\n=== Example 6: Document Information ===")
	doc5, err := docx.Open("example_table.docx")
	if err != nil {
		log.Fatalf("Error opening document: %v", err)
	}

	fmt.Printf("Paragraphs: %d\n", doc5.GetParagraphCount())
	fmt.Printf("Tables: %d\n", doc5.GetTableCount())
	fmt.Printf("Text content: %s\n", doc5.GetText()[:100]+"...")

	// Example 7: Delete content
	fmt.Println("\n=== Example 7: Deleting Content ===")
	doc6, err := docx.Open("example_replaced.docx")
	if err != nil {
		log.Fatalf("Error opening document: %v", err)
	}

	fmt.Printf("Paragraphs before deletion: %d\n", doc6.GetParagraphCount())
	if err := doc6.DeleteParagraph(0); err != nil {
		log.Fatalf("Error deleting paragraph: %v", err)
	}
	fmt.Printf("Paragraphs after deletion: %d\n", doc6.GetParagraphCount())

	if err := doc6.Save("example_deleted.docx"); err != nil {
		log.Fatalf("Error saving document: %v", err)
	}
	fmt.Println("Created: example_deleted.docx")

	fmt.Println("\n=== All examples completed successfully! ===")
}
