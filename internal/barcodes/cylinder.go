package barcodes

import (
	"image"
	"math"
)

// Codes on bottles and tubes (T-20260926-12). The 2D detectors sample on a
// perspective grid; on a cylinder the modules narrow along a cosine towards
// the edges, and when the camera looks down or up the rows bend into arcs.
// Off-centre the samples then miss the middle modules, and a code that is
// perfectly sharp is not read (see docs/decoder.md, testdata/cylinder).
//
// decodeCurved is the fallback for that: the spots that look like a 2D
// code are unwrapped under a few assumed bottle shapes and read again. It
// runs only when nothing else in the image was read.

// guessesPerCandidate says how many of the cylinderGuesses each of the
// largest 2D-code spots gets, largest first: a still photo can afford
// about a second when nothing else was read; a camera frame only the most
// likely shapes of the largest spot — the next frame comes a moment later.
var (
	guessesPerCandidate     = []int{30, 9, 9}
	guessesPerCandidateLive = []int{12}
)

// cylinderGuess is one assumed shape: the code covers span degrees of the
// circumference, its centre is turned by turn degrees away from the camera,
// and the camera looks down on the bottle by tilt degrees (negative: up).
type cylinderGuess struct{ span, turn, tilt float64 }

// cylinderGuesses in the order they are tried: the common cases first
// (code facing the camera, camera slightly above), then turned bottles.
// Measured on testdata/cylinder, the hits come from the first twelve and
// from the ±35° turns.
var cylinderGuesses = func() []cylinderGuess {
	var g []cylinderGuess
	for _, turn := range []float64{0, 20, -20, 35, -35} {
		for _, span := range []float64{70, 110} {
			for _, tilt := range []float64{15, 0, -15} {
				g = append(g, cylinderGuess{span, turn, tilt})
			}
		}
	}
	return g
}()

// surfaceSpan is the curvature behind a Surface hint: q = ½ of the
// bottle's width gives 2·asin(½) = 60°; a code wrapping a thin tube about
// 100° (q ≈ ¾).
var surfaceSpan = map[Surface]float64{SurfaceBottle: 60, SurfaceTube: 100}

// surfaceGuesses are the turns and tilts of one known curvature, in the
// order of cylinderGuesses.
func surfaceGuesses(span float64) []cylinderGuess {
	var g []cylinderGuess
	for _, turn := range []float64{0, 20, -20, 35, -35} {
		for _, tilt := range []float64{15, 0, -15} {
			g = append(g, cylinderGuess{span, turn, tilt})
		}
	}
	return g
}

func decodeCurved(img *image.Gray, fr frame, opts DecodeOptions) []Decoded {
	guesses, budget := cylinderGuesses, guessesPerCandidate
	if opts.Live {
		budget = guessesPerCandidateLive
	}
	switch opts.Surface {
	case SurfaceFlat:
		return nil
	case SurfaceBottle, SurfaceTube:
		// The user knows the curvature: every turn and tilt of that one,
		// which a camera frame can afford in full.
		guesses = surfaceGuesses(surfaceSpan[opts.Surface])
		if opts.Live {
			budget = []int{len(guesses)}
		}
	}
	boxes := matrixCandidates(img)
	for i, box := range boxes[:min(len(boxes), len(budget))] {
		win := box.Inset(-(max(box.Dx(), box.Dy())/3 + 8)).Intersect(img.Rect)
		for _, g := range guesses[:min(len(guesses), budget[i])] {
			c := fitCylinder(box, win, g)
			flat := c.unwrap(img)
			if flat == nil {
				continue
			}
			mapper := shifted{c.mapper(fr), -quietZone, -quietZone}
			if found := decodeOnce(withQuietZone(flat), mapper, curvedReaders); len(found) > 0 {
				return found
			}
		}
	}
	return nil
}

// cylinder maps between a window of the image and its unwrapped version.
// In the image a point of the surface at angle theta lies at
// x = cx + r·sin(theta); seen from above it drops by
// r·(cos(theta) − cos(turn))·sin(tilt) relative to the code's centre.
type cylinder struct {
	cx, r, turn, tilt float64
	thetaL, thetaR    float64 // angles of the window's left and right edge
	top, height       int     // the window's rows
}

// fitCylinder places the guessed shape so the code spans the busy box.
func fitCylinder(box, win image.Rectangle, g cylinderGuess) cylinder {
	rad := math.Pi / 180
	a, b := (g.turn-g.span/2)*rad, (g.turn+g.span/2)*rad
	r := float64(box.Dx()) / (math.Sin(b) - math.Sin(a))
	cx := float64(box.Min.X) - r*math.Sin(a)
	edge := func(x int) float64 {
		s := (float64(x) - cx) / r
		return math.Asin(math.Max(-maxSin, math.Min(maxSin, s)))
	}
	return cylinder{
		cx: cx, r: r, turn: g.turn * rad, tilt: g.tilt * rad,
		thetaL: edge(win.Min.X), thetaR: edge(win.Max.X),
		top: win.Min.Y, height: win.Dy(),
	}
}

// maxSin keeps the window off the silhouette of the bottle, where the
// surface is seen edge-on and one pixel covers many modules.
const maxSin = 0.97

// source is the image point of the unwrapped point (u, v).
func (c cylinder) source(u, v float64) (float64, float64) {
	theta := c.thetaL + u/c.r
	x := c.cx + c.r*math.Sin(theta)
	y := float64(c.top) + v + c.r*(math.Cos(theta)-math.Cos(c.turn))*math.Sin(c.tilt)
	return x, y
}

// unwrap renders the window flat: columns at equal arc length, rows
// straightened. Pixels outside the image are white (paper).
func (c cylinder) unwrap(img *image.Gray) *image.Gray {
	w := int(c.r * (c.thetaR - c.thetaL))
	if w < 16 || c.height < 16 || w > 4*img.Rect.Dx() {
		return nil
	}
	out := image.NewGray(image.Rect(0, 0, w, c.height))
	// x and the row shift depend only on the column: one sine and cosine
	// per column, not per pixel.
	b := img.Rect
	at := func(px, py int) float64 {
		if px < b.Min.X || py < b.Min.Y || px >= b.Max.X || py >= b.Max.Y {
			return 255
		}
		return float64(img.Pix[img.PixOffset(px, py)])
	}
	for u := 0; u < w; u++ {
		x, y0 := c.source(float64(u)+0.5, 0.5)
		x0 := int(math.Floor(x - 0.5))
		fx := x - 0.5 - float64(x0)
		for v := 0; v < c.height; v++ {
			y := y0 + float64(v) - 0.5
			yi := int(math.Floor(y))
			fy := y - float64(yi)
			top := at(x0, yi)*(1-fx) + at(x0+1, yi)*fx
			bottom := at(x0, yi+1)*(1-fx) + at(x0+1, yi+1)*fx
			out.Pix[v*out.Stride+u] = uint8(top*(1-fy) + bottom*fy + 0.5)
		}
	}
	return out
}

func (c cylinder) mapper(fr frame) pointMapper { return cylinderMapper{c, fr} }

type cylinderMapper struct {
	c  cylinder
	fr frame
}

func (m cylinderMapper) point(u, v float64) image.Point {
	return m.fr.point(m.c.source(u, v))
}

// shifted moves points by (dx, dy) before mapping them, e.g. to undo the
// quiet zone added around an image.
type shifted struct {
	m      pointMapper
	dx, dy float64
}

func (s shifted) point(x, y float64) image.Point { return s.m.point(x+s.dx, y+s.dy) }
