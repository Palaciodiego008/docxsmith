package docx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// HeaderFooterTestSuite defines the test suite for header/footer functionality
type HeaderFooterTestSuite struct {
	suite.Suite
	doc *Document
}

// SetupTest runs before each test
func (suite *HeaderFooterTestSuite) SetupTest() {
	suite.doc = New()
}

// TestSetHeader tests setting headers with various configurations
func (suite *HeaderFooterTestSuite) TestSetHeader() {
	tests := []struct {
		name        string
		hfType      HeaderFooterType
		content     string
		options     []HeaderFooterOption
		expectError bool
		errorMsg    string
	}{
		{
			name:    "default header with plain text",
			hfType:  HeaderTypeDefault,
			content: "Default Header",
			options: nil,
		},
		{
			name:    "first page header with bold text",
			hfType:  HeaderTypeFirst,
			content: "First Page Header",
			options: []HeaderFooterOption{WithHFBold()},
		},
		{
			name:    "even page header with italic and center alignment",
			hfType:  HeaderTypeEven,
			content: "Even Page Header",
			options: []HeaderFooterOption{WithHFItalic(), WithHFAlignment("center")},
		},
		{
			name:    "header with custom font size and color",
			hfType:  HeaderTypeDefault,
			content: "Styled Header",
			options: []HeaderFooterOption{WithHFFontSize("28"), WithHFTextColor("FF0000")},
		},
		{
			name:    "header with all formatting options",
			hfType:  HeaderTypeDefault,
			content: "Fully Styled Header",
			options: []HeaderFooterOption{
				WithHFBold(),
				WithHFItalic(),
				WithHFFontSize("32"),
				WithHFTextColor("0066CC"),
				WithHFAlignment("right"),
				WithHFFont("Arial"),
			},
		},
		{
			name:        "invalid header type",
			hfType:      HeaderFooterType("invalid"),
			content:     "Invalid Header",
			expectError: true,
			errorMsg:    "invalid header type",
		},
		{
			name:        "footer type used for header should fail",
			hfType:      FooterTypeDefault,
			content:     "Wrong Type",
			expectError: true,
			errorMsg:    "invalid header type",
		},
		{
			name:    "empty content header",
			hfType:  HeaderTypeDefault,
			content: "",
			options: nil,
		},
		{
			name:    "header with special characters",
			hfType:  HeaderTypeDefault,
			content: "Header with © ® ™ symbols",
			options: nil,
		},
		{
			name:    "very long header content",
			hfType:  HeaderTypeDefault,
			content: "This is a very long header content that spans multiple words and tests how the system handles lengthy text in headers",
			options: []HeaderFooterOption{WithHFAlignment("center")},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := suite.doc.SetHeader(tt.hfType, tt.content, tt.options...)

			if tt.expectError {
				assert.Error(suite.T(), err)
				if tt.errorMsg != "" {
					assert.Contains(suite.T(), err.Error(), tt.errorMsg)
				}
				return
			}

			require.NoError(suite.T(), err)
			assert.True(suite.T(), suite.doc.HasHeader(tt.hfType))

			header, err := suite.doc.GetHeader(tt.hfType)
			require.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.hfType, header.Type)
			assert.False(suite.T(), header.IsFooter)
			assert.NotEmpty(suite.T(), header.Paragraphs)
		})
	}
}

