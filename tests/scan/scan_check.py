#!/usr/bin/env python3
"""Scan test: generate codes with the CLI and read them back with zxing-cpp.

Unit tests check what the encoder is asked to do; only a decoder shows that
the result is readable. This is how the PDF417 row-indicator bug (about half
of all lengths unreadable) was found.

Usage: scan_check.py <barcode-cli> <work-dir>
Needs: pip install zxing-cpp pillow (task test:scan sets up a venv).
"""
import subprocess
import sys
from pathlib import Path

import zxingcpp
from PIL import Image

CLI, WORK = sys.argv[1], Path(sys.argv[2])
WORK.mkdir(parents=True, exist_ok=True)

# (type, input, expected decoded text) — EAN/UPC inputs lack the check digit
# on purpose; the decoder must see the completed code.
FIXED = [
    ("qr", "https://mlcgo.eu/produkte?x=Größe", None),
    ("datamatrix", "MLC-DM-12345 äöü", None),
    ("code128", "MLC-128-abc", None),
    ("code39", "CODE39-TEST", None),
    ("ean13", "400638133393", "4006381333931"),
    ("ean8", "9638507", "96385074"),
    ("upca", "03600029145", "036000291452"),
    ("itf", "12345678", None),
]
SWEEP_TYPES = ["aztec", "pdf417", "qr", "datamatrix"]
SWEEP_SOURCES = {"ascii": "abcdefghij" * 6, "umlaut": "äöüßÄÖÜ" * 9,
                 "digits": "0123456789" * 6, "mixed": "Größe 12 € x;y,z" * 4}


def decode(path: Path):
    # Read the file exactly as exported — no extra margin, so a missing
    # quiet zone shows up here (it was hidden before 1.6.0).
    img = Image.open(path).convert("L")
    found = zxingcpp.read_barcodes(img)
    return found[0].text if found else None


def check(btype: str, data: str, expected: str | None, name: str) -> bool:
    out = WORK / f"{name}.png"
    run = subprocess.run([CLI, "-type", btype, "-data", data, "-out", str(out)],
                         capture_output=True, text=True)
    if run.returncode != 0:
        print(f"✗ {name}: generation failed: {run.stderr.strip()[:120]}")
        return False
    got = decode(out)
    want = expected or data
    # UPC-A may be reported as EAN-13 with a leading 0.
    ok = got == want or (btype == "upca" and got == "0" + want)
    if not ok:
        print(f"✗ {name}: {data[:30]!r} decoded as {got!r}, want {want!r}")
    return ok


failed = total = 0
for btype, data, expected in FIXED:
    total += 1
    failed += not check(btype, data, expected, f"{btype}-fixed")
for btype in SWEEP_TYPES:
    for src_name, src in SWEEP_SOURCES.items():
        for n in range(1, 61):
            data = src[:n].strip() or "x"
            if btype == "aztec" and any(ord(c) > 0xFF for c in data):
                continue  # rejected on purpose: Aztec is ISO-8859-1 only (no ECI)
            total += 1
            failed += not check(btype, data, None, f"{btype}-{src_name}-{n}")

print(f"{total - failed}/{total} codes decoded")
sys.exit(1 if failed else 0)
