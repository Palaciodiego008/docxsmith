package operations

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
	"github.com/Palaciodiego008/docxsmith/pkg/pdf"
)

func TestSplitDOCXByParagraphs(t *testing.T) {
	tests := []struct {
		name           string
		totalParas     int
		ranges         []ParagraphRange
		expectedFiles  int
		expectedParas  []int // paragraphs in each output file
	}{
		{
			name:       "Split into 2 parts",
			totalParas: 10,
			ranges: []ParagraphRange{
				{Start: 0, End: 4},
				{Start: 5, End: 9},
			},
			expectedFiles: 2,
			expectedParas: []int{5, 5},
		},
		{
			name:       "Split into 3 unequal parts",
			totalParas: 10,
			ranges: []ParagraphRange{
				{Start: 0, End: 2},
				{Start: 3, End: 6},
				{Start: 7, End: 9},
			},
			expectedFiles: 3,
			expectedParas: []int{3, 4, 3},
		},
		{
			name:       "Single paragraph range",
			totalParas: 5,
			ranges: []ParagraphRange{
				{Start: 2, End: 2},
			},
			expectedFiles: 1,
			expectedParas: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create test document
			doc := docx.New()
			for i := 0; i < tt.totalParas; i++ {
				doc.AddParagraph(fmt.Sprintf("Paragraph %d", i+1))
			}

			inputPath := filepath.Join(tmpDir, "input.docx")
			if err := doc.Save(inputPath); err != nil {
				t.Fatalf("Failed to save test document: %v", err)
			}

			// Split document
			opts := SplitOptions{
				OutputPattern: "part{n}.docx",
				OutputDir:     tmpDir,
			}

			outputFiles, err := SplitDOCXByParagraphs(inputPath, tt.ranges, opts)
			if err != nil {
				t.Fatalf("Split failed: %v", err)
			}

			// Verify number of output files
			if len(outputFiles) != tt.expectedFiles {
				t.Errorf("Expected %d output files, got %d", tt.expectedFiles, len(outputFiles))
			}

			// Verify content of each file
			for i, outPath := range outputFiles {
				outDoc, err := docx.Open(outPath)
				if err != nil {
					t.Fatalf("Failed to open output file %s: %v", outPath, err)
				}

				if i < len(tt.expectedParas) {
					if outDoc.GetParagraphCount() != tt.expectedParas[i] {
						t.Errorf("File %d: expected %d paragraphs, got %d",
							i+1, tt.expectedParas[i], outDoc.GetParagraphCount())
					}
				}
			}
		})
	}
}

func TestSplitPDFByPages(t *testing.T) {
	tests := []struct {
		name          string
		totalPages    int
		ranges        []PageRange
		expectedFiles int
		expectedPages []int
	}{
		{
			name:       "Split into 2 parts",
			totalPages: 10,
			ranges: []PageRange{
				{Start: 0, End: 4},
				{Start: 5, End: 9},
			},
			expectedFiles: 2,
			expectedPages: []int{5, 5},
		},
		{
			name:       "Extract single page",
			totalPages: 5,
			ranges: []PageRange{
				{Start: 2, End: 2},
			},
			expectedFiles: 1,
			expectedPages: []int{1},
		},
		{
			name:       "Multiple ranges",
			totalPages: 10,
			ranges: []PageRange{
				{Start: 0, End: 1},
				{Start: 3, End: 5},
				{Start: 7, End: 9},
			},
			expectedFiles: 3,
			expectedPages: []int{2, 3, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create test PDF
			doc := pdf.New()
			for i := 0; i < tt.totalPages; i++ {
				page := doc.AddPage()
				page.AddText(fmt.Sprintf("Page %d", i+1), 20, 30, 12)
			}

			inputPath := filepath.Join(tmpDir, "input.pdf")
			if err := doc.Save(inputPath); err != nil {
				t.Fatalf("Failed to save test PDF: %v", err)
			}

			// Split PDF
			opts := SplitOptions{
				OutputPattern: "part{n}.pdf",
				OutputDir:     tmpDir,
			}

			outputFiles, err := SplitPDFByPages(inputPath, tt.ranges, opts)
			if err != nil {
				t.Fatalf("Split failed: %v", err)
			}

			// Verify number of output files
			if len(outputFiles) != tt.expectedFiles {
				t.Errorf("Expected %d output files, got %d", tt.expectedFiles, len(outputFiles))
			}

			// Verify content of each file
			for i, outPath := range outputFiles {
				outDoc, err := pdf.Open(outPath)
				if err != nil {
					t.Fatalf("Failed to open output file %s: %v", outPath, err)
				}

				if i < len(tt.expectedPages) {
					if outDoc.GetPageCount() != tt.expectedPages[i] {
						t.Errorf("File %d: expected %d pages, got %d",
							i+1, tt.expectedPages[i], outDoc.GetPageCount())
					}
				}
			}
		})
	}
}

