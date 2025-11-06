## Document Diff - Professional Comparison Tool

DocxSmith includes a sophisticated document comparison engine for tracking changes between document versions.

## Quick Start

```bash
# Compare two documents (HTML output)
./docxsmith diff -old version1.docx -new version2.docx -output changes.html

# Compare with markdown output
./docxsmith diff -old v1.docx -new v2.docx -format markdown -output changes.md

# Compare to terminal
./docxsmith diff -old v1.docx -new v2.docx -format text

# Ignore whitespace differences
./docxsmith diff -old v1.docx -new v2.docx -ignore-whitespace

# Ignore case differences
./docxsmith diff -old v1.docx -new v2.docx -ignore-case
```

## CLI Usage

### Basic Comparison

```bash
docxsmith diff [options]
```

**Required Flags:**
- `-old` - Path to old/original document
- `-new` - Path to new/modified document

**Optional Flags:**
- `-output` - Output file (default: stdout)
- `-format` - Output format: html, markdown, text (default: html)
- `-ignore-whitespace` - Ignore whitespace differences
- `-ignore-case` - Ignore case differences
- `-stats` - Show statistics (default: true)

### Output Formats

#### 1. HTML (Default)

Professional HTML output with colors and styling.

```bash
docxsmith diff -old v1.docx -new v2.docx -output report.html
```

**Features:**
- Color-coded changes (green=added, red=deleted, yellow=modified)
- Statistics dashboard
- Clean, professional styling
- Line numbers
- Legend
- Responsive design

#### 2. Markdown

GitHub-flavored markdown for documentation.

```bash
docxsmith diff -old v1.docx -new v2.docx -format markdown -output changes.md
```

**Features:**
- Markdown formatting
- Strikethrough for deletions
- Inline indicators (`+`, `-`, `~`)
- Statistics table
- Git-friendly

#### 3. Plain Text

Simple text output for terminals.

```bash
docxsmith diff -old v1.docx -new v2.docx -format text
```

**Features:**
- Clean terminal output
- Simple prefix indicators
- No special formatting
- Easy to parse

## Library Usage

### Basic Comparison

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/Palaciodiego008/docxsmith/pkg/diff"
)

func main() {
    // Configure options
    opts := diff.DefaultDiffOptions()
    opts.IgnoreWhitespace = true

    // Compare documents
    result, err := diff.CompareDOCX("old.docx", "new.docx", opts)
    if err != nil {
        log.Fatal(err)
    }

    // Render as HTML
    renderer := diff.NewHTMLRenderer(true)
    output, err := renderer.Render(result)
    if err != nil {
        log.Fatal(err)
    }

    // Save to file
    os.WriteFile("diff.html", []byte(output), 0644)

    // Print summary
    fmt.Printf("Total changes: %d\n", result.Stats.TotalChanges)
    fmt.Printf("Added lines: %d\n", result.Stats.AddedLines)
    fmt.Printf("Deleted lines: %d\n", result.Stats.DeletedLines)
}
```

### Custom Rendering

```go
// HTML with custom settings
htmlRenderer := diff.NewHTMLRenderer(true) // true = show stats
html, _ := htmlRenderer.Render(result)

// Markdown
mdRenderer := diff.NewMarkdownRenderer(true)
md, _ := mdRenderer.Render(result)

// Plain text
txtRenderer := diff.NewPlainTextRenderer(true, true) // showStats, colorOutput
txt, _ := txtRenderer.Render(result)
```

### Advanced Options

```go
opts := diff.DiffOptions{
    IgnoreWhitespace: true,  // Ignore whitespace
    IgnoreCase:       true,  // Case-insensitive
    ContextLines:     3,     // Lines of context
    MinChangeLength:  1,     // Minimum change length
}

result, err := diff.CompareDOCX("v1.docx", "v2.docx", opts)
```

### Accessing Change Details

```go
for _, change := range result.Changes {
    switch change.Type {
    case diff.DiffAdded:
        fmt.Printf("Added at line %d: %s\n", change.Position, change.New)
    case diff.DiffDeleted:
        fmt.Printf("Deleted at line %d: %s\n", change.Position, change.Old)
    case diff.DiffModified:
        fmt.Printf("Modified at line %d: %s -> %s\n",
            change.Position, change.Old, change.New)
    }
}
```

## Use Cases

### 1. Version Control

Track document changes between versions:

```bash
# Compare contract versions
docxsmith diff -old contract_v1.docx -new contract_v2.docx -output contract_changes.html

# Email changes to team
# Open contract_changes.html in browser and share
```

### 2. Review Process

Review edits made by others:

```bash
# Compare before and after edits
docxsmith diff -old draft.docx -new reviewed.docx -format markdown -output review.md

# Add to pull request or review system
```

### 3. Quality Assurance

Verify document modifications:

```bash
# Ensure only approved changes were made
docxsmith diff -old approved.docx -new final.docx -output qa_report.html
```

### 4. Audit Trail

Create audit trails for document changes:

```bash
#!/bin/bash
# Create audit trail for all versions

for i in {1..5}; do
    if [ $i -eq 1 ]; then
        continue
    fi

    prev=$((i-1))
    docxsmith diff \
        -old "version_$prev.docx" \
        -new "version_$i.docx" \
        -output "audit/changes_v${prev}_to_v${i}.html"
done
```

### 5. Automated Reporting

Generate change reports automatically:

```bash
#!/bin/bash
# Daily document diff report

DATE=$(date +%Y-%m-%d)
docxsmith diff \
    -old "archive/doc_yesterday.docx" \
    -new "current/doc_today.docx" \
    -output "reports/changes_$DATE.html" \
    -ignore-whitespace

