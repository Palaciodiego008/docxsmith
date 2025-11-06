#!/bin/bash

# DocxSmith Installation Script
# This script builds and installs the DocxSmith CLI tool

set -e

echo "🔨 DocxSmith - The Document Forge"
echo "=================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo "✓ Go version: $(go version)"
echo ""

# Build the binary
echo "📦 Building DocxSmith..."
go build -o docxsmith ./cmd/docxsmith

if [ $? -eq 0 ]; then
    echo "✓ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi

# Make it executable
chmod +x docxsmith

echo ""
echo "🎉 DocxSmith has been built successfully!"
echo ""
echo "To use DocxSmith:"
echo "  - Run from current directory: ./docxsmith <command>"
echo "  - Install globally: sudo mv docxsmith /usr/local/bin/"
echo "  - Or add to PATH: export PATH=\$PATH:\$(pwd)"
echo ""
echo "Quick start:"
echo "  ./docxsmith create -output sample.docx -text \"Hello World\""
echo "  ./docxsmith --help"
echo ""
echo "For more information, see README.md"
