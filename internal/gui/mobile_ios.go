//go:build ios

package gui

import "github.com/wailsapp/wails/v3/pkg/application"

// The share sheet takes file URLs: a PNG in the app's tmp directory is
// offered as an image.
const canShareFiles = true

func nativeShare(payload string) bool { application.IOS.Share(payload); return true }
func nativeOpenURL(url string) bool   { application.IOS.OpenURL(url); return true }
func nativeTorch(on bool) bool        { application.IOS.SetTorch(on); return true }
func nativeHaptic(kind string)        { application.IOS.Haptic(kind) }
