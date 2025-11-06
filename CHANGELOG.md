# Changelog

All notable changes to DocxSmith will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **PDF Support** - Full PDF document manipulation capabilities
  - Create new PDF documents with text and styling
  - Read and parse existing PDF files
  - Extract text from PDFs
  - Add text content with formatting (bold, italic, colors, sizes)
  - Table support in PDF generation
  - Metadata management (title, author, subject)
- **Format Conversion** - Convert between DOCX and PDF
  - DOCX to PDF conversion with formatting preservation
  - PDF to DOCX conversion for editing
  - Customizable conversion options (font size, family, page size)
- **New CLI Commands**
  - `pdf-create` - Create new PDF documents
  - `pdf-add` - Add content to PDF documents
  - `pdf-info` - Display PDF information
  - `pdf-extract` - Extract text from PDFs
  - `convert` - Convert between DOCX and PDF formats
- **New Packages**
  - `pkg/pdf` - PDF document manipulation
  - `pkg/converter` - Format conversion utilities

### Changed
- Refactored CLI architecture for better maintainability
  - Moved command logic from `cmd/docxsmith/main.go` to `internal/cli` package
  - Separated commands into logical files: `create.go`, `content.go`, `text.go`, `table.go`, `info.go`, `pdf.go`, `convert.go`
  - main.go now only 12 lines (minimal entry point)
  - Improved code organization following Go best practices
- Updated help text to include PDF and conversion commands
- Enhanced project structure for scalability

## [1.0.0] - 2025-11-05

### Added
- Initial release of DocxSmith
- Core document manipulation features:
  - Create new .docx documents
  - Open and read existing documents
  - Add paragraphs with rich formatting (bold, italic, colors, sizes, alignment)
  - Delete paragraphs and ranges
  - Find and replace text
  - Extract text content
- Table support:
  - Create tables
  - Set/get cell content
  - Add/delete rows
  - Delete tables
- CLI tool with commands:
  - `create` - Create new documents
  - `add` - Add content
  - `delete` - Remove content
  - `replace` - Replace text
  - `find` - Search for text
  - `extract` - Extract text
  - `table` - Table operations
  - `info` - Display document information
  - `clear` - Clear all content
- Comprehensive test suite with 100% core functionality coverage
- Full documentation and examples
- Makefile for common development tasks
- MIT License

### Features
- Zero dependencies for core functionality
- Uses Go standard library
- Clean and intuitive API
- Well-documented code
- Professional project structure

[1.0.0]: https://github.com/Palaciodiego008/docxsmith/releases/tag/v1.0.0
