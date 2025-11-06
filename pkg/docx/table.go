package docx

import (
	"encoding/xml"
	"fmt"
)

// Table represents a table in the document
type Table struct {
	XMLName xml.Name `xml:"tbl"`
	Props   *TblPr   `xml:"tblPr,omitempty"`
	Grid    *TblGrid `xml:"tblGrid,omitempty"`
	Rows    []TblRow `xml:"tr"`
}

// TblPr represents table properties
type TblPr struct {
	XMLName xml.Name `xml:"tblPr"`
	Style   *TblStyle `xml:"tblStyle,omitempty"`
	Width   *TblWidth `xml:"tblW,omitempty"`
}

// TblStyle represents table style
type TblStyle struct {
	XMLName xml.Name `xml:"tblStyle"`
	Val     string   `xml:"val,attr"`
}

// TblWidth represents table width
type TblWidth struct {
	XMLName xml.Name `xml:"tblW"`
	Type    string   `xml:"type,attr"`
	W       string   `xml:"w,attr"`
}

// TblGrid represents table grid/columns
type TblGrid struct {
	XMLName xml.Name      `xml:"tblGrid"`
	Cols    []TblGridCol  `xml:"gridCol"`
}

// TblGridCol represents a table column
type TblGridCol struct {
	XMLName xml.Name `xml:"gridCol"`
	W       string   `xml:"w,attr,omitempty"`
}

// TblRow represents a table row
type TblRow struct {
	XMLName xml.Name `xml:"tr"`
	Props   *TrPr    `xml:"trPr,omitempty"`
	Cells   []TblCell `xml:"tc"`
}

// TrPr represents row properties
type TrPr struct {
	XMLName xml.Name `xml:"trPr"`
}

// TblCell represents a table cell
type TblCell struct {
	XMLName xml.Name   `xml:"tc"`
	Props   *TcPr      `xml:"tcPr,omitempty"`
	Content []Paragraph `xml:"p"`
}

// TcPr represents cell properties
type TcPr struct {
	XMLName xml.Name `xml:"tcPr"`
	Width   *TblWidth `xml:"tcW,omitempty"`
}

// AddTable adds a new table to the document
func (d *Document) AddTable(rows, cols int) *Table {
	table := Table{
		Props: &TblPr{
			Width: &TblWidth{
				Type: "auto",
				W:    "0",
			},
		},
		Grid: &TblGrid{
			Cols: make([]TblGridCol, cols),
		},
		Rows: make([]TblRow, rows),
	}

	// Initialize rows and cells
	for i := 0; i < rows; i++ {
		table.Rows[i] = TblRow{
			Cells: make([]TblCell, cols),
		}
		for j := 0; j < cols; j++ {
			table.Rows[i].Cells[j] = TblCell{
				Content: []Paragraph{
					{
						Runs: []Run{
							{
								Text: []Text{
									{Space: "preserve", Content: ""},
								},
							},
						},
					},
				},
			}
		}
	}

	d.Body.Tables = append(d.Body.Tables, table)
	return &d.Body.Tables[len(d.Body.Tables)-1]
}

// SetCellText sets the text content of a cell
func (t *Table) SetCellText(row, col int, text string) error {
	if row < 0 || row >= len(t.Rows) {
		return fmt.Errorf("row index %d out of range", row)
	}
	if col < 0 || col >= len(t.Rows[row].Cells) {
		return fmt.Errorf("column index %d out of range", col)
	}

	cell := &t.Rows[row].Cells[col]
	if len(cell.Content) == 0 {
		cell.Content = []Paragraph{{}}
	}
	if len(cell.Content[0].Runs) == 0 {
		cell.Content[0].Runs = []Run{{}}
	}

	cell.Content[0].Runs[0].Text = []Text{
		{
			Space:   "preserve",
			Content: text,
		},
	}

	return nil
}

// GetCellText gets the text content of a cell
func (t *Table) GetCellText(row, col int) (string, error) {
	if row < 0 || row >= len(t.Rows) {
		return "", fmt.Errorf("row index %d out of range", row)
	}
	if col < 0 || col >= len(t.Rows[row].Cells) {
		return "", fmt.Errorf("column index %d out of range", col)
	}

	cell := t.Rows[row].Cells[col]
	var text string
	for _, p := range cell.Content {
		for _, r := range p.Runs {
			for _, t := range r.Text {
				text += t.Content
			}
		}
	}

	return text, nil
}

// AddRow adds a new row to the table
func (t *Table) AddRow() {
	if len(t.Rows) == 0 {
		return
	}

	cols := len(t.Rows[0].Cells)
	newRow := TblRow{
		Cells: make([]TblCell, cols),
	}

	for i := 0; i < cols; i++ {
		newRow.Cells[i] = TblCell{
			Content: []Paragraph{
				{
					Runs: []Run{
						{
							Text: []Text{
								{Space: "preserve", Content: ""},
							},
						},
					},
				},
			},
		}
	}

	t.Rows = append(t.Rows, newRow)
}

// DeleteRow deletes a row from the table
func (t *Table) DeleteRow(index int) error {
	if index < 0 || index >= len(t.Rows) {
		return fmt.Errorf("row index %d out of range", index)
	}

	t.Rows = append(t.Rows[:index], t.Rows[index+1:]...)
	return nil
}

// GetRowCount returns the number of rows in the table
func (t *Table) GetRowCount() int {
	return len(t.Rows)
}

// GetColumnCount returns the number of columns in the table
func (t *Table) GetColumnCount() int {
	if len(t.Rows) == 0 {
		return 0
	}
	return len(t.Rows[0].Cells)
}