echo "Change report generated: reports/changes_$DATE.html"
```

## Architecture

### Clean Separation of Concerns

```
pkg/diff/
├── diff.go          # Core diff algorithm (Myers LCS)
├── renderer.go      # Rendering interfaces & implementations
└── diff_test.go     # Table-driven tests
```

### Design Patterns Used

1. **Strategy Pattern** - Multiple renderers (HTML, Markdown, Text)
2. **Interface Segregation** - Clean Renderer interface
3. **Dependency Injection** - Options passed to constructors
4. **Factory Pattern** - Renderer creation
5. **Single Responsibility** - Each file has one purpose

### Key Components

**Differ Interface:**
```go
type Differ interface {
    Compare(old, new string) (*DiffResult, error)
}
```

**Renderer Interface:**
```go
type Renderer interface {
    Render(result *DiffResult) (string, error)
}
```

**Implementations:**
- `DocxDiffer` - Compares DOCX documents
- `HTMLRenderer` - Renders as HTML
- `MarkdownRenderer` - Renders as Markdown
- `PlainTextRenderer` - Renders as plain text

## Algorithm

DocxSmith uses a **simplified Myers diff algorithm** based on Longest Common Subsequence (LCS):

1. Extract text lines from both documents
2. Build dynamic programming table for LCS
3. Backtrack to identify changes
4. Classify changes (added, deleted, modified)
5. Calculate statistics
6. Render output

**Time Complexity:** O(n*m) where n, m are document lengths
**Space Complexity:** O(n*m) for DP table

## Best Practices

### 1. Choose the Right Format

- **HTML** - For sharing with non-technical users, presentations
- **Markdown** - For documentation, Git repos, technical reviews
- **Text** - For quick terminal checks, automation scripts

### 2. Use Ignore Options Wisely

```bash
# For code/technical documents
docxsmith diff -old v1.docx -new v2.docx

# For general text documents (ignore formatting)
docxsmith diff -old v1.docx -new v2.docx -ignore-whitespace -ignore-case
```

### 3. Automate Regular Comparisons

```bash
# Add to CI/CD pipeline
if ! docxsmith diff -old master_doc.docx -new branch_doc.docx -format text | grep "identical"; then
    echo "Document has changes - review required"
    exit 1
fi
```

### 4. Combine with Other Features

```bash
# Compare after template rendering
docxsmith template-render -template tmpl.docx -data old.json -output old_result.docx
docxsmith template-render -template tmpl.docx -data new.json -output new_result.docx
docxsmith diff -old old_result.docx -new new_result.docx -output template_changes.html
```

## Troubleshooting

### Large Documents

For very large documents:
- Use `-format text` for faster output
- Disable stats with `-stats=false`
- Consider splitting documents first

### No Changes Detected

If you expect changes but see none:
- Check if files are actually different
- Remove `-ignore-whitespace` and `-ignore-case` flags
- Verify correct file paths

### HTML Not Rendering Properly

- Ensure output file has `.html` extension
- Open in modern browser (Chrome, Firefox, Safari)
- Check file permissions

## Advanced Examples

### Compare Multiple Versions

```bash
#!/bin/bash
# Compare across multiple versions

versions=("v1" "v2" "v3" "v4")

for ((i=0; i<${#versions[@]}-1; i++)); do
    curr="${versions[$i]}"
    next="${versions[$i+1]}"

    docxsmith diff \
        -old "docs/${curr}.docx" \
        -new "docs/${next}.docx" \
        -output "diffs/${curr}_to_${next}.html"

    echo "Compared $curr → $next"
done
```

### Integration with Git

```bash
#!/bin/bash
# Git hook to show document changes

OLD_VERSION=$(git show HEAD:document.docx > old_temp.docx)
docxsmith diff -old old_temp.docx -new document.docx -output git_changes.html
rm old_temp.docx
```

## API Reference

### Types

```go
type DiffType int
const (
    DiffNone
    DiffAdded
    DiffDeleted
    DiffModified
)

type Change struct {
    Type     DiffType
    Old      string
    New      string
    Position int
    Context  string
}

type DiffResult struct {
    Changes      []Change
    Stats        DiffStats
    OldDocument  string
    NewDocument  string
}

type DiffStats struct {
    TotalChanges   int
    AddedLines     int
    DeletedLines   int
    ModifiedLines  int
    UnchangedLines int
}
```

### Functions

```go
// Compare two DOCX documents
func CompareDOCX(oldPath, newPath string, opts DiffOptions) (*DiffResult, error)

// Create renderers
func NewHTMLRenderer(showStats bool) *HTMLRenderer
func NewMarkdownRenderer(showStats bool) *MarkdownRenderer
func NewPlainTextRenderer(showStats, colorOutput bool) *PlainTextRenderer

// Render diff
func (r Renderer) Render(result *DiffResult) (string, error)
```

## Testing

Run the comprehensive test suite:

```bash
# Run diff tests
go test ./pkg/diff -v

# With coverage
go test ./pkg/diff -cover
```

**Test Coverage:**
- 7 test suites
- 20+ test cases (table-driven)
- All diff scenarios covered
- All renderers tested

## Performance

- **Small documents** (<100 paragraphs): < 100ms
- **Medium documents** (100-1000 paragraphs): < 500ms
- **Large documents** (1000+ paragraphs): < 2s

## Future Enhancements

Planned improvements:
- Word-level diff (not just line-level)
- Visual side-by-side comparison
- Track changes mode (modify DOCX with changes marked)
- PDF diff support
- Image diff detection

## Resources

- [Main Documentation](../README.md)
- [Usage Guide](../USAGE.md)
- [Architecture Overview](./ARCHITECTURE.md)
