package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hmsoft0815/mlcartifact/client"
	"github.com/mlcmcp/mlc_barcode/internal/barcodes"
	"github.com/mlcmcp/mlc_barcode/internal/qrformats"
	"github.com/mlcmcp/mlc_barcode/internal/version"
)

func main() {
	showVersion := flag.Bool("version", false, "Show version and exit")
	btype := flag.String("type", "qr", "Barcode type (qr, datamatrix, code128, code39, ean13, ean8, upca, itf)")
	data := flag.String("data", "", "Data to encode")
	customText := flag.String("custom-text", "", "Custom caption text to display below barcode")
	output := flag.String("out", "barcode.svg", "Output filename (.svg or .png)")
	width := flag.Int("width", 0, "Width of the barcode (0 for default)")
	height := flag.Int("height", 0, "Height of the barcode (0 for default)")
	showText := flag.Bool("text", false, "Show text below barcode (if supported)")
	fgColor := flag.String("fg", "black", "Foreground color (e.g. black, #ff0000)")
	bgColor := flag.String("bg", "white", "Background color (e.g. white, transparent, #ffffff)")

	// Structured QR Flags: WIFI
	wifiSSID := flag.String("wifi-ssid", "", "WIFI SSID (triggers WIFI QR)")
	wifiPass := flag.String("wifi-pass", "", "WIFI Password")
	wifiEnc := flag.String("wifi-enc", "WPA", "WIFI Encryption (WPA, WEP, nopass)")
	wifiHidden := flag.Bool("wifi-hidden", false, "WIFI Hidden SSID")

	// Structured QR Flags: vCard
	vcardFirst := flag.String("vcard-first", "", "vCard First Name (triggers vCard QR)")
	vcardLast := flag.String("vcard-last", "", "vCard Last Name")
	vcardEmail := flag.String("vcard-email", "", "vCard Email")
	vcardTel := flag.String("vcard-tel", "", "vCard Phone")

	// Structured QR Flags: Event
	eventSummary := flag.String("event-summary", "", "Event Summary (triggers Event QR)")
	eventStart := flag.String("event-start", "", "Event Start Time (YYYYMMDDTHHMMSS)")
	eventEnd := flag.String("event-end", "", "Event End Time")
	eventTZ := flag.String("event-tz", "", "Event TimeZone (e.g. Europe/Berlin)")

	// Structured QR Flags: EPC / GiroCode (SEPA-Überweisung)
	epcName := flag.String("epc-name", "", "EPC/GiroCode Beneficiary Name (triggers GiroCode QR)")
	epcIBAN := flag.String("epc-iban", "", "EPC/GiroCode IBAN")
	epcBIC := flag.String("epc-bic", "", "EPC/GiroCode BIC")
	epcAmount := flag.Float64("epc-amount", 0, "EPC/GiroCode Amount in EUR (e.g. 12.50)")
	epcRef := flag.String("epc-ref", "", "EPC/GiroCode Reference / Verwendungszweck")

	// Structured QR Flags: Crypto Wallet
	cryptoCoin := flag.String("crypto-coin", "", "Crypto Coin (btc, eth, sol, etc., triggers Crypto QR)")
	cryptoAddr := flag.String("crypto-addr", "", "Crypto Wallet Address")
	cryptoAmount := flag.Float64("crypto-amount", 0, "Crypto Payment Amount")
	cryptoMsg := flag.String("crypto-msg", "", "Crypto Payment Message / Memo")

	// Structured QR Flags: Geo Location (Maps)
	geoLat := flag.Float64("geo-lat", 0, "Geo Latitude (e.g. 52.5200, triggers Geo QR)")
	geoLon := flag.Float64("geo-lon", 0, "Geo Longitude (e.g. 13.4050)")
	geoQuery := flag.String("geo-query", "", "Geo Location Query / Name (e.g. 'Berlin')")

	// Structured QR Flags: Phone & SMS
	telNumber := flag.String("tel", "", "Telephone Number (triggers tel: QR)")
	smsNumber := flag.String("sms", "", "SMS Recipient Number (triggers smsto: QR)")
	smsText := flag.String("sms-text", "", "SMS Message Text")

	// Structured QR Flags: Email
	mailTo := flag.String("mail-to", "", "Email Recipient (triggers mailto: QR)")
	mailSubject := flag.String("mail-subject", "", "Email Subject Line")
	mailBody := flag.String("mail-body", "", "Email Body Text")

	// Artifact options
	saveArtifact := flag.Bool("artifact", false, "Save generated barcode to mlcartifact service")
	artifactAddr := flag.String("artifact-addr", os.Getenv("ARTIFACT_GRPC_ADDR"), "Address of the mlcartifact gRPC server")

	flag.Parse()

	if *showVersion {
		fmt.Printf("MLC Barcode CLI v%s\nAuthor: %s\n", version.Version, version.Author)
		return
	}

	dataStr := *data

	// Structured QR overrides
	if *epcIBAN != "" || *epcName != "" {
		dataStr = qrformats.FormatEPC(qrformats.EPCOptions{
			Name:      *epcName,
			IBAN:      *epcIBAN,
			BIC:       *epcBIC,
			Amount:    *epcAmount,
			Reference: *epcRef,
		})
	} else if *cryptoAddr != "" || *cryptoCoin != "" {
		dataStr = qrformats.FormatCrypto(qrformats.CryptoOptions{
			Coin:    *cryptoCoin,
			Address: *cryptoAddr,
			Amount:  *cryptoAmount,
			Message: *cryptoMsg,
		})
	} else if *geoLat != 0 || *geoLon != 0 || *geoQuery != "" {
		dataStr = qrformats.FormatGeo(qrformats.GeoOptions{
			Latitude:  *geoLat,
			Longitude: *geoLon,
			Query:     *geoQuery,
		})
	} else if *telNumber != "" {
		dataStr = qrformats.FormatTel(qrformats.TelOptions{
			PhoneNumber: *telNumber,
		})
	} else if *smsNumber != "" {
		dataStr = qrformats.FormatSMS(qrformats.SMSOptions{
			PhoneNumber: *smsNumber,
			Message:     *smsText,
		})
	} else if *mailTo != "" {
		dataStr = qrformats.FormatEmail(qrformats.EmailOptions{
			To:      *mailTo,
			Subject: *mailSubject,
			Body:    *mailBody,
		})
	} else if *wifiSSID != "" {
		dataStr = qrformats.FormatWifi(qrformats.WifiOptions{
			SSID:       *wifiSSID,
			Password:   *wifiPass,
			Encryption: *wifiEnc,
			Hidden:     *wifiHidden,
		})
	} else if *vcardFirst != "" || *vcardLast != "" {
		dataStr = qrformats.FormatVCard(qrformats.VCardOptions{
			FirstName: *vcardFirst,
			LastName:  *vcardLast,
			Email:     *vcardEmail,
			Phone:     *vcardTel,
		})
	} else if *eventSummary != "" {
		dataStr = qrformats.FormatVCalendar(qrformats.VCalendarOptions{
			Summary:   *eventSummary,
			StartTime: *eventStart,
			EndTime:   *eventEnd,
			TimeZone:  *eventTZ,
		})
	}

	if dataStr == "" {
		fmt.Println("Error: Data is required (or use structured QR flags like -wifi-ssid, -epc-iban, -crypto-addr, -geo-lat, -tel, -sms, -mail-to)")
		flag.Usage()
		os.Exit(1)
	}

	barcodeType := barcodes.BarcodeType(strings.ToLower(*btype))
	opts := barcodes.DefaultOptions(barcodeType)

	if *width > 0 {
		opts.Width = *width
	}
	if *height > 0 {
		opts.Height = *height
	}
	opts.ShowText = *showText
	opts.CustomText = *customText
	opts.ForegroundColor = *fgColor
	opts.BackgroundColor = *bgColor

	ext := strings.ToLower(filepath.Ext(*output))

	var err error
	var content []byte
	var mimeType string

	switch ext {
	case ".svg":
		var svg string
		svg, err = barcodes.GenerateSVG(barcodeType, dataStr, opts)
		content = []byte(svg)
		mimeType = "image/svg+xml"
	case ".png":
		content, err = barcodes.GeneratePNG(barcodeType, dataStr, opts)
		mimeType = "image/png"
	default:
		err = fmt.Errorf("unsupported output format: %s", ext)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating barcode: %v\n", err)
		os.Exit(1)
	}

	// Save to file
	err = os.WriteFile(*output, content, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s: %s\n", barcodeType, *output)

	// Optional: Save to artifact service
	if *saveArtifact {
		addr := *artifactAddr
		if addr == "" {
			addr = ":9590" // default
		}
		c, err := client.NewClientWithAddr(addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not connect to artifact server at %s: %v\n", addr, err)
		} else {
			resp, err := c.Write(context.Background(), *output, content,
				client.WithMimeType(mimeType),
				client.WithDescription(fmt.Sprintf("Generated %s barcode for: %s", barcodeType, dataStr)))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving artifact: %v\n", err)
			} else {
				fmt.Printf("Artifact successfully saved (ID: %s)\n", resp.Id)
			}
			c.Close()
		}
	}
}
