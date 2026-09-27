# Barcode reader — the decode pipeline

How MLC Barcode finds and reads codes in an image, from the entry points
(CLI, MCP, GUI with camera) down to the individual ZXing readers. Read this
before changing anything in `internal/barcodes/decode*.go`; the tests named
at the end guard every stage.

## Overview

```
 CLI -decode <file>        MCP decode_barcode        GUI "Prüfen" tab
 (image.Decode)            (base64 / path)           file · drag & drop · Ctrl+V
        │                          │                 camera frames · region selection
        │                          └────────┬────────────────┘
        │                    DecodeImageBytes: size and format limits
        └──────────────────────────┬────────┘
                                   ▼
                         barcodes.Decode(img)
                                   │
   1. toGray              8-bit luminance; JPEG: the Y plane as is
   2. shrink              long side > 2000 px → integer box average
   3. decodeScene         on the reduced image
        a. whole image, all readers  (QR multi, PDF417 multi, DataMatrix,
                                      Aztec, EAN/UPC, Code 128, Code 39, ITF)
        b. 2D candidates   → a window per candidate, DataMatrix + Aztec
        c. 1D paint-over   → found 1D codes whited out, 1D readers again
   4. fallbacks, only while nothing was found:
        scale2x → invert → decodeCurved (codes on bottles)
        → decodeScene at full resolution
                                   │
                                   ▼
             []Decoded {Type, Text, Points in the caller's pixels}
                                   │
                    qrformats.Parse(text) → kind + fields
          (securPharm pack codes, GiroCode/EPC, vCard, Wi-Fi, event, …)
```

## Entry points

| Caller | Code | Input | Notes |
|---|---|---|---|
| CLI | `cmd/barcode/decode.go` | image file | `-decode <file>`, prints `type<TAB>content` and parsed fields |
| MCP | `cmd/mcp-server/tools_decode.go` | base64; a file path only when the server is started with paths allowed | tool `decode_barcode`; content parsed into structured output |
| GUI | `internal/gui/app.go` `DecodeImage` | data URL from the frontend | returns image size and points so the UI can mark the codes |
| GUI read-back | `internal/gui/app.go` | every generated PNG | "Lesbar geprüft": each generated code is decoded again |

`DecodeImageBytes` (`internal/barcodes/decode_image.go`) is the guarded
door for untrusted input: at most 20 MB (`MaxImageBytes`), at most 40 MP
(`MaxImagePixels`, checked from the header before pixels are allocated),
PNG, JPEG, GIF or WebP. Errors are `InputError`s with the codes
`image_format` and `image_too_large`; nothing found is `nothing_found`. The
GUI shows them in German from `frontend/src/lib/i18n/errors.ts`, MCP and
CLI in English.

## Stage by stage

### 1. Luminance — `toGray`

Everything after this works on one `*image.Gray` at origin (0,0). JPEG
photos arrive as `*image.YCbCr`, whose Y plane already is the luminance and
is copied row by row; `*image.Gray` likewise. Anything else is drawn over
white, so transparent pixels count as paper, not as black.

### 2. Working size — `shrink`, `maxWorkSide = 2000`

A 12 MP phone photo has far more pixels than any code needs; at full size
one decode took about 12 s. Images whose long side exceeds 2000 px are
reduced by an integer factor k (box average of k×k pixels): a 4000×3000
photo becomes 2000×1500. A code that fills a tenth of the photo's width
keeps about 200 px, plenty for every reader.

### 3. `decodeScene`

#### a) Whole image, all readers

`decodeOnce` builds one ZXing `BinaryBitmap` (hybrid binarizer) and runs
all readers with `TRY_HARDER`. A 40 px white quiet zone is added first
(`withQuietZone`): codes cropped to their edge, as in a screenshot, are
not readable without it.

How many codes each reader reports per pass decides what the next steps
must do:

| Reader | Codes per pass | Detector |
|---|---|---|
| QR (`multi/qrcode`) | all | finder patterns anywhere in the image |
| PDF417 (`internal/pdf417decode`, ZXing port) | all | start/stop patterns, row by row |
| DataMatrix | one | white-rectangle search **from the centre outwards** |
| Aztec | one | bullseye search **from the centre outwards** |
| EAN-13/8, UPC-A, Code 128, Code 39, ITF | one | row scans from the middle; `TRY_HARDER` also rotates 90° |

