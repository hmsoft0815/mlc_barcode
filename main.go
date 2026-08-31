package main

import (
	"embed"
	_ "embed"
	"log"

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

func main() {
	barcodeApp := NewBarcodeApp()

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
	})

	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:           "MLC Barcode v" + version.Version,
		Width:           1200,
		Height:          850,
		MinWidth:        900,
		MinHeight:       650,
		DevToolsEnabled: isDebugBuild(),
		KeyBindings:     devtoolsKeyBindings(),
		Linux: application.LinuxWindow{
			Icon: appIcon,
		},
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(248, 249, 250),
		URL:              "/",
	})

	if err := wailsApp.Run(); err != nil {
		log.Fatal(err)
	}
}
