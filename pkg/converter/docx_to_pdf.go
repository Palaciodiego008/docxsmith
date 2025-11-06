package converter

import (
	"fmt"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
	"github.com/Palaciodiego008/docxsmith/pkg/pdf"
)

// DocxToPDF converts a DOCX document to PDF
type DocxToPDF struct {
	Options ConvertOptions
}

// NewDocxToPDF creates a new DOCX to PDF converter
func NewDocxToPDF(opts ConvertOptions) *DocxToPDF {
	return &DocxToPDF{
		Options: opts,
	}
}

// Convert converts a DOCX document to PDF
func (c *DocxToPDF) Convert(doc *docx.Document, outputPath string) error {
	pdfDoc := pdf.New()

	// Set metadata
	pdfDoc.SetMetadata("Converted from DOCX", "", "")

	// Add a page
	page := pdfDoc.AddPage()

	// Current Y position for content
	currentY := page.Margin.Top

	// Convert paragraphs
	for _, para := range doc.Body.Paragraphs {
		text := ""
		isBold := false
		isItalic := false
		fontSize := c.Options.FontSize
		color := "000000"

		// Extract text and styling from runs
		for _, run := range para.Runs {
			for _, t := range run.Text {
				text += t.Content
			}

			// Check for formatting
			if run.Props != nil {
				if run.Props.Bold != nil {
					isBold = true
				}
				if run.Props.Italic != nil {
					isItalic = true
				}
				if run.Props.Size != nil && run.Props.Size.Val != "" {
					// Size in DOCX is in half-points, convert to points
					var sz float64
					fmt.Sscanf(run.Props.Size.Val, "%f", &sz)
					fontSize = sz / 2
				}
				if run.Props.Color != nil && run.Props.Color.Val != "" {
					color = run.Props.Color.Val
				}
			}
		}

		if text != "" {
			style := pdf.TextStyle{
				FontSize:   fontSize,
				FontFamily: c.Options.FontFamily,
				Bold:       isBold,
				Italic:     isItalic,
				Color:      color,
			}

			page.AddTextStyled(text, page.Margin.Left, currentY, style)
			currentY += fontSize * 0.5 // Line spacing

			// Check if we need a new page
			if currentY > page.Height-page.Margin.Bottom {
				page = pdfDoc.AddPage()
				currentY = page.Margin.Top
			}
		}
	}

	// Convert tables
	for _, table := range doc.Body.Tables {
		// Check if we need a new page for the table
		estimatedTableHeight := float64(len(table.Rows)) * 8.0
		if currentY+estimatedTableHeight > page.Height-page.Margin.Bottom {
			page = pdfDoc.AddPage()
			currentY = page.Margin.Top
		}

		// Extract table data
		rows := [][]string{}
		for _, row := range table.Rows {
			cells := []string{}
			for _, cell := range row.Cells {
				cellText := ""
				for _, p := range cell.Content {
					for _, r := range p.Runs {
						for _, t := range r.Text {
							cellText += t.Content
						}
					}
				}
				cells = append(cells, cellText)
			}
			rows = append(rows, cells)
		}

		// Add table to PDF
		tableContent := pdf.TableContent{
			X:    page.Margin.Left,
			Y:    currentY,
			Rows: rows,
			HeaderStyle: &pdf.TextStyle{
				FontSize:   c.Options.FontSize,
				FontFamily: c.Options.FontFamily,
				Bold:       true,
			},
			CellStyle: &pdf.TextStyle{
				FontSize:   c.Options.FontSize,
				FontFamily: c.Options.FontFamily,
				Bold:       false,
			},
		}
		page.Content = append(page.Content, tableContent)

		currentY += estimatedTableHeight + 5 // Add some spacing after table
	}

	// Save PDF
	return pdfDoc.Save(outputPath)
}

// ConvertFile converts a DOCX file to PDF
func ConvertDocxToPDF(inputPath, outputPath string, opts ConvertOptions) error {
	// Open DOCX
	doc, err := docx.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open DOCX: %w", err)
	}

	// Convert
	converter := NewDocxToPDF(opts)
	return converter.Convert(doc, outputPath)
}
