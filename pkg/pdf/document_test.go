package pdf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	doc := New()
	if doc == nil {
		t.Fatal("New() returned nil")
	}
	if doc.Pages == nil {
		t.Fatal("New document has nil pages")
	}
	if len(doc.Pages) != 0 {
		t.Errorf("New document should have 0 pages, got %d", len(doc.Pages))
	}
}

func TestAddPage(t *testing.T) {
	doc := New()
	initialCount := doc.GetPageCount()

	page := doc.AddPage()
	if page == nil {
		t.Fatal("AddPage() returned nil")
	}

	if doc.GetPageCount() != initialCount+1 {
		t.Errorf("Expected %d pages, got %d", initialCount+1, doc.GetPageCount())
	}

	if page.Number != 1 {
		t.Errorf("Expected page number 1, got %d", page.Number)
	}
}

func TestAddText(t *testing.T) {
	doc := New()
	page := doc.AddPage()

	page.AddText("Test text", 20, 30, 12)

	if len(page.Content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(page.Content))
	}

	if tc, ok := page.Content[0].(TextContent); ok {
		if tc.Text != "Test text" {
			t.Errorf("Expected 'Test text', got '%s'", tc.Text)
		}
		if tc.X != 20 || tc.Y != 30 {
			t.Errorf("Expected position (20,30), got (%f,%f)", tc.X, tc.Y)
		}
		if tc.FontSize != 12 {
			t.Errorf("Expected font size 12, got %f", tc.FontSize)
		}
	} else {
		t.Error("Content is not TextContent")
	}
}

func TestGetPage(t *testing.T) {
	doc := New()
	doc.AddPage()
	doc.AddPage()

	page, err := doc.GetPage(0)
	if err != nil {
		t.Fatalf("GetPage(0) failed: %v", err)
	}
	if page.Number != 1 {
		t.Errorf("Expected page number 1, got %d", page.Number)
	}

	page, err = doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1) failed: %v", err)
	}
	if page.Number != 2 {
		t.Errorf("Expected page number 2, got %d", page.Number)
	}

	// Test out of range
	_, err = doc.GetPage(5)
	if err == nil {
		t.Error("Expected error for out of range index")
	}
}

func TestDeletePage(t *testing.T) {
	doc := New()
	doc.AddPage()
	doc.AddPage()
	doc.AddPage()

	if doc.GetPageCount() != 3 {
		t.Fatalf("Expected 3 pages initially, got %d", doc.GetPageCount())
	}

	err := doc.DeletePage(1)
	if err != nil {
		t.Fatalf("DeletePage failed: %v", err)
	}

	if doc.GetPageCount() != 2 {
		t.Errorf("Expected 2 pages after deletion, got %d", doc.GetPageCount())
	}

	// Verify page numbers are updated
	for i, page := range doc.Pages {
		if page.Number != i+1 {
			t.Errorf("Page %d has incorrect number %d", i, page.Number)
		}
	}
}

func TestGetText(t *testing.T) {
	doc := New()
	page := doc.AddPage()

	page.AddText("First text", 20, 30, 12)
	page.AddText("Second text", 20, 50, 12)

	text := page.GetText()
	expected := "First text Second text "
	if text != expected {
		t.Errorf("Expected '%s', got '%s'", expected, text)
	}
}

func TestGetAllText(t *testing.T) {
	doc := New()

	page1 := doc.AddPage()
	page1.AddText("Page 1 text", 20, 30, 12)

	page2 := doc.AddPage()
	page2.AddText("Page 2 text", 20, 30, 12)

	text := doc.GetAllText()
	if !contains(text, "Page 1 text") || !contains(text, "Page 2 text") {
		t.Errorf("GetAllText() doesn't contain expected text: %s", text)
	}
}

func TestSetMetadata(t *testing.T) {
	doc := New()
	doc.SetMetadata("Test Title", "Test Author", "Test Subject")

	if doc.Metadata.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", doc.Metadata.Title)
	}
	if doc.Metadata.Author != "Test Author" {
		t.Errorf("Expected author 'Test Author', got '%s'", doc.Metadata.Author)
	}
	if doc.Metadata.Subject != "Test Subject" {
		t.Errorf("Expected subject 'Test Subject', got '%s'", doc.Metadata.Subject)
	}
}

func TestSaveAndOpen(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.pdf")

	// Create and save
	doc := New()
	doc.SetMetadata("Test Document", "Test Author", "")
	page := doc.AddPage()
	page.AddText("Hello PDF", 20, 30, 12)

	err := doc.Save(testFile)
	if err != nil {
		t.Fatalf("Error saving PDF: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("PDF file was not created")
	}

	// Open and verify
	doc2, err := Open(testFile)
	if err != nil {
		t.Fatalf("Error opening PDF: %v", err)
	}

	if doc2.GetPageCount() != 1 {
		t.Errorf("Expected 1 page, got %d", doc2.GetPageCount())
	}

	text := doc2.GetAllText()
	if !contains(text, "Hello PDF") {
		t.Errorf("Expected text to contain 'Hello PDF', got: %s", text)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
