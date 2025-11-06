package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

func TestNew(t *testing.T) {
	doc := docx.New()
	tmpl := New(doc)

	if tmpl == nil {
		t.Fatal("New() returned nil")
	}
	if tmpl.doc == nil {
		t.Fatal("Template has nil document")
	}
}

func TestSimpleVariableReplacement(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     Data
		expected string
	}{
		{
			name:     "Single variable",
			template: "Hello {{.Name}}",
			data:     Data{"Name": "World"},
			expected: "Hello World",
		},
		{
			name:     "Multiple variables",
			template: "{{.FirstName}} {{.LastName}}",
			data:     Data{"FirstName": "John", "LastName": "Doe"},
			expected: "John Doe",
		},
		{
			name:     "Numeric variable",
			template: "Price: ${{.Price}}",
			data:     Data{"Price": 99.99},
			expected: "Price: $99.99",
		},
		{
			name:     "Variable with surrounding text",
			template: "Welcome {{.User}} to our system!",
			data:     Data{"User": "Alice"},
			expected: "Welcome Alice to our system!",
		},
		{
			name:     "Empty variable",
			template: "Value: {{.Value}}",
			data:     Data{"Value": ""},
			expected: "Value: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := docx.New()
			doc.AddParagraph(tt.template)

			tmpl := New(doc)
			opts := DefaultOptions()

			result, err := tmpl.Render(tt.data, opts)
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			if len(result.Body.Paragraphs) == 0 {
				t.Fatal("No paragraphs in result")
			}

			text := extractParagraphText(&result.Body.Paragraphs[0])
			if text != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, text)
			}
		})
	}
}

func TestMissingVariables(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		data        Data
		strictMode  bool
		defaultVal  string
		expectError bool
		expected    string
	}{
		{
			name:        "Missing var in strict mode",
			template:    "Hello {{.Name}}",
			data:        Data{},
			strictMode:  true,
			expectError: true,
		},
		{
			name:        "Missing var with default",
			template:    "Hello {{.Name}}",
			data:        Data{},
			strictMode:  false,
			defaultVal:  "Guest",
			expectError: false,
			expected:    "Hello Guest",
		},
		{
			name:        "Missing var with empty default",
			template:    "Hello {{.Name}}",
			data:        Data{},
			strictMode:  false,
			defaultVal:  "",
			expectError: false,
			expected:    "Hello ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := docx.New()
			doc.AddParagraph(tt.template)

			tmpl := New(doc)
			opts := RenderOptions{
				StrictMode:            tt.strictMode,
				DefaultValue:          tt.defaultVal,
				RemoveEmptyParagraphs: false,
			}

			result, err := tmpl.Render(tt.data, opts)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(result.Body.Paragraphs) == 0 {
				t.Fatal("No paragraphs in result")
			}

			text := extractParagraphText(&result.Body.Paragraphs[0])
			if text != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, text)
			}
		})
	}
}

