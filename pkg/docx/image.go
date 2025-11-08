package docx

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

// Drawing represents a drawing element in a run
type Drawing struct {
	XMLName xml.Name `xml:"drawing"`
	Inline  *Inline  `xml:"inline"`
}

// Inline represents an inline drawing
type Inline struct {
	XMLName    xml.Name    `xml:"inline"`
	DistT      string      `xml:"distT,attr,omitempty"`
	DistB      string      `xml:"distB,attr,omitempty"`
	DistL      string      `xml:"distL,attr,omitempty"`
	DistR      string      `xml:"distR,attr,omitempty"`
	Extent     *Extent     `xml:"extent"`
	EffectExt  *EffectExt  `xml:"effectExtent"`
	DocPr      *DocPr      `xml:"docPr"`
	CNvGraphic *CNvGraphic `xml:"cNvGraphicFramePr"`
	Graphic    *Graphic    `xml:"graphic"`
}

// Extent represents the size of the drawing
type Extent struct {
	XMLName xml.Name `xml:"extent"`
	Cx      string   `xml:"cx,attr"`
	Cy      string   `xml:"cy,attr"`
}

// EffectExt represents effect extents
type EffectExt struct {
	XMLName xml.Name `xml:"effectExtent"`
	L       string   `xml:"l,attr"`
	T       string   `xml:"t,attr"`
	R       string   `xml:"r,attr"`
	B       string   `xml:"b,attr"`
}

