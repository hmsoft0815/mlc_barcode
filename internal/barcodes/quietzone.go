package barcodes

import (
	"image"
	"image/color"
	"math"

	"github.com/boombuler/barcode"
)

// quietModules is the blank margin each symbology needs around it, in
// modules (x: left/right, y: top/bottom). Scanners find the code's edges
// by that margin; without it a code placed on a coloured or busy
// background often cannot be read.
func quietModules(btype BarcodeType) (x, y int) {
	switch btype {
	case TypeQR:
		return 4, 4 // ISO/IEC 18004
	case TypeDataMatrix:
		return 1, 1 // ISO/IEC 16022
	case TypeAztec:
		return 1, 1 // ISO/IEC 24778 needs none; one module costs little
	case TypePDF417:
		return 2, 2 // ISO/IEC 15438
	case TypeEAN13:
		return 11, 0 // GS1: 11 left, 7 right — the larger on both sides
	case TypeEAN8, TypeUPCA:
		return 9, 0 // GS1: 7 and 9
	default:
		return 10, 0 // Code 128, Code 39, ITF: at least 10 modules
	}
}

// fitWithQuietZone scales bc to the requested total size with the quiet
// zone inside it, so layouts built on a fixed size stay unchanged. Without
// a requested size one module is one pixel.
func fitWithQuietZone(bc barcode.Barcode, btype BarcodeType, opts BarcodeOptions) (barcode.Barcode, error) {
	qx, qy := 0, 0
	if !opts.NoQuietZone {
		qx, qy = quietModules(btype)
	}
	n := bc.Bounds()
	mx, my := qx, qy
	if opts.Width > 0 && opts.Height > 0 {
		mx = int(math.Round(float64(qx) * float64(opts.Width) / float64(n.Dx()+2*qx)))
		my = int(math.Round(float64(qy) * float64(opts.Height) / float64(n.Dy()+2*qy)))
		// barcode.Scale cannot shrink: a long Code 128 keeps its own width.
		var err error
		bc, err = barcode.Scale(bc, max(opts.Width-2*mx, n.Dx()), max(opts.Height-2*my, n.Dy()))
		if err != nil {
			return nil, err
		}
	}
	if mx == 0 && my == 0 {
		return bc, nil
	}
	return &margined{Barcode: bc, mx: mx, my: my}, nil
}

// margined adds a background-coloured border around a barcode. At reports
// white there, which the renderers treat as background.
type margined struct {
	barcode.Barcode
	mx, my int
}

func (m *margined) Bounds() image.Rectangle {
	b := m.Barcode.Bounds()
	return image.Rect(0, 0, b.Dx()+2*m.mx, b.Dy()+2*m.my)
}

func (m *margined) At(x, y int) color.Color {
	b := m.Barcode.Bounds()
	p := image.Pt(x-m.mx+b.Min.X, y-m.my+b.Min.Y)
	if !p.In(b) {
		return color.White
	}
	return m.Barcode.At(p.X, p.Y)
}