func TestConditionals(t *testing.T) {
	tests := []struct {
		name          string
		conditionVar  string
		conditionVal  interface{}
		trueContent   string
		falseContent  string
		expectedCount int
		expectedText  string
	}{
		{
			name:          "True boolean condition",
			conditionVar:  "IsActive",
			conditionVal:  true,
			trueContent:   "Active",
			falseContent:  "Inactive",
			expectedCount: 1,
			expectedText:  "Active",
		},
		{
			name:          "False boolean condition",
			conditionVar:  "IsActive",
			conditionVal:  false,
			trueContent:   "Active",
			falseContent:  "Inactive",
			expectedCount: 1,
			expectedText:  "Inactive",
		},
		{
			name:          "Non-empty string is true",
			conditionVar:  "Status",
			conditionVal:  "approved",
			trueContent:   "Approved",
			falseContent:  "Pending",
			expectedCount: 1,
			expectedText:  "Approved",
		},
		{
			name:          "Empty string is false",
			conditionVar:  "Status",
			conditionVal:  "",
			trueContent:   "Approved",
			falseContent:  "Pending",
			expectedCount: 1,
			expectedText:  "Pending",
		},
		{
			name:          "Non-zero number is true",
			conditionVar:  "Count",
			conditionVal:  5,
			trueContent:   "Has items",
			falseContent:  "No items",
			expectedCount: 1,
			expectedText:  "Has items",
		},
		{
			name:          "Zero number is false",
			conditionVar:  "Count",
			conditionVal:  0,
			trueContent:   "Has items",
			falseContent:  "No items",
			expectedCount: 1,
			expectedText:  "No items",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := docx.New()
			doc.AddParagraph("{{if ." + tt.conditionVar + "}}")
			doc.AddParagraph(tt.trueContent)
			doc.AddParagraph("{{else}}")
			doc.AddParagraph(tt.falseContent)
			doc.AddParagraph("{{end}}")

			data := Data{tt.conditionVar: tt.conditionVal}

			tmpl := New(doc)
			opts := DefaultOptions()

			result, err := tmpl.Render(data, opts)
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			if len(result.Body.Paragraphs) != tt.expectedCount {
				t.Errorf("Expected %d paragraphs, got %d", tt.expectedCount, len(result.Body.Paragraphs))
			}

			if len(result.Body.Paragraphs) > 0 {
				text := extractParagraphText(&result.Body.Paragraphs[0])
				if text != tt.expectedText {
					t.Errorf("Expected '%s', got '%s'", tt.expectedText, text)
				}
			}
		})
	}
}

func TestLoops(t *testing.T) {
	tests := []struct {
		name          string
		items         []map[string]interface{}
		expectedCount int
	}{
		{
			name: "Loop with 3 items",
			items: []map[string]interface{}{
				{"Name": "Item1", "Price": "10"},
				{"Name": "Item2", "Price": "20"},
				{"Name": "Item3", "Price": "30"},
			},
			expectedCount: 3,
		},
		{
			name:          "Loop with empty list",
			items:         []map[string]interface{}{},
			expectedCount: 0,
		},
		{
			name: "Loop with one item",
			items: []map[string]interface{}{
				{"Name": "Solo", "Price": "100"},
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := docx.New()
			doc.AddParagraph("{{range .Items}}")
			doc.AddParagraph("{{.Item.Name}}: ${{.Item.Price}}")
			doc.AddParagraph("{{end}}")

			data := Data{"Items": tt.items}

			tmpl := New(doc)
			opts := DefaultOptions()

			result, err := tmpl.Render(data, opts)
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			if len(result.Body.Paragraphs) != tt.expectedCount {
				t.Errorf("Expected %d paragraphs, got %d", tt.expectedCount, len(result.Body.Paragraphs))
			}

			// Verify content if there are items
			if tt.expectedCount > 0 && len(tt.items) > 0 {
				text := extractParagraphText(&result.Body.Paragraphs[0])
				expectedText := tt.items[0]["Name"].(string) + ": $" + tt.items[0]["Price"].(string)
				if text != expectedText {
					t.Errorf("Expected '%s', got '%s'", expectedText, text)
				}
			}
		})
	}
}

func TestGetVariables(t *testing.T) {
	tests := []struct {
		name             string
		paragraphs       []string
		expectedVarCount int
		expectedVars     []string
	}{
		{
			name:             "Single variable",
			paragraphs:       []string{"Hello {{.Name}}"},
			expectedVarCount: 1,
			expectedVars:     []string{"Name"},
		},
		{
			name:             "Multiple unique variables",
			paragraphs:       []string{"{{.First}} {{.Last}} {{.Email}}"},
			expectedVarCount: 3,
			expectedVars:     []string{"First", "Last", "Email"},
		},
		{
			name:             "Duplicate variables",
			paragraphs:       []string{"{{.Name}}", "Hello {{.Name}} again"},
			expectedVarCount: 1,
			expectedVars:     []string{"Name"},
		},
		{
			name:             "No variables",
			paragraphs:       []string{"Plain text"},
			expectedVarCount: 0,
			expectedVars:     []string{},
		},
		{
			name: "Mixed content",
			paragraphs: []string{
				"Title: {{.Title}}",
				"Plain text",
				"Author: {{.Author}}",
			},
			expectedVarCount: 2,
			expectedVars:     []string{"Title", "Author"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := docx.New()
			for _, p := range tt.paragraphs {
				doc.AddParagraph(p)
			}

			tmpl := New(doc)
			vars := tmpl.GetVariables()

			if len(vars) != tt.expectedVarCount {
				t.Errorf("Expected %d variables, got %d", tt.expectedVarCount, len(vars))
			}

			// Check that expected vars are present
			varMap := make(map[string]bool)
			for _, v := range vars {
				varMap[v] = true
			}

			for _, expected := range tt.expectedVars {
				if !varMap[expected] {
					t.Errorf("Expected variable '%s' not found", expected)
				}
			}
		})
	}
}

func TestLoadAndRenderToFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create template file
	templatePath := filepath.Join(tmpDir, "template.docx")
	doc := docx.New()
	doc.AddParagraph("Invoice for {{.Customer}}")
	doc.AddParagraph("Amount: ${{.Amount}}")

	err := doc.Save(templatePath)
	if err != nil {
		t.Fatalf("Failed to save template: %v", err)
	}

	// Load template
	tmpl, err := Load(templatePath)
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	// Render to file
	outputPath := filepath.Join(tmpDir, "output.docx")
	data := Data{
		"Customer": "John Doe",
		"Amount":   "1250.00",
	}

	err = tmpl.RenderToFile(data, outputPath, DefaultOptions())
	if err != nil {
		t.Fatalf("Failed to render to file: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Open and verify content
	resultDoc, err := docx.Open(outputPath)
	if err != nil {
		t.Fatalf("Failed to open result: %v", err)
	}

	if len(resultDoc.Body.Paragraphs) < 2 {
		t.Fatal("Expected at least 2 paragraphs in result")
	}

	text0 := extractParagraphText(&resultDoc.Body.Paragraphs[0])
	text1 := extractParagraphText(&resultDoc.Body.Paragraphs[1])

	if text0 != "Invoice for John Doe" {
		t.Errorf("Expected 'Invoice for John Doe', got '%s'", text0)
	}

	if text1 != "Amount: $1250.00" {
		t.Errorf("Expected 'Amount: $1250.00', got '%s'", text1)
	}
}

func TestComplexScenario(t *testing.T) {
	// Test combining variables, conditionals, and loops
	doc := docx.New()

	doc.AddParagraph("{{.CompanyName}}")
	doc.AddParagraph("Invoice for {{.CustomerName}}")
	doc.AddParagraph("")

	doc.AddParagraph("{{range .Items}}")
	doc.AddParagraph("- {{.Item.Name}}: ${{.Item.Price}}")
	doc.AddParagraph("{{end}}")

	doc.AddParagraph("")
	doc.AddParagraph("Total: ${{.Total}}")

	doc.AddParagraph("{{if .IsPaid}}")
	doc.AddParagraph("Status: PAID")
	doc.AddParagraph("{{else}}")
	doc.AddParagraph("Status: UNPAID")
	doc.AddParagraph("{{end}}")

	data := Data{
		"CompanyName":  "ACME Corp",
		"CustomerName": "Jane Smith",
		"Items": []map[string]interface{}{
			{"Name": "Product A", "Price": "100"},
			{"Name": "Product B", "Price": "200"},
		},
		"Total":  "300",
		"IsPaid": true,
	}

	tmpl := New(doc)
	opts := DefaultOptions()

	result, err := tmpl.Render(data, opts)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Verify structure
	expectedTexts := []string{
		"ACME Corp",
		"Invoice for Jane Smith",
		"- Product A: $100",
		"- Product B: $200",
		"Total: $300",
		"Status: PAID",
	}

	actualCount := 0
	for _, para := range result.Body.Paragraphs {
		text := extractParagraphText(&para)
		if text != "" {
			actualCount++
		}
	}

	if actualCount < len(expectedTexts) {
		t.Errorf("Expected at least %d non-empty paragraphs, got %d", len(expectedTexts), actualCount)
	}
}
