package template

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

// Template represents a document template
type Template struct {
	doc      *docx.Document
	filePath string
}

// Data represents template data
type Data map[string]interface{}

// RenderOptions holds rendering options
type RenderOptions struct {
	// StrictMode causes rendering to fail on missing variables
	StrictMode bool

	// DefaultValue is used when variable is missing (if not in strict mode)
	DefaultValue string

	// RemoveEmptyParagraphs removes paragraphs that become empty after rendering
	RemoveEmptyParagraphs bool
}

// DefaultOptions returns default rendering options
func DefaultOptions() RenderOptions {
	return RenderOptions{
		StrictMode:            false,
		DefaultValue:          "",
		RemoveEmptyParagraphs: true,
	}
}

// New creates a new template from a document
func New(doc *docx.Document) *Template {
	return &Template{
		doc: doc,
	}
}

// Load loads a template from a file
func Load(filePath string) (*Template, error) {
	doc, err := docx.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template: %w", err)
	}

	return &Template{
		doc:      doc,
		filePath: filePath,
	}, nil
}

// Render renders the template with the given data
func (t *Template) Render(data Data, opts RenderOptions) (*docx.Document, error) {
	// Clone the document to avoid modifying the original
	renderedDoc := t.doc.Clone()

	// Process all paragraphs
	for i := 0; i < len(renderedDoc.Body.Paragraphs); i++ {
		para := &renderedDoc.Body.Paragraphs[i]

		// Extract text from paragraph
		text := extractParagraphText(para)

		// Check for loop directive
		if strings.Contains(text, "{{range") {
			// Handle loop
			loopResult, consumed, err := t.processLoop(renderedDoc, i, data, opts)
			if err != nil {
				return nil, fmt.Errorf("error processing loop at paragraph %d: %w", i, err)
			}

			// Replace the loop paragraphs
			if consumed > 0 {
				// Remove original loop paragraphs
				newParas := append(renderedDoc.Body.Paragraphs[:i], renderedDoc.Body.Paragraphs[i+consumed:]...)
				// Insert rendered paragraphs
				renderedDoc.Body.Paragraphs = append(newParas[:i], append(loopResult, newParas[i:]...)...)
				i += len(loopResult) - 1
			}
			continue
		}

		// Check for conditional directive
		if strings.Contains(text, "{{if") {
			// Handle conditional
			condResult, consumed, err := t.processConditional(renderedDoc, i, data, opts)
			if err != nil {
				return nil, fmt.Errorf("error processing conditional at paragraph %d: %w", i, err)
			}

			if consumed > 0 {
				// Replace the conditional paragraphs
				newParas := append(renderedDoc.Body.Paragraphs[:i], renderedDoc.Body.Paragraphs[i+consumed:]...)
				if condResult != nil {
					renderedDoc.Body.Paragraphs = append(newParas[:i], append(condResult, newParas[i:]...)...)
					i += len(condResult) - 1
				} else {
					renderedDoc.Body.Paragraphs = newParas
					i--
				}
			}
			continue
		}

		// Replace variables in paragraph
		if err := t.replaceParagraphVariables(para, data, opts); err != nil {
			if opts.StrictMode {
				return nil, fmt.Errorf("error replacing variables in paragraph %d: %w", i, err)
			}
		}

		// Remove if empty and option is set
		if opts.RemoveEmptyParagraphs && isParagraphEmpty(para) {
			renderedDoc.Body.Paragraphs = append(
				renderedDoc.Body.Paragraphs[:i],
				renderedDoc.Body.Paragraphs[i+1:]...,
			)
			i--
		}
	}

	// Process tables
	for _, table := range renderedDoc.Body.Tables {
		if err := t.processTable(&table, data, opts); err != nil {
			return nil, fmt.Errorf("error processing table: %w", err)
		}
	}

	return renderedDoc, nil
}

// replaceParagraphVariables replaces variables in a paragraph
func (t *Template) replaceParagraphVariables(para *docx.Paragraph, data Data, opts RenderOptions) error {
	varPattern := regexp.MustCompile(`\{\{\.([a-zA-Z0-9_]+)\}\}`)

	for i := range para.Runs {
		for j := range para.Runs[i].Text {
			text := &para.Runs[i].Text[j]

			// Find all variables
			matches := varPattern.FindAllStringSubmatch(text.Content, -1)

			for _, match := range matches {
				if len(match) < 2 {
					continue
				}

				varName := match[1]
				placeholder := match[0]

				// Get value from data
				value, err := getValueFromData(data, varName)
				if err != nil {
					if opts.StrictMode {
						return fmt.Errorf("variable %s not found", varName)
					}
					value = opts.DefaultValue
				}

				// Replace in text
				text.Content = strings.ReplaceAll(text.Content, placeholder, fmt.Sprint(value))
			}
		}
	}

	return nil
}

// getValueFromData retrieves a value from the data map
func getValueFromData(data Data, key string) (interface{}, error) {
	// Support nested keys with dot notation
	keys := strings.Split(key, ".")

	var current interface{} = data
	for _, k := range keys {
		switch v := current.(type) {
		case map[string]interface{}:
			if val, ok := v[k]; ok {
				current = val
			} else {
				return nil, fmt.Errorf("key %s not found", k)
			}
		case Data:
			if val, ok := v[k]; ok {
				current = val
			} else {
				return nil, fmt.Errorf("key %s not found", k)
			}
		default:
			// Try reflection for struct fields
			rv := reflect.ValueOf(current)
			if rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
			if rv.Kind() == reflect.Struct {
				field := rv.FieldByName(k)
				if field.IsValid() {
					current = field.Interface()
				} else {
					return nil, fmt.Errorf("field %s not found", k)
				}
			} else {
				return nil, fmt.Errorf("cannot access key %s on type %T", k, current)
			}
		}
	}

	return current, nil
}

// extractParagraphText extracts all text from a paragraph
func extractParagraphText(para *docx.Paragraph) string {
	var text string
	for _, run := range para.Runs {
		for _, t := range run.Text {
			text += t.Content
		}
	}
	return text
}

// isParagraphEmpty checks if a paragraph is empty
func isParagraphEmpty(para *docx.Paragraph) bool {
	text := extractParagraphText(para)
	return strings.TrimSpace(text) == ""
}

// RenderToFile renders the template and saves to a file
func (t *Template) RenderToFile(data Data, outputPath string, opts RenderOptions) error {
	doc, err := t.Render(data, opts)
	if err != nil {
		return err
	}

	return doc.Save(outputPath)
}

// GetVariables returns all variables found in the template
func (t *Template) GetVariables() []string {
	varPattern := regexp.MustCompile(`\{\{\.([a-zA-Z0-9_]+)\}\}`)
	varSet := make(map[string]bool)

	// Check paragraphs
	for _, para := range t.doc.Body.Paragraphs {
		text := extractParagraphText(&para)
		matches := varPattern.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				varSet[match[1]] = true
			}
		}
	}

	// Check tables
	for _, table := range t.doc.Body.Tables {
		for _, row := range table.Rows {
			for _, cell := range row.Cells {
				for _, para := range cell.Content {
					text := extractParagraphText(&para)
					matches := varPattern.FindAllStringSubmatch(text, -1)
					for _, match := range matches {
						if len(match) >= 2 {
							varSet[match[1]] = true
						}
					}
				}
			}
		}
	}

	// Convert to slice
	variables := make([]string, 0, len(varSet))
	for v := range varSet {
		variables = append(variables, v)
	}

	return variables
}

// ParseInt safely parses an integer
func ParseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

// ParseBool safely parses a boolean
func ParseBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1" || s == "yes"
}
