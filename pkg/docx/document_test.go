package docx

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
	if doc.Body == nil {
		t.Fatal("New document has nil body")
	}
	if len(doc.files) == 0 {
		t.Fatal("New document has no default files")
	}
}

func TestAddParagraph(t *testing.T) {
	doc := New()
	initialCount := doc.GetParagraphCount()

	doc.AddParagraph("Test paragraph")
	if doc.GetParagraphCount() != initialCount+1 {
		t.Errorf("Expected %d paragraphs, got %d", initialCount+1, doc.GetParagraphCount())
	}

	text, err := doc.GetParagraphText(initialCount)
	if err != nil {
		t.Fatalf("Error getting paragraph text: %v", err)
	}
	if text != "Test paragraph" {
		t.Errorf("Expected 'Test paragraph', got '%s'", text)
	}
}

func TestAddParagraphWithOptions(t *testing.T) {
	doc := New()

	doc.AddParagraph("Bold text", WithBold())
	doc.AddParagraph("Italic text", WithItalic())
	doc.AddParagraph("Colored text", WithColor("FF0000"))
	doc.AddParagraph("Large text", WithSize("32"))
	doc.AddParagraph("Centered text", WithAlignment("center"))

	if doc.GetParagraphCount() != 5 {
		t.Errorf("Expected 5 paragraphs, got %d", doc.GetParagraphCount())
	}

	// Verify bold
	if doc.Body.Paragraphs[0].Runs[0].Props == nil || doc.Body.Paragraphs[0].Runs[0].Props.Bold == nil {
		t.Error("Bold option was not applied")
	}

	// Verify italic
	if doc.Body.Paragraphs[1].Runs[0].Props == nil || doc.Body.Paragraphs[1].Runs[0].Props.Italic == nil {
		t.Error("Italic option was not applied")
	}

	// Verify color
	if doc.Body.Paragraphs[2].Runs[0].Props == nil || doc.Body.Paragraphs[2].Runs[0].Props.Color == nil {
		t.Error("Color option was not applied")
	}
	if doc.Body.Paragraphs[2].Runs[0].Props.Color.Val != "FF0000" {
		t.Errorf("Expected color FF0000, got %s", doc.Body.Paragraphs[2].Runs[0].Props.Color.Val)
	}

	// Verify alignment
	if doc.Body.Paragraphs[4].Props == nil || doc.Body.Paragraphs[4].Props.Jc == nil {
		t.Error("Alignment option was not applied")
	}
	if doc.Body.Paragraphs[4].Props.Jc.Val != "center" {
		t.Errorf("Expected alignment center, got %s", doc.Body.Paragraphs[4].Props.Jc.Val)
	}
}

func TestAddParagraphAt(t *testing.T) {
	doc := New()
	doc.AddParagraph("First")
	doc.AddParagraph("Third")

	err := doc.AddParagraphAt(1, "Second")
	if err != nil {
		t.Fatalf("Error adding paragraph at index: %v", err)
	}

	if doc.GetParagraphCount() != 3 {
		t.Errorf("Expected 3 paragraphs, got %d", doc.GetParagraphCount())
	}

	texts := []string{"First", "Second", "Third"}
	for i, expected := range texts {
		text, _ := doc.GetParagraphText(i)
		if text != expected {
			t.Errorf("Paragraph %d: expected '%s', got '%s'", i, expected, text)
		}
	}
}

func TestDeleteParagraph(t *testing.T) {
	doc := New()
	doc.AddParagraph("First")
	doc.AddParagraph("Second")
	doc.AddParagraph("Third")

	err := doc.DeleteParagraph(1)
	if err != nil {
		t.Fatalf("Error deleting paragraph: %v", err)
	}

	if doc.GetParagraphCount() != 2 {
		t.Errorf("Expected 2 paragraphs, got %d", doc.GetParagraphCount())
	}

	text0, _ := doc.GetParagraphText(0)
	text1, _ := doc.GetParagraphText(1)

	if text0 != "First" || text1 != "Third" {
		t.Errorf("Wrong paragraphs after deletion: got '%s' and '%s'", text0, text1)
	}
}

