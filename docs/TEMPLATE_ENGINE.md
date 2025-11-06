# Template Engine Guide

DocxSmith includes a powerful template engine for generating documents dynamically from data.

## Quick Start

```bash
# 1. Create example template and data
docxsmith template-example -template invoice.docx -data data.json

# 2. Render the template
docxsmith template-render -template invoice.docx -data data.json -output result.docx

# 3. List variables in template
docxsmith template-variables -template invoice.docx
```

## Template Syntax

### 1. Variables

Replace placeholders with actual values from your data.

**Template:**
```
Customer: {{.CustomerName}}
Date: {{.Date}}
Total: ${{.Total}}
```

**Data (JSON):**
```json
{
  "CustomerName": "John Doe",
  "Date": "2025-11-05",
  "Total": "1250.00"
}
```

**Result:**
```
Customer: John Doe
Date: 2025-11-05
Total: $1250.00
```

### 2. Conditionals

Show/hide content based on conditions.

**Template:**
```
{{if .IsPaid}}
Status: PAID ✓
{{else}}
Status: UNPAID - Payment Required
{{end}}
```

**Data:**
```json
{
  "IsPaid": true
}
```

**Result:**
```
Status: PAID ✓
```

### Condition Evaluation

- `true` → true
- `false` → false
- Non-empty string → true
- Empty string "" → false
- Non-zero number → true
- Zero → false
- nil → false

### 3. Loops

Iterate over lists to generate repeated content.

**Template:**
```
{{range .Items}}
- {{.Item.Name}}: ${{.Item.Price}} (Qty: {{.Item.Quantity}})
{{end}}
```

**Data:**
```json
{
  "Items": [
    {"Name": "Product A", "Price": "500", "Quantity": "2"},
    {"Name": "Product B", "Price": "250", "Quantity": "1"}
  ]
}
```

**Result:**
```
- Product A: $500 (Qty: 2)
- Product B: $250 (Qty: 1)
```

### 4. Tables with Loops

Generate table rows dynamically from data.

**Template:**
```
Table with 3 rows:
Row 1: {{range .Items}}
Row 2: {{.Item.Name}} | {{.Item.Quantity}} | {{.Item.Price}}
Row 3: [remaining header or data rows]
```

The first row contains `{{range .Items}}`, the second row is the template that gets repeated for each item.

**Data:**
```json
{
  "Items": [
    {"Name": "Item 1", "Quantity": "2", "Price": "$10"},
    {"Name": "Item 2", "Quantity": "5", "Price": "$25"}
  ]
}
```

**Result:**
A table with rows for each item.

## Complete Example

### Invoice Template

**invoice_template.docx contains:**
```
{{.Title}}
Company: {{.CompanyName}}
Date: {{.Date}}

Bill To: {{.CustomerName}}

Items:
[Table with 3 rows]
  {{range .Items}}
  {{.Item.Name}} | {{.Item.Quantity}} | {{.Item.Price}}
  [headers]

Total: ${{.Total}}

{{if .IsPaid}}
✓ PAID - Thank you!
{{else}}
⚠ UNPAID - Payment Due
{{end}}

{{if .Notes}}
Notes: {{.Notes}}
{{end}}
```

### Data File

**invoice_data.json:**
```json
{
  "Title": "INVOICE #001",
  "CompanyName": "ACME Corporation",
  "Date": "2025-11-05",
  "CustomerName": "Jane Smith",
  "Items": [
    {
      "Name": "Consulting Services",
      "Quantity": "10 hours",
      "Price": "$150/hr"
    },
    {
      "Name": "Software License",
      "Quantity": "1",
      "Price": "$500"
    }
  ],
  "Total": "2,000.00",
  "IsPaid": false,
  "Notes": "Payment due within 30 days"
}
```

### Render Command

```bash
docxsmith template-render \
  -template invoice_template.docx \
  -data invoice_data.json \
  -output invoice_final.docx
```

## CLI Commands

### template-render

Render a template with data.

```bash
docxsmith template-render [options]
```

**Options:**
- `-template` - Template file path (required)
- `-data` - Data file (JSON or YAML) (required)
- `-output` - Output file path (required)
- `-strict` - Strict mode: fail on missing variables
- `-default` - Default value for missing variables
- `-keep-empty` - Keep empty paragraphs

**Examples:**
```bash
# Basic rendering
docxsmith template-render -template invoice.docx -data data.json -output result.docx

# Strict mode (fails on missing variables)
docxsmith template-render -template invoice.docx -data data.json -output result.docx -strict

# With default values
docxsmith template-render -template invoice.docx -data data.json -output result.docx -default "N/A"

# Using YAML data
docxsmith template-render -template report.docx -data data.yaml -output result.docx
```

### template-variables

List all variables in a template.

```bash
docxsmith template-variables -template invoice.docx
```

**Output:**
```
Variables found in template (6):
  - Title
  - CompanyName
  - Date
  - CustomerName
  - Total
  - IsPaid
```

### template-example

Create example template and data files.

```bash
docxsmith template-example [options]
```

**Options:**
- `-template` - Output template file (default: "template.docx")
- `-data` - Output data file (default: "data.json")
- `-format` - Data format: json or yaml (default: "json")

**Example:**
```bash
docxsmith template-example -template invoice.docx -data invoice.json
```

## Library Usage

