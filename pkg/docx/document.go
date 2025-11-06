package docx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// Document represents a .docx document structure
type Document struct {
	FilePath    string
	Body        *Body
	Styles      *Styles
	ContentTypes *ContentTypes
	Rels        *Relationships
	files       map[string][]byte // All files in the docx zip
}

// Body represents the document body
type Body struct {
	XMLName    xml.Name    `xml:"body"`
	Paragraphs []Paragraph `xml:"p"`
	Tables     []Table     `xml:"tbl"`
}

// Paragraph represents a paragraph in the document
type Paragraph struct {
	XMLName xml.Name `xml:"p"`
	Runs    []Run    `xml:"r"`
	Props   *PProps  `xml:"pPr,omitempty"`
}

// Run represents a text run
type Run struct {
	XMLName xml.Name `xml:"r"`
	Props   *RProps  `xml:"rPr,omitempty"`
	Text    []Text   `xml:"t"`
	Tab     *Tab     `xml:"tab,omitempty"`
	Break   *Break   `xml:"br,omitempty"`
}

// Text represents text content
type Text struct {
	XMLName   xml.Name `xml:"t"`
	Space     string   `xml:"space,attr,omitempty"`
	Content   string   `xml:",chardata"`
}

// PProps represents paragraph properties
type PProps struct {
	XMLName xml.Name `xml:"pPr"`
	Style   *PStyle  `xml:"pStyle,omitempty"`
	Jc      *Jc      `xml:"jc,omitempty"` // Justification
	Spacing *Spacing `xml:"spacing,omitempty"`
}

// RProps represents run properties
type RProps struct {
	XMLName xml.Name `xml:"rPr"`
	Bold    *Bold    `xml:"b,omitempty"`
	Italic  *Italic  `xml:"i,omitempty"`
	Size    *Size    `xml:"sz,omitempty"`
	Color   *Color   `xml:"color,omitempty"`
}

// Bold represents bold formatting
type Bold struct {
	XMLName xml.Name `xml:"b"`
}

// Italic represents italic formatting
type Italic struct {
	XMLName xml.Name `xml:"i"`
}

// Size represents font size
type Size struct {
	XMLName xml.Name `xml:"sz"`
	Val     string   `xml:"val,attr"`
}

// Color represents text color
type Color struct {
	XMLName xml.Name `xml:"color"`
	Val     string   `xml:"val,attr"`
}

// Tab represents a tab character
type Tab struct {
	XMLName xml.Name `xml:"tab"`
}

// Break represents a line break
type Break struct {
	XMLName xml.Name `xml:"br"`
}

// PStyle represents paragraph style
type PStyle struct {
	XMLName xml.Name `xml:"pStyle"`
	Val     string   `xml:"val,attr"`
}

// Jc represents text justification
type Jc struct {
	XMLName xml.Name `xml:"jc"`
	Val     string   `xml:"val,attr"` // left, center, right, both
}

// Spacing represents paragraph spacing
type Spacing struct {
	XMLName xml.Name `xml:"spacing"`
	Before  string   `xml:"before,attr,omitempty"`
	After   string   `xml:"after,attr,omitempty"`
	Line    string   `xml:"line,attr,omitempty"`
}

// Styles represents document styles
type Styles struct {
	XMLName xml.Name `xml:"styles"`
	// Add more style definitions as needed
}

// ContentTypes represents [Content_Types].xml
type ContentTypes struct {
	XMLName xml.Name `xml:"Types"`
	// Add content type definitions
}

// Relationships represents document relationships
type Relationships struct {
	XMLName xml.Name `xml:"Relationships"`
	// Add relationships
}

// GetText extracts all text from the document
func (d *Document) GetText() string {
	var texts []string
	for _, p := range d.Body.Paragraphs {
		for _, r := range p.Runs {
			for _, t := range r.Text {
				texts = append(texts, t.Content)
			}
		}
	}
	return strings.Join(texts, " ")
}

// FindText searches for text in the document and returns paragraph indices
func (d *Document) FindText(searchText string) []int {
	var indices []int
	searchLower := strings.ToLower(searchText)

	for i, p := range d.Body.Paragraphs {
		var paragraphText string
		for _, r := range p.Runs {
			for _, t := range r.Text {
				paragraphText += t.Content
			}
		}
		if strings.Contains(strings.ToLower(paragraphText), searchLower) {
			indices = append(indices, i)
		}
	}
	return indices
}

// GetParagraphText returns text from a specific paragraph
func (d *Document) GetParagraphText(index int) (string, error) {
	if index < 0 || index >= len(d.Body.Paragraphs) {
		return "", fmt.Errorf("paragraph index %d out of range", index)
	}

	var texts []string
	p := d.Body.Paragraphs[index]
	for _, r := range p.Runs {
		for _, t := range r.Text {
			texts = append(texts, t.Content)
		}
	}
	return strings.Join(texts, ""), nil
}

// readZipFile reads a file from the zip archive
func readZipFile(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return io.ReadAll(rc)
}

// saveZipFile saves data to the zip archive
func saveZipFile(w *zip.Writer, name string, data []byte) error {
	fw, err := w.Create(name)
	if err != nil {
		return err
	}
	_, err = fw.Write(data)
	return err
}