func TestSplitDOCXByCount(t *testing.T) {
	tests := []struct {
		name          string
		totalParas    int
		splitCount    int
		expectedFiles int
		minParasEach  int
	}{
		{
			name:          "Split 10 paragraphs into 2 parts",
			totalParas:    10,
			splitCount:    2,
			expectedFiles: 2,
			minParasEach:  5,
		},
		{
			name:          "Split 9 paragraphs into 3 parts",
			totalParas:    9,
			splitCount:    3,
			expectedFiles: 3,
			minParasEach:  3,
		},
		{
			name:          "Split 5 paragraphs into 2 parts",
			totalParas:    5,
			splitCount:    2,
			expectedFiles: 2,
			minParasEach:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create test document
			doc := docx.New()
			for i := 0; i < tt.totalParas; i++ {
				doc.AddParagraph(fmt.Sprintf("Paragraph %d", i+1))
			}

			inputPath := filepath.Join(tmpDir, "input.docx")
			if err := doc.Save(inputPath); err != nil {
				t.Fatalf("Failed to save test document: %v", err)
			}

			// Split document
			opts := DefaultSplitOptions()
			opts.OutputDir = tmpDir

			outputFiles, err := SplitDOCXByCount(inputPath, tt.splitCount, opts)
			if err != nil {
				t.Fatalf("Split failed: %v", err)
			}

			if len(outputFiles) != tt.expectedFiles {
				t.Errorf("Expected %d files, got %d", tt.expectedFiles, len(outputFiles))
			}

			// Verify total paragraphs are preserved
			totalParasAfter := 0
			for _, outPath := range outputFiles {
				outDoc, err := docx.Open(outPath)
				if err != nil {
					t.Fatalf("Failed to open output: %v", err)
				}
				totalParasAfter += outDoc.GetParagraphCount()
			}

			if totalParasAfter != tt.totalParas {
				t.Errorf("Total paragraphs not preserved: expected %d, got %d",
					tt.totalParas, totalParasAfter)
			}
		})
	}
}

func TestSplitPDFByCount(t *testing.T) {
	tests := []struct {
		name          string
		totalPages    int
		splitCount    int
		expectedFiles int
	}{
		{
			name:          "Split 10 pages into 2 parts",
			totalPages:    10,
			splitCount:    2,
			expectedFiles: 2,
		},
		{
			name:          "Split 12 pages into 3 parts",
			totalPages:    12,
			splitCount:    3,
			expectedFiles: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create test PDF
			doc := pdf.New()
			for i := 0; i < tt.totalPages; i++ {
				page := doc.AddPage()
				page.AddText(fmt.Sprintf("Page %d", i+1), 20, 30, 12)
			}

			inputPath := filepath.Join(tmpDir, "input.pdf")
			if err := doc.Save(inputPath); err != nil {
				t.Fatalf("Failed to save test PDF: %v", err)
			}

			// Split PDF
			opts := DefaultSplitOptions()
			opts.OutputDir = tmpDir

			outputFiles, err := SplitPDFByCount(inputPath, tt.splitCount, opts)
			if err != nil {
				t.Fatalf("Split failed: %v", err)
			}

			if len(outputFiles) != tt.expectedFiles {
				t.Errorf("Expected %d files, got %d", tt.expectedFiles, len(outputFiles))
			}

			// Verify total pages are preserved
			totalPagesAfter := 0
			for _, outPath := range outputFiles {
				outDoc, err := pdf.Open(outPath)
				if err != nil {
					t.Fatalf("Failed to open output: %v", err)
				}
				totalPagesAfter += outDoc.GetPageCount()
			}

			if totalPagesAfter != tt.totalPages {
				t.Errorf("Total pages not preserved: expected %d, got %d",
					tt.totalPages, totalPagesAfter)
			}
		})
	}
}

