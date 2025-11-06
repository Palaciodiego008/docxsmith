# Document Merge & Split Guide

DocxSmith provides powerful capabilities for merging and splitting both DOCX and PDF documents.

## Quick Start

```bash
# Merge documents
docxsmith merge -inputs doc1.docx,doc2.docx,doc3.docx -output combined.docx

# Split into 3 parts
docxsmith split -input large.pdf -count 3

# Split by page ranges
docxsmith split -input doc.pdf -pages "1-5,10-15,20-25"

# Split by headings (smart)
docxsmith split -input book.docx -by-heading -heading-level 1
```

## Merge Operations

### Merge DOCX Documents

Combine multiple DOCX files into a single document.

```bash
docxsmith merge \
  -inputs report1.docx,report2.docx,report3.docx \
  -output combined_report.docx
```

**Options:**
- `-inputs` - Comma-separated list of input files (required)
- `-output` - Output file path (required)
- `-page-breaks` - Add page breaks between documents (default: true)
- `-separator` - Add separator text between documents (default: false)
- `-separator-text` - Custom separator text (default: "---")

**Examples:**
```bash
# Basic merge
docxsmith merge -inputs doc1.docx,doc2.docx -output merged.docx

# Merge without page breaks
docxsmith merge -inputs doc1.docx,doc2.docx -output merged.docx -page-breaks=false

# Merge with separator
docxsmith merge -inputs doc1.docx,doc2.docx -output merged.docx -separator -separator-text "=== NEW SECTION ==="
```

### Merge PDF Documents

Combine multiple PDF files into a single PDF.

```bash
docxsmith merge -inputs file1.pdf,file2.pdf,file3.pdf -output combined.pdf
```

### Merge Information

Preview what will be merged without actually merging.

```bash
docxsmith merge-info -inputs doc1.docx,doc2.docx,doc3.docx
```

**Output:**
```
Merge Information (DOCX):
  Documents: 3
  Total Paragraphs: 45
  Total Tables: 2
```

## Split Operations

### Split by Page Ranges (PDF)

Extract specific pages or ranges from a PDF.

```bash
# Extract pages 1-5
docxsmith split -input document.pdf -pages "1-5"

# Extract multiple ranges
docxsmith split -input document.pdf -pages "1-5,10-15,20-25"

# Extract specific pages
docxsmith split -input document.pdf -pages "1,5,10,15"

# Mixed ranges and pages
docxsmith split -input document.pdf -pages "1-3,5,7-9,12"
```

**Output:**
- Creates separate PDF files for each range
- Files named: `part1.pdf`, `part2.pdf`, etc.

### Split into Equal Parts

Divide a document into N equal parts.

```bash
# Split DOCX into 3 parts
docxsmith split -input large.docx -count 3

# Split PDF into 5 parts
docxsmith split -input large.pdf -count 5

# Custom output pattern
docxsmith split -input report.pdf -count 4 -pattern "section{n}.pdf"
```

### Split by Headings (Smart Split - DOCX only)

Automatically split a DOCX document at each heading.

```bash
# Split at Heading 1 (chapters)
docxsmith split -input book.docx -by-heading -heading-level 1

# Split at Heading 2 (sections)
docxsmith split -input book.docx -by-heading -heading-level 2

# Custom output pattern with title
docxsmith split -input book.docx -by-heading -heading-level 1 -pattern "{title}.docx"
```

**Features:**
- Automatically detects heading styles
- Creates one file per heading
- Can use heading text in filename
- Preserves all content between headings

**Heading Levels:**
- Level 1: Main chapters
- Level 2: Sections
- Level 3-6: Subsections

### Custom Output Patterns

Control how output files are named:

```bash
# Using placeholders
docxsmith split -input doc.pdf -count 3 -pattern "chapter_{n}.pdf"
# Output: chapter_1.pdf, chapter_2.pdf, chapter_3.pdf

docxsmith split -input report.docx -count 2 -pattern "{base}_part{n}.docx"
# Output: report_part1.docx, report_part2.docx

# For heading splits
docxsmith split -input book.docx -by-heading -pattern "{title}.docx"
# Output: Introduction.docx, Chapter 1.docx, Conclusion.docx
```

**Placeholders:**
- `{n}` - Part number (1, 2, 3, ...)
- `{base}` - Original filename without extension
- `{title}` - Heading text (heading split only)

