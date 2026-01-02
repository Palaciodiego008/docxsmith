#!/bin/bash

echo "=== DocxSmith Conversion Test ==="
echo

# Build the tool
echo "Building docxsmith..."
go build -o docxsmith ./cmd/docxsmith
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi
echo "✅ Build successful"
echo

# Create test document
echo "Creating test DOCX..."
./docxsmith create -output test_conversion.docx -text "This is a test document for conversion testing."
./docxsmith add -input test_conversion.docx -output test_conversion.docx -text "Second paragraph with more content." -bold
echo "✅ Test document created"
echo

# Check for external tools
echo "Checking for external conversion tools..."
if command -v libreoffice &> /dev/null; then
    echo "✅ LibreOffice found"
    HAS_EXTERNAL=true
elif command -v pandoc &> /dev/null; then
    echo "✅ Pandoc found"
    HAS_EXTERNAL=true
else
    echo "⚠️  No external tools found (using built-in converters)"
    echo "   Install LibreOffice or Pandoc for better quality:"
    echo "   - Ubuntu/Debian: sudo apt-get install libreoffice-writer"
    echo "   - macOS: brew install libreoffice"
    HAS_EXTERNAL=false
fi
echo

# Test DOCX to PDF
echo "Testing DOCX → PDF conversion..."
./docxsmith convert -input test_conversion.docx -output test_conversion.pdf
if [ $? -eq 0 ] && [ -f test_conversion.pdf ]; then
    echo "✅ DOCX → PDF successful"
    ls -lh test_conversion.pdf
else
    echo "❌ DOCX → PDF failed"
fi
echo

# Test PDF to DOCX (only if external tools available)
if [ "$HAS_EXTERNAL" = true ]; then
    echo "Testing PDF → DOCX conversion..."
    ./docxsmith convert -input test_conversion.pdf -output test_conversion_back.docx
    if [ $? -eq 0 ] && [ -f test_conversion_back.docx ]; then
        echo "✅ PDF → DOCX successful"
        ls -lh test_conversion_back.docx
    else
        echo "❌ PDF → DOCX failed"
    fi
else
    echo "⚠️  Skipping PDF → DOCX (requires external tools)"
fi
echo

echo "=== Test Complete ==="
echo "Generated files:"
ls -lh test_conversion.* 2>/dev/null
