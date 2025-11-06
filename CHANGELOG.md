# Changelog

All notable changes to DocxSmith will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Document Diff** - Professional document comparison tool
  - Compare DOCX documents line-by-line
  - Multiple output formats (HTML, Markdown, Plain Text)
  - Beautiful HTML reports with color-coded changes
  - Options to ignore whitespace and case
  - Detailed statistics (added, deleted, modified lines)
  - Myers diff algorithm (LCS-based)
  - 89.4% test coverage
- **Merge & Split** - Advanced document operations
  - Merge multiple DOCX or PDF files
  - Split documents by page ranges
  - Split into N equal parts
  - Smart split by heading levels (DOCX)
  - Custom output patterns
  - Merge info preview
  - 82.1% test coverage
- **CI/CD Pipeline** - Professional automation
  - GitHub Actions workflows (CI, Release)
  - 4 CI jobs: Unit Tests, Lint, Build, Security
  - Unit tests with race detector
  - Code quality checks (fmt, vet)
  - Automated releases with binaries for 5 platforms
  - Pre-commit hooks for local validation
  - Dependabot for automatic updates
  - Security scanning (Gosec)
  - Coverage reporting
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
- **Major CLI Architecture Refactor**
  - Implemented Command Pattern for extensibility
  - Created common utilities module (`common.go`) to eliminate code duplication
  - Added interfaces for commands (`Command`, `Renderer`)
  - Separated concerns: each command in its own file
  - main.go remains minimal (12 lines)
  - Improved error handling with helper functions
  - Reusable flag definitions
  - DRY principle applied throughout
- **Improved Code Quality**
  - Eliminated duplicate validation code
  - Common flag helpers across commands
  - Consistent error messages
  - Better separation of concerns
- Enhanced project structure for maximum scalability
- Updated help text to include all new commands

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
