//go:build !ios && !android

package gui

// On the desktop the frontend opens links with the Wails runtime and has no
// share sheet, torch or vibration.
const canShareFiles = false

func nativeShare(string) bool   { return false }
func nativeOpenURL(string) bool { return false }
func nativeTorch(bool) bool     { return false }
func nativeHaptic(string)       {}
