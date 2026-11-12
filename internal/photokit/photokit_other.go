//go:build !darwin || ios || !cgo

package photokit

func supported() bool { return false }
func canSave() bool   { return false }

func pick(PickOptions) ([]Photo, error)        { return nil, ErrUnsupported }
func save(string, SaveOptions) (string, error) { return "", ErrUnsupported }
