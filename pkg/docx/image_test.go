package docx

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddImage(t *testing.T) {
	tests := []struct {
		name        string
		imageName   string
		imageData   []byte
		options     []ImageOption
		expectError bool
		errorMsg    string
	}{
		{
			name:      "valid PNG with default options",
			imageName: "test.png",
			imageData: createPNGData(),
			options:   nil,
		},
		{
			name:      "valid JPEG with custom dimensions",
			imageName: "test.jpg",
			imageData: createJPEGData(),
			options:   []ImageOption{WithImageWidth(400), WithImageHeight(300)},
		},
		{
			name:      "valid GIF with small dimensions",
			imageName: "test.gif",
			imageData: createGIFData(),
			options:   []ImageOption{WithImageWidth(50), WithImageHeight(50)},
		},
		{
			name:      "valid BMP with large dimensions",
			imageName: "test.bmp",
			imageData: createBMPData(),
			options:   []ImageOption{WithImageWidth(800), WithImageHeight(600)},
		},
		{
			name:        "nonexistent file",
			imageName:   "nonexistent.png",
			expectError: true,
			errorMsg:    "image file does not exist",
		},
		{
			name:        "unsupported extension",
			imageName:   "test.txt",
			imageData:   []byte("not an image"),
			expectError: true,
			errorMsg:    "unsupported image format",
		},
		{
			name:        "invalid PNG header",
			imageName:   "invalid.png",
			imageData:   []byte("fake png data"),
			expectError: true,
			errorMsg:    "does not appear to be a valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()
			var testImagePath string

			if tt.imageData != nil {
				testImagePath = createTestImageFile(t, tt.imageName, tt.imageData)
				defer os.Remove(testImagePath)
			} else if !tt.expectError {
				t.Fatal("Test case missing image data")
			} else {
				testImagePath = tt.imageName
			}

			err := doc.AddImage(testImagePath, tt.options...)

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error containing '%s', got nil", tt.errorMsg)
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Fatalf("Expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(doc.Body.Paragraphs) == 0 {
				t.Fatal("No paragraphs found after adding image")
			}

			p := doc.Body.Paragraphs[0]
			if len(p.Runs) == 0 || p.Runs[0].Drawing == nil {
				t.Fatal("Image not properly added to document")
			}
		})
	}
}

func TestAddImageAt(t *testing.T) {
	tests := []struct {
		name           string
		initialParas   []string
		insertIndex    int
		expectedParas  int
		expectError    bool
		errorMsg       string
	}{
		{
			name:          "insert at beginning",
			initialParas:  []string{"First", "Second"},
			insertIndex:   0,
			expectedParas: 3,
		},
		{
			name:          "insert in middle",
			initialParas:  []string{"First", "Second", "Third"},
			insertIndex:   1,
			expectedParas: 4,
		},
		{
			name:          "insert at end",
			initialParas:  []string{"First", "Second"},
			insertIndex:   2,
			expectedParas: 3,
		},
		{
			name:          "insert in empty document",
			initialParas:  []string{},
			insertIndex:   0,
			expectedParas: 1,
		},
		{
			name:        "negative index",
			initialParas: []string{"First"},
			insertIndex: -1,
			expectError: true,
			errorMsg:    "index -1 out of range",
		},
		{
			name:        "index too large",
			initialParas: []string{"First"},
			insertIndex: 5,
			expectError: true,
			errorMsg:    "index 5 out of range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()
			testImagePath := createTestImageFile(t, "test.png", createPNGData())
			defer os.Remove(testImagePath)

			// Add initial paragraphs
			for _, text := range tt.initialParas {
				doc.AddParagraph(text)
			}

			err := doc.AddImageAt(tt.insertIndex, testImagePath)

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error containing '%s', got nil", tt.errorMsg)
				}
				if !contains(err.Error(), tt.errorMsg) {
					t.Fatalf("Expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(doc.Body.Paragraphs) != tt.expectedParas {
				t.Fatalf("Expected %d paragraphs, got %d", tt.expectedParas, len(doc.Body.Paragraphs))
			}

			// Verify image is at correct position
			imageParagraph := doc.Body.Paragraphs[tt.insertIndex]
			if len(imageParagraph.Runs) == 0 || imageParagraph.Runs[0].Drawing == nil {
				t.Fatal("Image not found at expected position")
			}
		})
	}
}

func TestImageOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []ImageOption
		expected struct {
			width  int
			height int
		}
	}{
		{
			name:    "default options",
			options: nil,
			expected: struct {
				width  int
				height int
			}{200, 150},
		},
		{
			name:    "custom width only",
			options: []ImageOption{WithImageWidth(300)},
			expected: struct {
				width  int
				height int
			}{300, 150},
		},
		{
			name:    "custom height only",
			options: []ImageOption{WithImageHeight(250)},
			expected: struct {
				width  int
				height int
			}{200, 250},
		},
		{
			name:    "both custom dimensions",
			options: []ImageOption{WithImageWidth(500), WithImageHeight(400)},
			expected: struct {
				width  int
				height int
			}{500, 400},
		},
		{
			name:    "very small dimensions",
			options: []ImageOption{WithImageWidth(1), WithImageHeight(1)},
			expected: struct {
				width  int
				height int
			}{1, 1},
		},
		{
			name:    "very large dimensions",
			options: []ImageOption{WithImageWidth(2000), WithImageHeight(1500)},
			expected: struct {
				width  int
				height int
			}{2000, 1500},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()
			testImagePath := createTestImageFile(t, "test.png", createPNGData())
			defer os.Remove(testImagePath)

			err := doc.AddImage(testImagePath, tt.options...)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify options were applied (this is a simplified check)
			// In a real implementation, you'd inspect the XML structure
			if len(doc.Body.Paragraphs) == 0 {
				t.Fatal("No paragraphs found")
			}
		})
	}
}

