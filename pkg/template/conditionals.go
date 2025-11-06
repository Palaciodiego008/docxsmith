package template

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

// processConditional processes a {{if .Condition}}...{{end}} directive
func (t *Template) processConditional(doc *docx.Document, startIdx int, data Data, opts RenderOptions) ([]docx.Paragraph, int, error) {
	// Find the if directive
	startText := extractParagraphText(&doc.Body.Paragraphs[startIdx])
	ifPattern := regexp.MustCompile(`\{\{if\s+\.([a-zA-Z0-9_]+)\}\}`)
	matches := ifPattern.FindStringSubmatch(startText)

	if len(matches) < 2 {
		return nil, 0, fmt.Errorf("invalid if directive: %s", startText)
	}

	conditionName := matches[1]

	// Get the condition value
	conditionValue, err := getValueFromData(data, conditionName)
	if err != nil {
		if opts.StrictMode {
			return nil, 0, fmt.Errorf("condition variable %s not found", conditionName)
		}
		conditionValue = false
	}

	// Evaluate condition
	condition := evaluateCondition(conditionValue)

	// Find the end directive
	endIdx := -1
	elseIdx := -1
	for i := startIdx + 1; i < len(doc.Body.Paragraphs); i++ {
		text := extractParagraphText(&doc.Body.Paragraphs[i])
		if strings.Contains(text, "{{else}}") && elseIdx == -1 {
			elseIdx = i
		}
		if strings.Contains(text, "{{end}}") {
			endIdx = i
			break
		}
	}

	if endIdx == -1 {
		return nil, 0, fmt.Errorf("no matching {{end}} found for {{if}}")
	}

	result := []docx.Paragraph{}

	if condition {
		// Include content from startIdx+1 to elseIdx-1 (or endIdx-1 if no else)
		endContentIdx := endIdx
		if elseIdx != -1 {
			endContentIdx = elseIdx
		}

		for i := startIdx + 1; i < endContentIdx; i++ {
			para := cloneParagraph(&doc.Body.Paragraphs[i])
			// Replace variables in the paragraph
			if err := t.replaceParagraphVariables(&para, data, opts); err != nil {
				if opts.StrictMode {
					return nil, 0, err
				}
			}
			result = append(result, para)
		}
	} else if elseIdx != -1 {
		// Include content from elseIdx+1 to endIdx-1
		for i := elseIdx + 1; i < endIdx; i++ {
			para := cloneParagraph(&doc.Body.Paragraphs[i])
			// Replace variables in the paragraph
			if err := t.replaceParagraphVariables(&para, data, opts); err != nil {
				if opts.StrictMode {
					return nil, 0, err
				}
			}
			result = append(result, para)
		}
	}

	// Return result and number of paragraphs consumed
	consumed := endIdx - startIdx + 1
	return result, consumed, nil
}

// evaluateCondition evaluates a condition value to boolean
func evaluateCondition(value interface{}) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		return v != "" && v != "false" && v != "0"
	case int, int8, int16, int32, int64:
		return v != 0
	case uint, uint8, uint16, uint32, uint64:
		return v != 0
	case float32, float64:
		return v != 0.0
	default:
		// For other types, check if non-nil
		return true
	}
}

// processTable processes variables in table cells
func (t *Template) processTable(table *docx.Table, data Data, opts RenderOptions) error {
	// Check if table has range directive in first row
	if len(table.Rows) > 0 {
		firstRowText := ""
		if len(table.Rows[0].Cells) > 0 && len(table.Rows[0].Cells[0].Content) > 0 {
			firstRowText = extractParagraphText(&table.Rows[0].Cells[0].Content[0])
		}

		// Check for range directive
		if strings.Contains(firstRowText, "{{range") {
			return t.processTableLoop(table, data, opts)
		}
	}

	// Regular table - just replace variables in each cell
	for i := range table.Rows {
		for j := range table.Rows[i].Cells {
			for k := range table.Rows[i].Cells[j].Content {
				para := &table.Rows[i].Cells[j].Content[k]
				if err := t.replaceParagraphVariables(para, data, opts); err != nil {
					if opts.StrictMode {
						return err
					}
				}
			}
		}
	}

	return nil
}

// processTableLoop processes a range directive in a table
func (t *Template) processTableLoop(table *docx.Table, data Data, opts RenderOptions) error {
	if len(table.Rows) < 2 {
		return fmt.Errorf("table loop requires at least 2 rows (directive + template)")
	}

	// Parse range directive from first row
	firstRowText := extractParagraphText(&table.Rows[0].Cells[0].Content[0])
	rangePattern := regexp.MustCompile(`\{\{range\s+\.([a-zA-Z0-9_]+)\}\}`)
	matches := rangePattern.FindStringSubmatch(firstRowText)

	if len(matches) < 2 {
		return fmt.Errorf("invalid range directive in table: %s", firstRowText)
	}

	collectionName := matches[1]

	// Get the collection
	collection, err := getValueFromData(data, collectionName)
	if err != nil {
		if opts.StrictMode {
			return fmt.Errorf("collection %s not found", collectionName)
		}
		return nil
	}

	// Convert to slice
	collectionSlice, err := toSlice(collection)
	if err != nil {
		return fmt.Errorf("collection %s is not iterable: %w", collectionName, err)
	}

	// The second row is the template
	templateRow := table.Rows[1]

	// Remove first two rows (directive + template)
	table.Rows = table.Rows[2:]

	// Generate rows for each item
	newRows := []docx.TblRow{}

	for _, item := range collectionSlice {
		// Clone the template row
		newRow := cloneTableRow(&templateRow)

		// Replace variables in each cell
		for i := range newRow.Cells {
			for j := range newRow.Cells[i].Content {
				para := &newRow.Cells[i].Content[j]
				if err := t.replaceLoopVariables(para, item, opts); err != nil {
					if opts.StrictMode {
						return err
					}
				}
			}
		}

		newRows = append(newRows, newRow)
	}

	// Prepend new rows to existing rows
	table.Rows = append(newRows, table.Rows...)

	return nil
}

// cloneTableRow creates a deep copy of a table row
func cloneTableRow(row *docx.TblRow) docx.TblRow {
	newRow := docx.TblRow{
		Cells: make([]docx.TblCell, len(row.Cells)),
	}

	for i, cell := range row.Cells {
		newCell := docx.TblCell{
			Content: make([]docx.Paragraph, len(cell.Content)),
		}

		// Copy cell properties
		if cell.Props != nil {
			newCell.Props = &docx.TcPr{
				Width: cell.Props.Width,
			}
		}

		// Clone each paragraph
		for j, para := range cell.Content {
			newCell.Content[j] = cloneParagraph(&para)
		}

		newRow.Cells[i] = newCell
	}

	// Copy row properties
	if row.Props != nil {
		newRow.Props = &docx.TrPr{}
	}

	return newRow
}
