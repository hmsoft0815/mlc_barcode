package gui

// BarcodeRequest contains parameters for generating a single barcode.
type BarcodeRequest struct {
	Type            string `json:"type"`            // qr, datamatrix, code128, code39, ean13, ean8, upca, itf
	Data            string `json:"data"`            // Content / text to encode
	CustomText      string `json:"customText"`      // Optional custom caption text to show below barcode
	Width           int    `json:"width"`           // Output width in px (0 for default)
	Height          int    `json:"height"`          // Output height in px (0 for default)
	ShowText        bool   `json:"showText"`        // Show readable text below barcode
	ForegroundColor string `json:"foregroundColor"` // Foreground color (e.g. #000000)
	BackgroundColor string `json:"backgroundColor"` // Background color (e.g. #ffffff or transparent)
}

// BarcodeResult holds the generated barcode content and metadata.
type BarcodeResult struct {
	Type    string `json:"type"`
	Data    string `json:"data"`
	SVG     string `json:"svg,omitempty"`     // Raw SVG XML
	PNGData string `json:"pngData,omitempty"` // Base64 data URI (data:image/png;base64,...)
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// WifiInput contains options for Wi-Fi QR codes.
type WifiInput struct {
	SSID       string `json:"ssid"`
	Password   string `json:"password"`
	Encryption string `json:"encryption"` // WPA, WEP, nopass
	Hidden     bool   `json:"hidden"`
}

// VCardInput contains options for contact vCard 3.0 QR codes.
type VCardInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// EventInput contains options for iCal/vCalendar QR codes.
type EventInput struct {
	Summary   string `json:"summary"`
	StartTime string `json:"startTime"` // YYYYMMDDTHHMMSS
	EndTime   string `json:"endTime"`
	TimeZone  string `json:"timeZone"`
}

// BatchBarcodeRequest for bulk barcode generation.
type BatchBarcodeRequest struct {
	Type            string   `json:"type"`
	Lines           []string `json:"lines"`
	Width           int      `json:"width"`
	Height          int      `json:"height"`
	ShowText        bool     `json:"showText"`
	ForegroundColor string   `json:"foregroundColor"`
	BackgroundColor string   `json:"backgroundColor"`
}

// BatchItemResult holds the result for one item in a batch.
type BatchItemResult struct {
	Index   int    `json:"index"`
	Data    string `json:"data"`
	SVG     string `json:"svg,omitempty"`
	PNGData string `json:"pngData,omitempty"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// BatchBarcodeResponse summarizes batch generation results.
type BatchBarcodeResponse struct {
	Items      []BatchItemResult `json:"items"`
	Total      int               `json:"total"`
	ValidCount int               `json:"validCount"`
	ErrorCount int               `json:"errorCount"`
}

// BatchExportRequest for exporting generated batch items to a folder.
type BatchExportRequest struct {
	FolderPath   string            `json:"folderPath"`
	NamingScheme string            `json:"namingScheme"` // "index", "data_slug", "data_raw"
	Prefix       string            `json:"prefix"`       // Optional prefix, e.g. "barcode_"
	Format       string            `json:"format"`       // "svg" or "png"
	Items        []BatchItemResult `json:"items"`
}

// BatchExportResponse returns the outcome of folder export.
type BatchExportResponse struct {
	ExportedCount int      `json:"exportedCount"`
	SkippedCount  int      `json:"skippedCount"`
	SavedPaths    []string `json:"savedPaths"`
	Error         string   `json:"error,omitempty"`
}

// SaveSingleFileRequest for saving one barcode to disk.
type SaveSingleFileRequest struct {
	DefaultName string `json:"defaultName"`
	Format      string `json:"format"`  // "svg" or "png"
	Content     string `json:"content"` // Raw SVG or Base64 PNG
}
