package converter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Palaciodiego008/docxsmith/pkg/docx"
)

func TestExternalToolDetection(t *testing.T) {
	tests := []struct {
		name     string
		checkFn  func() bool
		toolName string
	}{
		{
			name:     "LibreOffice detection",
			checkFn:  IsLibreOfficeAvailable,
			toolName: "LibreOffice",
		},
		{
			name:     "Pandoc detection",
			checkFn:  IsPandocAvailable,
			toolName: "Pandoc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			available := tt.checkFn()
			t.Logf("%s available: %v", tt.toolName, available)
		})
	}
}

func TestExternalConversion(t *testing.T) {
	tests := []struct {
		name        string
		converter   string
		skipCheck   func() bool
		inputExt    string
		outputExt   string
		content     []string
		expectError bool
	}{
		{
			name:      "LibreOffice DOCX to PDF",
			converter: "libreoffice",
			skipCheck: func() bool { return !IsLibreOfficeAvailable() },
			inputExt:  ".docx",
			outputExt: ".pdf",
			content:   []string{"Test document", "Second paragraph"},
		},
		{
			name:      "Pandoc DOCX to PDF",
			converter: "pandoc",
			skipCheck: func() bool { return !IsPandocAvailable() },
			inputExt:  ".docx",
			outputExt: ".pdf",
			content:   []string{"Pandoc test", "Another line"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipCheck() {
				t.Skipf("%s not available", tt.converter)
			}

			tmpDir := t.TempDir()
			inputPath := filepath.Join(tmpDir, "test"+tt.inputExt)
			outputPath := filepath.Join(tmpDir, "test"+tt.outputExt)

			// Create test document
			doc := docx.New()
			for _, line := range tt.content {
				doc.AddParagraph(line)
			}
			if err := doc.Save(inputPath); err != nil {
				t.Fatalf("Failed to create test document: %v", err)
			}

			// Run conversion
			var err error
			switch tt.converter {
			case "libreoffice":
				err = ConvertWithLibreOffice(inputPath, outputPath)
			case "pandoc":
				err = ConvertWithPandoc(inputPath, outputPath)
			}

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Conversion failed: %v", err)
			}

			// Verify output
			info, err := os.Stat(outputPath)
			if os.IsNotExist(err) {
				t.Errorf("Output file not created: %s", outputPath)
			}
			if err == nil && info.Size() < 100 {
				t.Errorf("Output file too small: %d bytes", info.Size())
			}
		})
	}
}

func TestUnsupportedConversions(t *testing.T) {
	tests := []struct {
		name      string
		inputExt  string
		outputExt string
		skipCheck func() bool
	}{
		{
			name:      "PDF to DOCX with LibreOffice",
			inputExt:  ".pdf",
			outputExt: ".docx",
			skipCheck: func() bool { return !IsLibreOfficeAvailable() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipCheck() {
				t.Skip("External tool not available")
			}

			tmpDir := t.TempDir()
			inputPath := filepath.Join(tmpDir, "test"+tt.inputExt)
			outputPath := filepath.Join(tmpDir, "test"+tt.outputExt)

			// Create dummy PDF
			if err := os.WriteFile(inputPath, []byte("%PDF-1.4\n"), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			err := ConvertWithLibreOffice(inputPath, outputPath)
			if err == nil {
				t.Error("Expected error for unsupported conversion")
			}
		})
	}
}

func TestMissingToolErrors(t *testing.T) {
	tests := []struct {
		name      string
		converter func(string, string) error
		wantError string
	}{
		{
			name: "LibreOffice not found",
			converter: func(in, out string) error {
				originalPath := os.Getenv("PATH")
				defer os.Setenv("PATH", originalPath)
				os.Setenv("PATH", "")
				return ConvertWithLibreOffice(in, out)
			},
			wantError: "libreoffice not found",
		},
		{
			name: "Pandoc not found",
			converter: func(in, out string) error {
				originalPath := os.Getenv("PATH")
				defer os.Setenv("PATH", originalPath)
				os.Setenv("PATH", "")
				return ConvertWithPandoc(in, out)
			},
			wantError: "pandoc not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			inputPath := filepath.Join(tmpDir, "test.docx")
			outputPath := filepath.Join(tmpDir, "test.pdf")

			err := tt.converter(inputPath, outputPath)
			if err == nil {
				t.Error("Expected error when tool not found")
			}
		})
	}
}

func TestBuiltInConverter(t *testing.T) {
	tests := []struct {
		name      string
		content   []string
		options   ConvertOptions
		minSize   int64
		expectErr bool
	}{
		{
			name:    "Simple document",
			content: []string{"Simple text"},
			options: DefaultOptions(),
			minSize: 100,
		},
		{
			name:    "Multiple paragraphs",
			content: []string{"First", "Second", "Third"},
			options: DefaultOptions(),
			minSize: 100,
		},
		{
			name:    "Custom font size",
			content: []string{"Large text"},
			options: ConvertOptions{
				FontSize:   16,
				FontFamily: "Arial",
			},
			minSize: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			inputPath := filepath.Join(tmpDir, "test.docx")
			outputPath := filepath.Join(tmpDir, "test.pdf")

			// Create document
			doc := docx.New()
			for _, line := range tt.content {
				doc.AddParagraph(line)
			}
			if err := doc.Save(inputPath); err != nil {
				t.Fatalf("Failed to create document: %v", err)
			}

			// Convert
			err := ConvertDocxToPDF(inputPath, outputPath, tt.options)
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Conversion failed: %v", err)
			}

			// Verify
			info, err := os.Stat(outputPath)
			if err != nil {
				t.Fatalf("Output file not found: %v", err)
			}
			if info.Size() < tt.minSize {
				t.Errorf("Output too small: got %d bytes, want >= %d", info.Size(), tt.minSize)
			}
		})
	}
}
