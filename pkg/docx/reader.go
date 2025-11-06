package docx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// Open opens and reads a .docx file
func Open(filePath string) (*Document, error) {
	doc := &Document{
		FilePath: filePath,
		files:    make(map[string][]byte),
	}

	// Open the docx file (which is a zip archive)
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open docx file: %w", err)
	}
	defer r.Close()

	// Read all files from the zip
	var documentXML []byte
	for _, f := range r.File {
		data, err := readZipFile(f)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", f.Name, err)
		}
		doc.files[f.Name] = data

		// Parse the main document.xml
		if f.Name == "word/document.xml" {
			documentXML = data
		}
	}

	if documentXML == nil {
		return nil, fmt.Errorf("document.xml not found in docx file")
	}

	// Parse the XML document
	if err := doc.parseDocument(documentXML); err != nil {
		return nil, fmt.Errorf("failed to parse document.xml: %w", err)
	}

	return doc, nil
}

// parseDocument parses the main document.xml content
func (d *Document) parseDocument(data []byte) error {
	// Define the document structure with namespace
	type WDocument struct {
		XMLName xml.Name `xml:"http://schemas.openxmlformats.org/wordprocessingml/2006/main document"`
		Body    *Body    `xml:"body"`
	}

	var doc WDocument
	if err := xml.Unmarshal(data, &doc); err != nil {
		return err
	}

	if doc.Body == nil {
		d.Body = &Body{
			Paragraphs: []Paragraph{},
			Tables:     []Table{},
		}
	} else {
		d.Body = doc.Body
	}

	return nil
}

// ReadBytes reads a .docx file from bytes
func ReadBytes(data []byte) (*Document, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "docx-*.docx")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write data to temp file
	if _, err := tmpFile.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write to temp file: %w", err)
	}

	// Open the temp file
	return Open(tmpFile.Name())
}

// ReadFrom reads a .docx document from an io.Reader
func ReadFrom(r io.Reader) (*Document, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}
	return ReadBytes(data)
}
