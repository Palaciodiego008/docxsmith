# DocxSmith - The Document Forge

<p align="center">
  <strong>A powerful and elegant Go library and CLI tool for manipulating .docx files</strong>
</p>

## Features

- **Create** new .docx documents from scratch
- **Read** and parse existing .docx files
- **Modify** document content programmatically
- **Add** paragraphs with rich formatting (bold, italic, colors, sizes)
- **Delete** paragraphs or ranges of content
- **Find** and **replace** text throughout documents
- **Tables** support (create, modify, delete)
- **Extract** text content from documents
- **CLI tool** for command-line operations
- **Zero dependencies** for core functionality
- **Well-tested** with comprehensive test coverage

## Installation

### As a Library

```bash
go get github.com/Palaciodiego008/docxsmith
```

### As a CLI Tool

```bash
go install github.com/Palaciodiego008/docxsmith/cmd/docxsmith@latest
```

Or build from source:

```bash
git clone https://github.com/Palaciodiego008/docxsmith.git
cd docxsmith
go build -o docxsmith ./cmd/docxsmith
```

## Quick Start

### Using as a Library

```go
package main

import (
    "log"
    "github.com/Palaciodiego008/docxsmith/pkg/docx"
)

func main() {
    // Create a new document
    doc := docx.New()

    // Add content
    doc.AddParagraph("Welcome to DocxSmith!")
    doc.AddParagraph("This is bold text", docx.WithBold())
    doc.AddParagraph("This is colored text", docx.WithColor("FF0000"))

    // Save the document
    if err := doc.Save("output.docx"); err != nil {
        log.Fatal(err)
    }
}
```

### Using the CLI

```bash
# Create a new document
docxsmith create -output hello.docx -text "Hello, World!"

# Add content to an existing document
docxsmith add -input hello.docx -output hello2.docx -text "New paragraph" -bold

# Find text in a document
docxsmith find -input hello.docx -text "World"

# Replace text
docxsmith replace -input hello.docx -output hello3.docx -old "World" -new "DocxSmith"

# Extract text
docxsmith extract -input hello.docx

# Create a table
docxsmith table -input hello.docx -output table.docx -create -rows 3 -cols 4
```

## Library API

### Creating Documents

```go
// Create a new empty document
doc := docx.New()

// Create from an existing template
doc, err := docx.CreateFromTemplate("template.docx")

// Open an existing document
doc, err := docx.Open("existing.docx")
```

### Working with Paragraphs

```go
// Add a simple paragraph
doc.AddParagraph("Simple text")

// Add with formatting
doc.AddParagraph("Bold text", docx.WithBold())
doc.AddParagraph("Italic text", docx.WithItalic())
doc.AddParagraph("Colored text", docx.WithColor("0000FF"))
doc.AddParagraph("Large text", docx.WithSize("32"))
doc.AddParagraph("Centered text", docx.WithAlignment("center"))

// Combine multiple options
doc.AddParagraph("Fancy text",
    docx.WithBold(),
    docx.WithItalic(),
    docx.WithColor("FF0000"),
    docx.WithSize("28"))

// Add paragraph at specific position
doc.AddParagraphAt(2, "Inserted text")

// Delete a paragraph
doc.DeleteParagraph(0)

// Delete a range of paragraphs
doc.DeleteParagraphsRange(0, 5)
```

### Text Operations

```go
// Find text in document
indices := doc.FindText("search term")
// Returns slice of paragraph indices where text was found

// Replace all occurrences
count := doc.ReplaceText("old", "new")

// Replace in specific paragraph
doc.ReplaceTextInParagraph(2, "old", "new")

// Get all text content
text := doc.GetText()

// Get text from specific paragraph
text, err := doc.GetParagraphText(0)
```

### Working with Tables

```go
// Create a table
table := doc.AddTable(3, 4) // 3 rows, 4 columns

// Set cell content
table.SetCellText(0, 0, "Header 1")
table.SetCellText(0, 1, "Header 2")

// Get cell content
text, err := table.GetCellText(1, 1)

// Add a row
table.AddRow()

// Delete a row
table.DeleteRow(1)

// Get table dimensions
rows := table.GetRowCount()
cols := table.GetColumnCount()

// Delete entire table
doc.DeleteTable(0)
```

### Document Information

```go
// Get counts
paraCount := doc.GetParagraphCount()
tableCount := doc.GetTableCount()

// Clear all content
doc.Clear()

// Clone document
newDoc := doc.Clone()
```

