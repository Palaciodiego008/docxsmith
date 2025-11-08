package docx

import (
	"encoding/xml"
	"fmt"
)

// HeaderFooterType represents the type of header or footer
type HeaderFooterType string

const (
	HeaderTypeDefault HeaderFooterType = "header-default"
	HeaderTypeFirst   HeaderFooterType = "header-first"
	HeaderTypeEven    HeaderFooterType = "header-even"
	FooterTypeDefault HeaderFooterType = "footer-default"
	FooterTypeFirst   HeaderFooterType = "footer-first"
	FooterTypeEven    HeaderFooterType = "footer-even"
)

// HeaderFooter represents a header or footer element
type HeaderFooter struct {
	XMLName    xml.Name `xml:"hdr,omitempty"`
	Type       HeaderFooterType
	Paragraphs []Paragraph `xml:"p"`
	IsFooter   bool
}

// HeaderFooterManager interface defines operations for managing headers and footers
type HeaderFooterManager interface {
	SetHeader(hfType HeaderFooterType, content string, opts ...HeaderFooterOption) error
	SetFooter(hfType HeaderFooterType, content string, opts ...HeaderFooterOption) error
	GetHeader(hfType HeaderFooterType) (*HeaderFooter, error)
	GetFooter(hfType HeaderFooterType) (*HeaderFooter, error)
	RemoveHeader(hfType HeaderFooterType) error
	RemoveFooter(hfType HeaderFooterType) error
	HasHeader(hfType HeaderFooterType) bool
	HasFooter(hfType HeaderFooterType) bool
}

// HeaderFooterOption is a function type for configuring headers and footers
type HeaderFooterOption func(*HeaderFooterConfig)

// HeaderFooterConfig holds configuration for headers and footers
type HeaderFooterConfig struct {
	Alignment string
	Bold      bool
	Italic    bool
	Size      string
	Color     string
	Font      string
}

// HeaderFooterService implements HeaderFooterManager
type HeaderFooterService struct {
	document *Document
	headers  map[HeaderFooterType]*HeaderFooter
	footers  map[HeaderFooterType]*HeaderFooter
}

// NewHeaderFooterService creates a new header/footer service
func NewHeaderFooterService(doc *Document) HeaderFooterManager {
	return &HeaderFooterService{
		document: doc,
		headers:  make(map[HeaderFooterType]*HeaderFooter),
		footers:  make(map[HeaderFooterType]*HeaderFooter),
	}
}

// SetHeader sets a header with the specified type and content
func (hfs *HeaderFooterService) SetHeader(hfType HeaderFooterType, content string, opts ...HeaderFooterOption) error {
	if err := hfs.validateHeaderFooterType(hfType, false); err != nil {
		return err
	}

	config := hfs.applyOptions(opts...)
	header := hfs.createHeaderFooter(hfType, content, config, false)
	hfs.headers[hfType] = header

	return nil
}

// SetFooter sets a footer with the specified type and content
func (hfs *HeaderFooterService) SetFooter(hfType HeaderFooterType, content string, opts ...HeaderFooterOption) error {
	if err := hfs.validateHeaderFooterType(hfType, true); err != nil {
		return err
	}

	config := hfs.applyOptions(opts...)
	footer := hfs.createHeaderFooter(hfType, content, config, true)
	hfs.footers[hfType] = footer

	return nil
}

// GetHeader retrieves a header by type
func (hfs *HeaderFooterService) GetHeader(hfType HeaderFooterType) (*HeaderFooter, error) {
	header, exists := hfs.headers[hfType]
	if !exists {
		return nil, fmt.Errorf("header of type %s not found", hfType)
	}
	return header, nil
}

// GetFooter retrieves a footer by type
func (hfs *HeaderFooterService) GetFooter(hfType HeaderFooterType) (*HeaderFooter, error) {
	footer, exists := hfs.footers[hfType]
	if !exists {
		return nil, fmt.Errorf("footer of type %s not found", hfType)
	}
	return footer, nil
}

// RemoveHeader removes a header by type
func (hfs *HeaderFooterService) RemoveHeader(hfType HeaderFooterType) error {
	if !hfs.HasHeader(hfType) {
		return fmt.Errorf("header of type %s does not exist", hfType)
	}
	delete(hfs.headers, hfType)
	return nil
}

