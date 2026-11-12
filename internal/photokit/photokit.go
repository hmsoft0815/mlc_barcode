// Package photokit gives Go apps access to the system photo library:
// pick photos from it and save images into it.
//
// It does not depend on Wails or any other GUI framework, but it needs a
// running Cocoa application (NSApp) on macOS: the picker is a window, and
// every call hands its work to the main thread. A Wails app has one; a
// plain command-line program does not and gets ErrNoApp.
//
// Platforms:
//
//	macOS 13+   Pick (PHPickerViewController), Save (PhotoKit)
//	macOS 12    Save only; Pick returns ErrUnsupported
//	iOS         planned (ErrUnsupported for now)
//	others      ErrUnsupported — hide the buttons, see Supported
//
// Picking needs no permission: the picker runs outside the app and hands
// over only what the user chose. Saving asks once for "add photos only";
// the app's Info.plist needs NSPhotoLibraryAddUsageDescription, or macOS
// ends the app when Save is called.
package photokit // copy of mlc_opticscript/photokit, see ORIGIN.md

import "errors"

var (
	// ErrUnsupported: this platform or OS version has no photo library
	// access through this package.
	ErrUnsupported = errors.New("photokit: not supported on this platform")
	// ErrNoApp: no running Cocoa application to show the picker in.
	ErrNoApp = errors.New("photokit: no running application (NSApp)")
	// ErrCancelled: the user closed the picker without choosing.
	ErrCancelled = errors.New("photokit: cancelled")
	// ErrDenied: the user did not allow adding photos (System Settings →
	// Privacy & Security → Photos).
	ErrDenied = errors.New("photokit: access to the photo library denied")
)

// PickOptions control Pick.
type PickOptions struct {
	// Limit is the most photos the user may choose; 0 means no limit.
	Limit int
	// Dir is where the picked files are written. It must exist. Empty
	// means a new directory under os.TempDir(); the caller removes it.
	Dir string
}

// Photo is one picked photo, written to disk as the original file.
type Photo struct {
	// Path of the file in PickOptions.Dir. The file keeps the format the
	// library holds (HEIC, JPEG, PNG, RAW …) with its metadata.
	Path string
	// Name is the file name the library suggests (often IMG_1234.HEIC).
	Name string
	// Type is the uniform type identifier, e.g. "public.heic".
	Type string
}

// SaveOptions control Save.
type SaveOptions struct {
	// Name is the original file name stored with the photo; empty means
	// the base name of the path.
	Name string
}

// Supported reports whether Pick works here. Use it to show or hide a
// "From Photos" button.
func Supported() bool { return supported() }

// CanSave reports whether Save works here.
func CanSave() bool { return canSave() }

// Pick shows the system photo picker and blocks until the user is done.
// The chosen photos are written to opts.Dir in the order chosen.
// Returns ErrCancelled when nothing was chosen.
func Pick(opts PickOptions) ([]Photo, error) { return pick(opts) }

// Save adds the image file at path to the photo library as a new photo
// and returns its library identifier. The file is copied; the library
// never changes or replaces an existing photo.
func Save(path string, opts SaveOptions) (string, error) { return save(path, opts) }