### Saving Documents

```go
// Save to file
err := doc.Save("output.docx")

// Save to a different file
err := doc.SaveAs("copy.docx")

// Get document as bytes
data, err := doc.ToBytes()
```

## CLI Commands

### create - Create a new document

```bash
docxsmith create -output file.docx [-text "content"]
```

Options:
- `-output`: Output file path (required)
- `-text`: Initial text content (optional)

### add - Add content

```bash
docxsmith add -input in.docx -output out.docx -text "content" [options]
```

Options:
- `-input`: Input file path (required)
- `-output`: Output file path (required)
- `-text`: Text to add (required)
- `-at`: Insert at specific index (optional)
- `-bold`: Make text bold
- `-italic`: Make text italic
- `-size`: Font size (e.g., "24" for 12pt)
- `-color`: Text color (hex without #)
- `-align`: Alignment (left, center, right, both)

### delete - Delete content

```bash
docxsmith delete -input in.docx -output out.docx [options]
```

Options:
- `-input`: Input file path (required)
- `-output`: Output file path (required)
- `-paragraph`: Paragraph index to delete
- `-start` & `-end`: Delete range of paragraphs
- `-table`: Table index to delete

### replace - Replace text

```bash
docxsmith replace -input in.docx -output out.docx -old "text" -new "replacement"
```

Options:
- `-input`: Input file path (required)
- `-output`: Output file path (required)
- `-old`: Text to replace (required)
- `-new`: Replacement text (required)
- `-paragraph`: Only replace in specific paragraph

### find - Find text

```bash
docxsmith find -input file.docx -text "search"
```

Options:
- `-input`: Input file path (required)
- `-text`: Text to find (required)

### extract - Extract text

```bash
docxsmith extract -input file.docx [-output text.txt]
```

Options:
- `-input`: Input file path (required)
- `-output`: Output text file (optional, prints to stdout if omitted)

### table - Table operations

```bash
docxsmith table -input in.docx -output out.docx [options]
```

Options:
- `-input`: Input file path (required)
- `-output`: Output file path (required)
- `-create`: Create a new table
- `-rows`: Number of rows (default: 2)
- `-cols`: Number of columns (default: 2)
- `-set`: Set cell text (format: "tableIdx,row,col,text")

### info - Document information

```bash
docxsmith info -input file.docx
```

Options:
- `-input`: Input file path (required)

### clear - Clear all content

```bash
docxsmith clear -input in.docx -output out.docx
```

Options:
- `-input`: Input file path (required)
- `-output`: Output file path (required)

## Examples

See the [examples](./examples) directory for more comprehensive examples:

```bash
# Run the basic usage example
cd examples
go run basic_usage.go
```

This will generate several example documents demonstrating various features.

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run tests with verbose output:

```bash
go test -v ./pkg/docx
```

## Project Structure

```
docxsmith/
├── cmd/
│   └── docxsmith/          # CLI tool
│       └── main.go
├── pkg/
│   └── docx/               # Core library
│       ├── document.go     # Document structure
│       ├── reader.go       # Reading .docx files
│       ├── writer.go       # Writing .docx files
│       ├── operations.go   # Document operations
│       ├── table.go        # Table operations
│       ├── creator.go      # Document creation
│       ├── *_test.go       # Tests
├── examples/               # Usage examples
├── testdata/              # Test fixtures
├── go.mod
└── README.md
```

## How It Works

.docx files are actually ZIP archives containing XML files. DocxSmith:

1. Unzips the .docx file
2. Parses the XML content (mainly `word/document.xml`)
3. Manipulates the XML structure
4. Serializes back to XML
5. Repackages as a ZIP file with .docx extension

The library handles all the complexity of the Office Open XML format while providing a simple, intuitive API.

## Limitations

- Currently focuses on document content (paragraphs and tables)
- Advanced features like images, charts, and headers/footers are not yet supported
- Complex formatting and styles have limited support
- Does not preserve all metadata from original documents

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - feel free to use this project for any purpose.

## Author

Diego Palacio ([@Palaciodiego008](https://github.com/Palaciodiego008))

## Acknowledgments

- Built with Go's standard library
- Inspired by the need for simple .docx manipulation
- Name inspired by blacksmiths who forge powerful tools

---

**DocxSmith** - Forging documents with precision and elegance.