// RemoveFooter removes a footer by type
func (hfs *HeaderFooterService) RemoveFooter(hfType HeaderFooterType) error {
	if !hfs.HasFooter(hfType) {
		return fmt.Errorf("footer of type %s does not exist", hfType)
	}
	delete(hfs.footers, hfType)
	return nil
}

// HasHeader checks if a header exists
func (hfs *HeaderFooterService) HasHeader(hfType HeaderFooterType) bool {
	_, exists := hfs.headers[hfType]
	return exists
}

// HasFooter checks if a footer exists
func (hfs *HeaderFooterService) HasFooter(hfType HeaderFooterType) bool {
	_, exists := hfs.footers[hfType]
	return exists
}

// Private methods

func (hfs *HeaderFooterService) validateHeaderFooterType(hfType HeaderFooterType, isFooter bool) error {
	if isFooter {
		validTypes := []HeaderFooterType{FooterTypeDefault, FooterTypeFirst, FooterTypeEven}
		for _, validType := range validTypes {
			if hfType == validType {
				return nil
			}
		}
		return fmt.Errorf("invalid footer type: %s", hfType)
	} else {
		validTypes := []HeaderFooterType{HeaderTypeDefault, HeaderTypeFirst, HeaderTypeEven}
		for _, validType := range validTypes {
			if hfType == validType {
				return nil
			}
		}
		return fmt.Errorf("invalid header type: %s", hfType)
	}
}

func (hfs *HeaderFooterService) applyOptions(opts ...HeaderFooterOption) *HeaderFooterConfig {
	config := &HeaderFooterConfig{
		Alignment: "left",
		Font:      "Calibri",
		Size:      "22", // 11pt
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}

func (hfs *HeaderFooterService) createHeaderFooter(hfType HeaderFooterType, content string, config *HeaderFooterConfig, isFooter bool) *HeaderFooter {
	paragraph := hfs.createStyledParagraph(content, config)

	hf := &HeaderFooter{
		Type:       hfType,
		Paragraphs: []Paragraph{paragraph},
		IsFooter:   isFooter,
	}

	if isFooter {
		hf.XMLName = xml.Name{Local: "ftr"}
	}

	return hf
}

func (hfs *HeaderFooterService) createStyledParagraph(content string, config *HeaderFooterConfig) Paragraph {
	run := Run{
		Text: []Text{{
			Space:   "preserve",
			Content: content,
		}},
	}

	// Apply formatting
	if config.Bold || config.Italic || config.Size != "" || config.Color != "" || config.Font != "" {
		run.Props = &RProps{}

		if config.Bold {
			run.Props.Bold = &Bold{}
		}
		if config.Italic {
			run.Props.Italic = &Italic{}
		}
		if config.Size != "" {
			run.Props.Size = &Size{Val: config.Size}
		}
		if config.Color != "" {
			run.Props.Color = &Color{Val: config.Color}
		}
		if config.Font != "" {
			run.Props.RFonts = &RFonts{ASCII: config.Font}
		}
	}

	paragraph := Paragraph{
		Runs: []Run{run},
	}

	// Apply alignment
	if config.Alignment != "left" {
		paragraph.Props = &PProps{
			Jc: &Jc{Val: config.Alignment},
		}
	}

	return paragraph
}

// Option functions

// WithHFAlignment sets the text alignment for headers/footers
func WithHFAlignment(align string) HeaderFooterOption {
	return func(config *HeaderFooterConfig) {
		config.Alignment = align
	}
}

// WithHFBold makes the header/footer text bold
func WithHFBold() HeaderFooterOption {
	return func(config *HeaderFooterConfig) {
		config.Bold = true
	}
}

// WithHFItalic makes the header/footer text italic
func WithHFItalic() HeaderFooterOption {
	return func(config *HeaderFooterConfig) {
		config.Italic = true
	}
}

// WithHFFontSize sets the font size for headers/footers
func WithHFFontSize(size string) HeaderFooterOption {
	return func(config *HeaderFooterConfig) {
		config.Size = size
	}
}

// WithHFTextColor sets the text color for headers/footers
func WithHFTextColor(color string) HeaderFooterOption {
	return func(config *HeaderFooterConfig) {
		config.Color = color
	}
}

// WithHFFont sets the font family for headers/footers
func WithHFFont(font string) HeaderFooterOption {
	return func(config *HeaderFooterConfig) {
		config.Font = font
	}
}