// TestSetFooter tests setting footers with various configurations
func (suite *HeaderFooterTestSuite) TestSetFooter() {
	tests := []struct {
		name        string
		hfType      HeaderFooterType
		content     string
		options     []HeaderFooterOption
		expectError bool
		errorMsg    string
	}{
		{
			name:    "default footer with plain text",
			hfType:  FooterTypeDefault,
			content: "Default Footer",
			options: nil,
		},
		{
			name:    "first page footer with center alignment",
			hfType:  FooterTypeFirst,
			content: "Page 1",
			options: []HeaderFooterOption{WithHFAlignment("center")},
		},
		{
			name:    "even page footer with right alignment and bold",
			hfType:  FooterTypeEven,
			content: "Even Page Footer",
			options: []HeaderFooterOption{WithHFBold(), WithHFAlignment("right")},
		},
		{
			name:    "footer with page numbering placeholder",
			hfType:  FooterTypeDefault,
			content: "Page {PAGE} of {NUMPAGES}",
			options: []HeaderFooterOption{WithHFAlignment("center")},
		},
		{
			name:    "footer with company info",
			hfType:  FooterTypeDefault,
			content: "© 2024 Company Name. All rights reserved.",
			options: []HeaderFooterOption{WithHFAlignment("center"), WithHFFontSize("18")},
		},
		{
			name:        "invalid footer type",
			hfType:      HeaderFooterType("invalid"),
			content:     "Invalid Footer",
			expectError: true,
			errorMsg:    "invalid footer type",
		},
		{
			name:        "header type used for footer should fail",
			hfType:      HeaderTypeDefault,
			content:     "Wrong Type",
			expectError: true,
			errorMsg:    "invalid footer type",
		},
		{
			name:    "footer with multiple formatting",
			hfType:  FooterTypeDefault,
			content: "Confidential Document",
			options: []HeaderFooterOption{
				WithHFBold(),
				WithHFItalic(),
				WithHFTextColor("FF0000"),
				WithHFAlignment("center"),
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := suite.doc.SetFooter(tt.hfType, tt.content, tt.options...)

			if tt.expectError {
				assert.Error(suite.T(), err)
				if tt.errorMsg != "" {
					assert.Contains(suite.T(), err.Error(), tt.errorMsg)
				}
				return
			}

			require.NoError(suite.T(), err)
			assert.True(suite.T(), suite.doc.HasFooter(tt.hfType))

			footer, err := suite.doc.GetFooter(tt.hfType)
			require.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.hfType, footer.Type)
			assert.True(suite.T(), footer.IsFooter)
			assert.NotEmpty(suite.T(), footer.Paragraphs)
		})
	}
}

// TestHeaderFooterRetrieval tests getting headers and footers
func (suite *HeaderFooterTestSuite) TestHeaderFooterRetrieval() {
	tests := []struct {
		name         string
		setupHeaders map[HeaderFooterType]string
		setupFooters map[HeaderFooterType]string
		getHeader    HeaderFooterType
		getFooter    HeaderFooterType
		expectHeaderError bool
		expectFooterError bool
	}{
		{
			name: "retrieve existing header and footer",
			setupHeaders: map[HeaderFooterType]string{
				HeaderTypeDefault: "Test Header",
			},
			setupFooters: map[HeaderFooterType]string{
				FooterTypeDefault: "Test Footer",
			},
			getHeader: HeaderTypeDefault,
			getFooter: FooterTypeDefault,
		},
		{
			name: "retrieve multiple headers and footers",
			setupHeaders: map[HeaderFooterType]string{
				HeaderTypeDefault: "Default Header",
				HeaderTypeFirst:   "First Header",
				HeaderTypeEven:    "Even Header",
			},
			setupFooters: map[HeaderFooterType]string{
				FooterTypeDefault: "Default Footer",
				FooterTypeFirst:   "First Footer",
			},
			getHeader: HeaderTypeEven,
			getFooter: FooterTypeFirst,
		},
		{
			name:              "retrieve non-existent header",
			getHeader:         HeaderTypeFirst, // Use different type to avoid conflicts
			expectHeaderError: true,
		},
		{
			name:              "retrieve non-existent footer",
			getFooter:         FooterTypeFirst, // Use different type to avoid conflicts
			expectFooterError: true,
		},
		{
			name: "retrieve after setting all types",
			setupHeaders: map[HeaderFooterType]string{
				HeaderTypeDefault: "Default Header",
				HeaderTypeFirst:   "First Header",
				HeaderTypeEven:    "Even Header",
			},
			setupFooters: map[HeaderFooterType]string{
				FooterTypeDefault: "Default Footer",
				FooterTypeFirst:   "First Footer",
				FooterTypeEven:    "Even Footer",
			},
			getHeader: HeaderTypeFirst,
			getFooter: FooterTypeEven,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Use fresh document for each test
			doc := New()
			
			// Setup
			for hfType, content := range tt.setupHeaders {
				err := doc.SetHeader(hfType, content)
				require.NoError(suite.T(), err)
			}
			for hfType, content := range tt.setupFooters {
				err := doc.SetFooter(hfType, content)
				require.NoError(suite.T(), err)
			}

			// Test header retrieval
			if tt.getHeader != "" {
				header, err := doc.GetHeader(tt.getHeader)
				if tt.expectHeaderError {
					assert.Error(suite.T(), err)
					assert.Nil(suite.T(), header)
				} else {
					require.NoError(suite.T(), err)
					assert.NotNil(suite.T(), header)
					assert.Equal(suite.T(), tt.getHeader, header.Type)
				}
			}

			// Test footer retrieval
			if tt.getFooter != "" {
				footer, err := doc.GetFooter(tt.getFooter)
				if tt.expectFooterError {
					assert.Error(suite.T(), err)
					assert.Nil(suite.T(), footer)
				} else {
					require.NoError(suite.T(), err)
					assert.NotNil(suite.T(), footer)
					assert.Equal(suite.T(), tt.getFooter, footer.Type)
				}
			}
		})
	}
}

