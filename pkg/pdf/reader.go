package pdf

import (
	"fmt"

	"github.com/ledongthuc/pdf"
)

// Open opens and reads a PDF file
func Open(filePath string) (*Document, error) {
	doc := &Document{
		FilePath: filePath,
		Pages:    []*Page{},
		Metadata: &Metadata{
			Creator: "DocxSmith",
		},
	}

	// Open PDF file
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	// Get number of pages
	numPages := r.NumPage()

	// Read each page
	for i := 1; i <= numPages; i++ {
		page := &Page{
			Number:  i,
			Content: []Content{},
			Width:   210,
			Height:  297,
			Margin: Margin{
				Left:   20,
				Top:    20,
				Right:  20,
				Bottom: 20,
			},
		}

		// Get page content
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}

		// Extract text from page
		text, err := p.GetPlainText(nil)
		if err == nil && text != "" {
			// Add text as content
			textContent := TextContent{
				Text:       text,
				X:          20,
				Y:          20,
				FontSize:   12,
				FontFamily: "Arial",
				Bold:       false,
				Italic:     false,
				Color:      "000000",
			}
			page.Content = append(page.Content, textContent)
		}

		doc.Pages = append(doc.Pages, page)
	}

	return doc, nil
}

// ReadBytes reads a PDF from bytes
func ReadBytes(data []byte) (*Document, error) {
	// For now, this is a placeholder
	// In production, you would use a library that supports reading from bytes
	return nil, fmt.Errorf("reading PDF from bytes not yet implemented")
}
