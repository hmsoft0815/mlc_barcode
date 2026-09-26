// Package pdf417 is a copy of github.com/boombuler/barcode/pdf417 (v1.1.0,
// MIT licence, see LICENSE in this directory) with one fix.
//
// getLeftCodeWord computed the row indicator of cluster 0 as (rows-3)/3.
// ISO/IEC 15438 specifies (rows-1)/3, which the right-hand indicator
// already used. The two indicators then disagreed about the number of rows
// for many data lengths and scanners rejected the symbol: in a sweep of
// 1-60 characters (ASCII, umlauts, digits) about half of the codes could
// not be decoded. With the fix all of them decode (zxing-cpp).
//
// Drop this copy once the upstream package carries the fix.
package pdf417
