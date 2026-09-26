// Package pdf417 is a copy of github.com/boombuler/barcode/pdf417 (v1.1.0,
// MIT licence, see LICENSE in this directory) with two changes.
//
// getLeftCodeWord computed the row indicator of cluster 0 as (rows-3)/3.
// ISO/IEC 15438 specifies (rows-1)/3, which the right-hand indicator
// already used. The two indicators then disagreed about the number of rows
// for many data lengths and scanners rejected the symbol: in a sweep of
// 1-60 characters (ASCII, umlauts, digits) about half of the codes could
// not be decoded. With the fix all of them decode (zxing-cpp).
//
// Second, text outside ASCII is written as UTF-8 bytes but was not marked:
// readers assume ISO-8859-1 and show "GrÃ¶Ãe" for "Größe". highlevelEncode
// now starts such text with ECI 26 (UTF-8), codewords 927, 26 — found by
// the round trip against the ported reader in internal/pdf417decode.
//
// Drop this copy once the upstream package carries both.
package pdf417