### Basic Usage

```go
package main

import (
    "log"
    "github.com/Palaciodiego008/docxsmith/pkg/template"
)

func main() {
    // Load template
    tmpl, err := template.Load("invoice_template.docx")
    if err != nil {
        log.Fatal(err)
    }

    // Prepare data
    data := template.Data{
        "CustomerName": "John Doe",
        "Total": "1250.00",
        "IsPaid": true,
    }

    // Render
    opts := template.DefaultOptions()
    err = tmpl.RenderToFile(data, "output.docx", opts)
    if err != nil {
        log.Fatal(err)
    }
}
```

### With Custom Options

```go
opts := template.RenderOptions{
    StrictMode:            true,  // Fail on missing variables
    DefaultValue:          "N/A", // Default for missing vars
    RemoveEmptyParagraphs: true,  // Clean up empty paragraphs
}

doc, err := tmpl.Render(data, opts)
```

### Get Template Variables

```go
tmpl, _ := template.Load("invoice.docx")
variables := tmpl.GetVariables()

fmt.Println("Variables:", variables)
// Output: Variables: [CustomerName Total IsPaid Date]
```

### Complex Data Structures

```go
type Invoice struct {
    Number   string
    Customer string
    Items    []LineItem
    Total    float64
    Paid     bool
}

type LineItem struct {
    Description string
    Quantity    int
    Price       float64
}

invoice := Invoice{
    Number:   "INV-001",
    Customer: "ACME Corp",
    Items: []LineItem{
        {Description: "Service A", Quantity: 2, Price: 100.00},
        {Description: "Service B", Quantity: 1, Price: 250.00},
    },
    Total: 450.00,
    Paid:  false,
}

// Convert struct to template.Data
data := template.Data{
    "Number":   invoice.Number,
    "Customer": invoice.Customer,
    "Items":    invoice.Items,
    "Total":    invoice.Total,
    "Paid":     invoice.Paid,
}

tmpl.RenderToFile(data, "invoice.docx", template.DefaultOptions())
```

## Best Practices

### 1. Template Design

- **Keep it simple**: Start with basic variables, add complexity as needed
- **Test incrementally**: Render after each addition to catch errors early
- **Use meaningful names**: `{{.CustomerName}}` not `{{.CN}}`
- **Document your templates**: Add comments in separate sections

### 2. Data Preparation

- **Validate data**: Ensure all required fields are present
- **Format values**: Format numbers, dates before passing to template
- **Handle nulls**: Provide defaults for optional fields
- **Type consistency**: Keep data types consistent (all strings or typed)

### 3. Error Handling

- **Use strict mode** during development to catch missing variables
- **Provide defaults** in production for graceful degradation
- **Test with edge cases**: Empty lists, null values, special characters

### 4. Performance

- **Reuse templates**: Load once, render multiple times
- **Batch processing**: Render multiple documents in parallel
- **Monitor size**: Large templates with many loops can be slow

## Common Patterns

### Invoice/Receipt
```
Title, Company Info
Customer Details
{{range .Items}} - Line items
Total, Tax calculations
{{if .Paid}} - Payment status
```

### Report
```
Report Title and Date
{{range .Sections}} - Multiple sections
  Section Title
  {{range .Section.Data}} - Section data
```

### Letter/Email
```
Date, Recipient
Dear {{.Name}},

Body with {{.Variables}}

{{if .Urgent}}
URGENT: ...
{{end}}

Signature
```

### Certificate
```
Certificate of {{.Type}}
Awarded to: {{.RecipientName}}
Date: {{.Date}}
{{if .Honors}} - With Honors
```

## Troubleshooting

### Variable Not Found

**Error:** `variable CustomerName not found`

**Solutions:**
- Check data has the exact key name (case-sensitive)
- Use `-default` flag to provide fallback value
- Enable strict mode during development

### Loop Not Rendering

**Problem:** `{{range .Items}}` produces no output

**Solutions:**
- Verify data contains "Items" key
- Check Items is a slice/array, not empty
- Ensure {{end}} directive is present

### Conditional Always Shows Same Branch

**Problem:** `{{if .IsPaid}}` always shows one branch

**Solutions:**
- Check boolean value is actually bool type
- Remember: empty string and 0 are falsy
- Test condition evaluation separately

### Table Loop Issues

**Problem:** Table rows not generated correctly

**Solutions:**
- First row: `{{range .Items}}`
- Second row: template with `{{.Item.Field}}`
- Ensure template row has correct number of columns

## Advanced Features

### Nested Data Access

```go
data := template.Data{
    "Company": map[string]interface{}{
        "Name": "ACME",
        "Address": "123 Main St",
    },
}

// Template: {{.Company.Name}}
```

### Loop Index

```
{{range .Items}}
Item {{.Index}}: {{.Item.Name}}
{{end}}
```

### Multiple Conditions

```
{{if .IsUrgent}}
URGENT
{{end}}

{{if .IsConfidential}}
CONFIDENTIAL
{{end}}
```

## Examples

See the `/examples` directory for complete working examples:
- `invoice_template.docx` - Invoice template
- `invoice_data.json` - Sample invoice data
- `invoice_rendered.docx` - Rendered result

## Support

For issues or questions:
- Check the [README](../README.md)
- See [USAGE](../USAGE.md) for general CLI usage
- Open an issue on GitHub