### Output Directory

Specify where to save split files:

```bash
docxsmith split -input large.pdf -count 5 -dir output/chapters/
```

## Library Usage

### Merging Documents

```go
package main

import (
    "log"
    "github.com/Palaciodiego008/docxsmith/pkg/operations"
)

func main() {
    // Merge DOCX files
    inputs := []string{"doc1.docx", "doc2.docx", "doc3.docx"}
    opts := operations.DefaultMergeOptions()
    opts.AddPageBreaks = true
    opts.AddSeparator = true
    opts.SeparatorText = "=== SECTION ==="

    err := operations.MergeDOCX(inputs, "combined.docx", opts)
    if err != nil {
        log.Fatal(err)
    }

    // Merge PDF files
    pdfInputs := []string{"file1.pdf", "file2.pdf"}
    err = operations.MergePDF(pdfInputs, "combined.pdf")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Splitting Documents

```go
// Split by paragraph ranges
ranges := []operations.ParagraphRange{
    {Start: 0, End: 10},
    {Start: 11, End: 20},
    {Start: 21, End: 30},
}

opts := operations.DefaultSplitOptions()
opts.OutputPattern = "chapter{n}.docx"

files, err := operations.SplitDOCXByParagraphs("book.docx", ranges, opts)

// Split into equal parts
files, err := operations.SplitDOCXByCount("large.docx", 5, opts)

// Split by headings
files, err := operations.SplitDOCXByHeadings("book.docx", 1, opts)

// Split PDF by pages
pageRanges := []operations.PageRange{
    {Start: 0, End: 9},    // Pages 1-10
    {Start: 10, End: 19},  // Pages 11-20
}
files, err := operations.SplitPDFByPages("doc.pdf", pageRanges, opts)
```

### Getting Merge Information

```go
// Get DOCX merge info
info, err := operations.GetMergeDOCXInfo([]string{"doc1.docx", "doc2.docx"})
fmt.Printf("Will merge %d documents with %d total paragraphs\n",
    info.TotalDocuments, info.TotalParagraphs)

// Get PDF merge info
info, err := operations.GetMergePDFInfo([]string{"file1.pdf", "file2.pdf"})
fmt.Printf("Will merge %d PDFs with %d total pages\n",
    info.TotalDocuments, info.TotalPages)
```

## Use Cases

### 1. Combine Monthly Reports

```bash
# Merge all monthly reports into annual report
docxsmith merge \
  -inputs jan.docx,feb.docx,mar.docx,apr.docx,may.docx,jun.docx,jul.docx,aug.docx,sep.docx,oct.docx,nov.docx,dec.docx \
  -output annual_report_2025.docx \
  -page-breaks \
  -separator \
  -separator-text "=== END OF MONTH ==="
```

### 2. Split Book into Chapters

```bash
# Split book by Heading 1 (chapters)
docxsmith split \
  -input complete_book.docx \
  -by-heading \
  -heading-level 1 \
  -pattern "{title}.docx" \
  -dir chapters/

# Output:
# chapters/Introduction.docx
# chapters/Chapter 1: Getting Started.docx
# chapters/Chapter 2: Advanced Topics.docx
# chapters/Conclusion.docx
```

### 3. Extract PDF Pages

```bash
# Extract executive summary (pages 1-3)
docxsmith split -input report.pdf -pages "1-3" -pattern "executive_summary.pdf"

# Extract multiple sections
docxsmith split \
  -input manual.pdf \
  -pages "1-10,20-30,50-60" \
  -pattern "section{n}.pdf"
```

### 4. Split Large Document for Email

```bash
# Split 100-page document into 4 parts
docxsmith split -input large_report.pdf -count 4 -pattern "report_part{n}.pdf"
```

### 5. Batch Process with Shell Script

```bash
#!/bin/bash
# Merge all DOCX files in current directory

files=$(ls *.docx | tr '\n' ',' | sed 's/,$//')
docxsmith merge -inputs "$files" -output all_combined.docx

echo "Merged $(ls *.docx | wc -l) documents into all_combined.docx"
```

## Advanced Examples

### Conditional Merge

```bash
#!/bin/bash
# Only merge files containing specific keyword

