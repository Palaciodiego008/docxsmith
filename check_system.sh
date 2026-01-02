#!/bin/bash

echo "=== DocxSmith System Check ==="
echo

# Check Go installation
echo "Go Version:"
go version
echo

# Check external tools
echo "External Conversion Tools:"
if command -v libreoffice &> /dev/null; then
    echo "✅ LibreOffice: $(libreoffice --version | head -1)"
else
    echo "❌ LibreOffice: Not installed"
    echo "   Install: sudo apt-get install libreoffice-writer (Ubuntu/Debian)"
    echo "           brew install libreoffice (macOS)"
fi

if command -v pandoc &> /dev/null; then
    echo "✅ Pandoc: $(pandoc --version | head -1)"
else
    echo "❌ Pandoc: Not installed"
    echo "   Install: sudo apt-get install pandoc (Ubuntu/Debian)"
    echo "           brew install pandoc (macOS)"
fi

if command -v ocrmypdf &> /dev/null; then
    echo "✅ OCRmyPDF: $(ocrmypdf --version)"
else
    echo "⚠️  OCRmyPDF: Not installed (optional, for scanned PDFs)"
    echo "   Install: sudo apt-get install ocrmypdf (Ubuntu/Debian)"
fi
echo

# Check DocxSmith build
echo "DocxSmith Build:"
if [ -f "./docxsmith" ]; then
    echo "✅ Binary exists: $(ls -lh docxsmith | awk '{print $5}')"
    echo "   Last modified: $(stat -c %y docxsmith 2>/dev/null || stat -f %Sm docxsmith 2>/dev/null)"
else
    echo "❌ Binary not found"
    echo "   Build with: go build -o docxsmith ./cmd/docxsmith"
fi
echo

# Check dependencies
echo "Go Dependencies:"
go list -m all | grep -E "(gofpdf|pdf)" || echo "⚠️  Dependencies not downloaded"
echo

# System resources
echo "System Resources:"
if command -v free &> /dev/null; then
    echo "Memory: $(free -h | grep Mem | awk '{print $3 " used / " $2 " total"}')"
else
    echo "Memory: $(vm_stat | grep "Pages free" | awk '{print $3}' | sed 's/\.//')KB free"
fi
echo "Disk: $(df -h . | tail -1 | awk '{print $4 " available"}')"
echo

echo "=== Recommendations ==="
if ! command -v libreoffice &> /dev/null && ! command -v pandoc &> /dev/null; then
    echo "⚠️  Install LibreOffice for best conversion quality"
    echo "   sudo apt-get install libreoffice-writer"
else
    echo "✅ System ready for PDF conversion"
fi
