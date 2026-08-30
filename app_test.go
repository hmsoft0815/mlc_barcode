package main

import (
	"os"
	"testing"
)

func TestBarcodeApp_Formatters(t *testing.T) {
	app := NewBarcodeApp()

	// Test Wifi
	wifiStr := app.FormatWifi(WifiInput{
		SSID:       "MyWiFi",
		Password:   "secret123",
		Encryption: "WPA",
	})
	if wifiStr != "WIFI:T:WPA;S:MyWiFi;P:secret123;;" {
		t.Errorf("unexpected wifi format: %s", wifiStr)
	}

	// Test vCard
	vcardStr := app.FormatVCard(VCardInput{
		FirstName: "Max",
		LastName:  "Mustermann",
		Email:     "max@beispiel.de",
		Phone:     "+4912345678",
	})
	if vcardStr == "" || len(vcardStr) < 20 {
		t.Errorf("unexpected vcard format: %s", vcardStr)
	}

	// Test Event
	eventStr := app.FormatEvent(EventInput{
		Summary:   "Meeting",
		StartTime: "20260901T100000",
		EndTime:   "20260901T110000",
		TimeZone:  "Europe/Berlin",
	})
	if eventStr == "" || len(eventStr) < 20 {
		t.Errorf("unexpected event format: %s", eventStr)
	}
}

func TestBarcodeApp_GenerateBarcode(t *testing.T) {
	app := NewBarcodeApp()

	// Valid QR code
	res, err := app.GenerateBarcode(BarcodeRequest{
		Type: "qr",
		Data: "https://mlcgo.eu",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Success || res.SVG == "" || res.PNGData == "" {
		t.Fatalf("expected successful barcode generation, got: %+v", res)
	}

	// Valid EAN13
	resEan, err := app.GenerateBarcode(BarcodeRequest{
		Type: "ean13",
		Data: "4012345678901",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resEan.Success || resEan.SVG == "" {
		t.Fatalf("expected valid EAN13 generation, got: %+v", resEan)
	}

	// Empty data should return failure
	resEmpty, err := app.GenerateBarcode(BarcodeRequest{
		Type: "qr",
		Data: "",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resEmpty.Success {
		t.Errorf("expected failure for empty data")
	}

	// Invalid barcode data
	resInvalid, err := app.GenerateBarcode(BarcodeRequest{
		Type: "ean13",
		Data: "INVALID_LETTERS",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resInvalid.Success {
		t.Errorf("expected failure for invalid ean13 data")
	}
}

func TestBarcodeApp_GenerateBatchAndExport(t *testing.T) {
	app := NewBarcodeApp()

	batchRes, err := app.GenerateBatch(BatchBarcodeRequest{
		Type: "code128",
		Lines: []string{
			"ITEM-001",
			"ITEM-002",
			"ITEM-003",
		},
	})
	if err != nil {
		t.Fatalf("unexpected batch error: %v", err)
	}
	if batchRes.Total != 3 || batchRes.ValidCount != 3 || batchRes.ErrorCount != 0 {
		t.Fatalf("unexpected batch counts: %+v", batchRes)
	}

	// Test Export to temp folder
	tmpDir, err := os.MkdirTemp("", "mlc_barcode_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	exportRes, err := app.ExportBatchToFolder(BatchExportRequest{
		FolderPath:   tmpDir,
		NamingScheme: "data_slug",
		Prefix:       "test_",
		Format:       "svg",
		Items:        batchRes.Items,
	})
	if err != nil {
		t.Fatalf("unexpected export error: %v", err)
	}
	if exportRes.ExportedCount != 3 {
		t.Errorf("expected 3 exported files, got %d", exportRes.ExportedCount)
	}

	// Check that files exist
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to read temp dir: %v", err)
	}
	if len(files) != 3 {
		t.Errorf("expected 3 files on disk, found %d", len(files))
	}
}

func TestSanitizeFilename(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"https://mlcgo.eu/test?a=1&b=2", "https_mlcgo_eu_test_a_1_b_2"},
		{"  Hello World!  ", "Hello_World"},
		{"", "barcode"},
	}

	for _, c := range cases {
		actual := sanitizeFilename(c.input)
		if actual != c.expected {
			t.Errorf("sanitizeFilename(%q) = %q, expected %q", c.input, actual, c.expected)
		}
	}
}
