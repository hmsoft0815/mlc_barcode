package qrformats

import (
	"fmt"
	"net/url"
	"strings"
)

// GeoOptions holds latitude, longitude and optional search query for Map QR codes.
// Standard: RFC 5870 (geo: URI scheme).
type GeoOptions struct {
	Latitude  float64
	Longitude float64
	Query     string // Optional location name or search query (e.g. "Eiffel Tower")
}

// FormatGeo returns an RFC 5870 geo URI (e.g. "geo:52.5200,13.4050" or "geo:52.5200,13.4050?q=Brandenburg+Gate").
func FormatGeo(opts GeoOptions) string {
	query := strings.TrimSpace(opts.Query)
	if query != "" {
		return fmt.Sprintf("geo:%f,%f?q=%s", opts.Latitude, opts.Longitude, url.QueryEscape(query))
	}
	return fmt.Sprintf("geo:%f,%f", opts.Latitude, opts.Longitude)
}
