package photokit

import (
	"errors"
	"runtime"
	"testing"
)

// Without a running Cocoa app (as in `go test`) Pick must say so instead
// of waiting for a main thread that never answers.
func TestPickWithoutApp(t *testing.T) {
	_, err := Pick(PickOptions{Dir: t.TempDir()})
	switch {
	case !Supported():
		if !errors.Is(err, ErrUnsupported) {
			t.Fatalf("unsupported platform: got %v, want ErrUnsupported", err)
		}
	case runtime.GOOS == "darwin":
		if !errors.Is(err, ErrNoApp) {
			t.Fatalf("macOS without NSApp: got %v, want ErrNoApp", err)
		}
	}
}

func TestSaveMissingFile(t *testing.T) {
	if !CanSave() {
		if _, err := Save("x.jpg", SaveOptions{}); !errors.Is(err, ErrUnsupported) {
			t.Fatalf("got %v, want ErrUnsupported", err)
		}
		return
	}
	if _, err := Save(t.TempDir()+"/missing.jpg", SaveOptions{}); err == nil {
		t.Fatal("saving a missing file must fail before PhotoKit is asked")
	}
}
