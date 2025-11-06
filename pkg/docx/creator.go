package docx

import (
	"archive/zip"
)

// New creates a new empty document
func New() *Document {
	return &Document{
		Body: &Body{
			Paragraphs: []Paragraph{},
			Tables:     []Table{},
		},
		files: getDefaultDocxFiles(),
	}
}

// getDefaultDocxFiles returns the minimum required files for a valid .docx
func getDefaultDocxFiles() map[string][]byte {
	files := make(map[string][]byte)

	// [Content_Types].xml
	files["[Content_Types].xml"] = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
	<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
	<Default Extension="xml" ContentType="application/xml"/>
	<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`)

	// _rels/.rels
	files["_rels/.rels"] = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
	<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`)

	// word/_rels/document.xml.rels
	files["word/_rels/document.xml.rels"] = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
</Relationships>`)

	return files
}

// CreateFromTemplate creates a new document based on a template
func CreateFromTemplate(templatePath string) (*Document, error) {
	return Open(templatePath)
}

// Clone creates a deep copy of the document
func (d *Document) Clone() *Document {
	newDoc := &Document{
		FilePath: d.FilePath,
		Body: &Body{
			Paragraphs: make([]Paragraph, len(d.Body.Paragraphs)),
			Tables:     make([]Table, len(d.Body.Tables)),
		},
		files: make(map[string][]byte),
	}

	// Copy paragraphs
	copy(newDoc.Body.Paragraphs, d.Body.Paragraphs)

	// Copy tables
	copy(newDoc.Body.Tables, d.Body.Tables)

	// Copy files
	for k, v := range d.files {
		newDoc.files[k] = append([]byte(nil), v...)
	}

	return newDoc
}

// CreateMinimalDocx creates a minimal valid .docx file for testing
func CreateMinimalDocx(outputPath string) error {
	doc := New()
	doc.AddParagraph("This is a test document created by DocxSmith.")
	doc.AddParagraph("")
	doc.AddParagraph("It contains multiple paragraphs for testing purposes.")

	return doc.Save(outputPath)
}

// WriteToZip writes the document to an open zip.Writer (useful for streaming)
func (d *Document) WriteToZip(w *zip.Writer) error {
	// Marshal the document
	documentXML, err := d.marshalDocument()
	if err != nil {
		return err
	}

	// Update the document.xml
	d.files["word/document.xml"] = documentXML

	// Write all files to the zip
	for name, data := range d.files {
		if err := saveZipFile(w, name, data); err != nil {
			return err
		}
	}

	return nil
}
