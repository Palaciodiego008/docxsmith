package template

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

// processLoop processes a {{range .Items}}...{{end}} directive
func (t *Template) processLoop(doc *docx.Document, startIdx int, data Data, opts RenderOptions) ([]docx.Paragraph, int, error) {
	result := []docx.Paragraph{}

	// Find the range directive
	startText := extractParagraphText(&doc.Body.Paragraphs[startIdx])
	rangePattern := regexp.MustCompile(`\{\{range\s+\.([a-zA-Z0-9_]+)\}\}`)
	matches := rangePattern.FindStringSubmatch(startText)

	if len(matches) < 2 {
		return nil, 0, fmt.Errorf("invalid range directive: %s", startText)
	}

	collectionName := matches[1]

	// Get the collection from data
	collection, err := getValueFromData(data, collectionName)
	if err != nil {
		if opts.StrictMode {
			return nil, 0, fmt.Errorf("collection %s not found", collectionName)
		}
		return result, 1, nil // Return empty result
	}

	// Find the end directive
	endIdx := -1
	for i := startIdx + 1; i < len(doc.Body.Paragraphs); i++ {
		text := extractParagraphText(&doc.Body.Paragraphs[i])
		if strings.Contains(text, "{{end}}") {
			endIdx = i
			break
		}
	}

	if endIdx == -1 {
		return nil, 0, fmt.Errorf("no matching {{end}} found for {{range}}")
	}

	// Get template paragraphs (between start and end)
	templateParas := doc.Body.Paragraphs[startIdx+1 : endIdx]

	// Convert collection to slice
	collectionSlice, err := toSlice(collection)
	if err != nil {
		return nil, 0, fmt.Errorf("collection %s is not iterable: %w", collectionName, err)
	}

	// Iterate over collection
	for idx, item := range collectionSlice {
		// Create data context for this iteration
		itemData := Data{
			"Index": idx,
			"Item":  item,
		}

		// Merge with parent data
		for k, v := range data {
			if k != collectionName {
				itemData[k] = v
			}
		}

		// Render each template paragraph with item data
		for _, templatePara := range templateParas {
			// Clone paragraph
			newPara := cloneParagraph(&templatePara)

			// Replace {{.Item.Field}} with actual values
			if err := t.replaceLoopVariables(&newPara, item, opts); err != nil {
				if opts.StrictMode {
					return nil, 0, err
				}
			}

			// Also replace {{.Index}}
			if err := t.replaceParagraphVariables(&newPara, itemData, opts); err != nil {
				if opts.StrictMode {
					return nil, 0, err
				}
			}

			result = append(result, newPara)
		}
	}

	// Return result and number of paragraphs consumed (start + templates + end)
	consumed := endIdx - startIdx + 1
	return result, consumed, nil
}

// replaceLoopVariables replaces {{.Item.Field}} variables
func (t *Template) replaceLoopVariables(para *docx.Paragraph, item interface{}, opts RenderOptions) error {
	itemPattern := regexp.MustCompile(`\{\{\.Item\.([a-zA-Z0-9_]+)\}\}`)

	for i := range para.Runs {
		for j := range para.Runs[i].Text {
			text := &para.Runs[i].Text[j]

			// Find all item variables
			matches := itemPattern.FindAllStringSubmatch(text.Content, -1)

			for _, match := range matches {
				if len(match) < 2 {
					continue
				}

				fieldName := match[1]
				placeholder := match[0]

				// Get field value from item
				value, err := getFieldValue(item, fieldName)
				if err != nil {
					if opts.StrictMode {
						return fmt.Errorf("field %s not found in item", fieldName)
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

// getFieldValue gets a field value from a struct or map
func getFieldValue(item interface{}, fieldName string) (interface{}, error) {
	// If item is a map
	if m, ok := item.(map[string]interface{}); ok {
		if val, exists := m[fieldName]; exists {
			return val, nil
		}
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	// If item is Data
	if d, ok := item.(Data); ok {
		if val, exists := d[fieldName]; exists {
			return val, nil
		}
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	// Use reflection for structs
	rv := reflect.ValueOf(item)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Struct {
		field := rv.FieldByName(fieldName)
		if field.IsValid() {
			return field.Interface(), nil
		}
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	return nil, fmt.Errorf("item is not a struct or map")
}

// toSlice converts various types to a slice
func toSlice(v interface{}) ([]interface{}, error) {
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		result := make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			result[i] = rv.Index(i).Interface()
		}
		return result, nil

	default:
		return nil, fmt.Errorf("value is not a slice or array")
	}
}

// cloneParagraph creates a deep copy of a paragraph
func cloneParagraph(p *docx.Paragraph) docx.Paragraph {
	newPara := docx.Paragraph{
		Runs: make([]docx.Run, len(p.Runs)),
	}

	// Copy runs
	for i, run := range p.Runs {
		newRun := docx.Run{
			Text: make([]docx.Text, len(run.Text)),
		}

		// Copy text
		copy(newRun.Text, run.Text)

		// Copy properties
		if run.Props != nil {
			newRun.Props = &docx.RProps{
				Bold:   run.Props.Bold,
				Italic: run.Props.Italic,
				Size:   run.Props.Size,
				Color:  run.Props.Color,
			}
		}

		newPara.Runs[i] = newRun
	}

	// Copy properties
	if p.Props != nil {
		newPara.Props = &docx.PProps{
			Style:   p.Props.Style,
			Jc:      p.Props.Jc,
			Spacing: p.Props.Spacing,
		}
	}

	return newPara
}
