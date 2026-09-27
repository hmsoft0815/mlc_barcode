//go:build !ios

package main

// installCameraGrant: on the desktop the window's Permissions option does
// it (see main.go); nothing to install.
func installCameraGrant() {}