So after this pass QR and PDF417 are complete; DataMatrix and Aztec are
found only near the image centre, and of 1D codes only the first one.

#### b) 2D candidates — `matrixCandidates` + windows

The DataMatrix and Aztec detectors find a code reliably **if the centre of
the image they are given lies on the code**, even among text and lines —
and not otherwise (measured: a 300 px window finds a 180 px DataMatrix
with its centre up to 60 px off; a window centred beside the code finds
nothing). A fixed grid of windows would need a step smaller than the
smallest code and costs too much. Instead the image is searched for spots
that look like a 2D code:

1. **Cells.** The image is divided into cells, 160 along the long side
   (`cellsPerSide`, at least 6 px, `minCell`).
2. **Changes in both directions.** Per cell the dark/light changes
   (grey step > 40, `edgeContrast`) are counted along rows and along
   columns. A cell is *busy* when both reach 0.04 per pixel
   (`busyDensity`). A 2D code changes often in both directions; plain
   ground, ruled lines and the bars of a 1D code change in one direction
   only.
3. **Joining.** Inside a code some cells show no change — a large module,
   Aztec's bullseye. Busy cells are therefore grown by 2 cells
   (`joinCells`, `dilate`) so the parts of one code form one area;
   without this a code fell apart into four or six small areas, none of
   them centred on the code.
4. **Areas.** Connected areas (8-neighbourhood) of at least 4 cells
   (`minCandidateArea`) become bounding boxes, largest first, at most 40
   (`maxCandidates`).
