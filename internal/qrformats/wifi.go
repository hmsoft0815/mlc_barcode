package qrformats

import (
	"fmt"
	"strings"
)

// WifiOptions holds configuration for a WIFI QR code
type WifiOptions struct {
	SSID       string
	Password   string
	Encryption string // WPA, WEP, or nopass
	Hidden     bool
}

// FormatWifi returns a formatted WIFI string: WIFI:T:WPA;S:SSID;P:PASSWORD;H:true;;
func FormatWifi(opts WifiOptions) string {
	enc := opts.Encryption
	if enc == "" {
		enc = "WPA"
	}
	
	hidden := ""
	if opts.Hidden {
		hidden = "H:true;"
	}

	// Escape special characters in SSID and Password if needed (;, :, \)
	ssid := escapeWifi(opts.SSID)
	pass := escapeWifi(opts.Password)

	return fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;%s;", enc, ssid, pass, hidden)
}

func escapeWifi(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, ";", "\\;")
	s = strings.ReplaceAll(s, ":", "\\:")
	s = strings.ReplaceAll(s, ",", "\\,")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}
