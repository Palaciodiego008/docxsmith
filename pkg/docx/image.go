package docx

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Drawing represents a drawing element in a run
type Drawing struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/wordprocessingml/2006/main drawing"`
	Inline  *Inline  `xml:"http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing inline"`
}

// Inline represents an inline drawing
type Inline struct {
	XMLName    xml.Name    `xml:"http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing inline"`
	DistT      string      `xml:"distT,attr,omitempty"`
	DistB      string      `xml:"distB,attr,omitempty"`
	DistL      string      `xml:"distL,attr,omitempty"`
	DistR      string      `xml:"distR,attr,omitempty"`
	Extent     *Extent     `xml:"http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing extent"`
	EffectExt  *EffectExt  `xml:"http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing effectExtent"`
	DocPr      *DocPr      `xml:"http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing docPr"`
	CNvGraphic *CNvGraphic `xml:"http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing cNvGraphicFramePr"`
	Graphic    *Graphic    `xml:"http://schemas.openxmlformats.org/drawingml/2006/main graphic"`
}

// Extent represents the size of the drawing
type Extent struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing extent"`
	Cx      string   `xml:"cx,attr"`
	Cy      string   `xml:"cy,attr"`
}

// EffectExt represents effect extents
type EffectExt struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing effectExtent"`
	L       string   `xml:"l,attr"`
	T       string   `xml:"t,attr"`
	R       string   `xml:"r,attr"`
	B       string   `xml:"b,attr"`
}

// DocPr represents document properties
type DocPr struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing docPr"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

// CNvGraphic represents graphic frame properties
type CNvGraphic struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing cNvGraphicFramePr"`
}

// Graphic represents the graphic element
type Graphic struct {
	XMLName     xml.Name     `xml:"http://schemas.openxmlformats.org/drawingml/2006/main graphic"`
	GraphicData *GraphicData `xml:"http://schemas.openxmlformats.org/drawingml/2006/main graphicData"`
}

// GraphicData represents graphic data
type GraphicData struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/main graphicData"`
	URI     string   `xml:"uri,attr"`
	Pic     *Pic     `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture pic"`
}

// Pic represents a picture
type Pic struct {
	XMLName  xml.Name  `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture pic"`
	NvPicPr  *NvPicPr  `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture nvPicPr"`
	BlipFill *BlipFill `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture blipFill"`
	SpPr     *SpPr     `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture spPr"`
}

// NvPicPr represents non-visual picture properties
type NvPicPr struct {
	XMLName  xml.Name   `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture nvPicPr"`
	CNvPr    *CNvPr     `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture cNvPr"`
	CNvPicPr *CNvPicPr2 `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture cNvPicPr"`
}

// CNvPr represents common non-visual properties
type CNvPr struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture cNvPr"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

// CNvPicPr2 represents picture-specific non-visual properties
type CNvPicPr2 struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture cNvPicPr"`
}

// BlipFill represents the image fill
type BlipFill struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture blipFill"`
	Blip    *Blip    `xml:"http://schemas.openxmlformats.org/drawingml/2006/main blip"`
	Stretch *Stretch `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture stretch"`
}

// Blip represents the image reference
type Blip struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/main blip"`
	Embed   string   `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships embed,attr"`
}

// Stretch represents stretch properties
type Stretch struct {
	XMLName  xml.Name  `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture stretch"`
	FillRect *FillRect `xml:"http://schemas.openxmlformats.org/drawingml/2006/main fillRect"`
}

// FillRect represents fill rectangle
type FillRect struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/main fillRect"`
}

// SpPr represents shape properties
type SpPr struct {
	XMLName  xml.Name  `xml:"http://schemas.openxmlformats.org/drawingml/2006/picture spPr"`
	Xfrm     *Xfrm     `xml:"http://schemas.openxmlformats.org/drawingml/2006/main xfrm"`
	PrstGeom *PrstGeom `xml:"http://schemas.openxmlformats.org/drawingml/2006/main prstGeom"`
}

// Xfrm represents transform properties
type Xfrm struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/main xfrm"`
	Ext     *XfrmExt `xml:"http://schemas.openxmlformats.org/drawingml/2006/main ext"`
}

