# DocxSmith Usage Guide

Comprehensive guide to using DocxSmith for document manipulation.

## Table of Contents
- [Getting Started](#getting-started)
- [DOCX Operations](#docx-operations)
- [PDF Operations](#pdf-operations)
- [Format Conversion](#format-conversion)
- [Library Usage](#library-usage)

## Getting Started

### Build the CLI

```bash
go build -o docxsmith ./cmd/docxsmith
```

### Run Tests

```bash
go test ./...
```

## DOCX Operations

### Create a New Document

```bash
./docxsmith create -output mydoc.docx -text "Hello World"
```

### Add Content with Formatting

```bash
# Bold text
./docxsmith add -input mydoc.docx -output mydoc.docx -text "Bold paragraph" -bold

# Colored text
./docxsmith add -input mydoc.docx -output mydoc.docx -text "Red text" -color FF0000

# Large centered text
./docxsmith add -input mydoc.docx -output mydoc.docx -text "Title" -size 32 -align center

# Combine multiple styles
./docxsmith add -input mydoc.docx -output mydoc.docx -text "Fancy" -bold -italic -size 24 -color 0000FF
```

### Work with Tables

```bash
# Create a table
./docxsmith table -input mydoc.docx -output mydoc.docx -create -rows 4 -cols 3

# Set cell content
./docxsmith table -input mydoc.docx -output mydoc.docx -set "0,0,0,Name"
./docxsmith table -input mydoc.docx -output mydoc.docx -set "0,0,1,Age"
```

### Find and Replace

```bash
# Find text
./docxsmith find -input mydoc.docx -text "Hello"

# Replace text globally
./docxsmith replace -input mydoc.docx -output mydoc.docx -old "Hello" -new "Hi"

# Replace in specific paragraph
./docxsmith replace -input mydoc.docx -output mydoc.docx -old "text" -new "content" -paragraph 2
```

### Extract and View Info

```bash
# Extract all text
./docxsmith extract -input mydoc.docx

# Extract to file
./docxsmith extract -input mydoc.docx -output extracted.txt

# View document information
./docxsmith info -input mydoc.docx
```

## PDF Operations

### Create a New PDF

```bash
./docxsmith pdf-create -output myfile.pdf -text "Hello PDF World" -title "My First PDF" -author "John Doe"
```

### Add Content to PDF

```bash
# Add regular text
./docxsmith pdf-add -input myfile.pdf -output myfile.pdf -text "New paragraph"

# Add bold text
./docxsmith pdf-add -input myfile.pdf -output myfile.pdf -text "Important!" -bold

# Add colored text with custom size
./docxsmith pdf-add -input myfile.pdf -output myfile.pdf -text "Red Alert" -color FF0000 -size 16

# Combine styling
./docxsmith pdf-add -input myfile.pdf -output myfile.pdf -text "Notice" -bold -italic -size 14
```

### Extract PDF Content

```bash
# Extract text to console
./docxsmith pdf-extract -input document.pdf

# Extract text to file
./docxsmith pdf-extract -input document.pdf -output content.txt
```

### View PDF Information

```bash
./docxsmith pdf-info -input document.pdf
```

## Format Conversion

### Convert DOCX to PDF

```bash
# Basic conversion
./docxsmith convert -input document.docx -output document.pdf

# With custom options
./docxsmith convert -input document.docx -output document.pdf \
  -font-size 12 \
  -font-family "Times" \
  -page-size A4
```

### Convert PDF to DOCX

```bash
# Basic conversion
./docxsmith convert -input document.pdf -output document.docx

# With custom options
./docxsmith convert -input document.pdf -output document.docx \
  -font-size 11 \
  -font-family "Arial"
```

## Library Usage

### Working with DOCX in Code

```go
package main

import (
    "log"
    "github.com/Palaciodiego008/docxsmith/pkg/docx"
)

func main() {
    // Create document
    doc := docx.New()

    // Add content
    doc.AddParagraph("Title", docx.WithBold(), docx.WithSize("32"))
    doc.AddParagraph("Subtitle", docx.WithItalic(), docx.WithColor("666666"))
    doc.AddParagraph("Body text...")

    // Add table
    table := doc.AddTable(3, 2)
    table.SetCellText(0, 0, "Header 1")
    table.SetCellText(0, 1, "Header 2")

    // Save
    if err := doc.Save("output.docx"); err != nil {
        log.Fatal(err)
    }
}
```

### Working with PDF in Code

```go
package main

import (
    "log"
    "github.com/Palaciodiego008/docxsmith/pkg/pdf"
)

func main() {
    // Create PDF
    pdfDoc := pdf.New()
    pdfDoc.SetMetadata("My Document", "Author", "Subject")

    // Add page and content
    page := pdfDoc.AddPage()

    // Add title
    titleStyle := pdf.TextStyle{
        FontSize:   18,
        FontFamily: "Arial",
        Bold:       true,
        Color:      "000000",
    }
    page.AddTextStyled("Document Title", 20, 30, titleStyle)

    // Add body text
    page.AddText("This is body text", 20, 50, 12)

    // Save
    if err := pdfDoc.Save("output.pdf"); err != nil {
        log.Fatal(err)
    }
}
```

### Converting Between Formats

```go
package main

import (
    "log"
    "github.com/Palaciodiego008/docxsmith/pkg/converter"
)

func main() {
    // Set conversion options
    opts := converter.DefaultOptions()
    opts.FontSize = 12
    opts.FontFamily = "Arial"
    opts.PageSize = "A4"

    // Convert DOCX to PDF
    err := converter.ConvertDocxToPDF("input.docx", "output.pdf", opts)
    if err != nil {
        log.Fatal(err)
    }

    // Convert PDF to DOCX
    err = converter.ConvertPDFToDocx("input.pdf", "output.docx", opts)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Advanced Examples

### Batch Processing

```bash
#!/bin/bash
# Convert all DOCX files in a directory to PDF

for file in *.docx; do
    output="${file%.docx}.pdf"
    ./docxsmith convert -input "$file" -output "$output"
    echo "Converted: $file -> $output"
done
```

### Document Pipeline

```bash
#!/bin/bash
# Create a document, modify it, and convert to PDF

# Create initial document
./docxsmith create -output doc.docx -text "Report"

# Add sections
./docxsmith add -input doc.docx -output doc.docx -text "Introduction" -bold -size 20
./docxsmith add -input doc.docx -output doc.docx -text "This is the intro paragraph."

./docxsmith add -input doc.docx -output doc.docx -text "Methodology" -bold -size 20
./docxsmith add -input doc.docx -output doc.docx -text "We used various methods..."

# Add a table
./docxsmith table -input doc.docx -output doc.docx -create -rows 3 -cols 2

# Convert to PDF
./docxsmith convert -input doc.docx -output final-report.pdf

echo "Report generated: final-report.pdf"
```

## Tips and Best Practices

1. **Always specify output paths** - Most commands require both input and output files
2. **Test with small documents first** - Verify your commands work before processing large files
3. **Use meaningful filenames** - Help track your document versions
4. **Check document info** - Use `info` and `pdf-info` to verify document properties
5. **Backup important files** - Keep originals before making modifications
6. **Chain operations carefully** - Test each step in your pipeline individually
7. **Use version control** - Track your document changes with git

## Troubleshooting

### Common Issues

**Build fails:**
```bash
# Update dependencies
go mod tidy
go build ./cmd/docxsmith
```

**Cannot open file:**
- Check file permissions
- Verify file path is correct
- Ensure file exists

**Conversion quality issues:**
- Try adjusting font size and family options
- Check that source document is well-formatted
- Verify PDF/DOCX is not corrupted

## Getting Help

```bash
# General help
./docxsmith help

# Command-specific help
./docxsmith create -help
./docxsmith pdf-create -help
./docxsmith convert -help
```

## Resources

- [README.md](./README.md) - Full documentation
- [CHANGELOG.md](./CHANGELOG.md) - Version history
- [CONTRIBUTING.md](./CONTRIBUTING.md) - How to contribute
- [Examples](./examples/) - Code examples
