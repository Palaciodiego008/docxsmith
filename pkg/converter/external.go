package converter

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// ConvertWithLibreOffice uses LibreOffice for conversion
func ConvertWithLibreOffice(inputPath, outputPath string) error {
	inputExt := strings.ToLower(filepath.Ext(inputPath))
	outputExt := strings.ToLower(filepath.Ext(outputPath))

	// Check if LibreOffice is available
	if _, err := exec.LookPath("libreoffice"); err != nil {
		return fmt.Errorf("libreoffice not found. Install with: sudo apt-get install libreoffice (Ubuntu/Debian) or brew install libreoffice (macOS)")
	}

	var filterName string
	switch {
	case inputExt == ".docx" && outputExt == ".pdf":
		filterName = "writer_pdf_Export"
	case inputExt == ".pdf" && outputExt == ".docx":
		return fmt.Errorf("PDF to DOCX conversion requires OCR. Use: ConvertPDFToDocxOCR")
	default:
		return fmt.Errorf("unsupported conversion: %s to %s", inputExt, outputExt)
	}

	outputDir := filepath.Dir(outputPath)
	
	cmd := exec.Command("libreoffice",
		"--headless",
		"--convert-to", strings.TrimPrefix(outputExt, "."),
		"--outdir", outputDir,
		inputPath,
	)

	if filterName != "" {
		cmd.Args = append(cmd.Args[:2], append([]string{"--convert-to", strings.TrimPrefix(outputExt, ".") + ":" + filterName}, cmd.Args[4:]...)...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("conversion failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// ConvertWithPandoc uses Pandoc for conversion
func ConvertWithPandoc(inputPath, outputPath string) error {
	if _, err := exec.LookPath("pandoc"); err != nil {
		return fmt.Errorf("pandoc not found. Install with: sudo apt-get install pandoc (Ubuntu/Debian) or brew install pandoc (macOS)")
	}

	cmd := exec.Command("pandoc", inputPath, "-o", outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("conversion failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// IsLibreOfficeAvailable checks if LibreOffice is installed
func IsLibreOfficeAvailable() bool {
	_, err := exec.LookPath("libreoffice")
	return err == nil
}

// IsPandocAvailable checks if Pandoc is installed
func IsPandocAvailable() bool {
	_, err := exec.LookPath("pandoc")
	return err == nil
}
