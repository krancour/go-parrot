package common

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Status of the camera settings

// CameraSettingsState ...
// TODO: Document this
type CameraSettingsState interface {
	// RLock blocks until a read lock is obtained. This permits callers to procede
	// with querying any or all attributes of the camera settings state without
	// worry that some attributes will be overwritten as others are read. i.e. It
	// permits the possibility of taking an atomic snapshop of camera settings
	// state. Note that use of this function is not obligatory for applications
	// that do not require such guarantees. Callers MUST call RUnlock() or else
	// camera settings state will never resume updating.
	RLock()
	// RUnlock releases a read lock on the camera settings state. See RLock().
	RUnlock()
	// FOV returns the horizontal field of vision in degrees. A boolean value is
	// also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	FOV() (float32, bool)
	// PanMax returns the maximum (right) pan in degrees. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	PanMax() (float32, bool)
	// PanMin returns the minimum (left) pan in degrees. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	PanMin() (float32, bool)
	// TiltMax returns the maximum (up) tilt in degrees. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	TiltMax() (float32, bool)
	// TiltMin returns the minimum (down) tilt in degrees. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	TiltMin() (float32, bool)
}

type cameraSettingsState struct {
	fov     *float32
	panMax  *float32
	panMin  *float32
	tiltMax *float32
	tiltMin *float32
	lock    sync.RWMutex
}

func (c *cameraSettingsState) ID() uint8 {
	return 15
}

func (c *cameraSettingsState) Name() string {
	return "CameraSettingsState"
}

func (c *cameraSettingsState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"CameraSettingsChanged",
			[]interface{}{
				float32(0), // fov,
				float32(0), // panMax,
				float32(0), // panMin,
				float32(0), // tiltMax,
				float32(0), // tiltMin,
			},
			c.cameraSettingsChanged,
		),
	}
}

// cameraSettingsChanged is invoked by the device at connection time.
func (c *cameraSettingsState) cameraSettingsChanged(args []interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.fov = ptr.ToFloat32(args[0].(float32))
	c.panMax = ptr.ToFloat32(args[1].(float32))
	c.panMin = ptr.ToFloat32(args[2].(float32))
	c.tiltMax = ptr.ToFloat32(args[3].(float32))
	c.tiltMin = ptr.ToFloat32(args[4].(float32))
	log.WithField(
		"fov", *c.fov,
	).WithField(
		"panMax", *c.panMax,
	).WithField(
		"panMin", *c.panMin,
	).WithField(
		"tiltMax", *c.tiltMax,
	).WithField(
		"tiltMin", *c.tiltMin,
	).Debug("camera settings changed")
	return nil
}

func (c *cameraSettingsState) RLock() {
	c.lock.RLock()
}

func (c *cameraSettingsState) RUnlock() {
	c.lock.RUnlock()
}

func (c *cameraSettingsState) FOV() (float32, bool) {
	if c.fov == nil {
		return 0, false
	}
	return *c.fov, true
}

func (c *cameraSettingsState) PanMax() (float32, bool) {
	if c.panMax == nil {
		return 0, false
	}
	return *c.panMax, true
}

func (c *cameraSettingsState) PanMin() (float32, bool) {
	if c.panMin == nil {
		return 0, false
	}
	return *c.panMin, true
}

func (c *cameraSettingsState) TiltMax() (float32, bool) {
	if c.tiltMax == nil {
		return 0, false
	}
	return *c.tiltMax, true
}

func (c *cameraSettingsState) TiltMin() (float32, bool) {
	if c.tiltMin == nil {
		return 0, false
	}
	return *c.tiltMin, true
}
