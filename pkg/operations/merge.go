package operations

import (
	"fmt"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
	"github.com/Palaciodiego008/docxsmith/pkg/pdf"
)

// MergeOptions holds options for merging documents
type MergeOptions struct {
	// AddPageBreaks adds page breaks between documents
	AddPageBreaks bool

	// AddSeparator adds a separator paragraph between documents
	AddSeparator bool

	// SeparatorText is the text to use as separator
	SeparatorText string

	// PreserveFormatting attempts to preserve source formatting
	PreserveFormatting bool
}

// DefaultMergeOptions returns default merge options
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{
		AddPageBreaks:      true,
		AddSeparator:       false,
		SeparatorText:      "---",
		PreserveFormatting: true,
	}
}

// MergeDOCX merges multiple DOCX documents into one
func MergeDOCX(inputPaths []string, outputPath string, opts MergeOptions) error {
	if len(inputPaths) == 0 {
		return fmt.Errorf("no input files provided")
	}

	// Create a new document for the result
	result := docx.New()

	// Process each input document
	for i, path := range inputPaths {
		doc, err := docx.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", path, err)
		}

		// Add separator before document (except first)
		if i > 0 && opts.AddSeparator {
			result.AddParagraph(opts.SeparatorText)
			result.AddParagraph("")
		}

		// Copy all paragraphs
		for _, para := range doc.Body.Paragraphs {
			result.Body.Paragraphs = append(result.Body.Paragraphs, para)
		}

		// Copy all tables
		for _, table := range doc.Body.Tables {
			result.Body.Tables = append(result.Body.Tables, table)
		}

		// Add page break after document (except last)
		if i < len(inputPaths)-1 && opts.AddPageBreaks {
			// Add empty paragraph as page break placeholder
			result.AddParagraph("")
		}
	}

	// Save the merged document
	return result.Save(outputPath)
}

// MergePDF merges multiple PDF documents into one
func MergePDF(inputPaths []string, outputPath string) error {
	if len(inputPaths) == 0 {
		return fmt.Errorf("no input files provided")
	}

	// Create a new PDF document
	result := pdf.New()

	// Process each input PDF
	for _, path := range inputPaths {
		doc, err := pdf.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", path, err)
		}

		// Copy all pages
		for _, page := range doc.Pages {
			newPage := result.AddPage()

			// Copy page properties
			newPage.Width = page.Width
			newPage.Height = page.Height
			newPage.Margin = page.Margin

			// Copy content
			newPage.Content = append(newPage.Content, page.Content...)
		}
	}

	// Save the merged PDF
	return result.Save(outputPath)
}

// MergeDocuments is a convenience function that detects file type and merges accordingly
func MergeDocuments(inputPaths []string, outputPath string, opts MergeOptions) error {
	if len(inputPaths) == 0 {
		return fmt.Errorf("no input files provided")
	}

	// Detect file type from first input
	firstPath := inputPaths[0]
	if len(firstPath) < 4 {
		return fmt.Errorf("invalid file path: %s", firstPath)
	}

	ext := firstPath[len(firstPath)-4:]

	switch ext {
	case "docx":
		return MergeDOCX(inputPaths, outputPath, opts)
	case ".pdf":
		return MergePDF(inputPaths, outputPath)
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}
}

// MergeInfo holds information about a merge operation
type MergeInfo struct {
	TotalDocuments  int
	TotalPages      int
	TotalParagraphs int
	TotalTables     int
}

// GetMergeDOCXInfo returns information about what would be merged
func GetMergeDOCXInfo(inputPaths []string) (*MergeInfo, error) {
	info := &MergeInfo{
		TotalDocuments: len(inputPaths),
	}

	for _, path := range inputPaths {
		doc, err := docx.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s: %w", path, err)
		}

		info.TotalParagraphs += doc.GetParagraphCount()
		info.TotalTables += doc.GetTableCount()
	}

	return info, nil
}

// GetMergePDFInfo returns information about what would be merged
func GetMergePDFInfo(inputPaths []string) (*MergeInfo, error) {
	info := &MergeInfo{
		TotalDocuments: len(inputPaths),
	}

	for _, path := range inputPaths {
		doc, err := pdf.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s: %w", path, err)
		}

		info.TotalPages += doc.GetPageCount()
	}

	return info, nil
}