// XfrmExt represents transform extent (different from drawing extent)
type XfrmExt struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/main ext"`
	Cx      string   `xml:"cx,attr"`
	Cy      string   `xml:"cy,attr"`
}

// PrstGeom represents preset geometry
type PrstGeom struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/drawingml/2006/main prstGeom"`
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
	// Check if file exists first
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("image file does not exist: %s", imagePath)
	}

	// Read image file once
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return fmt.Errorf("failed to read image file: %v", err)
	}

	// Validate image file
	if err := d.validateImageFile(imagePath, imageData); err != nil {
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
	p, err := d.createImageParagraph(imagePath, imageData, options)
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

	// Check if file exists first
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("image file does not exist: %s", imagePath)
	}

	// Read image file once
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return fmt.Errorf("failed to read image file: %v", err)
	}

	// Validate image file
	if err := d.validateImageFile(imagePath, imageData); err != nil {
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
	p, err := d.createImageParagraph(imagePath, imageData, options)
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

// validateImageFile validates the image format and content
func (d *Document) validateImageFile(imagePath string, imageData []byte) error {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(imagePath))
	supportedFormats := []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".svg", ".ico", ".tiff", ".tif", ".heic", ".heif"}

	if !slices.Contains(supportedFormats, ext) {
		return fmt.Errorf("unsupported image format: %s", ext)
	}

	// Validate file content - check if it's actually an image
	// Read enough bytes to validate all formats (WebP needs 12 bytes)
	headerSize := min(12, len(imageData))
	if headerSize < 8 {
		return fmt.Errorf("image file too small to validate: %s", imagePath)
	}

	// Check for common image magic numbers
	header := imageData[:headerSize]
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
	case ".webp":
		// WebP signature: RIFF[size]WEBP (need to read 12 bytes)
		return len(header) >= 12 &&
			header[0] == 0x52 && header[1] == 0x49 && // "RI"
			header[2] == 0x46 && header[3] == 0x46 && // "FF"
			header[8] == 0x57 && header[9] == 0x45 && // "WE"
			header[10] == 0x42 && header[11] == 0x50 // "BP"
	case ".svg":
		// SVG should start with <?xml declaration or <svg tag
		if len(header) == 0 {
			return false
		}
		headerStr := string(header)
		trimmed := strings.TrimSpace(headerStr)
		return strings.HasPrefix(trimmed, "<?xml") || strings.HasPrefix(trimmed, "<svg")
	case ".ico":
		// ICO signature: 00 00 01 00
		return len(header) >= 4 &&
			header[0] == 0x00 && header[1] == 0x00 &&
			header[2] == 0x01 && header[3] == 0x00
	case ".tiff", ".tif":
		// TIFF signature: "II" (little-endian) or "MM" (big-endian)
		return len(header) >= 4 &&
			((header[0] == 0x49 && header[1] == 0x49) || // "II"
				(header[0] == 0x4D && header[1] == 0x4D)) // "MM"
	case ".heic", ".heif":
		// HEIC/HEIF signature: starts with ftyp box
		// Format: [size:4 bytes]["ftyp":4 bytes][brand:4 bytes]
		return len(header) >= 8 &&
			header[4] == 0x66 && header[5] == 0x74 && // "ft"
			header[6] == 0x79 && header[7] == 0x70 // "yp"
	default:
		return false
	}
}

