package diff

import (
	"fmt"
	"strings"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

// DiffType represents the type of change
type DiffType int

const (
	DiffNone DiffType = iota
	DiffAdded
	DiffDeleted
	DiffModified
)

func (d DiffType) String() string {
	switch d {
	case DiffAdded:
		return "added"
	case DiffDeleted:
		return "deleted"
	case DiffModified:
		return "modified"
	default:
		return "unchanged"
	}
}

// Change represents a single change in the diff
type Change struct {
	Type     DiffType
	Old      string
	New      string
	Position int // Paragraph or line number
	Context  string
}

// DiffResult represents the result of comparing two documents
type DiffResult struct {
	Changes     []Change
	Stats       DiffStats
	OldDocument string
	NewDocument string
}

// DiffStats holds statistics about the diff
type DiffStats struct {
	TotalChanges   int
	AddedLines     int
	DeletedLines   int
	ModifiedLines  int
	UnchangedLines int
}

// Differ is the interface for diff implementations
type Differ interface {
	Compare(old, new string) (*DiffResult, error)
}

// DocxDiffer compares DOCX documents
type DocxDiffer struct {
	options DiffOptions
}

// DiffOptions holds options for diff operations
type DiffOptions struct {
	// IgnoreWhitespace ignores whitespace differences
	IgnoreWhitespace bool

	// IgnoreCase ignores case differences
	IgnoreCase bool

	// ContextLines number of context lines to show around changes
	ContextLines int

	// MinChangeLength minimum length to consider a change
	MinChangeLength int
}

// DefaultDiffOptions returns default diff options
func DefaultDiffOptions() DiffOptions {
	return DiffOptions{
		IgnoreWhitespace: false,
		IgnoreCase:       false,
		ContextLines:     3,
		MinChangeLength:  1,
	}
}

// NewDocxDiffer creates a new DOCX differ
func NewDocxDiffer(opts DiffOptions) *DocxDiffer {
	return &DocxDiffer{
		options: opts,
	}
}

// Compare compares two DOCX documents
func (d *DocxDiffer) Compare(oldPath, newPath string) (*DiffResult, error) {
	// Open documents
	oldDoc, err := docx.Open(oldPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open old document: %w", err)
	}

	newDoc, err := docx.Open(newPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open new document: %w", err)
	}

	// Extract text from paragraphs
	oldLines := extractLines(oldDoc)
	newLines := extractLines(newDoc)

	// Compute diff
	changes := d.computeDiff(oldLines, newLines)

	// Calculate stats
	stats := calculateStats(changes)

	return &DiffResult{
		Changes:     changes,
		Stats:       stats,
		OldDocument: oldPath,
		NewDocument: newPath,
	}, nil
}

// extractLines extracts text lines from a document
func extractLines(doc *docx.Document) []string {
	lines := []string{}
	for _, para := range doc.Body.Paragraphs {
		text := ""
		for _, run := range para.Runs {
			for _, t := range run.Text {
				text += t.Content
			}
		}
		lines = append(lines, text)
	}
	return lines
}

// computeDiff computes the diff between two sets of lines
func (d *DocxDiffer) computeDiff(oldLines, newLines []string) []Change {
	changes := []Change{}

	// Use Myers diff algorithm (simplified implementation)
	oldLen := len(oldLines)
	newLen := len(newLines)

	// Create a DP table for LCS (Longest Common Subsequence)
	dp := make([][]int, oldLen+1)
	for i := range dp {
		dp[i] = make([]int, newLen+1)
	}

	// Fill DP table
	for i := 1; i <= oldLen; i++ {
		for j := 1; j <= newLen; j++ {
			if d.linesEqual(oldLines[i-1], newLines[j-1]) {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	// Backtrack to find changes
	i, j := oldLen, newLen
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && d.linesEqual(oldLines[i-1], newLines[j-1]) {
			// No change
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			// Addition
			changes = append([]Change{{
				Type:     DiffAdded,
				New:      newLines[j-1],
				Position: j - 1,
			}}, changes...)
			j--
		} else if i > 0 {
			// Deletion
			changes = append([]Change{{
				Type:     DiffDeleted,
				Old:      oldLines[i-1],
				Position: i - 1,
			}}, changes...)
			i--
		}
	}

	return changes
}

// linesEqual checks if two lines are equal considering options
func (d *DocxDiffer) linesEqual(line1, line2 string) bool {
	if d.options.IgnoreWhitespace {
		line1 = strings.TrimSpace(line1)
		line2 = strings.TrimSpace(line2)
	}

	if d.options.IgnoreCase {
		line1 = strings.ToLower(line1)
		line2 = strings.ToLower(line2)
	}

	return line1 == line2
}

// calculateStats calculates statistics from changes
func calculateStats(changes []Change) DiffStats {
	stats := DiffStats{}

	for _, change := range changes {
		stats.TotalChanges++
		switch change.Type {
		case DiffAdded:
			stats.AddedLines++
		case DiffDeleted:
			stats.DeletedLines++
		case DiffModified:
			stats.ModifiedLines++
		}
	}

	return stats
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// CompareDOCX is a convenience function to compare two DOCX files
func CompareDOCX(oldPath, newPath string, opts DiffOptions) (*DiffResult, error) {
	differ := NewDocxDiffer(opts)
	return differ.Compare(oldPath, newPath)
}
