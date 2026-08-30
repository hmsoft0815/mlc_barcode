package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mlcmcp/mlc_barcode/internal/barcodes"
	"github.com/mlcmcp/mlc_barcode/internal/qrformats"
	"github.com/mlcmcp/mlc_barcode/internal/version"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// BarcodeApp is the main service struct exposed to the Wails v3 frontend.
type BarcodeApp struct{}

// NewBarcodeApp creates a new instance of BarcodeApp.
func NewBarcodeApp() *BarcodeApp {
	return &BarcodeApp{}
}

// GetVersion returns the current application version.
func (a *BarcodeApp) GetVersion() string {
	return version.Version
}

// FormatWifi formats Wi-Fi connection parameters into standard QR data format.
func (a *BarcodeApp) FormatWifi(opts WifiInput) string {
	enc := opts.Encryption
	if enc == "" {
		enc = "WPA"
	}
	return qrformats.FormatWifi(qrformats.WifiOptions{
		SSID:       opts.SSID,
		Password:   opts.Password,
		Encryption: enc,
	})
}

// FormatVCard formats contact parameters into vCard 3.0 standard QR data format.
func (a *BarcodeApp) FormatVCard(opts VCardInput) string {
	return qrformats.FormatVCard(qrformats.VCardOptions{
		FirstName: opts.FirstName,
		LastName:  opts.LastName,
		Email:     opts.Email,
		Phone:     opts.Phone,
	})
}

// FormatEvent formats event parameters into iCal/vCalendar QR data format.
func (a *BarcodeApp) FormatEvent(opts EventInput) string {
	return qrformats.FormatVCalendar(qrformats.VCalendarOptions{
		Summary:   opts.Summary,
		StartTime: opts.StartTime,
		EndTime:   opts.EndTime,
		TimeZone:  opts.TimeZone,
	})
}

