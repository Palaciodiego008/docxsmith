package operations

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
	"github.com/Palaciodiego008/docxsmith/pkg/pdf"
)

func TestMergeDOCX(t *testing.T) {
	tests := []struct {
		name          string
		numDocs       int
		parasPerDoc   int
		expectedParas int
		addPageBreaks bool
		addSeparator  bool
		separatorText string
	}{
		{
			name:          "Merge 2 documents with page breaks",
			numDocs:       2,
			parasPerDoc:   3,
			expectedParas: 7, // 3 + 3 + 1 page break
			addPageBreaks: true,
			addSeparator:  false,
		},
		{
			name:          "Merge 3 documents without page breaks",
			numDocs:       3,
			parasPerDoc:   2,
			expectedParas: 6, // 2 + 2 + 2
			addPageBreaks: false,
			addSeparator:  false,
		},
		{
			name:          "Merge with separator",
			numDocs:       2,
			parasPerDoc:   2,
			expectedParas: 6, // 2 + 1 separator + 1 empty + 2
			addPageBreaks: false,
			addSeparator:  true,
			separatorText: "===",
		},
		{
			name:          "Merge single document",
			numDocs:       1,
			parasPerDoc:   5,
			expectedParas: 5,
			addPageBreaks: false,
			addSeparator:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			inputFiles := []string{}

			// Create test documents
			for i := 0; i < tt.numDocs; i++ {
				doc := docx.New()
				for j := 0; j < tt.parasPerDoc; j++ {
					doc.AddParagraph(fmt.Sprintf("Doc%d Para%d", i+1, j+1))
				}

				path := filepath.Join(tmpDir, fmt.Sprintf("doc%d.docx", i+1))
				if err := doc.Save(path); err != nil {
					t.Fatalf("Failed to save test document: %v", err)
				}
				inputFiles = append(inputFiles, path)
			}

			// Merge documents
			outputPath := filepath.Join(tmpDir, "merged.docx")
			opts := MergeOptions{
				AddPageBreaks:      tt.addPageBreaks,
				AddSeparator:       tt.addSeparator,
				SeparatorText:      tt.separatorText,
				PreserveFormatting: true,
			}

			err := MergeDOCX(inputFiles, outputPath, opts)
			if err != nil {
				t.Fatalf("Merge failed: %v", err)
			}

			// Verify merged document
			merged, err := docx.Open(outputPath)
			if err != nil {
				t.Fatalf("Failed to open merged document: %v", err)
			}

			if merged.GetParagraphCount() != tt.expectedParas {
				t.Errorf("Expected %d paragraphs, got %d", tt.expectedParas, merged.GetParagraphCount())
			}
		})
	}
}

func TestMergePDF(t *testing.T) {
	tests := []struct {
		name          string
		numDocs       int
		pagesPerDoc   int
		expectedPages int
	}{
		{
			name:          "Merge 2 PDFs",
			numDocs:       2,
			pagesPerDoc:   3,
			expectedPages: 6,
		},
		{
			name:          "Merge 3 PDFs",
			numDocs:       3,
			pagesPerDoc:   2,
			expectedPages: 6,
		},
		{
			name:          "Merge single PDF",
			numDocs:       1,
			pagesPerDoc:   5,
			expectedPages: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			inputFiles := []string{}

			// Create test PDFs
			for i := 0; i < tt.numDocs; i++ {
				doc := pdf.New()
				for j := 0; j < tt.pagesPerDoc; j++ {
					page := doc.AddPage()
					page.AddText(fmt.Sprintf("PDF%d Page%d", i+1, j+1), 20, 30, 12)
				}

				path := filepath.Join(tmpDir, fmt.Sprintf("pdf%d.pdf", i+1))
				if err := doc.Save(path); err != nil {
					t.Fatalf("Failed to save test PDF: %v", err)
				}
				inputFiles = append(inputFiles, path)
			}

			// Merge PDFs
			outputPath := filepath.Join(tmpDir, "merged.pdf")
			err := MergePDF(inputFiles, outputPath)
			if err != nil {
				t.Fatalf("Merge failed: %v", err)
			}

			// Verify merged PDF
			merged, err := pdf.Open(outputPath)
			if err != nil {
				t.Fatalf("Failed to open merged PDF: %v", err)
			}

			if merged.GetPageCount() != tt.expectedPages {
				t.Errorf("Expected %d pages, got %d", tt.expectedPages, merged.GetPageCount())
			}
		})
	}
}

func TestMergeDOCXErrors(t *testing.T) {
	tests := []struct {
		name        string
		inputPaths  []string
		expectError bool
	}{
		{
			name:        "No input files",
			inputPaths:  []string{},
			expectError: true,
		},
		{
			name:        "Non-existent file",
			inputPaths:  []string{"nonexistent.docx"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputPath := filepath.Join(tmpDir, "output.docx")

			err := MergeDOCX(tt.inputPaths, outputPath, DefaultMergeOptions())

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGetMergeDOCXInfo(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test documents
	doc1 := docx.New()
	doc1.AddParagraph("Para 1")
	doc1.AddParagraph("Para 2")
	doc1.AddTable(2, 2)

	doc2 := docx.New()
	doc2.AddParagraph("Para 3")
	doc2.AddTable(3, 3)

	path1 := filepath.Join(tmpDir, "doc1.docx")
	path2 := filepath.Join(tmpDir, "doc2.docx")

	doc1.Save(path1)
	doc2.Save(path2)

	// Get merge info
	info, err := GetMergeDOCXInfo([]string{path1, path2})
	if err != nil {
		t.Fatalf("Failed to get merge info: %v", err)
	}

	if info.TotalDocuments != 2 {
		t.Errorf("Expected 2 documents, got %d", info.TotalDocuments)
	}
	if info.TotalParagraphs != 3 {
		t.Errorf("Expected 3 paragraphs, got %d", info.TotalParagraphs)
	}
	if info.TotalTables != 2 {
		t.Errorf("Expected 2 tables, got %d", info.TotalTables)
	}
}

func TestGetMergePDFInfo(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test PDFs
	pdf1 := pdf.New()
	pdf1.AddPage()
	pdf1.AddPage()

	pdf2 := pdf.New()
	pdf2.AddPage()
	pdf2.AddPage()
	pdf2.AddPage()

	path1 := filepath.Join(tmpDir, "pdf1.pdf")
	path2 := filepath.Join(tmpDir, "pdf2.pdf")

	pdf1.Save(path1)
	pdf2.Save(path2)

	// Get merge info
	info, err := GetMergePDFInfo([]string{path1, path2})
	if err != nil {
		t.Fatalf("Failed to get merge info: %v", err)
	}

	if info.TotalDocuments != 2 {
		t.Errorf("Expected 2 documents, got %d", info.TotalDocuments)
	}
	if info.TotalPages != 5 {
		t.Errorf("Expected 5 pages, got %d", info.TotalPages)
	}
}
