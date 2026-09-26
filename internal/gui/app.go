package gui

import (
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

// FormatEPC formats SEPA credit transfer parameters into EPC-QR-Code (GiroCode).
func (a *BarcodeApp) FormatEPC(opts EPCInput) string {
	return qrformats.FormatEPC(qrformats.EPCOptions{
		Name:      opts.Name,
		IBAN:      opts.IBAN,
		BIC:       opts.BIC,
		Amount:    opts.Amount,
		Reference: opts.Reference,
		Purpose:   opts.Purpose,
	})
}

// FormatCrypto formats cryptocurrency payment parameters into standard URI format.
func (a *BarcodeApp) FormatCrypto(opts CryptoInput) string {
	return qrformats.FormatCrypto(qrformats.CryptoOptions{
		Coin:    opts.Coin,
		Address: opts.Address,
		Amount:  opts.Amount,
		Label:   opts.Label,
		Message: opts.Message,
	})
}

// FormatGeo formats latitude, longitude and search query into RFC 5870 geo URI.
func (a *BarcodeApp) FormatGeo(opts GeoInput) string {
	return qrformats.FormatGeo(qrformats.GeoOptions{
		Latitude:  opts.Latitude,
		Longitude: opts.Longitude,
		Query:     opts.Query,
	})
}

// FormatTel formats telephone dialing number into RFC 3966 tel URI.
func (a *BarcodeApp) FormatTel(opts TelInput) string {
	return qrformats.FormatTel(qrformats.TelOptions{
		PhoneNumber: opts.PhoneNumber,
	})
}

// FormatSMS formats phone number and message into smsto URI.
func (a *BarcodeApp) FormatSMS(opts SMSInput) string {
	return qrformats.FormatSMS(qrformats.SMSOptions{
		PhoneNumber: opts.PhoneNumber,
		Message:     opts.Message,
	})
}

// FormatEmail formats email address, subject and body into RFC 6068 mailto URI.
func (a *BarcodeApp) FormatEmail(opts EmailInput) string {
	return qrformats.FormatEmail(qrformats.EmailOptions{
		To:      opts.To,
		Subject: opts.Subject,
		Body:    opts.Body,
	})
}

// ValidateBarcode checks length and check digit of EAN-13, EAN-8 and UPC-A
// input while the user types.
func (a *BarcodeApp) ValidateBarcode(barcodeType, data string) RetailValidation {
	btype := barcodes.BarcodeType(strings.ToLower(barcodeType))
	if !barcodes.IsRetail(btype) {
		return RetailValidation{Valid: true}
	}
	c := barcodes.CheckRetail(btype, strings.TrimSpace(data))
	return RetailValidation{
		Applies:         true,
		Valid:           c.Valid,
		Reason:          c.Reason,
		Code:            c.Code,
		CheckDigitAdded: c.CheckDigitAdded,
		Given:           c.Given,
		Expected:        c.Expected,
		MinLength:       c.Lengths[0],
		MaxLength:       c.Lengths[1],
	}
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
			Type:      string(btype),
			Data:      data,
			Success:   false,
			Error:     "Data cannot be empty",
			ErrorInfo: ErrorInfo{ErrorCode: barcodes.ErrEmpty},
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
	opts.FontSize = req.FontSize
	opts.CustomText = req.CustomText
	if req.ForegroundColor != "" {
		opts.ForegroundColor = req.ForegroundColor
	}
	if req.BackgroundColor != "" {
		opts.BackgroundColor = req.BackgroundColor
	}

	svgStr, err := barcodes.GenerateSVG(btype, data, opts)
	if err != nil {
		return BarcodeResult{
			Type:      string(btype),
			Data:      data,
			Success:   false,
			Error:     err.Error(),
			ErrorInfo: errorInfo(err),
		}, nil
	}

	pngBytes, err := barcodes.GeneratePNG(btype, data, opts)
	var pngDataURI string
	if err == nil && len(pngBytes) > 0 {
		pngDataURI = "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes)
	}

	return BarcodeResult{
		Type:     string(btype),
		Data:     data,
		SVG:      svgStr,
		PNGData:  pngDataURI,
		Success:  true,
		ReadBack: readBack(btype, data, pngBytes),
	}, nil
}

// readBack decodes the generated PNG and compares it with the input.
func readBack(btype barcodes.BarcodeType, data string, png []byte) string {
	if btype == barcodes.TypePDF417 {
		return "unsupported"
	}
	want := data
	if barcodes.IsRetail(btype) {
		want = barcodes.CheckRetail(btype, data).Code
	}
	found, _, err := barcodes.DecodeImageBytes(png)
	if err != nil || len(found) == 0 {
		return "unreadable"
	}
	for _, d := range found {
		if d.Type == btype && d.Text == want {
			return "ok"
		}
	}
	return "mismatch"
}

// DecodeImage reads the barcodes in an image (PNG, JPEG, GIF, WebP), given
// as base64 or data URI, and splits known payloads into fields.
func (a *BarcodeApp) DecodeImage(imageBase64 string) DecodeImageResult {
	if i := strings.Index(imageBase64, ","); strings.HasPrefix(imageBase64, "data:") && i > 0 {
		imageBase64 = imageBase64[i+1:]
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(imageBase64))
	if err != nil {
		return DecodeImageResult{Error: "invalid image data", ErrorInfo: ErrorInfo{ErrorCode: barcodes.ErrImageFormat}}
	}
	found, size, err := barcodes.DecodeImageBytes(raw)
	res := DecodeImageResult{Width: size.X, Height: size.Y, Codes: []DecodedCode{}}
	if err != nil {
		res.Error, res.ErrorInfo = err.Error(), errorInfo(err)
		return res
	}
	for _, d := range found {
		c := DecodedCode{Type: string(d.Type), Text: d.Text}
		if p := qrformats.Parse(d.Text); p.Kind != "text" {
			c.Content = &p
		}
		for _, pt := range d.Points {
			c.Points = append(c.Points, Point{X: pt.X, Y: pt.Y})
		}
		res.Codes = append(res.Codes, c)
	}
	res.Success = true
	return res
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
	opts.FontSize = req.FontSize
	if req.ForegroundColor != "" {
		opts.ForegroundColor = req.ForegroundColor
	}
	if req.BackgroundColor != "" {
		opts.BackgroundColor = req.BackgroundColor
	}

	var results []BatchItemResult
	validCount := 0
	errorCount := 0

	for i, line := range req.Lines {
		data := strings.TrimSpace(line)
		itemIndex := i + 1

		if data == "" {
			continue
		}

		svgStr, err := barcodes.GenerateSVG(btype, data, opts)
		if err != nil {
			results = append(results, BatchItemResult{
				Index:     itemIndex,
				Data:      data,
				Success:   false,
				Error:     err.Error(),
				ErrorInfo: errorInfo(err),
			})
			errorCount++
			continue
		}

		pngBytes, err := barcodes.GeneratePNG(btype, data, opts)
		var pngDataURI string
		if err == nil && len(pngBytes) > 0 {
			pngDataURI = "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes)
		}

		results = append(results, BatchItemResult{
			Index:   itemIndex,
			Data:    data,
			SVG:     svgStr,
			PNGData: pngDataURI,
			Success: true,
		})
		validCount++
	}

	return BatchBarcodeResponse{
		Items:      results,
		Total:      len(results),
		ValidCount: validCount,
		ErrorCount: errorCount,
	}, nil
}

// PickExportFolder opens a native folder selection dialog.
func (a *BarcodeApp) PickExportFolder() (string, error) {
	dialog := application.Get().Dialog.OpenFile()
	dialog.SetMessage("Zielordner für Barcode-Dateien auswählen")
	dialog.CanChooseDirectories(true)
	dialog.CanChooseFiles(false)

	folderPath, err := dialog.PromptForSingleSelection()
	if err != nil {
		return "", err
	}
	return folderPath, nil
}

// ExportBatchToFolder exports generated batch barcodes into a target directory.
func (a *BarcodeApp) ExportBatchToFolder(req BatchExportRequest) (BatchExportResponse, error) {
	if req.FolderPath == "" {
		return BatchExportResponse{Error: "Zielordner ist nicht angegeben"}, nil
	}

	if err := os.MkdirAll(req.FolderPath, 0755); err != nil {
		return BatchExportResponse{Error: fmt.Sprintf("Ordner konnte nicht erstellt werden: %v", err)}, nil
	}

	ext := ".png"
	if strings.ToLower(req.Format) == "svg" {
		ext = ".svg"
	}

	var saved []string
	skipped := 0
	exported := 0

	for i, item := range req.Items {
		if !item.Success {
			skipped++
			continue
		}

		var filename string
		switch req.NamingScheme {
		case "index":
			filename = fmt.Sprintf("%s%03d%s", req.Prefix, i+1, ext)
		case "data_raw":
			slug := sanitizeFilename(item.Data)
			filename = fmt.Sprintf("%s%s%s", req.Prefix, slug, ext)
		case "data_slug":
			fallthrough
		default:
			slug := sanitizeFilename(item.Data)
			filename = fmt.Sprintf("%s%03d_%s%s", req.Prefix, i+1, slug, ext)
		}

		outPath := filepath.Join(req.FolderPath, filename)

		var writeErr error
		if ext == ".svg" && item.SVG != "" {
			writeErr = os.WriteFile(outPath, []byte(item.SVG), 0644)
		} else if ext == ".png" && item.PNGData != "" {
			// Extract base64 payload from data URI
			dataURI := item.PNGData
			if idx := strings.Index(dataURI, ","); idx != -1 {
				dataURI = dataURI[idx+1:]
			}
			pngBytes, decodeErr := base64.StdEncoding.DecodeString(dataURI)
			if decodeErr == nil {
				writeErr = os.WriteFile(outPath, pngBytes, 0644)
			} else {
				writeErr = decodeErr
			}
		}

		if writeErr != nil {
			skipped++
		} else {
			exported++
			saved = append(saved, outPath)
		}
	}

	return BatchExportResponse{
		ExportedCount: exported,
		SkippedCount:  skipped,
		SavedPaths:    saved,
	}, nil
}

// SaveSingleFile opens a native save-file dialog and writes the given SVG or PNG content.
func (a *BarcodeApp) SaveSingleFile(req SaveSingleFileRequest) (string, error) {
	dialog := application.Get().Dialog.SaveFile()
	dialog.SetMessage("Barcode speichern unter...")
	if req.Format == "svg" {
		dialog.AddFilter("SVG Vektorgrafik (*.svg)", "*.svg")
		dialog.SetFilename(req.DefaultName)
	} else {
		dialog.AddFilter("PNG Bild (*.png)", "*.png")
		dialog.SetFilename(req.DefaultName)
	}

	targetPath, err := dialog.PromptForSingleSelection()
	if err != nil {
		return "", err
	}
	if targetPath == "" {
		return "", nil // User cancelled
	}

	if req.Format == "svg" {
		if err := os.WriteFile(targetPath, []byte(req.Content), 0644); err != nil {
			return "", fmt.Errorf("failed to write svg file: %w", err)
		}
	} else {
		// PNG data URI
		dataURI := req.Content
		if idx := strings.Index(dataURI, ","); idx != -1 {
			dataURI = dataURI[idx+1:]
		}
		pngBytes, err := base64.StdEncoding.DecodeString(dataURI)
		if err != nil {
			return "", fmt.Errorf("failed to decode png base64: %w", err)
		}
		if err := os.WriteFile(targetPath, pngBytes, 0644); err != nil {
			return "", fmt.Errorf("failed to write png file: %w", err)
		}
	}

	return targetPath, nil
}

// CopyToClipboard copies text (e.g. raw SVG or string) to the native OS clipboard.
func (a *BarcodeApp) CopyToClipboard(text string) bool {
	app := application.Get()
	if app == nil || app.Clipboard == nil {
		return false
	}
	return app.Clipboard.SetText(text)
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

func sanitizeFilename(s string) string {
	clean := nonAlphanumericRegex.ReplaceAllString(s, "_")
	clean = strings.Trim(clean, "_")
	if len(clean) > 40 {
		clean = clean[:40]
	}
	if clean == "" {
		clean = "item"
	}
	return clean
}