// GenerateBarcode generates SVG and PNG representations of a single barcode.
func (a *BarcodeApp) GenerateBarcode(req BarcodeRequest) (BarcodeResult, error) {
	btype := barcodes.BarcodeType(strings.ToLower(req.Type))
	if btype == "" {
		btype = barcodes.TypeQR
	}

	data := strings.TrimSpace(req.Data)
	if data == "" {
		return BarcodeResult{
			Type:    string(btype),
			Data:    data,
			Success: false,
			Error:   "Data cannot be empty",
		}, nil
	}

	opts := barcodes.DefaultOptions(btype)
	if req.Width > 0 {
		opts.Width = req.Width
	}
	if req.Height > 0 {
		opts.Height = req.Height
	}
	opts.ShowText = req.ShowText
	if req.ForegroundColor != "" {
		opts.ForegroundColor = req.ForegroundColor
	}
	if req.BackgroundColor != "" {
		opts.BackgroundColor = req.BackgroundColor
	}

	svgStr, err := barcodes.GenerateSVG(btype, data, opts)
	if err != nil {
		return BarcodeResult{
			Type:    string(btype),
			Data:    data,
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	pngBytes, err := barcodes.GeneratePNG(btype, data, opts)
	var pngDataURI string
	if err == nil && len(pngBytes) > 0 {
		pngDataURI = "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes)
	}

	return BarcodeResult{
		Type:    string(btype),
		Data:    data,
		SVG:     svgStr,
		PNGData: pngDataURI,
		Success: true,
	}, nil
}

// GenerateBatch processes multiple lines of text and generates barcodes for each.
func (a *BarcodeApp) GenerateBatch(req BatchBarcodeRequest) (BatchBarcodeResponse, error) {
	btype := barcodes.BarcodeType(strings.ToLower(req.Type))
	if btype == "" {
		btype = barcodes.TypeQR
	}

	opts := barcodes.DefaultOptions(btype)
	if req.Width > 0 {
		opts.Width = req.Width
	}
	if req.Height > 0 {
		opts.Height = req.Height
	}
	opts.ShowText = req.ShowText
	if req.ForegroundColor != "" {
		opts.ForegroundColor = req.ForegroundColor
	}
	if req.BackgroundColor != "" {
		opts.BackgroundColor = req.BackgroundColor
	}

	items := make([]BatchItemResult, 0, len(req.Lines))
	validCount := 0
	errorCount := 0

	for i, line := range req.Lines {
		data := strings.TrimSpace(line)
		if data == "" {
			continue
		}

		item := BatchItemResult{
			Index: i + 1,
			Data:  data,
		}

		svgStr, err := barcodes.GenerateSVG(btype, data, opts)
		if err != nil {
			item.Success = false
			item.Error = err.Error()
			errorCount++
		} else {
			item.Success = true
			item.SVG = svgStr

			pngBytes, pngErr := barcodes.GeneratePNG(btype, data, opts)
			if pngErr == nil && len(pngBytes) > 0 {
				item.PNGData = "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes)
			}
			validCount++
		}

		items = append(items, item)
	}

	return BatchBarcodeResponse{
		Items:      items,
		Total:      len(items),
		ValidCount: validCount,
		ErrorCount: errorCount,
	}, nil
}

// PickTextFile opens a native open-file dialog for .txt or .csv files and reads all lines.
func (a *BarcodeApp) PickTextFile() (string, []string, error) {
	filePath, err := application.Get().Dialog.OpenFile().
		SetTitle("Select text or CSV file for batch barcodes").
		AddFilter("Text & CSV Files (*.txt, *.csv)", "*.txt;*.csv").
		AddFilter("All Files (*.*)", "*.*").
		PromptForSingleSelection()

	if err != nil || filePath == "" {
		return "", nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return filePath, nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			lines = append(lines, text)
		}
	}

	if err := scanner.Err(); err != nil {
		return filePath, lines, fmt.Errorf("error reading file: %w", err)
	}

	return filePath, lines, nil
}

// PickExportFolder opens a native folder selection dialog.
func (a *BarcodeApp) PickExportFolder() (string, error) {
	folderPath, err := application.Get().Dialog.OpenFile().
		SetTitle("Select Destination Folder for Barcode Export").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		PromptForSingleSelection()

	if err != nil {
		return "", err
	}
	return folderPath, nil
}

// sanitizeFilename turns an arbitrary string into a safe file name.
var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

func sanitizeFilename(s string) string {
	s = strings.TrimSpace(s)
	s = nonAlphanumericRegex.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	if len(s) > 40 {
		s = s[:40]
	}
	if s == "" {
		s = "barcode"
	}
	return s
}

// ExportBatchToFolder exports generated batch barcodes into a target directory.
func (a *BarcodeApp) ExportBatchToFolder(req BatchExportRequest) (BatchExportResponse, error) {
	if req.FolderPath == "" {
		return BatchExportResponse{}, fmt.Errorf("folder path cannot be empty")
	}

	if err := os.MkdirAll(req.FolderPath, 0755); err != nil {
		return BatchExportResponse{}, fmt.Errorf("failed to create destination folder: %w", err)
	}

	format := strings.ToLower(req.Format)
	if format != "png" {
		format = "svg"
	}

	prefix := req.Prefix
	savedPaths := make([]string, 0, len(req.Items))
	exportedCount := 0
	skippedCount := 0

	for i, item := range req.Items {
		if !item.Success {
			skippedCount++
			continue
		}

		var filename string
		switch req.NamingScheme {
		case "data_slug":
			filename = fmt.Sprintf("%s%03d_%s.%s", prefix, i+1, sanitizeFilename(item.Data), format)
		case "data_raw":
			filename = fmt.Sprintf("%s%s.%s", prefix, sanitizeFilename(item.Data), format)
		default: // "index"
			filename = fmt.Sprintf("%s%03d.%s", prefix, i+1, format)
		}

		destPath := filepath.Join(req.FolderPath, filename)

		var contentBytes []byte
		if format == "svg" {
			contentBytes = []byte(item.SVG)
		} else {
			// PNG Data URI
			pngBase64 := item.PNGData
			if idx := strings.Index(pngBase64, ","); idx != -1 {
				pngBase64 = pngBase64[idx+1:]
			}
			var err error
			contentBytes, err = base64.StdEncoding.DecodeString(pngBase64)
			if err != nil {
				skippedCount++
				continue
			}
		}

		if err := os.WriteFile(destPath, contentBytes, 0644); err != nil {
			skippedCount++
			continue
		}

		savedPaths = append(savedPaths, destPath)
		exportedCount++
	}

	return BatchExportResponse{
		ExportedCount: exportedCount,
		SkippedCount:  skippedCount,
		SavedPaths:    savedPaths,
	}, nil
}

// SaveSingleFile opens a native save-file dialog and writes the given SVG or PNG content.
func (a *BarcodeApp) SaveSingleFile(req SaveSingleFileRequest) (string, error) {
	format := strings.ToLower(req.Format)
	if format != "png" {
		format = "svg"
	}

	defaultName := req.DefaultName
	if defaultName == "" {
		defaultName = "barcode." + format
	}
	if !strings.HasSuffix(strings.ToLower(defaultName), "."+format) {
		defaultName = defaultName + "." + format
	}

	filterName := "SVG Image (*.svg)"
	filterPattern := "*.svg"
	if format == "png" {
		filterName = "PNG Image (*.png)"
		filterPattern = "*.png"
	}

	destPath, err := application.Get().Dialog.SaveFile().
		SetMessage("Save Barcode").
		SetFilename(defaultName).
		AddFilter(filterName, filterPattern).
		PromptForSingleSelection()

	if err != nil || destPath == "" {
		return "", err // user cancelled
	}

	var contentBytes []byte
	if format == "svg" {
		contentBytes = []byte(req.Content)
	} else {
		pngBase64 := req.Content
		if idx := strings.Index(pngBase64, ","); idx != -1 {
			pngBase64 = pngBase64[idx+1:]
		}
		var decodeErr error
		contentBytes, decodeErr = base64.StdEncoding.DecodeString(pngBase64)
		if decodeErr != nil {
			return "", fmt.Errorf("failed to decode PNG data: %w", decodeErr)
		}
	}

	if err := os.WriteFile(destPath, contentBytes, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return destPath, nil
}

// CopyToClipboard copies text (e.g. raw SVG or string) to the native OS clipboard.
func (a *BarcodeApp) CopyToClipboard(text string) bool {
	return application.Get().Clipboard.SetText(text)
}