func TestParsePageRanges(t *testing.T) {
	tests := []struct {
		name          string
		rangeStr      string
		maxPages      int
		expectedCount int
		expectError   bool
	}{
		{
			name:          "Single range",
			rangeStr:      "1-5",
			maxPages:      10,
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Multiple ranges",
			rangeStr:      "1-3,5-7,9-10",
			maxPages:      10,
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:          "Single pages",
			rangeStr:      "1,3,5,7",
			maxPages:      10,
			expectedCount: 4,
			expectError:   false,
		},
		{
			name:          "Mixed ranges and single pages",
			rangeStr:      "1-3,5,7-9",
			maxPages:      10,
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:        "Invalid range - out of bounds",
			rangeStr:    "1-20",
			maxPages:    10,
			expectError: true,
		},
		{
			name:        "Invalid format",
			rangeStr:    "abc",
			maxPages:    10,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ranges, err := ParsePageRanges(tt.rangeStr, tt.maxPages)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(ranges) != tt.expectedCount {
				t.Errorf("Expected %d ranges, got %d", tt.expectedCount, len(ranges))
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			input:    "File/With\\Slashes",
			expected: "File_With_Slashes",
		},
		{
			input:    "File:With*Invalid?Chars",
			expected: "File_With_Invalid_Chars",
		},
		{
			input:    "  Extra   Spaces  ",
			expected: "Extra Spaces",
		},
		{
			input:    "Normal-Filename.txt",
			expected: "Normal-Filename.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestSplitDOCXByHeadings(t *testing.T) {
	tmpDir := t.TempDir()

	// Create document with headings
	doc := docx.New()
	doc.AddParagraph("Chapter 1", docx.WithStyle("Heading1"))
	doc.AddParagraph("Content for chapter 1")
	doc.AddParagraph("More content")

	doc.AddParagraph("Chapter 2", docx.WithStyle("Heading1"))
	doc.AddParagraph("Content for chapter 2")

	doc.AddParagraph("Chapter 3", docx.WithStyle("Heading1"))
	doc.AddParagraph("Content for chapter 3")
	doc.AddParagraph("More content for 3")

	inputPath := filepath.Join(tmpDir, "book.docx")
	if err := doc.Save(inputPath); err != nil {
		t.Fatalf("Failed to save test document: %v", err)
	}

	// Split by headings
	opts := SplitOptions{
		OutputPattern: "chapter{n}.docx",
		OutputDir:     tmpDir,
	}

	outputFiles, err := SplitDOCXByHeadings(inputPath, 1, opts)
	if err != nil {
		t.Fatalf("Split by headings failed: %v", err)
	}

	// Should create 3 files (one per heading)
	if len(outputFiles) != 3 {
		t.Errorf("Expected 3 output files, got %d", len(outputFiles))
	}

	// Verify each file exists
	for i, outPath := range outputFiles {
		if _, err := os.Stat(outPath); os.IsNotExist(err) {
			t.Errorf("Output file %d does not exist: %s", i+1, outPath)
		}
	}
}

func TestSplitErrors(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T) (string, []PageRange, SplitOptions)
		expectError bool
	}{
		{
			name: "Invalid range - start > end",
			setupFunc: func(t *testing.T) (string, []PageRange, SplitOptions) {
				tmpDir := t.TempDir()
				doc := pdf.New()
				doc.AddPage()
				doc.AddPage()
				path := filepath.Join(tmpDir, "test.pdf")
				doc.Save(path)

				ranges := []PageRange{{Start: 1, End: 0}}
				return path, ranges, DefaultSplitOptions()
			},
			expectError: true,
		},
		{
			name: "Invalid range - out of bounds",
			setupFunc: func(t *testing.T) (string, []PageRange, SplitOptions) {
				tmpDir := t.TempDir()
				doc := pdf.New()
				doc.AddPage()
				path := filepath.Join(tmpDir, "test.pdf")
				doc.Save(path)

				ranges := []PageRange{{Start: 0, End: 10}}
				return path, ranges, DefaultSplitOptions()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, ranges, opts := tt.setupFunc(t)
			_, err := SplitPDFByPages(path, ranges, opts)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestSplitByCountErrors(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{
			name:  "Zero count",
			count: 0,
		},
		{
			name:  "Negative count",
			count: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			doc := docx.New()
			doc.AddParagraph("Test")
			path := filepath.Join(tmpDir, "test.docx")
			doc.Save(path)

			_, err := SplitDOCXByCount(path, tt.count, DefaultSplitOptions())
			if err == nil {
				t.Error("Expected error for invalid count")
			}
		})
	}
}
