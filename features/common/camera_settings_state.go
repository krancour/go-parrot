package common

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Status of the camera settings

// CameraSettingsState ...
// TODO: Document this
type CameraSettingsState interface {
	lock.ReadLockable
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
	sync.RWMutex
	fov     *float32
	panMax  *float32
	panMin  *float32
	tiltMax *float32
	tiltMin *float32
}

func newCameraSettingsState() *cameraSettingsState {
	return &cameraSettingsState{}
}

func (c *cameraSettingsState) ClassID() uint8 {
	return 15
}

func (c *cameraSettingsState) ClassName() string {
	return "CameraSettingsState"
}

func (c *cameraSettingsState) D2CCommands(
	log *log.Entry,
) []arcommands.D2CCommand {
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
			log,
		),
	}
}

// cameraSettingsChanged is invoked by the device at connection time.
func (c *cameraSettingsState) cameraSettingsChanged(
	args []interface{},
	log *log.Entry,
) error {
	c.Lock()
	defer c.Unlock()
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
