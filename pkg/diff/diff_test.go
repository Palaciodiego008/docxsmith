package diff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

func TestCompareDOCX(t *testing.T) {
	tests := []struct {
		name            string
		oldLines        []string
		newLines        []string
		expectedAdded   int
		expectedDeleted int
		expectedTotal   int
	}{
		{
			name:            "Identical documents",
			oldLines:        []string{"Line 1", "Line 2", "Line 3"},
			newLines:        []string{"Line 1", "Line 2", "Line 3"},
			expectedAdded:   0,
			expectedDeleted: 0,
			expectedTotal:   0,
		},
		{
			name:            "One line added",
			oldLines:        []string{"Line 1", "Line 2"},
			newLines:        []string{"Line 1", "Line 2", "Line 3"},
			expectedAdded:   1,
			expectedDeleted: 0,
			expectedTotal:   1,
		},
		{
			name:            "One line deleted",
			oldLines:        []string{"Line 1", "Line 2", "Line 3"},
			newLines:        []string{"Line 1", "Line 3"},
			expectedAdded:   0,
			expectedDeleted: 1,
			expectedTotal:   1,
		},
		{
			name:            "Multiple changes",
			oldLines:        []string{"Old Line 1", "Line 2", "Old Line 3"},
			newLines:        []string{"New Line 1", "Line 2", "New Line 3"},
			expectedAdded:   2,
			expectedDeleted: 2,
			expectedTotal:   4,
		},
		{
			name:            "All lines different",
			oldLines:        []string{"A", "B", "C"},
			newLines:        []string{"X", "Y", "Z"},
			expectedAdded:   3,
			expectedDeleted: 3,
			expectedTotal:   6,
		},
		{
			name:            "Empty old document",
			oldLines:        []string{},
			newLines:        []string{"Line 1", "Line 2"},
			expectedAdded:   2,
			expectedDeleted: 0,
			expectedTotal:   2,
		},
		{
			name:            "Empty new document",
			oldLines:        []string{"Line 1", "Line 2"},
			newLines:        []string{},
			expectedAdded:   0,
			expectedDeleted: 2,
			expectedTotal:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create test documents
			oldDoc := docx.New()
			for _, line := range tt.oldLines {
				oldDoc.AddParagraph(line)
			}
			oldPath := filepath.Join(tmpDir, "old.docx")
			if err := oldDoc.Save(oldPath); err != nil {
				t.Fatalf("Failed to save old doc: %v", err)
			}

			newDoc := docx.New()
			for _, line := range tt.newLines {
				newDoc.AddParagraph(line)
			}
			newPath := filepath.Join(tmpDir, "new.docx")
			if err := newDoc.Save(newPath); err != nil {
				t.Fatalf("Failed to save new doc: %v", err)
			}

			// Compare
			opts := DefaultDiffOptions()
			result, err := CompareDOCX(oldPath, newPath, opts)
			if err != nil {
				t.Fatalf("Compare failed: %v", err)
			}

			// Verify stats
			if result.Stats.TotalChanges != tt.expectedTotal {
				t.Errorf("Expected %d total changes, got %d", tt.expectedTotal, result.Stats.TotalChanges)
			}
			if result.Stats.AddedLines != tt.expectedAdded {
				t.Errorf("Expected %d added lines, got %d", tt.expectedAdded, result.Stats.AddedLines)
			}
			if result.Stats.DeletedLines != tt.expectedDeleted {
				t.Errorf("Expected %d deleted lines, got %d", tt.expectedDeleted, result.Stats.DeletedLines)
			}
		})
	}
}

