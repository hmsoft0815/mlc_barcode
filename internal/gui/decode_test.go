package gui

import (
	"strings"
	"testing"

	"github.com/mlcmcp/mlc_barcode/internal/qrformats"
)

func TestReadBack(t *testing.T) {
	a := &BarcodeApp{}
	for _, tt := range []struct {
		req  BarcodeRequest
		want string
	}{
		{BarcodeRequest{Type: "qr", Data: "https://mlcgo.eu"}, "ok"},
		{BarcodeRequest{Type: "ean13", Data: "400638133393"}, "ok"}, // check digit completed
		{BarcodeRequest{Type: "upca", Data: "036000291452"}, "ok"},
		{BarcodeRequest{Type: "pdf417", Data: "BOARDING PASS"}, "unsupported"},
		// light grey on white: a scanner cannot read it — the GUI must warn
		{BarcodeRequest{Type: "qr", Data: "kaum sichtbar", ForegroundColor: "#f4f4f4", BackgroundColor: "#ffffff"}, "unreadable"},
	} {
		res, err := a.GenerateBarcode(tt.req)
		if err != nil || !res.Success {
			t.Fatalf("%+v: %v %s", tt.req, err, res.Error)
		}
		if res.ReadBack != tt.want {
			t.Errorf("%s %q: readBack %q, want %q", tt.req.Type, tt.req.Data, res.ReadBack, tt.want)
		}
	}
}

func TestDecodeImage(t *testing.T) {
	a := &BarcodeApp{}
	gen, _ := a.GenerateBarcode(BarcodeRequest{Type: "qr", Data: qrformats.FormatEPC(qrformats.EPCOptions{
		Name: "Muster GmbH", IBAN: "DE89370400440532013000", Amount: 12.5, Reference: "R-1"})})

	res := a.DecodeImage(gen.PNGData) // data URI straight from the generator
	if !res.Success || len(res.Codes) != 1 {
		t.Fatalf("decode: %+v", res)
	}
	c := res.Codes[0]
	if c.Type != "qr" || c.Content == nil || c.Content.Kind != "epc" || c.Content.Fields["iban_valid"] != "true" {
		t.Errorf("unexpected result %+v", c)
	}
	if res.Width == 0 || len(c.Points) == 0 {
		t.Errorf("size or points missing: %dx%d %v", res.Width, res.Height, c.Points)
	}

	bad := a.DecodeImage("data:image/png;base64," + strings.Repeat("A", 40))
	if bad.Success || bad.ErrorCode != "image_format" {
		t.Errorf("garbage accepted: %+v", bad)
	}
}
