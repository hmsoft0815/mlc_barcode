"""Build the emulator camera image for tests/android/emu-test.ps1 -CameraImage.

The emulator (-camera-back imagefile:) scales a 1600x1200 image to full
height; the app's portrait preview then shows the strip x≈129..771, and the
scan frame covers about x 193..707, y 300..680 (measured with a quadrant
calibration image on the mlcb-test AVD, Pixel 7). This puts the code of a
test photo — a bottle with the code in the middle — centred in that frame.

    python3 tests/android/camera_image.py <photo.jpg> <out.jpg>
"""
import sys
from PIL import Image, ImageFilter

W, H = 1600, 1200
CENTER = (450, 520)          # centre of the scan frame in image coordinates
CODE_PX = 300                # code size wanted inside the frame

src = Image.open(sys.argv[1]).convert("RGB")
# The test photos (internal/barcodes/testdata/cylinder) show the code in the
# middle at ~35 % of the width.
scale = CODE_PX / (src.width * 0.35)
im = src.resize((round(src.width * scale), round(src.height * scale)))
x0, y0 = CENTER[0] - im.width // 2, CENTER[1] - im.height // 2

# Background: the bottle without label (a row near the top of the photo),
# stretched to the full height, continued sideways with the edge colours.
row = im.crop((0, int(im.height * 0.08), im.width, int(im.height * 0.08) + 2))
canvas = Image.new("RGB", (W, H), row.getpixel((0, 0)))
right = Image.new("RGB", (W - x0 - im.width, H), row.getpixel((im.width - 1, 0)))
canvas.paste(right, (x0 + im.width, 0))
canvas.paste(row.resize((im.width, H)).filter(ImageFilter.GaussianBlur(1)), (x0, 0))
canvas.paste(im, (x0, y0))
canvas.save(sys.argv[2], quality=92)