// createImageParagraph creates a paragraph containing an image
func (d *Document) createImageParagraph(imagePath string, imageData []byte, options *ImageOptions) (*Paragraph, error) {
	// Generate relationship ID
	relID := fmt.Sprintf("rId%d", d.getNextRelationshipID())

	// Get the next image ID once to ensure consistency
	imageID := d.getNextImageID()
	imageIDStr := strconv.Itoa(imageID)

	// Store image data in document files
	imageExt := strings.ToLower(filepath.Ext(imagePath))
	imageFileName := fmt.Sprintf("word/media/image%d%s", imageID, imageExt)
	if d.files == nil {
		d.files = make(map[string][]byte)
	}
	d.files[imageFileName] = imageData

	// Update Content Types to register the image extension
	d.registerImageContentType(imageExt)

	// Update relationships to add the image relationship
	d.addImageRelationship(relID, imageFileName)

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
				ID:   imageIDStr,
				Name: fmt.Sprintf("Picture %d", imageID),
			},
			CNvGraphic: &CNvGraphic{},
			Graphic: &Graphic{
				GraphicData: &GraphicData{
					URI: "http://schemas.openxmlformats.org/drawingml/2006/picture",
					Pic: &Pic{
						NvPicPr: &NvPicPr{
							CNvPr: &CNvPr{
								ID:   imageIDStr,
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

// getNextRelationshipID returns the next available relationship ID and increments the counter
func (d *Document) getNextRelationshipID() int {
	id := d.nextRelationshipID
	d.nextRelationshipID++
	return id
}

// getNextImageID returns the next available image ID and increments the counter
func (d *Document) getNextImageID() int {
	id := d.nextImageID
	d.nextImageID++
	return id
}

// initializeImageID sets the nextImageID based on existing images in the document
func (d *Document) initializeImageID() {
	d.nextImageID = d.GetImageCount() + 1
}

// initializeRelationshipID sets the nextRelationshipID based on existing relationships in the document
func (d *Document) initializeRelationshipID() {
	maxRelID := 0

	// Check word/_rels/document.xml.rels for existing relationships
	if relsData, exists := d.files["word/_rels/document.xml.rels"]; exists {
		relsStr := string(relsData)

		// Use regex to find all relationship IDs (e.g., rId1, rId2, rId100)
		re := regexp.MustCompile(`\brId(\d+)\b`)
		matches := re.FindAllStringSubmatch(relsStr, -1)

		for _, match := range matches {
			if len(match) > 1 {
				if id, err := strconv.Atoi(match[1]); err == nil && id > maxRelID {
					maxRelID = id
				}
			}
		}
	}

	d.nextRelationshipID = maxRelID + 1
}

// registerImageContentType adds or updates the content type for an image extension
func (d *Document) registerImageContentType(ext string) {
	// Map of image extensions to MIME types
	mimeTypes := map[string]string{
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".ico":  "image/x-icon",
		".tiff": "image/tiff",
		".tif":  "image/tiff",
		".heic": "image/heic",
		".heif": "image/heif",
	}

	mimeType, exists := mimeTypes[ext]
	if !exists {
		mimeType = "image/png" // Default fallback
	}

	// Get current content types
	contentTypesData, ok := d.files["[Content_Types].xml"]
	if !ok {
		// Initialize with default if not exists
		contentTypesData = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
	<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
	<Default Extension="xml" ContentType="application/xml"/>
	<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`)
	}

	contentTypesStr := string(contentTypesData)
	extWithoutDot := strings.TrimPrefix(ext, ".")

	// Check if this extension is already registered
	extensionEntry := fmt.Sprintf(`Extension="%s"`, extWithoutDot)
	if strings.Contains(contentTypesStr, extensionEntry) {
		return // Already registered
	}

	// Add the new Default entry before the closing </Types> tag
	newEntry := fmt.Sprintf(`	<Default Extension="%s" ContentType="%s"/>`, extWithoutDot, mimeType)
	contentTypesStr = strings.Replace(contentTypesStr, "</Types>", newEntry+"\n</Types>", 1)

	d.files["[Content_Types].xml"] = []byte(contentTypesStr)
}

// addImageRelationship adds a relationship entry for an image
func (d *Document) addImageRelationship(relID, imagePath string) {
	// Get current relationships
	relsData, ok := d.files["word/_rels/document.xml.rels"]
	if !ok {
		// Initialize with default if not exists
		relsData = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
</Relationships>`)
	}

	relsStr := string(relsData)

	// Check if this relationship already exists
	if strings.Contains(relsStr, relID) {
		return // Already exists
	}

	// Extract target path (remove "word/" prefix for the relationship)
	target := strings.TrimPrefix(imagePath, "word/")

	// Add the new Relationship entry before the closing </Relationships> tag
	newRel := fmt.Sprintf(`	<Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="%s"/>`, relID, target)
	relsStr = strings.Replace(relsStr, "</Relationships>", newRel+"\n</Relationships>", 1)

	d.files["word/_rels/document.xml.rels"] = []byte(relsStr)
}

// GetImageAsBase64 returns an image as base64 string (utility function)
func GetImageAsBase64(imagePath string) (string, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
