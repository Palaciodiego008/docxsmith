package operations

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
	"github.com/Palaciodiego008/docxsmith/pkg/pdf"
)

// SplitOptions holds options for splitting documents
type SplitOptions struct {
	// OutputPattern is the pattern for output files (e.g., "chapter_{n}.docx")
	OutputPattern string

	// OutputDir is the directory for output files
	OutputDir string
}

// DefaultSplitOptions returns default split options
func DefaultSplitOptions() SplitOptions {
	return SplitOptions{
		OutputPattern: "part_{n}",
		OutputDir:     ".",
	}
}

// SplitDOCXByParagraphs splits a DOCX document by paragraph ranges
func SplitDOCXByParagraphs(inputPath string, ranges []ParagraphRange, opts SplitOptions) ([]string, error) {
	doc, err := docx.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open document: %w", err)
	}

	outputFiles := []string{}
	totalParagraphs := doc.GetParagraphCount()

	for i, r := range ranges {
		// Validate range
		if r.Start < 0 || r.End >= totalParagraphs || r.Start > r.End {
			return nil, fmt.Errorf("invalid range [%d:%d], document has %d paragraphs", r.Start, r.End, totalParagraphs)
		}

		// Create new document with paragraphs in range
		newDoc := docx.New()
		for j := r.Start; j <= r.End; j++ {
			newDoc.Body.Paragraphs = append(newDoc.Body.Paragraphs, doc.Body.Paragraphs[j])
		}

		// Generate output filename
		ext := filepath.Ext(inputPath)
		base := strings.TrimSuffix(filepath.Base(inputPath), ext)
		pattern := strings.ReplaceAll(opts.OutputPattern, "{n}", fmt.Sprintf("%d", i+1))
		pattern = strings.ReplaceAll(pattern, "{base}", base)

		if !strings.HasSuffix(pattern, ext) {
			pattern += ext
		}

		outputPath := filepath.Join(opts.OutputDir, pattern)

		// Save split document
		if err := newDoc.Save(outputPath); err != nil {
			return nil, fmt.Errorf("failed to save split document: %w", err)
		}

		outputFiles = append(outputFiles, outputPath)
	}

	return outputFiles, nil
}

// SplitPDFByPages splits a PDF document by page ranges
func SplitPDFByPages(inputPath string, ranges []PageRange, opts SplitOptions) ([]string, error) {
	doc, err := pdf.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}

	outputFiles := []string{}
	totalPages := doc.GetPageCount()

	for i, r := range ranges {
		// Validate range
		if r.Start < 0 || r.End >= totalPages || r.Start > r.End {
			return nil, fmt.Errorf("invalid page range [%d:%d], document has %d pages", r.Start, r.End, totalPages)
		}

		// Create new PDF with pages in range
		newDoc := pdf.New()
		for j := r.Start; j <= r.End; j++ {
			page, err := doc.GetPage(j)
			if err != nil {
				return nil, fmt.Errorf("failed to get page %d: %w", j, err)
			}

			newPage := newDoc.AddPage()
			newPage.Width = page.Width
			newPage.Height = page.Height
			newPage.Margin = page.Margin
			newPage.Content = append(newPage.Content, page.Content...)
		}

		// Generate output filename
		ext := filepath.Ext(inputPath)
		base := strings.TrimSuffix(filepath.Base(inputPath), ext)
		pattern := strings.ReplaceAll(opts.OutputPattern, "{n}", fmt.Sprintf("%d", i+1))
		pattern = strings.ReplaceAll(pattern, "{base}", base)

		if !strings.HasSuffix(pattern, ext) {
			pattern += ext
		}

		outputPath := filepath.Join(opts.OutputDir, pattern)

		// Save split PDF
		if err := newDoc.Save(outputPath); err != nil {
			return nil, fmt.Errorf("failed to save split PDF: %w", err)
		}

		outputFiles = append(outputFiles, outputPath)
	}

	return outputFiles, nil
}

// SplitDOCXByCount splits a DOCX into N equal parts
func SplitDOCXByCount(inputPath string, count int, opts SplitOptions) ([]string, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be positive")
	}

	doc, err := docx.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open document: %w", err)
	}

	totalParagraphs := doc.GetParagraphCount()
	if totalParagraphs == 0 {
		return nil, fmt.Errorf("document has no paragraphs")
	}

	// Calculate paragraphs per part
	parasPerPart := totalParagraphs / count
	if parasPerPart == 0 {
		parasPerPart = 1
	}

	ranges := []ParagraphRange{}
	start := 0

	for i := 0; i < count && start < totalParagraphs; i++ {
		end := start + parasPerPart - 1
		if i == count-1 {
			// Last part gets remaining paragraphs
			end = totalParagraphs - 1
		}
		if end >= totalParagraphs {
			end = totalParagraphs - 1
		}

		ranges = append(ranges, ParagraphRange{Start: start, End: end})
		start = end + 1
	}

	return SplitDOCXByParagraphs(inputPath, ranges, opts)
}

