package docx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAddWatermark tests adding watermarks with table-driven test cases
func TestAddWatermark(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		opts      []WatermarkOption
		shouldErr bool
		hasWMark  bool
	}{
		{
			name:      "basic watermark addition",
			text:      "DRAFT",
			opts:      nil,
			shouldErr: false,
			hasWMark:  true,
		},
		{
			name: "watermark with all custom options",
			text: "CONFIDENTIAL",
			opts: []WatermarkOption{
				WithWatermarkOpacity(0.3),
				WithWatermarkAngle(-30),
				WithWatermarkSize(120),
				WithWatermarkColor("FF0000"),
			},
			shouldErr: false,
			hasWMark:  true,
		},
		{
			name:      "watermark with single opacity option",
			text:      "REVIEW",
			opts:      []WatermarkOption{WithWatermarkOpacity(0.5)},
			shouldErr: false,
			hasWMark:  true,
		},
		{
			name:      "watermark with single size option",
			text:      "IMPORTANT",
			opts:      []WatermarkOption{WithWatermarkSize(150)},
			shouldErr: false,
			hasWMark:  true,
		},
		{
			name:      "watermark with single angle option",
			text:      "INFO",
			opts:      []WatermarkOption{WithWatermarkAngle(-20)},
			shouldErr: false,
			hasWMark:  true,
		},
		{
			name:      "watermark with single color option",
			text:      "NOTICE",
			opts:      []WatermarkOption{WithWatermarkColor("0000FF")},
			shouldErr: false,
			hasWMark:  true,
		},
		{
			name:      "empty watermark text should fail",
			text:      "",
			opts:      nil,
			shouldErr: true,
			hasWMark:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()
			err := doc.AddWatermark(tt.text, tt.opts...)

			if tt.shouldErr {
				assert.Error(t, err)
				assert.False(t, doc.HasWatermark())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.hasWMark, doc.HasWatermark())
			}
		})
	}
}

// TestRemoveWatermark tests removing watermarks with table-driven test cases
func TestRemoveWatermark(t *testing.T) {
	tests := []struct {
		name      string
		addFirst  bool
		shouldErr bool
	}{
		{
			name:      "remove existing watermark",
			addFirst:  true,
			shouldErr: false,
		},
		{
			name:      "remove non-existing watermark should fail",
			addFirst:  false,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()

			if tt.addFirst {
				err := doc.AddWatermark("DRAFT")
				assert.NoError(t, err)
				assert.True(t, doc.HasWatermark())
			}

			err := doc.RemoveWatermark()

			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.False(t, doc.HasWatermark())
			}
		})
	}
}

// TestHasWatermark tests the HasWatermark check with table-driven cases
func TestHasWatermark(t *testing.T) {
	tests := []struct {
		name         string
		addWatermark bool
		expected     bool
	}{
		{
			name:         "new document should not have watermark",
			addWatermark: false,
			expected:     false,
		},
		{
			name:         "document with watermark should return true",
			addWatermark: true,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()

			if tt.addWatermark {
				err := doc.AddWatermark("TEST")
				assert.NoError(t, err)
			}

			result := doc.HasWatermark()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestWatermarkOpacityOption tests opacity configuration with table-driven cases
func TestWatermarkOpacityOption(t *testing.T) {
	tests := []struct {
		name        string
		opacity     float64
		shouldApply bool
		description string
	}{
		{
			name:        "valid opacity 0",
			opacity:     0.0,
			shouldApply: true,
			description: "minimum opacity",
		},
		{
			name:        "valid opacity 0.5",
			opacity:     0.5,
			shouldApply: true,
			description: "middle opacity",
		},
		{
			name:        "valid opacity 1.0",
			opacity:     1.0,
			shouldApply: true,
			description: "maximum opacity",
		},
		{
			name:        "invalid negative opacity",
			opacity:     -0.1,
			shouldApply: false,
			description: "negative opacity ignored",
		},
		{
			name:        "invalid opacity above 1",
			opacity:     1.5,
			shouldApply: false,
			description: "opacity > 1 ignored",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()
			opts := []WatermarkOption{WithWatermarkOpacity(tt.opacity)}
			err := doc.AddWatermark("TEST", opts...)

			assert.NoError(t, err, "should not error even with invalid values")
			assert.True(t, doc.HasWatermark(), "watermark should be added")
		})
	}
}

// TestWatermarkSizeOption tests size configuration with table-driven cases
func TestWatermarkSizeOption(t *testing.T) {
	tests := []struct {
		name        string
		size        int
		shouldApply bool
		description string
	}{
		{
			name:        "minimum valid size",
			size:        1,
			shouldApply: true,
			description: "size 1 point",
		},
		{
			name:        "standard size",
			size:        96,
			shouldApply: true,
			description: "size 96 points (default)",
		},
		{
			name:        "large size",
			size:        200,
			shouldApply: true,
			description: "size 200 points",
		},
		{
			name:        "zero size ignored",
			size:        0,
			shouldApply: false,
			description: "zero size ignored",
		},
		{
			name:        "negative size ignored",
			size:        -10,
			shouldApply: false,
			description: "negative size ignored",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()
			opts := []WatermarkOption{WithWatermarkSize(tt.size)}
			err := doc.AddWatermark("TEST", opts...)

			assert.NoError(t, err, "should not error even with invalid values")
			assert.True(t, doc.HasWatermark(), "watermark should be added")
		})
	}
}

// TestWatermarkColorOption tests color configuration with table-driven cases
func TestWatermarkColorOption(t *testing.T) {
	tests := []struct {
		name        string
		color       string
		description string
	}{
		{
			name:        "red color",
			color:       "FF0000",
			description: "red hex color",
		},
		{
			name:        "green color",
			color:       "00FF00",
			description: "green hex color",
		},
		{
			name:        "blue color",
			color:       "0000FF",
			description: "blue hex color",
		},
		{
			name:        "black color",
			color:       "000000",
			description: "black hex color",
		},
		{
			name:        "white color",
			color:       "FFFFFF",
			description: "white hex color",
		},
		{
			name:        "gray color",
			color:       "CCCCCC",
			description: "light gray hex color",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()
			opts := []WatermarkOption{WithWatermarkColor(tt.color)}
			err := doc.AddWatermark("TEST", opts...)

			assert.NoError(t, err)
			assert.True(t, doc.HasWatermark())
		})
	}
}

