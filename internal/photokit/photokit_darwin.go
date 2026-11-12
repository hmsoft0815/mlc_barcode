//go:build darwin && !ios && cgo

package photokit

/*
#cgo CFLAGS: -x objective-c -fobjc-arc -mmacosx-version-min=12.0
#cgo LDFLAGS: -framework Cocoa -framework Photos -framework PhotosUI -framework UniformTypeIdentifiers
#include <stdint.h>
#include <stdlib.h>

int  pkPickAvailable(void);
int  pkAppRunning(void);
void pkPick(uintptr_t h, int limit, const char *dir);
void pkSave(uintptr_t h, const char *path, const char *name);
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/cgo"
	"strings"
	"unsafe"
)

// Status codes shared with photokit_darwin.m.
const (
	stOK          = 0
	stCancelled   = 1
	stError       = 2
	stDenied      = 3
	stUnsupported = 4
)

type result struct {
	status int
	data   string
	msg    string
}

// pkDone is called once per request from Objective-C, on any thread.
//
//export pkDone
func pkDone(h C.uintptr_t, status C.int, data *C.char, msg *C.char) {
	hd := cgo.Handle(h)
	ch := hd.Value().(chan result)
	hd.Delete()
	ch <- result{int(status), C.GoString(data), C.GoString(msg)}
}

// call starts an Objective-C request and waits for its pkDone.
func call(start func(h C.uintptr_t)) result {
	ch := make(chan result, 1)
	start(C.uintptr_t(cgo.NewHandle(ch)))
	return <-ch
}

func (r result) err() error {
	switch r.status {
	case stOK:
		return nil
	case stCancelled:
		return ErrCancelled
	case stDenied:
		return ErrDenied
	case stUnsupported:
		return ErrUnsupported
	}
	if r.msg == "" {
		r.msg = "unknown error"
	}
	return errors.New("photokit: " + r.msg)
}

func supported() bool { return C.pkPickAvailable() != 0 }
func canSave() bool   { return true }

func pick(o PickOptions) ([]Photo, error) {
	if !supported() {
		return nil, ErrUnsupported
	}
	if C.pkAppRunning() == 0 {
		return nil, ErrNoApp
	}
	dir := o.Dir
	if dir == "" {
		d, err := os.MkdirTemp("", "photokit-")
		if err != nil {
			return nil, err
		}
		dir = d
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	limit := max(o.Limit, 0)
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	r := call(func(h C.uintptr_t) { C.pkPick(h, C.int(limit), cdir) })
	if err := r.err(); err != nil {
		return nil, err
	}
	var out []Photo
	for line := range strings.SplitSeq(r.data, "\n") {
		f := strings.Split(line, "\t")
		if len(f) == 3 {
			out = append(out, Photo{Path: f[0], Name: f[1], Type: f[2]})
		}
	}
	// Some photos loaded, others not (e.g. still in iCloud and offline):
	// hand over what is there and say what is missing.
	if r.msg != "" {
		return out, fmt.Errorf("photokit: not every photo could be loaded: %s", r.msg)
	}
	return out, nil
}

func save(path string, o SaveOptions) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(abs); err != nil {
		return "", err
	}
	name := o.Name
	if name == "" {
		name = filepath.Base(abs)
	}
	cpath, cname := C.CString(abs), C.CString(name)
	defer C.free(unsafe.Pointer(cpath))
	defer C.free(unsafe.Pointer(cname))
	r := call(func(h C.uintptr_t) { C.pkSave(h, cpath, cname) })
	return r.data, r.err()
}
