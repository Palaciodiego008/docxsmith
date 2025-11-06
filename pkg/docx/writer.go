package docx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"os"
)

// Save saves the document to a file
func (d *Document) Save(filePath string) error {
	// Create output file
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Create zip writer
	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// Marshal the body back to XML
	documentXML, err := d.marshalDocument()
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	// Update the document.xml in files map
	d.files["word/document.xml"] = documentXML

	// Write all files back to the zip
	for name, data := range d.files {
		if err := saveZipFile(zipWriter, name, data); err != nil {
			return fmt.Errorf("failed to save file %s: %w", name, err)
		}
	}

	return nil
}

// SaveAs saves the document to a new file
func (d *Document) SaveAs(filePath string) error {
	return d.Save(filePath)
}

// marshalDocument marshals the document body to XML
func (d *Document) marshalDocument() ([]byte, error) {
	// Define the document structure with namespace
	type WBody struct {
		XMLName    xml.Name    `xml:"w:body"`
		Paragraphs []Paragraph `xml:"w:p"`
		Tables     []Table     `xml:"w:tbl"`
	}

	type WDocument struct {
		XMLName xml.Name `xml:"w:document"`
		Xmlns   string   `xml:"xmlns:w,attr"`
		XmlnsR  string   `xml:"xmlns:r,attr"`
		Body    WBody    `xml:"w:body"`
	}

	doc := WDocument{
		Xmlns:  "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
		XmlnsR: "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
		Body: WBody{
			Paragraphs: d.Body.Paragraphs,
			Tables:     d.Body.Tables,
		},
	}

	// Marshal with proper XML header
	output, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}

	// Add XML header
	xmlHeader := []byte(xml.Header)
	return append(xmlHeader, output...), nil
}

// ToBytes returns the document as bytes
func (d *Document) ToBytes() ([]byte, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "docx-*.docx")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Save to temp file
	if err := d.Save(tmpPath); err != nil {
		return nil, err
	}

	// Read back the file
	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read temp file: %w", err)
	}

	return data, nil
}