5. **Windows.** Each box gets a margin of a third of its size plus 8 px
   (the detector needs the code's edge and quiet zone), is clipped to the
   image and decoded with the DataMatrix and Aztec readers only
   (`matrixReaders`). A window whose area already contains all points of
   a code found before is skipped (`alreadyFound`).

The top and bottom edges of 1D barcodes also count as busy; such
candidates cost one cheap window and find nothing.

#### c) Further 1D codes — `decodeLinearRest`

The 1D readers report one code per pass. Every 1D code found so far is
painted white on a copy of the image and the 1D readers (`linearReaders`)
run again, until a pass brings nothing new (at most 12 rounds,
`maxLinearCodes`). This costs nothing when an image has no 1D code.

`paintOverLinear` knows only the result points, which lie on the line
that was read. It extends the area across the bars — up and down for a
horizontal code, left and right for a turned one — as long as the pixels
still match the read line (mean grey difference < 40, `lineMatch`), then
adds a small margin and paints it white.

A density-based 1D candidate search was tried first and dropped: ruled
lines and text baselines look exactly like a barcode turned by 90°, and on
a lined background the whole image became one candidate.

### 4. Fallbacks

Only when stage 3 found nothing, from cheap to expensive; the first that
finds anything wins:

1. **Doubled** (`scale2x`, nearest neighbour) — dense codes on very few
   pixels, e.g. a small screenshot.
2. **Inverted** (`invert`) — light codes on dark ground.
3. **Codes on bottles and tubes** — `decodeCurved` (`cylinder.go`), see
   below.
4. **Full resolution** — `decodeScene` on the unreduced image, only if
   stage 2 had shrunk it: small codes in a large photo.

### Codes on bottles — `decodeCurved`

The 2D detectors sample on a perspective grid. On a cylinder the modules
narrow along a cosine towards the edges, and when the camera looks down
(or up) the rows bend into arcs; off-centre the samples miss the middle
modules and a sharp code is not read. `decodeCurved` takes the largest
2D-code candidates (stage 3b) and unwraps each window under assumed
shapes (`cylinderGuesses`): the code covers 70° or 110° of the
circumference, is turned 0°, ±20° or ±35°, the camera looks 15° down,
level or 15° up. For a guess, `fitCylinder` places the shape so the code
spans the candidate box — a point at angle θ lies at
x = cx + R·sin θ and drops by R·(cos θ − cos θ₀)·sin(tilt) — and `unwrap`
renders the window flat: columns at equal arc length, rows straightened.
DataMatrix, Aztec and QR then read the flat image; result points are
mapped back through the same shape (`cylinderMapper`), so the marks sit
on the photo.

Cost: 11–18 ms per guess. A photo gets 30 guesses for the largest
candidate and 9 for the next two (`guessesPerCandidate`) — at most about
a second, and only when nothing else was read. A camera frame
(`DecodeOptions.Live`, GUI `DecodeCameraFrame`) gets the 12 most likely
guesses of the largest candidate; the next frame follows anyway.

**Shape switch.** In the mobile scan view the user can say what the code
sits on (`DecodeOptions.Surface`, left of the camera image): *auto*,
*flat pack* (skips the bottle search), *bottle* — the code about half as
wide as the bottle — or *tube* — the code wraps far around. In the image
a code of width q times the bottle's covers s = 2·asin(q) of the
circumference: q = ½ gives 60°, q ≈ ¾ about 100° (`surfaceSpan`). With a
known curvature a camera frame tries every turn and tilt of it (15
guesses) and reads what the full search of a still photo reads; without
the hint a turned DataMatrix at 60° is missed in a frame
(`TestDecodeSurfaceHint`).

2D reader sets get the grey image as gozxing's planar luminance source
(`binaryBitmap`) instead of the generic per-pixel conversion; the 1D
readers keep the generic one, which can rotate for vertical codes.

### Coordinates — `frame`

Each reader sees a different image: reduced, padded with a quiet zone,
cropped to a window, doubled. A `frame{ox, oy, f}` maps a point p of the
image a reader saw to the caller's image as `p·f + (ox, oy)`;
`crop(at)` and `scaled(s)` derive the frame of a window or an enlarged
copy, `inverse` maps back (used to paint over 1D codes). All `Points` in
`Decoded` are therefore in the pixels of the image the caller passed in,
and the GUI draws its marks directly over the picture.

### Results

Codes are merged by (type, text); the same code found in several passes
or windows is reported once. A 1D result that lies inside a 2D code read
from the same image is dropped (`dropInside2D`): a row through a QR
code's modules can match a short 1D pattern — UPC-E, with its weak check,
was read from an event QR code. An EAN-13 with a leading 0 is reported as
the 12-digit UPC-A it most likely was made as. `qrformats.Parse` then
splits known payloads into fields. It looks at securPharm pack codes
first, before any trimming, because their control characters (GS, RS,
EOT) matter: GS1 (AIs 01/17/10/21/710) and IFA format 06 (9N/1T/D/S),
with PZN, PPN and GTIN check digits verified. The classic German pack
barcode — Code 39 `-12345678` (PZN8) or `-1234567` (PZN7, shown as PZN8
with a leading 0) — is recognised too; it carries only the PZN, no batch
or expiry (for that text see github.com/hmsoft0815/mlc_expiry).

## GUI camera and region selection

`frontend/src/lib/components/Checker.svelte`, same `DecodeImage` binding
as files.

- **Live camera** (`getUserMedia`, rear camera preferred): one frame
  every 300 ms (`SCAN_INTERVAL_MS`), the next only after the decoder
  answered, so a slow device never builds a queue. Frames go through
  `DecodeCameraFrame` (shorter bottle search, see above) together with
  the shape switch (auto / flat / bottle / tube) and the code switch,
  both remembered per device.
- **Front camera:** its frames are sharpened first (`DecodeOptions.Sharpen`,
  unsharp mask σ 2, amount 1.5) and every third one is also sent mirrored.
  Front cameras have a fixed focus for faces (about 30–50 cm): a code held
  closer is blurred, one held at that distance is small. Sharpening reads
  one blur step more (a 120 px DataMatrix at σ 2.0, `TestDecodeSharpenFrontCamera`),
  no more — the scan view says so and points to the rear camera.
- **Code switch** (top left, `DecodeOptions.Family`): *all codes*,
  *square* — QR, DataMatrix, Aztec, square guide frame, bottle search on —
  or *barcode* — 1D codes and PDF417, a wide 3:1 guide frame, no 2D
  candidates and no bottle search. Only those readers run: on a
  1600×900 frame one pass took 180 ms with all readers and 16 ms with the
  2D ones for a DataMatrix, 128 ms and 71 ms with the 1D ones for an
  EAN-13 (`TestDecodeFamily` checks that each kind finds its own codes
  and ignores the other). Frames are sent as
  JPEG with at most 1600 px on the long side (`FRAME_MAX_SIDE`).
- **Guide frame:** every other frame only the area inside the drawn
  guide frame (18 % inset, `GUIDE_INSET`, must match `.camera-frame`) is
  sent, at full camera resolution: small codes get more pixels there than
  in the reduced whole frame. The preview is not cropped (no
  `object-fit: cover`), otherwise the drawn frame and the decoded area
  differ.
- The first frame with a code stops the camera and stays as the checked
  image with its marks. Switching tabs, loading a file or "Kamera beenden"
  also stops it. "Kamera wechseln" cycles through the cameras; the front
  camera is shown mirrored.
- **Region selection** ("Ausschnitt wählen") on a still image: the user
  drags a rectangle, only that part is decoded, and the points are shifted
  back onto the whole image.
- iOS: works in the WKWebView of the Wails app (tested on iPadOS 26.6);
  `NSCameraUsageDescription` in `build/ios/Info.plist` is required.

## Measurements

On the development machine (16 threads); a phone is slower by a factor of
roughly 3–5.

| Case | Before | Now |
|---|---|---|
| 12 MP photo, one DataMatrix | 11.9 s, not found | 0.8 s, found |
| Camera frame 1600×900 (benchmark) | — | 0.21 s |
| DataMatrix / Aztec away from the centre (8 positions each) | not found | found |
| 3, 4, 6 DataMatrix in one photo | none found | all found |
| Shipping label: 2 × Code 128 + DataMatrix | one Code 128 missing | all found |
| ZXing blackbox datamatrix-1 / -2 (real photos) | 23/23, 18/18 | 23/23, 18/18 |
| ZXing blackbox aztec-2 | 7/22 | 7/22 (ZXing's own test expects fewer) |
| Codes on bottles (OpticScript set, 16 images per type: 30°–120° of the circumference, turned 0/35°, camera level/20° from above) | DataMatrix 3, QR 5, EAN-13 8 | DataMatrix 11, QR 14, EAN-13 8 (camera frames: 9, 12) |
| … still not read | | 2D: 90°/120° turned 35°, 120° from above; EAN-13: everything turned beyond 30° and 120° (no 1D unwrapping yet) |

## Tests

| Test | File | Guards |
|---|---|---|
| `TestDecodeAnywhereInImage` | `decode_scene_test.go` | DataMatrix, Aztec, QR, PDF417, EAN-13 at 9 positions on a lined background |
| `TestDecodeSeveralPacks` | `decode_scene_test.go` | 3, 4 and 6 pharma DataMatrix codes in one photo |
| `TestDecodeShippingLabel` | `decode_scene_test.go` | two Code 128 and a DataMatrix |
| `TestDecodePhonePhotoTime` | `decode_scene_test.go` | 12 MP photo under 2 s (skipped with `-short`) |
| `BenchmarkDecodeCameraFrame` | `decode_scene_test.go` | time per camera frame |
| `TestDecodeBlackbox` | `decode_blackbox_test.go` | real photos from ZXing (`testdata/zxing`, Apache-2.0) and codes on bottles made with mlc OpticScript (`testdata/cylinder`, `task testdata:cylinder`); `minFound` per set may only rise |
| round trip | `decode_test.go` | every symbology we generate is read back with the same content |
| exported codes | `quietzone_test.go` | codes are readable exactly as exported, without the decoder's extra quiet zone |
| `task test:scan` | `tests/scan` | an independent reader (zxing-cpp) reads what we generate |

## Known limits and ideas

- **Aztec photos:** 7 of 22 ZXing blackbox photos (perspective, glare).
  The Aztec detector itself is the limit, not the pipeline.
- **Candidate count:** at most 40 per image. A photo full of small print
  can produce many busy areas; the code is usually the largest, and areas
  are tried largest first.
- **Strong perspective or rotation** of 1D codes can defeat
  `paintOverLinear`'s line match; the code is then found again and
  painting stops after `maxLinearCodes` rounds — no wrong result, only
  time.
- **Codes on bottles and tubes** (T-20260926-12): `decodeCurved` reads
  most 2D codes; strongly curved and turned ones (code over 90° of the
  circumference, turned 35°) are still lost — the shape guesses do not
  reach them, and near the silhouette one pixel covers several modules.
  1D codes have no unwrapping yet. Glare across a DataMatrix's solid "L"
  edge destroys the information; only turning the bottle helps there.
- **Aztec with non-Latin-1 text** needs ECI (ticket T-20260926-09).
- **Next step if needed:** run the candidate windows in parallel
  goroutines; the readers are independent per window.
