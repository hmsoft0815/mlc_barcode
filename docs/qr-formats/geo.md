# Geo Location QR Format

Opens a set of coordinates in the phone's map application — Google Maps, Apple Maps or OpenStreetMap, whichever the device is set to.

Standard: **RFC 5870** (`geo:` URI).

## Parameters

| CLI Flag | MCP Parameter | Description |
|----------|---------------|-------------|
| `-geo-lat` | `latitude` | **Required.** Latitude, e.g. `52.5200`. |
| `-geo-lon` | `longitude` | **Required.** Longitude, e.g. `13.4050`. |
| `-geo-query` | `query` | Optional place name, appended as `?q=`. |

## Usage Examples

### CLI Example
```bash
./bin/barcode -geo-lat 52.5200 -geo-lon 13.4050 -geo-query "Brandenburger Tor" -out ort.png
```

### MCP Tool Example
```json
{
  "name": "generate_geo_qr",
  "arguments": {
    "latitude": 52.5200,
    "longitude": 13.4050,
    "query": "Brandenburger Tor"
  }
}
```

## Notes

Coordinates are decimal degrees, with a dot as the decimal separator and a leading minus for south and west. Degrees-minutes-seconds is not accepted.

When `query` is set, most applications show the name and centre on it; the coordinates remain what actually locates the point.
