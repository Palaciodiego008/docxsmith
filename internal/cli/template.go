package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
	"github.com/Palaciodiego008/docxsmith/pkg/template"
	"gopkg.in/yaml.v3"
)

// HandleTemplateRender handles the template render command
func HandleTemplateRender(args []string) {
	fs := flag.NewFlagSet("template-render", flag.ExitOnError)
	templatePath := fs.String("template", "", "Template file path (required)")
	dataPath := fs.String("data", "", "Data file path (JSON or YAML) (required)")
	output := fs.String("output", "", "Output file path (required)")
	strict := fs.Bool("strict", false, "Strict mode - fail on missing variables")
	defaultVal := fs.String("default", "", "Default value for missing variables")
	keepEmpty := fs.Bool("keep-empty", false, "Keep empty paragraphs")
	fs.Parse(args)

	if *templatePath == "" || *dataPath == "" || *output == "" {
		fmt.Fprintln(os.Stderr, "Error: -template, -data, and -output are required")
		fs.Usage()
		os.Exit(1)
	}

	// Load template
	tmpl, err := template.Load(*templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading template: %v\n", err)
		os.Exit(1)
	}

	// Load data
	data, err := loadDataFile(*dataPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading data: %v\n", err)
		os.Exit(1)
	}

	// Configure options
	opts := template.RenderOptions{
		StrictMode:            *strict,
		DefaultValue:          *defaultVal,
		RemoveEmptyParagraphs: !*keepEmpty,
	}

	// Render
	err = tmpl.RenderToFile(data, *output, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Template rendered successfully: %s\n", *output)
}

// HandleTemplateVariables handles the template-variables command
func HandleTemplateVariables(args []string) {
	fs := flag.NewFlagSet("template-variables", flag.ExitOnError)
	templatePath := fs.String("template", "", "Template file path (required)")
	fs.Parse(args)

	if *templatePath == "" {
		fmt.Fprintln(os.Stderr, "Error: -template is required")
		fs.Usage()
		os.Exit(1)
	}

	// Load template
	tmpl, err := template.Load(*templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading template: %v\n", err)
		os.Exit(1)
	}

	// Get variables
	variables := tmpl.GetVariables()

	if len(variables) == 0 {
		fmt.Println("No variables found in template")
		return
	}

	fmt.Printf("Variables found in template (%d):\n", len(variables))
	for _, v := range variables {
		fmt.Printf("  - %s\n", v)
	}
}

// HandleTemplateExample handles the template-example command
func HandleTemplateExample(args []string) {
	fs := flag.NewFlagSet("template-example", flag.ExitOnError)
	outputTemplate := fs.String("template", "template.docx", "Output template file")
	outputData := fs.String("data", "data.json", "Output data file")
	format := fs.String("format", "json", "Data format (json or yaml)")
	fs.Parse(args)

	// Create example template
	fmt.Println("Creating example template...")

	doc, err := createExampleTemplate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating template: %v\n", err)
		os.Exit(1)
	}

	if err := doc.Save(*outputTemplate); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Template created: %s\n", *outputTemplate)

	// Create example data
	exampleData := map[string]interface{}{
		"Title":        "Invoice",
		"CompanyName":  "ACME Corp",
		"Date":         "2025-11-05",
		"CustomerName": "John Doe",
		"Total":        "$1,250.00",
		"IsPaid":       false,
		"Items": []map[string]interface{}{
			{"Name": "Product A", "Quantity": "2", "Price": "$500.00"},
			{"Name": "Product B", "Quantity": "1", "Price": "$250.00"},
		},
	}

	var dataBytes []byte
	if *format == "yaml" {
		dataBytes, err = yaml.Marshal(exampleData)
	} else {
		dataBytes, err = json.MarshalIndent(exampleData, "", "  ")
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating data: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*outputData, dataBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Data created: %s\n", *outputData)
	fmt.Println("\nTo render the template, run:")
	fmt.Printf("  docxsmith template-render -template %s -data %s -output result.docx\n", *outputTemplate, *outputData)
}

// loadDataFile loads data from JSON or YAML file
func loadDataFile(path string) (template.Data, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	var result template.Data

	// Try JSON first
	err = json.Unmarshal(data, &result)
	if err == nil {
		return result, nil
	}

	// Try YAML
	err = yaml.Unmarshal(data, &result)
	if err == nil {
		return result, nil
	}

	return nil, fmt.Errorf("failed to parse data file as JSON or YAML")
}

// createExampleTemplate creates an example template document
func createExampleTemplate() (*docx.Document, error) {
	doc := docx.New()

	// Title
	doc.AddParagraph("{{.Title}}", docx.WithBold(), docx.WithSize("32"), docx.WithAlignment("center"))
	doc.AddParagraph("")

	// Company info
	doc.AddParagraph("{{.CompanyName}}", docx.WithBold())
	doc.AddParagraph("Date: {{.Date}}")
	doc.AddParagraph("")

	// Customer
	doc.AddParagraph("Customer: {{.CustomerName}}")
	doc.AddParagraph("")

	// Items table with loop
	doc.AddParagraph("Items:")
	table := doc.AddTable(3, 3)

	// Table headers in first row
	table.SetCellText(0, 0, "{{range .Items}}")
	table.SetCellText(0, 1, "")
	table.SetCellText(0, 2, "")

	// Template row
	table.SetCellText(1, 0, "{{.Item.Name}}")
	table.SetCellText(1, 1, "{{.Item.Quantity}}")
	table.SetCellText(1, 2, "{{.Item.Price}}")

	// Additional rows will be generated from data
	table.SetCellText(2, 0, "Name")
	table.SetCellText(2, 1, "Qty")
	table.SetCellText(2, 2, "Price")

	doc.AddParagraph("")

	// Total
	doc.AddParagraph("Total: {{.Total}}", docx.WithBold(), docx.WithSize("24"))
	doc.AddParagraph("")

	// Conditional
	doc.AddParagraph("{{if .IsPaid}}")
	doc.AddParagraph("Status: PAID", docx.WithColor("00FF00"), docx.WithBold())
	doc.AddParagraph("{{else}}")
	doc.AddParagraph("Status: UNPAID", docx.WithColor("FF0000"), docx.WithBold())
	doc.AddParagraph("{{end}}")

	return doc, nil
}