// DocPr represents document properties
type DocPr struct {
	XMLName xml.Name `xml:"docPr"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

// CNvGraphic represents graphic frame properties
type CNvGraphic struct {
	XMLName xml.Name `xml:"cNvGraphicFramePr"`
}

// Graphic represents the graphic element
type Graphic struct {
	XMLName     xml.Name     `xml:"graphic"`
	GraphicData *GraphicData `xml:"graphicData"`
}

// GraphicData represents graphic data
type GraphicData struct {
	XMLName xml.Name `xml:"graphicData"`
	URI     string   `xml:"uri,attr"`
	Pic     *Pic     `xml:"pic"`
}

// Pic represents a picture
type Pic struct {
	XMLName  xml.Name  `xml:"pic"`
	NvPicPr  *NvPicPr  `xml:"nvPicPr"`
	BlipFill *BlipFill `xml:"blipFill"`
	SpPr     *SpPr     `xml:"spPr"`
}

// NvPicPr represents non-visual picture properties
type NvPicPr struct {
	XMLName  xml.Name   `xml:"nvPicPr"`
	CNvPr    *CNvPr     `xml:"cNvPr"`
	CNvPicPr *CNvPicPr2 `xml:"cNvPicPr"`
}

// CNvPr represents common non-visual properties
type CNvPr struct {
	XMLName xml.Name `xml:"cNvPr"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

// CNvPicPr2 represents picture-specific non-visual properties
type CNvPicPr2 struct {
	XMLName xml.Name `xml:"cNvPicPr"`
}

// BlipFill represents the image fill
type BlipFill struct {
	XMLName xml.Name `xml:"blipFill"`
	Blip    *Blip    `xml:"blip"`
	Stretch *Stretch `xml:"stretch"`
}

// Blip represents the image reference
type Blip struct {
	XMLName xml.Name `xml:"blip"`
	Embed   string   `xml:"embed,attr"`
}

// Stretch represents stretch properties
type Stretch struct {
	XMLName  xml.Name  `xml:"stretch"`
	FillRect *FillRect `xml:"fillRect"`
}

// FillRect represents fill rectangle
type FillRect struct {
	XMLName xml.Name `xml:"fillRect"`
}

// SpPr represents shape properties
type SpPr struct {
	XMLName  xml.Name  `xml:"spPr"`
	Xfrm     *Xfrm     `xml:"xfrm"`
	PrstGeom *PrstGeom `xml:"prstGeom"`
}

// Xfrm represents transform properties
type Xfrm struct {
	XMLName xml.Name `xml:"xfrm"`
	Ext     *XfrmExt `xml:"ext"`
}

// XfrmExt represents transform extent (different from drawing extent)
type XfrmExt struct {
	XMLName xml.Name `xml:"ext"`
	Cx      string   `xml:"cx,attr"`
	Cy      string   `xml:"cy,attr"`
}

// PrstGeom represents preset geometry
type PrstGeom struct {
	XMLName xml.Name `xml:"prstGeom"`
	Prst    string   `xml:"prst,attr"`
}

// ImageOptions holds configuration for image insertion
type ImageOptions struct {
	Width  int // Width in pixels
	Height int // Height in pixels
}

// ImageOption is a function type for configuring images
type ImageOption func(*ImageOptions)

// WithImageWidth sets the image width in pixels
func WithImageWidth(width int) ImageOption {
	return func(opts *ImageOptions) {
		opts.Width = width
	}
}

// WithImageHeight sets the image height in pixels
func WithImageHeight(height int) ImageOption {
	return func(opts *ImageOptions) {
		opts.Height = height
	}
}

// AddImage adds an image to the document
func (d *Document) AddImage(imagePath string, opts ...ImageOption) error {
	// Validate image file
	if err := d.validateImageFile(imagePath); err != nil {
		return err
	}

	// Apply options
	options := &ImageOptions{
		Width:  200, // Default width
		Height: 150, // Default height
	}
	for _, opt := range opts {
		opt(options)
	}

	// Create image paragraph
	p, err := d.createImageParagraph(imagePath, options)
	if err != nil {
		return err
	}

	// Add to document
	d.Body.Paragraphs = append(d.Body.Paragraphs, *p)
	return nil
}

// AddImageAt inserts an image at a specific paragraph index
func (d *Document) AddImageAt(index int, imagePath string, opts ...ImageOption) error {
	if index < 0 || index > len(d.Body.Paragraphs) {
		return fmt.Errorf("index %d out of range", index)
	}

	// Validate image file
	if err := d.validateImageFile(imagePath); err != nil {
		return err
	}

	// Apply options
	options := &ImageOptions{
		Width:  200,
		Height: 150,
	}
	for _, opt := range opts {
		opt(options)
	}

	// Create image paragraph
	p, err := d.createImageParagraph(imagePath, options)
	if err != nil {
		return err
	}

	// Insert at index
	d.Body.Paragraphs = append(
		d.Body.Paragraphs[:index],
		append([]Paragraph{*p}, d.Body.Paragraphs[index:]...)...,
	)

	return nil
}

// GetImageCount returns the number of images in the document
func (d *Document) GetImageCount() int {
	count := 0
	for _, p := range d.Body.Paragraphs {
		for _, r := range p.Runs {
			if r.Drawing != nil {
				count++
			}
		}
	}
	return count
}

// validateImageFile checks if the file exists and is a supported image format
func (d *Document) validateImageFile(imagePath string) error {
	// Check if file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("image file does not exist: %s", imagePath)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(imagePath))
	supportedFormats := []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".svg", ".ico", ".tiff", ".tif", ".heic", ".heif"}

	if !slices.Contains(supportedFormats, ext) {
		return fmt.Errorf("unsupported image format: %s", ext)
	}

	// Basic file content validation - check if it's actually an image
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()

	// Read first few bytes to check magic numbers
	header := make([]byte, 8)
	_, err = file.Read(header)
	if err != nil {
		return fmt.Errorf("failed to read image header: %v", err)
	}

	// Check for common image magic numbers
	if !isValidImageHeader(header, ext) {
		return fmt.Errorf("file does not appear to be a valid %s image", ext)
	}

	return nil
}