// TestHeaderFooterRemoval tests removing headers and footers
func (suite *HeaderFooterTestSuite) TestHeaderFooterRemoval() {
	tests := []struct {
		name           string
		setupHeaders   []HeaderFooterType
		setupFooters   []HeaderFooterType
		removeHeader   HeaderFooterType
		removeFooter   HeaderFooterType
		expectHeaderError bool
		expectFooterError bool
	}{
		{
			name:         "remove existing header and footer",
			setupHeaders: []HeaderFooterType{HeaderTypeDefault},
			setupFooters: []HeaderFooterType{FooterTypeDefault},
			removeHeader: HeaderTypeDefault,
			removeFooter: FooterTypeDefault,
		},
		{
			name: "remove one of multiple headers",
			setupHeaders: []HeaderFooterType{
				HeaderTypeDefault,
				HeaderTypeFirst,
				HeaderTypeEven,
			},
			removeHeader: HeaderTypeFirst,
		},
		{
			name:              "remove non-existent header",
			removeHeader:      HeaderTypeFirst, // Use different type
			expectHeaderError: true,
		},
		{
			name:              "remove non-existent footer",
			removeFooter:      FooterTypeFirst, // Use different type
			expectFooterError: true,
		},
		{
			name: "remove all headers and footers",
			setupHeaders: []HeaderFooterType{
				HeaderTypeDefault,
				HeaderTypeFirst,
				HeaderTypeEven,
			},
			setupFooters: []HeaderFooterType{
				FooterTypeDefault,
				FooterTypeFirst,
				FooterTypeEven,
			},
			removeHeader: HeaderTypeDefault,
			removeFooter: FooterTypeDefault,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Use fresh document for each test
			doc := New()
			
			// Setup
			for _, hfType := range tt.setupHeaders {
				err := doc.SetHeader(hfType, "Test Header")
				require.NoError(suite.T(), err)
			}
			for _, hfType := range tt.setupFooters {
				err := doc.SetFooter(hfType, "Test Footer")
				require.NoError(suite.T(), err)
			}

			// Test header removal
			if tt.removeHeader != "" {
				err := doc.RemoveHeader(tt.removeHeader)
				if tt.expectHeaderError {
					assert.Error(suite.T(), err)
				} else {
					require.NoError(suite.T(), err)
					assert.False(suite.T(), doc.HasHeader(tt.removeHeader))
				}
			}

			// Test footer removal
			if tt.removeFooter != "" {
				err := doc.RemoveFooter(tt.removeFooter)
				if tt.expectFooterError {
					assert.Error(suite.T(), err)
				} else {
					require.NoError(suite.T(), err)
					assert.False(suite.T(), doc.HasFooter(tt.removeFooter))
				}
			}
		})
	}
}