func TestDiffOptions(t *testing.T) {
	tests := []struct {
		name          string
		oldLines      []string
		newLines      []string
		options       DiffOptions
		expectedTotal int
	}{
		{
			name:     "Ignore whitespace",
			oldLines: []string{"  Line 1  ", "Line 2"},
			newLines: []string{"Line 1", "Line 2"},
			options: DiffOptions{
				IgnoreWhitespace: true,
				IgnoreCase:       false,
			},
			expectedTotal: 0, // Should be identical with whitespace ignored
		},
		{
			name:     "Don't ignore whitespace",
			oldLines: []string{"  Line 1  ", "Line 2"},
			newLines: []string{"Line 1", "Line 2"},
			options: DiffOptions{
				IgnoreWhitespace: false,
				IgnoreCase:       false,
			},
			expectedTotal: 2, // Should show differences
		},
		{
			name:     "Ignore case",
			oldLines: []string{"HELLO", "WORLD"},
			newLines: []string{"hello", "world"},
			options: DiffOptions{
				IgnoreWhitespace: false,
				IgnoreCase:       true,
			},
			expectedTotal: 0, // Should be identical with case ignored
		},
		{
			name:     "Don't ignore case",
			oldLines: []string{"HELLO", "WORLD"},
			newLines: []string{"hello", "world"},
			options: DiffOptions{
				IgnoreWhitespace: false,
				IgnoreCase:       false,
			},
			expectedTotal: 4, // Should show differences
		},
		{
			name:     "Combined options",
			oldLines: []string{"  HELLO  ", "  WORLD  "},
			newLines: []string{"hello", "world"},
			options: DiffOptions{
				IgnoreWhitespace: true,
				IgnoreCase:       true,
			},
			expectedTotal: 0, // Should be identical
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create test documents
			oldDoc := docx.New()
			for _, line := range tt.oldLines {
				oldDoc.AddParagraph(line)
			}
			oldPath := filepath.Join(tmpDir, "old.docx")
			oldDoc.Save(oldPath)

			newDoc := docx.New()
			for _, line := range tt.newLines {
				newDoc.AddParagraph(line)
			}
			newPath := filepath.Join(tmpDir, "new.docx")
			newDoc.Save(newPath)

			// Compare with options
			result, err := CompareDOCX(oldPath, newPath, tt.options)
			if err != nil {
				t.Fatalf("Compare failed: %v", err)
			}

			if result.Stats.TotalChanges != tt.expectedTotal {
				t.Errorf("Expected %d total changes, got %d", tt.expectedTotal, result.Stats.TotalChanges)
			}
		})
	}
}

func TestRenderers(t *testing.T) {
	// Create test diff result
	result := &DiffResult{
		Changes: []Change{
			{Type: DiffAdded, New: "New line 1", Position: 0},
			{Type: DiffDeleted, Old: "Old line 2", Position: 1},
			{Type: DiffAdded, New: "New line 3", Position: 2},
		},
		Stats: DiffStats{
			TotalChanges: 3,
			AddedLines:   2,
			DeletedLines: 1,
		},
		OldDocument: "test_old.docx",
		NewDocument: "test_new.docx",
	}

	tests := []struct {
		name          string
		renderer      Renderer
		shouldContain []string
	}{
		{
			name:     "HTML Renderer",
			renderer: NewHTMLRenderer(true),
			shouldContain: []string{
				"<!DOCTYPE html>",
				"Document Comparison",
				"added",
				"deleted",
				"New line 1",
				"Old line 2",
			},
		},
		{
			name:     "Markdown Renderer",
			renderer: NewMarkdownRenderer(true),
			shouldContain: []string{
				"# Document Comparison",
				"## Statistics",
				"## Changes",
				"New line 1",
				"Old line 2",
			},
		},
		{
			name:     "PlainText Renderer",
			renderer: NewPlainTextRenderer(true, false),
			shouldContain: []string{
				"Document Comparison",
				"Statistics:",
				"Changes:",
				"New line 1",
				"Old line 2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tt.renderer.Render(result)
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			for _, expected := range tt.shouldContain {
				if !strings.Contains(output, expected) {
					t.Errorf("Output should contain '%s' but doesn't", expected)
				}
			}
		})
	}
}

func TestDiffTypeString(t *testing.T) {
	tests := []struct {
		diffType DiffType
		expected string
	}{
		{DiffNone, "unchanged"},
		{DiffAdded, "added"},
		{DiffDeleted, "deleted"},
		{DiffModified, "modified"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.diffType.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestEmptyDiff(t *testing.T) {
	tmpDir := t.TempDir()

	// Create identical documents
	doc1 := docx.New()
	doc1.AddParagraph("Same content")
	doc1.AddParagraph("More same content")

	doc2 := docx.New()
	doc2.AddParagraph("Same content")
	doc2.AddParagraph("More same content")

	path1 := filepath.Join(tmpDir, "doc1.docx")
	path2 := filepath.Join(tmpDir, "doc2.docx")

	doc1.Save(path1)
	doc2.Save(path2)

	// Compare
	result, err := CompareDOCX(path1, path2, DefaultDiffOptions())
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}

	if result.Stats.TotalChanges != 0 {
		t.Errorf("Expected 0 changes for identical documents, got %d", result.Stats.TotalChanges)
	}

	if len(result.Changes) != 0 {
		t.Errorf("Expected 0 change items, got %d", len(result.Changes))
	}
}

func TestHTMLRendererOutput(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "diff.html")

	result := &DiffResult{
		Changes: []Change{
			{Type: DiffAdded, New: "Added line", Position: 0},
		},
		Stats: DiffStats{
			TotalChanges: 1,
			AddedLines:   1,
		},
		OldDocument: "old.docx",
		NewDocument: "new.docx",
	}

	renderer := NewHTMLRenderer(true)
	output, err := renderer.Render(result)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Save to file
	if err := os.WriteFile(outputPath, []byte(output), 0644); err != nil {
		t.Fatalf("Failed to write HTML: %v", err)
	}

	// Verify file exists and is not empty
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Output file doesn't exist: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Output file is empty")
	}
}