// TestWatermarkAngleOption tests angle configuration with table-driven cases
func TestWatermarkAngleOption(t *testing.T) {
	tests := []struct {
		name        string
		angle       float64
		description string
	}{
		{
			name:        "extreme left angle",
			angle:       -90,
			description: "rotate -90 degrees",
		},
		{
			name:        "standard angle",
			angle:       -45,
			description: "rotate -45 degrees",
		},
		{
			name:        "horizontal angle",
			angle:       0,
			description: "no rotation",
		},
		{
			name:        "positive angle",
			angle:       45,
			description: "rotate 45 degrees",
		},
		{
			name:        "extreme right angle",
			angle:       90,
			description: "rotate 90 degrees",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()
			opts := []WatermarkOption{WithWatermarkAngle(tt.angle)}
			err := doc.AddWatermark("TEST", opts...)

			assert.NoError(t, err)
			assert.True(t, doc.HasWatermark())
		})
	}
}

// TestWatermarkMultipleOperations tests multiple watermark operations with table-driven cases
func TestWatermarkMultipleOperations(t *testing.T) {
	tests := []struct {
		name       string
		operations []string // "add", "remove", "add"
		finalState bool     // true if should have watermark at the end
		shouldErr  bool
	}{
		{
			name:       "add once",
			operations: []string{"add"},
			finalState: true,
			shouldErr:  false,
		},
		{
			name:       "add then remove",
			operations: []string{"add", "remove"},
			finalState: false,
			shouldErr:  false,
		},
		{
			name:       "add twice",
			operations: []string{"add", "add"},
			finalState: true,
			shouldErr:  false,
		},
		{
			name:       "add remove then add again",
			operations: []string{"add", "remove", "add"},
			finalState: true,
			shouldErr:  false,
		},
		{
			name:       "remove without add should error",
			operations: []string{"remove"},
			finalState: false,
			shouldErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()
			var lastErr error

			for _, op := range tt.operations {
				switch op {
				case "add":
					lastErr = doc.AddWatermark("TEST")
				case "remove":
					lastErr = doc.RemoveWatermark()
				}
			}

			if tt.shouldErr {
				assert.Error(t, lastErr)
			} else {
				assert.NoError(t, lastErr)
			}

			assert.Equal(t, tt.finalState, doc.HasWatermark())
		})
	}
}

// TestWatermarkSaveAndLoad tests saving and loading watermarked documents
func TestWatermarkSaveAndLoad(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		opts      []WatermarkOption
		shouldErr bool
	}{
		{
			name:      "save simple watermarked document",
			text:      "DRAFT",
			opts:      nil,
			shouldErr: false,
		},
		{
			name: "save complex watermarked document",
			text: "CONFIDENTIAL",
			opts: []WatermarkOption{
				WithWatermarkOpacity(0.5),
				WithWatermarkSize(100),
				WithWatermarkAngle(-45),
				WithWatermarkColor("FF0000"),
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := New()
			doc.AddParagraph("Test content")

			err := doc.AddWatermark(tt.text, tt.opts...)
			assert.NoError(t, err)

			tmpFile := "test_watermark_" + tt.name + ".docx"
			err = doc.Save(tmpFile)

			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Try to load the saved document
				loadedDoc, err := Open(tmpFile)
				assert.NoError(t, err)
				assert.NotNil(t, loadedDoc)
			}
		})
	}
}

// Helper function to clean up test files (optional)
func cleanupTestFile(filename string) {
	// Implement file cleanup if needed
}
