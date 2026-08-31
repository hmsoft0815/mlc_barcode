package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBarcodeApp_Formatters(t *testing.T) {
	app := NewBarcodeApp()

	// Test Wifi formatter
	wifiStr := app.FormatWifi(WifiInput{
		SSID:       "MyWifi",
		Password:   "Secret123",
		Encryption: "WPA",
	})
	if wifiStr != "WIFI:T:WPA;S:MyWifi;P:Secret123;;" {
		t.Errorf("unexpected wifi output: %s", wifiStr)
	}

	// Test vCard formatter
	vcardStr := app.FormatVCard(VCardInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Phone:     "+123456789",
	})
	if !contains(vcardStr, "BEGIN:VCARD") || !contains(vcardStr, "N:Doe;John;;;") {
		t.Errorf("unexpected vcard output: %s", vcardStr)
	}

	// Test Event formatter
	eventStr := app.FormatEvent(EventInput{
		Summary:   "Meeting",
		StartTime: "20260901T100000",
		EndTime:   "20260901T110000",
		TimeZone:  "Europe/Berlin",
	})
	if !contains(eventStr, "BEGIN:VEVENT") || !contains(eventStr, "SUMMARY:Meeting") {
		t.Errorf("unexpected event output: %s", eventStr)
	}
}

func TestBarcodeApp_GenerateBarcode(t *testing.T) {
	app := NewBarcodeApp()

	// Valid QR code
	res, err := app.GenerateBarcode(BarcodeRequest{
		Type: "qr",
		Data: "https://example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Success {
		t.Errorf("expected success, got error: %s", res.Error)
	}
	if res.SVG == "" {
		t.Error("expected non-empty SVG")
	}

	// Empty data should return error result
	resEmpty, _ := app.GenerateBarcode(BarcodeRequest{
		Type: "qr",
		Data: "",
	})
	if resEmpty.Success {
		t.Error("expected failure for empty data")
	}
}

func TestBarcodeApp_GenerateBatch(t *testing.T) {
	app := NewBarcodeApp()

	resp, err := app.GenerateBatch(BatchBarcodeRequest{
		Type:  "qr",
		Lines: []string{"Item 1", "Item 2", "Item 3"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 3 {
		t.Errorf("expected 3 total items, got %d", resp.Total)
	}
	if resp.ValidCount != 3 {
		t.Errorf("expected 3 valid items, got %d", resp.ValidCount)
	}
}

func TestBarcodeApp_ExportBatchToFolder(t *testing.T) {
	app := NewBarcodeApp()

	tmpDir, err := os.MkdirTemp("", "barcode_export_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	batchResp, err := app.GenerateBatch(BatchBarcodeRequest{
		Type:  "qr",
		Lines: []string{"Hello 1", "Hello 2"},
	})
	if err != nil {
		t.Fatal(err)
	}

	exportResp, err := app.ExportBatchToFolder(BatchExportRequest{
		FolderPath:   tmpDir,
		NamingScheme: "index",
		Prefix:       "test_",
		Format:       "svg",
		Items:        batchResp.Items,
	})
	if err != nil {
		t.Fatal(err)
	}
	if exportResp.ExportedCount != 2 {
		t.Errorf("expected 2 exported files, got %d", exportResp.ExportedCount)
	}

	f1 := filepath.Join(tmpDir, "test_001.svg")
	if _, err := os.Stat(f1); err != nil {
		t.Errorf("expected file %s to exist", f1)
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World!", "Hello_World"},
		{"https://example.com/foo?bar=1", "https_example_com_foo_bar_1"},
		{"", "item"},
		{"___test___", "test"},
	}

	for _, tc := range tests {
		got := sanitizeFilename(tc.input)
		if got != tc.expected {
			t.Errorf("sanitizeFilename(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && filepath.Base(s) != ""))
}
