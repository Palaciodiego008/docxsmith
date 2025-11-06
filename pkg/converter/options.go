package converter

// ConvertOptions holds options for conversion
type ConvertOptions struct {
	// PageSize specifies the page size (A4, Letter, Legal)
	PageSize string

	// Orientation specifies page orientation (Portrait, Landscape)
	Orientation string

	// FontSize specifies the default font size
	FontSize float64

	// FontFamily specifies the default font family
	FontFamily string

	// Margins specifies page margins in mm (left, top, right, bottom)
	Margins [4]float64
}

// DefaultOptions returns default conversion options
func DefaultOptions() ConvertOptions {
	return ConvertOptions{
		PageSize:    "A4",
		Orientation: "Portrait",
		FontSize:    12,
		FontFamily:  "Arial",
		Margins:     [4]float64{20, 20, 20, 20}, // left, top, right, bottom
	}
}