func TestDeleteParagraphsRange(t *testing.T) {
	doc := New()
	for i := 1; i <= 5; i++ {
		doc.AddParagraph("Paragraph " + string(rune('0'+i)))
	}

	err := doc.DeleteParagraphsRange(1, 3)
	if err != nil {
		t.Fatalf("Error deleting range: %v", err)
	}

	if doc.GetParagraphCount() != 2 {
		t.Errorf("Expected 2 paragraphs, got %d", doc.GetParagraphCount())
	}

	text0, _ := doc.GetParagraphText(0)
	text1, _ := doc.GetParagraphText(1)

	if text0 != "Paragraph 1" || text1 != "Paragraph 5" {
		t.Errorf("Wrong paragraphs after range deletion: got '%s' and '%s'", text0, text1)
	}
}

func TestReplaceText(t *testing.T) {
	doc := New()
	doc.AddParagraph("Hello world")
	doc.AddParagraph("Hello again")

	count := doc.ReplaceText("Hello", "Hi")
	if count != 2 {
		t.Errorf("Expected 2 replacements, got %d", count)
	}

	text0, _ := doc.GetParagraphText(0)
	text1, _ := doc.GetParagraphText(1)

	if text0 != "Hi world" || text1 != "Hi again" {
		t.Errorf("Replacement failed: got '%s' and '%s'", text0, text1)
	}
}

func TestReplaceTextInParagraph(t *testing.T) {
	doc := New()
	doc.AddParagraph("Hello world")
	doc.AddParagraph("Hello again")

	count, err := doc.ReplaceTextInParagraph(0, "Hello", "Hi")
	if err != nil {
		t.Fatalf("Error replacing text: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 replacement, got %d", count)
	}

	text0, _ := doc.GetParagraphText(0)
	text1, _ := doc.GetParagraphText(1)

	if text0 != "Hi world" {
		t.Errorf("Expected 'Hi world', got '%s'", text0)
	}
	if text1 != "Hello again" {
		t.Errorf("Second paragraph should be unchanged, got '%s'", text1)
	}
}

func TestFindText(t *testing.T) {
	doc := New()
	doc.AddParagraph("This is a test")
	doc.AddParagraph("Another line")
	doc.AddParagraph("Test again")

	indices := doc.FindText("test")
	if len(indices) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(indices))
	}

	expected := []int{0, 2}
	for i, idx := range indices {
		if idx != expected[i] {
			t.Errorf("Expected index %d, got %d", expected[i], idx)
		}
	}
}

func TestGetText(t *testing.T) {
	doc := New()
	doc.AddParagraph("First paragraph")
	doc.AddParagraph("Second paragraph")

	text := doc.GetText()
	if text != "First paragraph Second paragraph" {
		t.Errorf("Unexpected text: %s", text)
	}
}

func TestClear(t *testing.T) {
	doc := New()
	doc.AddParagraph("Test")
	doc.AddTable(2, 2)

	doc.Clear()

	if doc.GetParagraphCount() != 0 {
		t.Errorf("Expected 0 paragraphs after clear, got %d", doc.GetParagraphCount())
	}
	if doc.GetTableCount() != 0 {
		t.Errorf("Expected 0 tables after clear, got %d", doc.GetTableCount())
	}
}

func TestSaveAndOpen(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.docx")

	// Create and save
	doc := New()
	doc.AddParagraph("Test content")
	doc.AddParagraph("More content", WithBold())

	err := doc.Save(testFile)
	if err != nil {
		t.Fatalf("Error saving document: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Document file was not created")
	}

	// Open and verify
	doc2, err := Open(testFile)
	if err != nil {
		t.Fatalf("Error opening document: %v", err)
	}

	if doc2.GetParagraphCount() != 2 {
		t.Errorf("Expected 2 paragraphs, got %d", doc2.GetParagraphCount())
	}

	text, _ := doc2.GetParagraphText(0)
	if text != "Test content" {
		t.Errorf("Expected 'Test content', got '%s'", text)
	}
}

func TestClone(t *testing.T) {
	doc := New()
	doc.AddParagraph("Original")

	cloned := doc.Clone()
	cloned.AddParagraph("Cloned")

	if doc.GetParagraphCount() == cloned.GetParagraphCount() {
		t.Error("Clone should not affect original document")
	}
}