// isValidImageHeader checks if the file header matches the expected format
func isValidImageHeader(header []byte, ext string) bool {
	switch ext {
	case ".png":
		// PNG signature: 89 50 4E 47 0D 0A 1A 0A
		return len(header) >= 8 &&
			header[0] == 0x89 && header[1] == 0x50 &&
			header[2] == 0x4E && header[3] == 0x47
	case ".jpg", ".jpeg":
		// JPEG signature: FF D8 FF
		return len(header) >= 3 &&
			header[0] == 0xFF && header[1] == 0xD8 && header[2] == 0xFF
	case ".gif":
		// GIF signature: GIF87a or GIF89a
		return len(header) >= 6 &&
			string(header[0:3]) == "GIF" &&
			(string(header[3:6]) == "87a" || string(header[3:6]) == "89a")
	case ".bmp":
		// BMP signature: BM
		return len(header) >= 2 && header[0] == 0x42 && header[1] == 0x4D
	default:
		return false
	}
}

// createImageParagraph creates a paragraph containing an image
func (d *Document) createImageParagraph(imagePath string, options *ImageOptions) (*Paragraph, error) {
	// Read image file
	imageData, err := d.readImageFile(imagePath)
	if err != nil {
		return nil, err
	}

	// Generate relationship ID
	relID := fmt.Sprintf("rId%d", d.getNextRelationshipID())

	// Store image data in document files
	imageFileName := fmt.Sprintf("word/media/image%d%s", d.getNextImageID(), filepath.Ext(imagePath))
	if d.files == nil {
		d.files = make(map[string][]byte)
	}
	d.files[imageFileName] = imageData

	// Convert pixels to EMUs (English Metric Units)
	// 1 pixel = 9525 EMUs at 96 DPI
	widthEMU := strconv.Itoa(options.Width * 9525)
	heightEMU := strconv.Itoa(options.Height * 9525)

	// Create drawing structure
	drawing := &Drawing{
		Inline: &Inline{
			DistT: "0",
			DistB: "0",
			DistL: "0",
			DistR: "0",
			Extent: &Extent{
				Cx: widthEMU,
				Cy: heightEMU,
			},
			EffectExt: &EffectExt{
				L: "0",
				T: "0",
				R: "0",
				B: "0",
			},
			DocPr: &DocPr{
				ID:   strconv.Itoa(d.getNextImageID()),
				Name: fmt.Sprintf("Picture %d", d.getNextImageID()),
			},
			CNvGraphic: &CNvGraphic{},
			Graphic: &Graphic{
				GraphicData: &GraphicData{
					URI: "http://schemas.openxmlformats.org/drawingml/2006/picture",
					Pic: &Pic{
						NvPicPr: &NvPicPr{
							CNvPr: &CNvPr{
								ID:   strconv.Itoa(d.getNextImageID()),
								Name: filepath.Base(imagePath),
							},
							CNvPicPr: &CNvPicPr2{},
						},
						BlipFill: &BlipFill{
							Blip: &Blip{
								Embed: relID,
							},
							Stretch: &Stretch{
								FillRect: &FillRect{},
							},
						},
						SpPr: &SpPr{
							Xfrm: &Xfrm{
								Ext: &XfrmExt{
									Cx: widthEMU,
									Cy: heightEMU,
								},
							},
							PrstGeom: &PrstGeom{
								Prst: "rect",
							},
						},
					},
				},
			},
		},
	}

	// Create paragraph with image
	p := &Paragraph{
		Runs: []Run{
			{
				Drawing: drawing,
			},
		},
	}

	return p, nil
}

// readImageFile reads an image file and returns its data
func (d *Document) readImageFile(imagePath string) ([]byte, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %v", err)
	}

	return data, nil
}

// getNextRelationshipID returns the next available relationship ID
func (d *Document) getNextRelationshipID() int {
	// Simple implementation - in a real scenario, you'd track existing relationships
	return len(d.files) + 1
}

// getNextImageID returns the next available image ID
func (d *Document) getNextImageID() int {
	return d.GetImageCount() + 1
}

// GetImageAsBase64 returns an image as base64 string (utility function)
func GetImageAsBase64(imagePath string) (string, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