// TestHeaderFooterOptions tests various formatting options
func (suite *HeaderFooterTestSuite) TestHeaderFooterOptions() {
	tests := []struct {
		name     string
		options  []HeaderFooterOption
		validate func(*testing.T, *HeaderFooter)
	}{
		{
			name:    "bold option",
			options: []HeaderFooterOption{WithHFBold()},
			validate: func(t *testing.T, hf *HeaderFooter) {
				require.NotEmpty(t, hf.Paragraphs)
				require.NotEmpty(t, hf.Paragraphs[0].Runs)
				run := hf.Paragraphs[0].Runs[0]
				assert.NotNil(t, run.Props)
				assert.NotNil(t, run.Props.Bold)
			},
		},
		{
			name:    "italic option",
			options: []HeaderFooterOption{WithHFItalic()},
			validate: func(t *testing.T, hf *HeaderFooter) {
				require.NotEmpty(t, hf.Paragraphs)
				require.NotEmpty(t, hf.Paragraphs[0].Runs)
				run := hf.Paragraphs[0].Runs[0]
				assert.NotNil(t, run.Props)
				assert.NotNil(t, run.Props.Italic)
			},
		},
		{
			name:    "font size option",
			options: []HeaderFooterOption{WithHFFontSize("32")},
			validate: func(t *testing.T, hf *HeaderFooter) {
				require.NotEmpty(t, hf.Paragraphs)
				require.NotEmpty(t, hf.Paragraphs[0].Runs)
				run := hf.Paragraphs[0].Runs[0]
				assert.NotNil(t, run.Props)
				assert.NotNil(t, run.Props.Size)
				assert.Equal(t, "32", run.Props.Size.Val)
			},
		},
		{
			name:    "text color option",
			options: []HeaderFooterOption{WithHFTextColor("FF0000")},
			validate: func(t *testing.T, hf *HeaderFooter) {
				require.NotEmpty(t, hf.Paragraphs)
				require.NotEmpty(t, hf.Paragraphs[0].Runs)
				run := hf.Paragraphs[0].Runs[0]
				assert.NotNil(t, run.Props)
				assert.NotNil(t, run.Props.Color)
				assert.Equal(t, "FF0000", run.Props.Color.Val)
			},
		},
		{
			name:    "center alignment option",
			options: []HeaderFooterOption{WithHFAlignment("center")},
			validate: func(t *testing.T, hf *HeaderFooter) {
				require.NotEmpty(t, hf.Paragraphs)
				paragraph := hf.Paragraphs[0]
				assert.NotNil(t, paragraph.Props)
				assert.NotNil(t, paragraph.Props.Jc)
				assert.Equal(t, "center", paragraph.Props.Jc.Val)
			},
		},
		{
			name:    "font option",
			options: []HeaderFooterOption{WithHFFont("Arial")},
			validate: func(t *testing.T, hf *HeaderFooter) {
				require.NotEmpty(t, hf.Paragraphs)
				require.NotEmpty(t, hf.Paragraphs[0].Runs)
				run := hf.Paragraphs[0].Runs[0]
				assert.NotNil(t, run.Props)
				assert.NotNil(t, run.Props.RFonts)
				assert.Equal(t, "Arial", run.Props.RFonts.ASCII)
			},
		},
		{
			name: "multiple options combined",
			options: []HeaderFooterOption{
				WithHFBold(),
				WithHFItalic(),
				WithHFFontSize("24"),
				WithHFTextColor("0066CC"),
				WithHFAlignment("right"),
			},
			validate: func(t *testing.T, hf *HeaderFooter) {
				require.NotEmpty(t, hf.Paragraphs)
				require.NotEmpty(t, hf.Paragraphs[0].Runs)
				
				run := hf.Paragraphs[0].Runs[0]
				paragraph := hf.Paragraphs[0]
				
				// Check run properties
				assert.NotNil(t, run.Props)
				assert.NotNil(t, run.Props.Bold)
				assert.NotNil(t, run.Props.Italic)
				assert.NotNil(t, run.Props.Size)
				assert.Equal(t, "24", run.Props.Size.Val)
				assert.NotNil(t, run.Props.Color)
				assert.Equal(t, "0066CC", run.Props.Color.Val)
				
				// Check paragraph properties
				assert.NotNil(t, paragraph.Props)
				assert.NotNil(t, paragraph.Props.Jc)
				assert.Equal(t, "right", paragraph.Props.Jc.Val)
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := suite.doc.SetHeader(HeaderTypeDefault, "Test Content", tt.options...)
			require.NoError(suite.T(), err)

			header, err := suite.doc.GetHeader(HeaderTypeDefault)
			require.NoError(suite.T(), err)

			tt.validate(suite.T(), header)
		})
	}
}

