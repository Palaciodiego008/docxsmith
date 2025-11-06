package converter

import (
	"fmt"
	"strings"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
	"github.com/Palaciodiego008/docxsmith/pkg/pdf"
)

// PDFToDocx converts a PDF document to DOCX
type PDFToDocx struct {
	Options ConvertOptions
}

// NewPDFToDocx creates a new PDF to DOCX converter
func NewPDFToDocx(opts ConvertOptions) *PDFToDocx {
	return &PDFToDocx{
		Options: opts,
	}
}

// Convert converts a PDF document to DOCX
func (c *PDFToDocx) Convert(pdfDoc *pdf.Document, outputPath string) error {
	docxDoc := docx.New()

	// Process each page
	for _, page := range pdfDoc.Pages {
		// Process content
		for _, content := range page.Content {
			switch c := content.(type) {
			case pdf.TextContent:
				// Convert text content
				var opts []docx.ParagraphOption

				if c.Bold {
					opts = append(opts, docx.WithBold())
				}
				if c.Italic {
					opts = append(opts, docx.WithItalic())
				}
				if c.Color != "" && c.Color != "000000" {
					opts = append(opts, docx.WithColor(c.Color))
				}

				// Convert font size (PDF points to DOCX half-points)
				if c.FontSize > 0 {
					sizeStr := fmt.Sprintf("%.0f", c.FontSize*2)
					opts = append(opts, docx.WithSize(sizeStr))
				}

				// Split by lines
				lines := strings.Split(c.Text, "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						docxDoc.AddParagraph(line, opts...)
					}
				}

			case pdf.TableContent:
				// Convert table
				if len(c.Rows) > 0 {
					// Find maximum column count across all rows
					maxCols := 0
					for _, row := range c.Rows {
						if len(row) > maxCols {
							maxCols = len(row)
						}
					}

					table := docxDoc.AddTable(len(c.Rows), maxCols)

					for i, row := range c.Rows {
						for j, cell := range row {
							table.SetCellText(i, j, cell)
						}
					}
				}
			}
		}
	}

	// Save DOCX
	return docxDoc.Save(outputPath)
}

// ConvertFile converts a PDF file to DOCX
func ConvertPDFToDocx(inputPath, outputPath string, opts ConvertOptions) error {
	// Open PDF
	pdfDoc, err := pdf.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open PDF: %w", err)
	}

	// Convert
	converter := NewPDFToDocx(opts)
	return converter.Convert(pdfDoc, outputPath)
}
