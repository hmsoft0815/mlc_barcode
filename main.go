package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/mlcmcp/mlc_barcode/internal/gui"
	"github.com/mlcmcp/mlc_barcode/internal/version"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// assets embeds the compiled frontend bundle into the binary.
//
//go:embed all:frontend/dist
var assets embed.FS

// appIcon embeds the application window icon.
//
//go:embed build/appicon.png
var appIcon []byte

// startUI selects the UI the window opens with. The Windows test build
// "…-ipad.exe" sets it to "mobile" (-ldflags -X main.startUI=mobile): the
// mobile UI in an iPad-shaped window, for debugging it at the desk; the
// camera is the webcam, share sheet, torch and vibration fall back.
var startUI = "desktop"

func main() {
	barcodeApp := gui.NewBarcodeApp()

	wailsApp := application.New(application.Options{
		Name:        "MLC Barcode",
		Description: "Cross-platform Barcode & QR Code Generator",
		Icon:        appIcon,
		Services: []application.Service{
			application.NewService(barcodeApp),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		// The camera preview is a <video> in the page. Without inline
		// playback WKWebView shows a still frame and plays the stream in its
		// own fullscreen player instead.
		IOS: application.IOSOptions{
			EnableInlineMediaPlayback:       true,
			EnableAutoplayWithoutUserAction: true,
			BackgroundColour:                application.NewRGB(20, 20, 22), // dark theme body, no white flash
		},
	})

	title, url := "MLC Barcode v"+version.Version, "/"
	width, height, minWidth, minHeight := 1200, 850, 900, 650
	if startUI == "mobile" {
		// iPad portrait (768×1024 points) — fits a 1080p screen.
		title += " — iPad-Ansicht (Test)"
		url = "/?ui=mobile"
		width, height, minWidth, minHeight = 768, 1024, 360, 560
	}

	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:           title,
		Width:           width,
		Height:          height,
		MinWidth:        minWidth,
		MinHeight:       minHeight,
		DevToolsEnabled: isDebugBuild(),
		// Camera (scan view) and location ("use my position") for our own
		// page without the webview asking on every start; WebView2 and
		// WebKitGTK apply this, iOS gets camera_grant_ios.go.
		Permissions: map[application.PermissionType]application.Permission{
			application.PermissionCamera:      application.PermissionAllow,
			application.PermissionGeolocation: application.PermissionAllow,
		},
		KeyBindings: devtoolsKeyBindings(),
		Linux: application.LinuxWindow{
			Icon: appIcon,
		},
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(248, 249, 250),
		URL:              url,
	})

	installCameraGrant()
	if err := wailsApp.Run(); err != nil {
		log.Fatal(err)
	}
}
