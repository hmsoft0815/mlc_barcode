//go:build android

package gui

import "github.com/wailsapp/wails/v3/pkg/application"

// Android refuses file:// URIs in share intents (FileUriExposedException);
// sharing a picture needs a FileProvider in the Java part first. Until then
// the frontend shares the text.
const canShareFiles = false

func nativeShare(payload string) bool { application.Android.Share(payload); return true }
func nativeOpenURL(url string) bool   { application.Android.OpenURL(url); return true }
func nativeTorch(on bool) bool        { application.Android.SetTorch(on); return true }
func nativeHaptic(kind string)        { application.Android.Haptic(kind) }
