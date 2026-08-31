//go:build !devtools && !debug

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func isDebugBuild() bool { return false }

func devtoolsKeyBindings() map[string]func(application.Window) { return nil }
