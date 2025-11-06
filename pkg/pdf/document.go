package pdf

import (
	"fmt"
)

// Document represents a PDF document structure
type Document struct {
	FilePath string
	Pages    []*Page
	Metadata *Metadata
}

// Page represents a single page in the PDF
type Page struct {
	Number    int
	Content   []Content
	Width     float64
	Height    float64
	Margin    Margin
}

// Content represents content on a page (text, image, table, etc.)
type Content interface {
	Type() string
}

// TextContent represents text content
type TextContent struct {
	Text       string
	X, Y       float64
	FontSize   float64
	FontFamily string
	Bold       bool
	Italic     bool
	Color      string
}

func (t TextContent) Type() string { return "text" }

// TableContent represents a table
type TableContent struct {
	X, Y        float64
	Rows        [][]string
	ColumnWidth []float64
	HeaderStyle *TextStyle
	CellStyle   *TextStyle
}

func (t TableContent) Type() string { return "table" }

// ImageContent represents an image
type ImageContent struct {
	Path   string
	X, Y   float64
	Width  float64
	Height float64
}

func (i ImageContent) Type() string { return "image" }

// TextStyle represents text styling
type TextStyle struct {
	FontSize   float64
	FontFamily string
	Bold       bool
	Italic     bool
	Color      string
	Align      string
}

// Margin represents page margins
type Margin struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

// Metadata represents PDF metadata
type Metadata struct {
	Title    string
	Author   string
	Subject  string
	Keywords string
	Creator  string
}

// New creates a new empty PDF document
func New() *Document {
	return &Document{
		Pages: []*Page{},
		Metadata: &Metadata{
			Creator: "DocxSmith",
		},
	}
}

// AddPage adds a new page to the document
func (d *Document) AddPage() *Page {
	page := &Page{
		Number:  len(d.Pages) + 1,
		Content: []Content{},
		Width:   210, // A4 width in mm
		Height:  297, // A4 height in mm
		Margin: Margin{
			Left:   20,
			Top:    20,
			Right:  20,
			Bottom: 20,
		},
	}
	d.Pages = append(d.Pages, page)
	return page
}

// GetPageCount returns the number of pages
func (d *Document) GetPageCount() int {
	return len(d.Pages)
}

// GetPage returns a page by index (0-based)
func (d *Document) GetPage(index int) (*Page, error) {
	if index < 0 || index >= len(d.Pages) {
		return nil, fmt.Errorf("page index %d out of range", index)
	}
	return d.Pages[index], nil
}

// DeletePage removes a page by index
func (d *Document) DeletePage(index int) error {
	if index < 0 || index >= len(d.Pages) {
		return fmt.Errorf("page index %d out of range", index)
	}
	d.Pages = append(d.Pages[:index], d.Pages[index+1:]...)

	// Update page numbers
	for i := range d.Pages {
		d.Pages[i].Number = i + 1
	}
	return nil
}

// AddText adds text content to a page
func (p *Page) AddText(text string, x, y, fontSize float64) {
	content := TextContent{
		Text:       text,
		X:          x,
		Y:          y,
		FontSize:   fontSize,
		FontFamily: "Arial",
		Bold:       false,
		Italic:     false,
		Color:      "000000",
	}
	p.Content = append(p.Content, content)
}

// AddTextStyled adds styled text content to a page
func (p *Page) AddTextStyled(text string, x, y float64, style TextStyle) {
	content := TextContent{
		Text:       text,
		X:          x,
		Y:          y,
		FontSize:   style.FontSize,
		FontFamily: style.FontFamily,
		Bold:       style.Bold,
		Italic:     style.Italic,
		Color:      style.Color,
	}
	p.Content = append(p.Content, content)
}

// GetText extracts all text from the page
func (p *Page) GetText() string {
	var result string
	for _, content := range p.Content {
		if tc, ok := content.(TextContent); ok {
			result += tc.Text + " "
		}
	}
	return result
}

// GetAllText extracts all text from the document
func (d *Document) GetAllText() string {
	var result string
	for _, page := range d.Pages {
		result += page.GetText()
	}
	return result
}

// SetMetadata sets the document metadata
func (d *Document) SetMetadata(title, author, subject string) {
	if d.Metadata == nil {
		d.Metadata = &Metadata{}
	}
	d.Metadata.Title = title
	d.Metadata.Author = author
	d.Metadata.Subject = subject
	d.Metadata.Creator = "DocxSmith"
}