// TestHeaderFooterService tests the service implementation directly
func (suite *HeaderFooterTestSuite) TestHeaderFooterService() {
	service := NewHeaderFooterService(suite.doc)

	// Test interface compliance
	assert.Implements(suite.T(), (*HeaderFooterManager)(nil), service)

	// Test service methods
	err := service.SetHeader(HeaderTypeDefault, "Service Header")
	assert.NoError(suite.T(), err)

	assert.True(suite.T(), service.HasHeader(HeaderTypeDefault))
	assert.False(suite.T(), service.HasHeader(HeaderTypeFirst))

	header, err := service.GetHeader(HeaderTypeDefault)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), header)
}

// TestComplexScenarios tests complex real-world scenarios
func (suite *HeaderFooterTestSuite) TestComplexScenarios() {
	tests := []struct {
		name     string
		scenario func(*testing.T, *Document)
	}{
		{
			name: "professional document with all header/footer types",
			scenario: func(t *testing.T, doc *Document) {
				// Set all header types
				err := doc.SetHeader(HeaderTypeDefault, "Company Name", WithHFBold(), WithHFAlignment("center"))
				require.NoError(t, err)
				
				err = doc.SetHeader(HeaderTypeFirst, "DRAFT", WithHFItalic(), WithHFTextColor("FF0000"))
				require.NoError(t, err)
				
				err = doc.SetHeader(HeaderTypeEven, "Even Page Header", WithHFAlignment("left"))
				require.NoError(t, err)

				// Set all footer types
				err = doc.SetFooter(FooterTypeDefault, "Page {PAGE}", WithHFAlignment("center"))
				require.NoError(t, err)
				
				err = doc.SetFooter(FooterTypeFirst, "© 2024 Company", WithHFAlignment("center"))
				require.NoError(t, err)
				
				err = doc.SetFooter(FooterTypeEven, "Confidential", WithHFAlignment("right"))
				require.NoError(t, err)

				// Verify all are set
				assert.True(t, doc.HasHeader(HeaderTypeDefault))
				assert.True(t, doc.HasHeader(HeaderTypeFirst))
				assert.True(t, doc.HasHeader(HeaderTypeEven))
				assert.True(t, doc.HasFooter(FooterTypeDefault))
				assert.True(t, doc.HasFooter(FooterTypeFirst))
				assert.True(t, doc.HasFooter(FooterTypeEven))
			},
		},
		{
			name: "update existing headers and footers",
			scenario: func(t *testing.T, doc *Document) {
				// Set initial header
				err := doc.SetHeader(HeaderTypeDefault, "Initial Header")
				require.NoError(t, err)

				// Update with new content and formatting
				err = doc.SetHeader(HeaderTypeDefault, "Updated Header", WithHFBold(), WithHFTextColor("0066CC"))
				require.NoError(t, err)

				// Verify update
				header, err := doc.GetHeader(HeaderTypeDefault)
				require.NoError(t, err)
				
				// Check content is updated
				assert.NotEmpty(t, header.Paragraphs)
				assert.NotEmpty(t, header.Paragraphs[0].Runs)
				
				// Check formatting is applied
				run := header.Paragraphs[0].Runs[0]
				assert.NotNil(t, run.Props)
				assert.NotNil(t, run.Props.Bold)
				assert.NotNil(t, run.Props.Color)
				assert.Equal(t, "0066CC", run.Props.Color.Val)
			},
		},
		{
			name: "remove and re-add headers/footers",
			scenario: func(t *testing.T, doc *Document) {
				// Add header
				err := doc.SetHeader(HeaderTypeDefault, "Original Header")
				require.NoError(t, err)
				assert.True(t, doc.HasHeader(HeaderTypeDefault))

				// Remove header
				err = doc.RemoveHeader(HeaderTypeDefault)
				require.NoError(t, err)
				assert.False(t, doc.HasHeader(HeaderTypeDefault))

				// Re-add with different content
				err = doc.SetHeader(HeaderTypeDefault, "New Header", WithHFBold())
				require.NoError(t, err)
				assert.True(t, doc.HasHeader(HeaderTypeDefault))

				// Verify new content
				header, err := doc.GetHeader(HeaderTypeDefault)
				require.NoError(t, err)
				assert.NotEmpty(t, header.Paragraphs)
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			doc := New() // Fresh document for each scenario
			tt.scenario(suite.T(), doc)
		})
	}
}

// Run the test suite
func TestHeaderFooterTestSuite(t *testing.T) {
	suite.Run(t, new(HeaderFooterTestSuite))
}
