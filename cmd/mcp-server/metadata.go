package main

import (
	"github.com/mlcmcp/mlc_barcode/internal/barcodes"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const serverInstructions = `Generates barcodes and QR codes as SVG (default) or PNG.

- Plain codes: generate_barcode with type qr, datamatrix, aztec, pdf417, code128, code39, ean13, ean8, upca or itf.
- Errors say which character or how much data is the problem and what to do instead (e.g. use code128 for lower case, aztec for more text) — correct the call accordingly.
- For structured QR payloads use the dedicated tool instead of hand-building the text: generate_epc_qr (SEPA transfer / GiroCode), generate_wifi_qr, generate_vcard_qr, generate_event_qr, generate_crypto_qr, generate_geo_qr, generate_tel_qr, generate_sms_qr, generate_email_qr.
- EAN-13, EAN-8 and UPC-A may be given without check digit; it is computed. A wrong check digit is rejected and the error names the right one.
- text:true adds a caption showing the encoded content; caption:"…" sets your own caption text instead (e.g. a product name under an EAN or a name under a vCard); font_size sets its size.
- structuredContent.encoded_data is exactly what the code contains, e.g. the EAN including check digit or the generated vCard.
- decode_barcode reads codes from an image (path or image_base64) and reports symbology and content — use it to verify a generated code or to read one the user provides.`

var toolTitles = map[string]string{
	"generate_barcode":   "Barcode or QR code",
	"generate_epc_qr":    "GiroCode (SEPA transfer)",
	"generate_wifi_qr":   "Wi-Fi access QR code",
	"generate_vcard_qr":  "Business card QR code (vCard)",
	"generate_event_qr":  "Calendar event QR code",
	"generate_crypto_qr": "Crypto payment QR code",
	"generate_geo_qr":    "Map location QR code",
	"generate_tel_qr":    "Phone call QR code",
	"generate_sms_qr":    "SMS QR code",
	"generate_email_qr":  "Email QR code",
}

// barcodeOutput is the structuredContent of every successful tool call.
type barcodeOutput struct {
	BarcodeType string `json:"barcode_type"`
	Format      string `json:"format"`
	MimeType    string `json:"mime_type"`
	EncodedData string `json:"encoded_data"`
	ArtifactID  string `json:"artifact_id,omitempty"`
}

var barcodeOutputSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"barcode_type": map[string]any{"type": "string", "description": "Symbology that was generated, e.g. qr or ean13"},
		"format":       map[string]any{"type": "string", "enum": []string{"svg", "png"}},
		"mime_type":    map[string]any{"type": "string", "description": "image/svg+xml or image/png"},
		"encoded_data": map[string]any{"type": "string", "description": "Exact payload inside the code (EAN/UPC including check digit, generated vCard/iCalendar/EPC text)"},
		"artifact_id":  map[string]any{"type": "string", "description": "ID in mlcartifact when save_artifact was set"},
	},
	"required": []string{"barcode_type", "format", "mime_type", "encoded_data"},
}

// addBarcodeTool registers a tool with its display title and the shared
// output schema.
func addBarcodeTool(s *mcp.Server, t *mcp.Tool, h mcp.ToolHandlerFor[map[string]any, any]) {
	t.Title = toolTitles[t.Name]
	t.OutputSchema = barcodeOutputSchema
	mcp.AddTool(s, t, h)
}

func barcodeTypeNames() []string {
	names := make([]string, len(barcodes.AllTypes))
	for i, t := range barcodes.AllTypes {
		names[i] = string(t)
	}
	return names
}
