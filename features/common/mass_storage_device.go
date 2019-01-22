package common

// MassStorageDevice ...
// TODO: Document this
type MassStorageDevice interface {
	// ID returns the ID of the mass storage device.
	ID() uint8
	// Name returns the mass storage device name. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	Name() (string, bool)
	// Size returns the size of the mass storage device in megabytes. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	Size() (uint32, bool)
	// UsedSize returns the used size of the mass storage device in megabytes. A
	// boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	UsedSize() (uint32, bool)
	// Plugged returns a boolean indicating whether the mass storage device is
	// currently plugged in. A boolean value is also returned, indicating whether
	// the first value was reported by the device (true) or a default value
	// (false). This permits callers to distinguish real zero values from default
	// zero values.
	Plugged() (bool, bool)
	// Full returns a boolean indicating whether the mass storage is full. A
	// boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	Full() (bool, bool)
	// Internal returns a boolean indicating whether the mass storage device is
	// internal. A boolean value is also returned, indicating whether the first
	// value was reported by the device (true) or a default value (false). This
	// permits callers to distinguish real zero values from default zero values.
	Internal() (bool, bool)
	// PhotoCount returns the number of photos in mass storage. A boolean value is
	// also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	PhotoCount() (uint16, bool)
	// VideoCount returns the numnber of videos in mass storage. A boolean value
	// is also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	VideoCount() (uint16, bool)
	// PudCount returns the number of puds in mass storage. A boolean value is
	// also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	PudCount() (uint16, bool)
	// CrashLogCount returns the number of crash logs in mass storage. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	CrashLogCount() (uint16, bool)
	// RawPhotoCount returns the number of raw photos in mass stroage. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	RawPhotoCount() (uint16, bool)
	// CurrentRunMassStorageID returns the mass storage ID for the current run. A
	// boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	CurrentRunMassStorageID() (uint8, bool)
	// CurrentRunPhotoCount returns the number of photos in mass storage related
	// to the current run. A boolean value is also returned, indicating whether
	// the first value was reported by the device (true) or a default value
	// (false). This permits callers to distinguish real zero values from default
	// zero values.
	CurrentRunPhotoCount() (uint16, bool)
	// CurrentRunVideoCount returns the number of photos in mass storage related
	// to the current run. A boolean value is also returned, indicating whether
	// the first value was reported by the device (true) or a default value
	// (false). This permits callers to distinguish real zero values from default
	// zero values.
	CurrentRunVideoCount() (uint16, bool)
	// CurrentRunRawPhotoCount returns the number of raw photos in mass storage
	// related to the current run. A boolean value is also returned, indicating
	// whether the first value was reported by the device (true) or a default
	// value (false). This permits callers to distinguish real zero values from
	// default zero values.
	CurrentRunRawPhotoCount() (uint16, bool)
}

type massStorageDevice struct {
	id                      uint8
	name                    *string
	size                    *uint32
	usedSize                *uint32
	plugged                 *bool
	full                    *bool
	internal                *bool
	photoCount              *uint16
	videoCount              *uint16
	pudCount                *uint16 // TODO: What is a pud?
	crashLogCount           *uint16
	rawPhotoCount           *uint16
	currentRunMassStorageID *uint8
	currentRunPhotoCount    *uint16
	currentRunVideoCount    *uint16
	currentRunRawPhotoCount *uint16
}

func (m *massStorageDevice) ID() uint8 {
	return m.id
}

func (m *massStorageDevice) Name() (string, bool) {
	if m.name == nil {
		return "", false
	}
	return *m.name, true
}

func (m *massStorageDevice) Size() (uint32, bool) {
	if m.size == nil {
		return 0, false
	}
	return *m.size, true
}

func (m *massStorageDevice) UsedSize() (uint32, bool) {
	if m.usedSize == nil {
		return 0, false
	}
	return *m.usedSize, true
}

func (m *massStorageDevice) Plugged() (bool, bool) {
	if m.plugged == nil {
		return false, false
	}
	return *m.plugged, true
}

func (m *massStorageDevice) Full() (bool, bool) {
	if m.full == nil {
		return false, false
	}
	return *m.full, true
}

func (m *massStorageDevice) Internal() (bool, bool) {
	if m.internal == nil {
		return false, false
	}
	return *m.internal, true
}

func (m *massStorageDevice) PhotoCount() (uint16, bool) {
	if m.photoCount == nil {
		return 0, false
	}
	return *m.photoCount, true
}

func (m *massStorageDevice) VideoCount() (uint16, bool) {
	if m.videoCount == nil {
		return 0, false
	}
	return *m.videoCount, true
}

func (m *massStorageDevice) PudCount() (uint16, bool) {
	if m.pudCount == nil {
		return 0, false
	}
	return *m.pudCount, true
}

func (m *massStorageDevice) CrashLogCount() (uint16, bool) {
	if m.crashLogCount == nil {
		return 0, false
	}
	return *m.crashLogCount, true
}

func (m *massStorageDevice) RawPhotoCount() (uint16, bool) {
	if m.rawPhotoCount == nil {
		return 0, false
	}
	return *m.rawPhotoCount, true
}

func (m *massStorageDevice) CurrentRunMassStorageID() (uint8, bool) {
	if m.currentRunMassStorageID == nil {
		return 0, false
	}
	return *m.currentRunMassStorageID, true
}

func (m *massStorageDevice) CurrentRunPhotoCount() (uint16, bool) {
	if m.currentRunPhotoCount == nil {
		return 0, false
	}
	return *m.currentRunPhotoCount, true
}

func (m *massStorageDevice) CurrentRunVideoCount() (uint16, bool) {
	if m.currentRunVideoCount == nil {
		return 0, false
	}
	return *m.currentRunVideoCount, true
}

func (m *massStorageDevice) CurrentRunRawPhotoCount() (uint16, bool) {
	if m.currentRunRawPhotoCount == nil {
		return 0, false
	}
	return *m.currentRunRawPhotoCount, true
}
