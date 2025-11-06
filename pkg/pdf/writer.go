package pdf

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

// Save saves the PDF document to a file
func (d *Document) Save(filePath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Set metadata
	if d.Metadata != nil {
		pdf.SetTitle(d.Metadata.Title, false)
		pdf.SetAuthor(d.Metadata.Author, false)
		pdf.SetSubject(d.Metadata.Subject, false)
		pdf.SetCreator(d.Metadata.Creator, false)
	}

	// Process each page
	for _, page := range d.Pages {
		pdf.AddPage()

		// Set margins
		pdf.SetMargins(page.Margin.Left, page.Margin.Top, page.Margin.Right)
		pdf.SetAutoPageBreak(true, page.Margin.Bottom)

		// Render content
		for _, content := range page.Content {
			switch c := content.(type) {
			case TextContent:
				renderText(pdf, c)
			case TableContent:
				renderTable(pdf, c)
			}
		}
	}

	// Save to file
	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		return fmt.Errorf("failed to save PDF: %w", err)
	}

	return nil
}

// renderText renders text content
func renderText(pdf *gofpdf.Fpdf, tc TextContent) {
	// Set font style
	style := ""
	if tc.Bold {
		style += "B"
	}
	if tc.Italic {
		style += "I"
	}

	// Set font
	fontFamily := tc.FontFamily
	if fontFamily == "" {
		fontFamily = "Arial"
	}
	pdf.SetFont(fontFamily, style, tc.FontSize)

	// Set text color
	if tc.Color != "" && tc.Color != "000000" {
		r, g, b := hexToRGB(tc.Color)
		pdf.SetTextColor(r, g, b)
	} else {
		pdf.SetTextColor(0, 0, 0)
	}

	// Set position and write text
	pdf.SetXY(tc.X, tc.Y)
	pdf.Cell(0, tc.FontSize*0.35, tc.Text)
}

// renderTable renders a table
func renderTable(pdf *gofpdf.Fpdf, tc TableContent) {
	pdf.SetXY(tc.X, tc.Y)

	// Calculate column widths if not provided
	colWidths := tc.ColumnWidth
	if len(colWidths) == 0 && len(tc.Rows) > 0 && len(tc.Rows[0]) > 0 {
		numCols := len(tc.Rows[0])
		availableWidth := 170.0 // A4 width minus margins
		colWidth := availableWidth / float64(numCols)
		colWidths = make([]float64, numCols)
		for i := range colWidths {
			colWidths[i] = colWidth
		}
	}

	// Render rows
	for i, row := range tc.Rows {
		for j, cell := range row {
			if j >= len(colWidths) {
				break
			}

			// Use header style for first row, cell style for others
			if i == 0 && tc.HeaderStyle != nil {
				pdf.SetFont(tc.HeaderStyle.FontFamily, "B", tc.HeaderStyle.FontSize)
				pdf.SetFillColor(200, 200, 200) // Light gray background
			} else if tc.CellStyle != nil {
				style := ""
				if tc.CellStyle.Bold {
					style = "B"
				}
				pdf.SetFont(tc.CellStyle.FontFamily, style, tc.CellStyle.FontSize)
				pdf.SetFillColor(255, 255, 255) // White background
			} else {
				pdf.SetFont("Arial", "", 10)
				pdf.SetFillColor(255, 255, 255)
			}

			// Draw cell with border
			pdf.CellFormat(colWidths[j], 8, cell, "1", 0, "L", true, 0, "")
		}
		pdf.Ln(-1) // New line
	}
}

// hexToRGB converts hex color to RGB
func hexToRGB(hex string) (int, int, int) {
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

// SaveAs saves the document to a new file
func (d *Document) SaveAs(filePath string) error {
	return d.Save(filePath)
}
