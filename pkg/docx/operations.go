package docx

import (
	"fmt"
	"strings"
)

// AddParagraph adds a new paragraph to the document
func (d *Document) AddParagraph(text string, opts ...ParagraphOption) {
	p := Paragraph{
		Runs: []Run{
			{
				Text: []Text{
					{
						Space:   "preserve",
						Content: text,
					},
				},
			},
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(&p)
	}

	d.Body.Paragraphs = append(d.Body.Paragraphs, p)
}

// AddParagraphAt inserts a paragraph at a specific index
func (d *Document) AddParagraphAt(index int, text string, opts ...ParagraphOption) error {
	if index < 0 || index > len(d.Body.Paragraphs) {
		return fmt.Errorf("index %d out of range", index)
	}

	p := Paragraph{
		Runs: []Run{
			{
				Text: []Text{
					{
						Space:   "preserve",
						Content: text,
					},
				},
			},
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(&p)
	}

	// Insert at index
	d.Body.Paragraphs = append(
		d.Body.Paragraphs[:index],
		append([]Paragraph{p}, d.Body.Paragraphs[index:]...)...,
	)

	return nil
}

// DeleteParagraph removes a paragraph by index
func (d *Document) DeleteParagraph(index int) error {
	if index < 0 || index >= len(d.Body.Paragraphs) {
		return fmt.Errorf("paragraph index %d out of range", index)
	}

	d.Body.Paragraphs = append(
		d.Body.Paragraphs[:index],
		d.Body.Paragraphs[index+1:]...,
	)

	return nil
}

// DeleteParagraphsRange deletes multiple paragraphs from start to end (inclusive)
func (d *Document) DeleteParagraphsRange(start, end int) error {
	if start < 0 || end >= len(d.Body.Paragraphs) || start > end {
		return fmt.Errorf("invalid range [%d:%d]", start, end)
	}

	d.Body.Paragraphs = append(
		d.Body.Paragraphs[:start],
		d.Body.Paragraphs[end+1:]...,
	)

	return nil
}

// ReplaceText replaces all occurrences of old text with new text
func (d *Document) ReplaceText(oldText, newText string) int {
	count := 0
	for i := range d.Body.Paragraphs {
		for j := range d.Body.Paragraphs[i].Runs {
			for k := range d.Body.Paragraphs[i].Runs[j].Text {
				text := &d.Body.Paragraphs[i].Runs[j].Text[k]
				if strings.Contains(text.Content, oldText) {
					text.Content = strings.ReplaceAll(text.Content, oldText, newText)
					count++
				}
			}
		}
	}
	return count
}

// ReplaceTextInParagraph replaces text in a specific paragraph
func (d *Document) ReplaceTextInParagraph(index int, oldText, newText string) error {
	if index < 0 || index >= len(d.Body.Paragraphs) {
		return fmt.Errorf("paragraph index %d out of range", index)
	}

	p := &d.Body.Paragraphs[index]
	for j := range p.Runs {
		for k := range p.Runs[j].Text {
			text := &p.Runs[j].Text[k]
			text.Content = strings.ReplaceAll(text.Content, oldText, newText)
		}
	}

	return nil
}

// Clear removes all paragraphs and tables from the document
func (d *Document) Clear() {
	d.Body.Paragraphs = []Paragraph{}
	d.Body.Tables = []Table{}
}

// GetParagraphCount returns the number of paragraphs
func (d *Document) GetParagraphCount() int {
	return len(d.Body.Paragraphs)
}

// GetTableCount returns the number of tables
func (d *Document) GetTableCount() int {
	return len(d.Body.Tables)
}

// DeleteTable removes a table by index
func (d *Document) DeleteTable(index int) error {
	if index < 0 || index >= len(d.Body.Tables) {
		return fmt.Errorf("table index %d out of range", index)
	}

	d.Body.Tables = append(
		d.Body.Tables[:index],
		d.Body.Tables[index+1:]...,
	)

	return nil
}

// ParagraphOption is a function type for configuring paragraphs
type ParagraphOption func(*Paragraph)

// WithBold makes the paragraph text bold
func WithBold() ParagraphOption {
	return func(p *Paragraph) {
		for i := range p.Runs {
			if p.Runs[i].Props == nil {
				p.Runs[i].Props = &RProps{}
			}
			p.Runs[i].Props.Bold = &Bold{}
		}
	}
}

// WithItalic makes the paragraph text italic
func WithItalic() ParagraphOption {
	return func(p *Paragraph) {
		for i := range p.Runs {
			if p.Runs[i].Props == nil {
				p.Runs[i].Props = &RProps{}
			}
			p.Runs[i].Props.Italic = &Italic{}
		}
	}
}

// WithSize sets the font size (in half-points, e.g., 24 = 12pt)
func WithSize(size string) ParagraphOption {
	return func(p *Paragraph) {
		for i := range p.Runs {
			if p.Runs[i].Props == nil {
				p.Runs[i].Props = &RProps{}
			}
			p.Runs[i].Props.Size = &Size{Val: size}
		}
	}
}

// WithColor sets the text color (hex without #, e.g., "FF0000" for red)
func WithColor(color string) ParagraphOption {
	return func(p *Paragraph) {
		for i := range p.Runs {
			if p.Runs[i].Props == nil {
				p.Runs[i].Props = &RProps{}
			}
			p.Runs[i].Props.Color = &Color{Val: color}
		}
	}
}

// WithAlignment sets paragraph alignment ("left", "center", "right", "both")
func WithAlignment(align string) ParagraphOption {
	return func(p *Paragraph) {
		if p.Props == nil {
			p.Props = &PProps{}
		}
		p.Props.Jc = &Jc{Val: align}
	}
}

// WithStyle sets a paragraph style
func WithStyle(styleName string) ParagraphOption {
	return func(p *Paragraph) {
		if p.Props == nil {
			p.Props = &PProps{}
		}
		p.Props.Style = &PStyle{Val: styleName}
	}
}