merged_files=""
for file in *.docx; do
    if docxsmith find -input "$file" -text "APPROVED" > /dev/null 2>&1; then
        merged_files="$merged_files,$file"
    fi
done

merged_files="${merged_files#,}"  # Remove leading comma
docxsmith merge -inputs "$merged_files" -output approved_docs.docx
```

### Split and Convert Pipeline

```bash
#!/bin/bash
# Split DOCX by chapters and convert each to PDF

docxsmith split -input book.docx -by-heading -heading-level 1 -dir chapters/

for file in chapters/*.docx; do
    pdf_file="${file%.docx}.pdf"
    docxsmith convert -input "$file" -output "$pdf_file"
done

echo "Split and converted book into individual chapter PDFs"
```

## Best Practices

### Merging

1. **Check compatibility**: Ensure documents have similar formatting
2. **Preview first**: Use `merge-info` to see what will be merged
3. **Backup originals**: Keep original files before merging
4. **Use separators**: For clarity when merging different sources
5. **Page breaks**: Essential for professional multi-document merges

### Splitting

1. **Plan ranges**: Know what content goes where
2. **Use meaningful patterns**: Name output files descriptively
3. **Heading splits**: Best for structured documents (books, reports)
4. **Test with small docs**: Verify split logic before processing large files
5. **Organize output**: Use `-dir` to keep split files organized

## Troubleshooting

### Merge Issues

**Problem:** Merged document has formatting issues

**Solution:**
- Check source documents have consistent styles
- Try merging smaller batches first
- Use `-separator` to clearly delineate sources

**Problem:** Page breaks not appearing

**Solution:**
- Ensure `-page-breaks` flag is set to true
- Some viewers may not show page breaks clearly

### Split Issues

**Problem:** `split -by-heading` produces no files

**Solution:**
- Verify document has proper heading styles (Heading1, Heading2, etc.)
- Check heading level matches document structure
- Try different heading levels (1-6)

**Problem:** Page range out of bounds

**Solution:**
- Check document page count first with `info` or `pdf-info`
- Remember pages are 1-indexed in commands (page 1 is first page)
- Verify range format: "1-5" not "0-4"

## Performance Tips

1. **Large merges**: Process in batches if merging 100+ documents
2. **PDF splits**: Page-based splits are faster than content-based
3. **DOCX splits**: Heading-based splits are more reliable than count-based
4. **Parallel processing**: Use shell scripts to process multiple operations concurrently

## Examples Directory

Run the examples to see merge & split in action:

```bash
cd examples
go run merge_split_examples.go
```

## API Reference

### Merge Functions

```go
// Merge DOCX files
func MergeDOCX(inputPaths []string, outputPath string, opts MergeOptions) error

// Merge PDF files
func MergePDF(inputPaths []string, outputPath string) error

// Auto-detect and merge
func MergeDocuments(inputPaths []string, outputPath string, opts MergeOptions) error

// Get merge information
func GetMergeDOCXInfo(inputPaths []string) (*MergeInfo, error)
func GetMergePDFInfo(inputPaths []string) (*MergeInfo, error)
```

### Split Functions

```go
// Split DOCX by paragraph ranges
func SplitDOCXByParagraphs(inputPath string, ranges []ParagraphRange, opts SplitOptions) ([]string, error)

// Split PDF by page ranges
func SplitPDFByPages(inputPath string, ranges []PageRange, opts SplitOptions) ([]string, error)

// Split into N equal parts
func SplitDOCXByCount(inputPath string, count int, opts SplitOptions) ([]string, error)
func SplitPDFByCount(inputPath string, count int, opts SplitOptions) ([]string, error)

// Smart split by headings (DOCX only)
func SplitDOCXByHeadings(inputPath string, headingLevel int, opts SplitOptions) ([]string, error)

// Parse page range string
func ParsePageRanges(rangeStr string, maxPages int) ([]PageRange, error)
```

## Supported Formats

| Operation | DOCX | PDF |
|-----------|------|-----|
| Merge     | ✅   | ✅  |
| Split by Count | ✅ | ✅ |
| Split by Pages | ❌ | ✅ |
| Split by Headings | ✅ | ❌ |

## Resources

- [README](../README.md) - Main documentation
- [USAGE](../USAGE.md) - General usage guide
- [Examples](../examples/) - Code examples
