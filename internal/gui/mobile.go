package gui

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Native helpers for the mobile UI (share sheet, torch, haptics, opening
// links). The platform parts live in mobile_ios.go, mobile_android.go and
// mobile_desktop.go; on the desktop they do nothing and report false, so
// the frontend can fall back to its own means.

// Platform reports the OS the app runs on (runtime.GOOS). The frontend
// shows the mobile UI on "ios" and "android".
func (a *BarcodeApp) Platform() string {
	return runtime.GOOS
}

// ShareText opens the system share sheet with a text or link. It reports
// false where there is none.
func (a *BarcodeApp) ShareText(text string) bool {
	return nativeShare(sharePayload("text", text))
}

// ShareImage writes a PNG (base64 or data URL) to the app's temporary
// directory and opens the share sheet with the file, so it can be sent as
// a picture. It reports false where that is not possible.
func (a *BarcodeApp) ShareImage(pngBase64, name string) bool {
	if !canShareFiles {
		return false
	}
	if i := strings.Index(pngBase64, ","); strings.HasPrefix(pngBase64, "data:") && i > 0 {
		pngBase64 = pngBase64[i+1:]
	}
	data, err := base64.StdEncoding.DecodeString(pngBase64)
	if err != nil {
		return false
	}
	path := filepath.Join(os.TempDir(), sanitizeFilename(name)+".png")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return false
	}
	return nativeShare(sharePayload("url", "file://"+path))
}

// OpenExternal opens a link in the system browser (or the app registered
// for it, e.g. mail). It reports false where the frontend must do it.
func (a *BarcodeApp) OpenExternal(url string) bool {
	return nativeOpenURL(url)
}

// SetTorch switches the camera light. It reports false where there is none.
func (a *BarcodeApp) SetTorch(on bool) bool {
	return nativeTorch(on)
}

// Haptic gives a short vibration: "success", "warning", "error" or
// "selection".
func (a *BarcodeApp) Haptic(kind string) {
	nativeHaptic(kind)
}

func sharePayload(key, value string) string {
	b, _ := json.Marshal(map[string]string{key: value})
	return string(b)
}
