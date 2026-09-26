// A code printed on the label of a bottle or tube: CONTENT (a clean,
// square code image) is wrapped around a vertical cylinder and lit like a
// photo — darker towards the edges, a glare stripe, sensor noise, blur.
// Generates the test set internal/barcodes/testdata/cylinder (see
// `task testdata:cylinder`); runs with mlc OpticScript's mlcos-run.
//!INPUT: CONTENT
//!OUTPUT: OUTPUT
//!PARAM: SPAN:number=60,min=10,max=160
//!PARAM: TURN:number=0,min=-60,max=60
//!PARAM: CODEW:number=320,min=80,max=700
//!PARAM: TILT:number=20,min=0,max=60
//!PARAM: GLARE:number=0.5,min=0,max=1
//!PARAM: GLAREPOS:number=999,min=-170,max=999
//!PARAM: SHADE:number=0.55,min=0,max=0.9
//!PARAM: NOISE:number=0.03,min=0,max=0.3
//!PARAM: BLUR:number=0.8,min=0,max=4
//!PARAM: QUALITY:integer=75,min=30,max=100

// SPAN      degrees of the circumference the code covers (small = almost
//           flat, 90+ = strongly curved); the bottle's radius follows
//           from it, so the code keeps its size and only the curvature
//           changes
// TURN      degrees the code is turned away from the camera
// CODEW     width of the code in the image (unrolled, px)
// TILT      degrees the camera looks down on the bottle: horizontal lines
//           on the label become arcs (the near middle lower than the
//           sides) and the code is foreshortened vertically
// GLARE     strength of the glare stripe, GLAREPOS its angle relative to
//           the code's centre (the default lies beside the code)
// SHADE     darkening towards the edges (0 = flat light)
// NOISE     sensor noise, sigma on the 0..1 scale
// QUALITY   JPEG quality of the output, as a phone would store it

const W = 800, H = 600;
const rad = (d) => (d * Math.PI) / 180;
const R = CODEW / rad(SPAN);
// A turned bottle: the photographer still centres the code, so the axis
// moves the other way and the code is seen at an angle.
const cx = W / 2 - R * Math.sin(rad(TURN));
const xAt = (theta) => cx + R * Math.sin(theta);

// Table and bottle body (a white label).
const base = Engine.createColoredImage(W, H, "#8a8478");
const left = Math.max(0, Math.round(cx - R)), right = Math.min(W, Math.round(cx + R));
base.fillRectangle(left, 0, right - left, H, "#f6f5f0");

// The code, unrolled: as wide as the arc it covers, as high as its aspect
// ratio says (square for 2D codes, flat for EAN).
const code = Engine.loadImage(CONTENT);
const cols = 33, rows = 3;
const width = CODEW;
const height = ((width * code.height) / code.width) * Math.cos(rad(TILT));
const top = (H - height) / 2;
// Seen from above, a point of the surface at angle theta lies R·cos(theta)
// closer to the camera than the axis and drops by that times sin(TILT);
// relative to the code's centre (theta = TURN) the sides rise.
const drop = (theta) => R * (Math.cos(theta) - Math.cos(rad(TURN))) * Math.sin(rad(TILT));
const nodes = [];
for (let r = 0; r < rows; r++) {
  const y = top + (height * r) / (rows - 1);
  const row = [];
  for (let c = 0; c < cols; c++) {
    const theta = rad(TURN) + (c / (cols - 1) - 0.5) * rad(SPAN);
    row.push(xAt(theta) / W, (y + drop(theta)) / H);
  }
  nodes.push(row);
}
base.stampGrid(code, { rows, cols, nodes }, Interp.Bicubic);
code.free();

// Light falling off towards the edges: multiply with a gradient across the
// body; brightness follows the cosine of the surface angle.
const bodyW = right - left;
const shadeCv = Engine.createCanvas(bodyW, H);
const stops = [], colors = [];
for (let i = 0; i <= 16; i++) {
  const t = i / 16;
  const xr = (left + t * bodyW - cx) / R; // -1 … 1 across the whole bottle
  const theta = Math.asin(Math.max(-1, Math.min(1, xr)));
  const v = Math.round(255 * (1 - SHADE * (1 - Math.cos(theta))));
  const hex = v.toString(16).padStart(2, "0");
  stops.push(t);
  colors.push(`#${hex}${hex}${hex}`);
}
shadeCv.linearGradient(px(0, 0), px(bodyW, 0), colors, stops);
shadeCv.drawPathStr(`M 0,0 L ${bodyW},0 L ${bodyW},${H} L 0,${H} Z`, false);
const shade = shadeCv.toImage();
base.blendAt(shade, px(left, 0), 1.0, Blend.Multiply);
shade.free();
shadeCv.free();

// Glare: a soft white stripe along the bottle.
if (GLARE > 0) {
  const pos = GLAREPOS === 999 ? -(SPAN / 2 + 12) : GLAREPOS; // just left of the code
  const gx = xAt(rad(TURN + pos));
  const gw = Math.max(8, R * 0.18);
  const glareCv = Engine.createCanvas(Math.round(2 * gw), H);
  glareCv.linearGradient(px(0, 0), px(2 * gw, 0), ["#00000000", "#ffffffff", "#00000000"], [0, 0.5, 1]);
  glareCv.drawPathStr(`M 0,0 L ${2 * gw},0 L ${2 * gw},${H} L 0,${H} Z`, false);
  const glare = glareCv.toImage();
  base.blendAt(glare, px(Math.round(gx - gw), 0), GLARE, Blend.Screen);
  glare.free();
  glareCv.free();
}

if (BLUR > 0) base.gaussianBlur(BLUR);
if (NOISE > 0) base.addNoise({ type: "gaussian", sigma: NOISE, color: false });
base.save(OUTPUT, { format: "jpeg", quality: QUALITY });