// SplitPDFByCount splits a PDF into N equal parts
func SplitPDFByCount(inputPath string, count int, opts SplitOptions) ([]string, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be positive")
	}

	doc, err := pdf.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}

	totalPages := doc.GetPageCount()
	if totalPages == 0 {
		return nil, fmt.Errorf("PDF has no pages")
	}

	// Calculate pages per part
	pagesPerPart := totalPages / count
	if pagesPerPart == 0 {
		pagesPerPart = 1
	}

	ranges := []PageRange{}
	start := 0

	for i := 0; i < count && start < totalPages; i++ {
		end := start + pagesPerPart - 1
		if i == count-1 {
			// Last part gets remaining pages
			end = totalPages - 1
		}
		if end >= totalPages {
			end = totalPages - 1
		}

		ranges = append(ranges, PageRange{Start: start, End: end})
		start = end + 1
	}

	return SplitPDFByPages(inputPath, ranges, opts)
}

// SplitDOCXByHeadings splits a DOCX by heading levels (smart split)
func SplitDOCXByHeadings(inputPath string, headingLevel int, opts SplitOptions) ([]string, error) {
	doc, err := docx.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open document: %w", err)
	}

	// Find paragraphs with heading style
	headingIndices := []int{}
	for i, para := range doc.Body.Paragraphs {
		if isHeading(&para, headingLevel) {
			headingIndices = append(headingIndices, i)
		}
	}

	if len(headingIndices) == 0 {
		return nil, fmt.Errorf("no headings found at level %d", headingLevel)
	}

	// Create ranges between headings
	ranges := []ParagraphRange{}
	for i := 0; i < len(headingIndices); i++ {
		start := headingIndices[i]
		end := doc.GetParagraphCount() - 1

		if i < len(headingIndices)-1 {
			end = headingIndices[i+1] - 1
		}

		ranges = append(ranges, ParagraphRange{Start: start, End: end})
	}

	// Use heading text in filename if possible
	outputFiles := []string{}
	for i, r := range ranges {
		newDoc := docx.New()
		for j := r.Start; j <= r.End; j++ {
			newDoc.Body.Paragraphs = append(newDoc.Body.Paragraphs, doc.Body.Paragraphs[j])
		}

		// Try to get heading text for filename
		headingText := ""
		if r.Start < doc.GetParagraphCount() {
			text, _ := doc.GetParagraphText(r.Start)
			headingText = sanitizeFilename(text)
			if len(headingText) > 50 {
				headingText = headingText[:50]
			}
		}

		// Generate filename
		ext := filepath.Ext(inputPath)
		base := strings.TrimSuffix(filepath.Base(inputPath), ext)

		var pattern string
		if headingText != "" {
			pattern = strings.ReplaceAll(opts.OutputPattern, "{n}", fmt.Sprintf("%d", i+1))
			pattern = strings.ReplaceAll(pattern, "{base}", base)
			pattern = strings.ReplaceAll(pattern, "{title}", headingText)
		} else {
			pattern = strings.ReplaceAll(opts.OutputPattern, "{n}", fmt.Sprintf("%d", i+1))
			pattern = strings.ReplaceAll(pattern, "{base}", base)
		}

		if !strings.HasSuffix(pattern, ext) {
			pattern += ext
		}

		outputPath := filepath.Join(opts.OutputDir, pattern)

		if err := newDoc.Save(outputPath); err != nil {
			return nil, fmt.Errorf("failed to save split document: %w", err)
		}

		outputFiles = append(outputFiles, outputPath)
	}

	return outputFiles, nil
}

// ParagraphRange represents a range of paragraphs
type ParagraphRange struct {
	Start int
	End   int
}

// PageRange represents a range of pages
type PageRange struct {
	Start int
	End   int
}

// isHeading checks if a paragraph is a heading of the specified level
func isHeading(para *docx.Paragraph, level int) bool {
	if para.Props == nil || para.Props.Style == nil {
		return false
	}

	styleName := strings.ToLower(para.Props.Style.Val)
	expectedStyle := fmt.Sprintf("heading%d", level)

	return strings.Contains(styleName, expectedStyle) || styleName == expectedStyle
}

// sanitizeFilename removes invalid characters from a filename
func sanitizeFilename(s string) string {
	// Remove invalid filename characters
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := s

	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}

	// Replace multiple spaces with single space
	result = strings.Join(strings.Fields(result), " ")

	return strings.TrimSpace(result)
}

// ParsePageRanges parses page range strings like "1-5,7,9-12"
func ParsePageRanges(rangeStr string, maxPages int) ([]PageRange, error) {
	ranges := []PageRange{}
	parts := strings.Split(rangeStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			// Range like "1-5"
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}

			var start, end int
			if _, err := fmt.Sscanf(rangeParts[0], "%d", &start); err != nil {
				return nil, fmt.Errorf("invalid start page: %s", rangeParts[0])
			}
			if _, err := fmt.Sscanf(rangeParts[1], "%d", &end); err != nil {
				return nil, fmt.Errorf("invalid end page: %s", rangeParts[1])
			}

			// Convert to 0-indexed
			start--
			end--

			if start < 0 || end >= maxPages || start > end {
				return nil, fmt.Errorf("invalid range [%d:%d], document has %d pages", start+1, end+1, maxPages)
			}

			ranges = append(ranges, PageRange{Start: start, End: end})
		} else {
			// Single page like "7"
			var page int
			if _, err := fmt.Sscanf(part, "%d", &page); err != nil {
				return nil, fmt.Errorf("invalid page number: %s", part)
			}

			// Convert to 0-indexed
			page--

			if page < 0 || page >= maxPages {
				return nil, fmt.Errorf("page %d out of range, document has %d pages", page+1, maxPages)
			}

			ranges = append(ranges, PageRange{Start: page, End: page})
		}
	}

	return ranges, nil
}