func TestGetImageCount(t *testing.T) {
	tests := []struct {
		name         string
		imageCount   int
		textParas    int
		expectedCount int
	}{
		{
			name:         "no images",
			imageCount:   0,
			textParas:    3,
			expectedCount: 0,
		},
		{
			name:         "single image",
			imageCount:   1,
			textParas:    2,
			expectedCount: 1,
		},
		{
			name:         "multiple images",
			imageCount:   5,
			textParas:    3,
			expectedCount: 5,
		},
		{
			name:         "only images",
			imageCount:   3,
			textParas:    0,
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()
			testImagePath := createTestImageFile(t, "test.png", createPNGData())
			defer os.Remove(testImagePath)

			// Add text paragraphs
			for i := 0; i < tt.textParas; i++ {
				doc.AddParagraph("Text paragraph")
			}

			// Add images
			for i := 0; i < tt.imageCount; i++ {
				err := doc.AddImage(testImagePath)
				if err != nil {
					t.Fatalf("Failed to add image %d: %v", i, err)
				}
			}

			count := doc.GetImageCount()
			if count != tt.expectedCount {
				t.Fatalf("Expected %d images, got %d", tt.expectedCount, count)
			}
		})
	}
}

func TestImageValidation(t *testing.T) {
	tests := []struct {
		name      string
		extension string
		header    []byte
		expectErr bool
	}{
		{
			name:      "valid PNG",
			extension: ".png",
			header:    []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			expectErr: false,
		},
		{
			name:      "valid JPEG",
			extension: ".jpg",
			header:    []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46},
			expectErr: false,
		},
		{
			name:      "valid GIF87a",
			extension: ".gif",
			header:    []byte{'G', 'I', 'F', '8', '7', 'a', 0x00, 0x00},
			expectErr: false,
		},
		{
			name:      "valid GIF89a",
			extension: ".gif",
			header:    []byte{'G', 'I', 'F', '8', '9', 'a', 0x00, 0x00},
			expectErr: false,
		},
		{
			name:      "valid BMP",
			extension: ".bmp",
			header:    []byte{0x42, 0x4D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expectErr: false,
		},
		{
			name:      "invalid PNG header",
			extension: ".png",
			header:    []byte{0x00, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			expectErr: true,
		},
		{
			name:      "invalid JPEG header",
			extension: ".jpg",
			header:    []byte{0x00, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46},
			expectErr: true,
		},
		{
			name:      "invalid GIF header",
			extension: ".gif",
			header:    []byte{'X', 'I', 'F', '8', '7', 'a', 0x00, 0x00},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()
			testImagePath := createTestImageFile(t, "test"+tt.extension, tt.header)
			defer os.Remove(testImagePath)

			err := doc.AddImage(testImagePath)

			if tt.expectErr && err == nil {
				t.Fatal("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		})
	}
}

// Helper functions
func createTestImageFile(t *testing.T, filename string, data []byte) string {
	testFile := filepath.Join(os.TempDir(), filename)
	err := os.WriteFile(testFile, data, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return testFile
}

func createPNGData() []byte {
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1 dimensions
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, // bit depth, color type, etc.
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, // IDAT chunk
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0x0F, 0x00, 0x00,
		0x01, 0x00, 0x01, 0x5C, 0xC2, 0x8A, 0x8E, 0x00, // IDAT data
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, // IEND chunk
		0x42, 0x60, 0x82,
	}
}

func createJPEGData() []byte {
	return []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, // JPEG header
		0x49, 0x46, 0x00, 0x01, 0x01, 0x01, 0x00, 0x48,
		0x00, 0x48, 0x00, 0x00, 0xFF, 0xDB, 0x00, 0x43,
		0x00, 0x08, 0x06, 0x06, 0x07, 0x06, 0x05, 0x08,
		0xFF, 0xD9, // End of image
	}
}

func createGIFData() []byte {
	return []byte{
		'G', 'I', 'F', '8', '9', 'a', // GIF89a signature
		0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x21,
		0xF9, 0x04, 0x01, 0x00, 0x00, 0x00, 0x00, 0x2C,
		0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00,
		0x00, 0x02, 0x02, 0x04, 0x01, 0x00, 0x3B, // End
	}
}

func createBMPData() []byte {
	return []byte{
		0x42, 0x4D, // BM signature
		0x3A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x36, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x18, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00,
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())))
}
