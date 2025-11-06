package docx

import (
	"testing"
)

func TestAddTable(t *testing.T) {
	doc := New()
	table := doc.AddTable(3, 4)

	if table == nil {
		t.Fatal("AddTable returned nil")
	}

	if doc.GetTableCount() != 1 {
		t.Errorf("Expected 1 table, got %d", doc.GetTableCount())
	}

	if table.GetRowCount() != 3 {
		t.Errorf("Expected 3 rows, got %d", table.GetRowCount())
	}

	if table.GetColumnCount() != 4 {
		t.Errorf("Expected 4 columns, got %d", table.GetColumnCount())
	}
}

func TestSetCellText(t *testing.T) {
	doc := New()
	table := doc.AddTable(2, 2)

	err := table.SetCellText(0, 0, "Hello")
	if err != nil {
		t.Fatalf("Error setting cell text: %v", err)
	}

	text, err := table.GetCellText(0, 0)
	if err != nil {
		t.Fatalf("Error getting cell text: %v", err)
	}

	if text != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", text)
	}
}

func TestSetCellTextOutOfRange(t *testing.T) {
	doc := New()
	table := doc.AddTable(2, 2)

	err := table.SetCellText(5, 0, "Test")
	if err == nil {
		t.Error("Expected error for out of range row, got nil")
	}

	err = table.SetCellText(0, 5, "Test")
	if err == nil {
		t.Error("Expected error for out of range column, got nil")
	}
}

func TestGetCellText(t *testing.T) {
	doc := New()
	table := doc.AddTable(2, 2)

	table.SetCellText(0, 0, "A1")
	table.SetCellText(0, 1, "B1")
	table.SetCellText(1, 0, "A2")
	table.SetCellText(1, 1, "B2")

	tests := []struct {
		row, col int
		expected string
	}{
		{0, 0, "A1"},
		{0, 1, "B1"},
		{1, 0, "A2"},
		{1, 1, "B2"},
	}

	for _, tt := range tests {
		text, err := table.GetCellText(tt.row, tt.col)
		if err != nil {
			t.Fatalf("Error getting cell [%d,%d]: %v", tt.row, tt.col, err)
		}
		if text != tt.expected {
			t.Errorf("Cell [%d,%d]: expected '%s', got '%s'", tt.row, tt.col, tt.expected, text)
		}
	}
}

func TestAddRow(t *testing.T) {
	doc := New()
	table := doc.AddTable(2, 3)

	initialRows := table.GetRowCount()
	table.AddRow()

	if table.GetRowCount() != initialRows+1 {
		t.Errorf("Expected %d rows, got %d", initialRows+1, table.GetRowCount())
	}

	// Verify new row has correct number of columns
	newRowIdx := table.GetRowCount() - 1
	if len(table.Rows[newRowIdx].Cells) != 3 {
		t.Errorf("New row should have 3 columns, got %d", len(table.Rows[newRowIdx].Cells))
	}
}

func TestDeleteRow(t *testing.T) {
	doc := New()
	table := doc.AddTable(3, 2)

	table.SetCellText(0, 0, "Row 0")
	table.SetCellText(1, 0, "Row 1")
	table.SetCellText(2, 0, "Row 2")

	err := table.DeleteRow(1)
	if err != nil {
		t.Fatalf("Error deleting row: %v", err)
	}

	if table.GetRowCount() != 2 {
		t.Errorf("Expected 2 rows, got %d", table.GetRowCount())
	}

	text0, _ := table.GetCellText(0, 0)
	text1, _ := table.GetCellText(1, 0)

	if text0 != "Row 0" || text1 != "Row 2" {
		t.Errorf("Wrong rows after deletion: got '%s' and '%s'", text0, text1)
	}
}

func TestDeleteRowOutOfRange(t *testing.T) {
	doc := New()
	table := doc.AddTable(2, 2)

	err := table.DeleteRow(5)
	if err == nil {
		t.Error("Expected error for out of range row, got nil")
	}
}

func TestDeleteTable(t *testing.T) {
	doc := New()
	doc.AddTable(2, 2)
	doc.AddTable(3, 3)

	if doc.GetTableCount() != 2 {
		t.Fatalf("Expected 2 tables, got %d", doc.GetTableCount())
	}

	err := doc.DeleteTable(0)
	if err != nil {
		t.Fatalf("Error deleting table: %v", err)
	}

	if doc.GetTableCount() != 1 {
		t.Errorf("Expected 1 table after deletion, got %d", doc.GetTableCount())
	}

	// Verify remaining table is the second one
	if doc.Body.Tables[0].GetRowCount() != 3 {
		t.Error("Wrong table remained after deletion")
	}
}

func TestTableWithParagraphsAndTables(t *testing.T) {
	doc := New()
	doc.AddParagraph("Header")
	doc.AddTable(2, 2)
	doc.AddParagraph("Middle")
	doc.AddTable(3, 3)
	doc.AddParagraph("Footer")

	if doc.GetParagraphCount() != 3 {
		t.Errorf("Expected 3 paragraphs, got %d", doc.GetParagraphCount())
	}

	if doc.GetTableCount() != 2 {
		t.Errorf("Expected 2 tables, got %d", doc.GetTableCount())
	}
}

func TestEmptyTableOperations(t *testing.T) {
	doc := New()
	table := doc.AddTable(0, 0)

	if table.GetRowCount() != 0 {
		t.Errorf("Expected 0 rows, got %d", table.GetRowCount())
	}

	if table.GetColumnCount() != 0 {
		t.Errorf("Expected 0 columns, got %d", table.GetColumnCount())
	}

	// AddRow on empty table should do nothing
	table.AddRow()
	if table.GetRowCount() != 0 {
		t.Error("AddRow on empty table should not add rows")
	}
}
