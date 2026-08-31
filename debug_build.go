//go:build devtools || debug

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func isDebugBuild() bool { return true }

// devtoolsKeyBindings enables F12 in debug builds.
func devtoolsKeyBindings() map[string]func(application.Window) {
	return map[string]func(application.Window){
		"F12": func(w application.Window) { w.OpenDevTools() },
	}
}
