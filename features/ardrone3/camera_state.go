package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Camera state

// CameraState ...
// TODO: Document this
type CameraState interface {
	lock.ReadLockable
	// Tilt returns camera tilt in degrees. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	Tilt() (float32, bool)
	// Pan returns camera pan in degrees. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	Pan() (float32, bool)
	// CenterTilt returns the camera tilt value, in degrees, that is considered
	// "center." This can later be used to re-center the camera if/when desired. A
	// boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	CenterTilt() (float32, bool)
	// CenterPan returns the camera pan value, in degrees, that is considered
	// "center." This can later be used to re-center the camera if/when desired. A
	// boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	CenterPan() (float32, bool)
	// MaxTiltVelocity returns the the maximum velocity, in degrees per second, at
	// which camera tilt may be adjusted. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	MaxTiltVelocity() (float32, bool)
	// MaxPanVelocity returns the maximum velocity, in degrees per second, at
	// which camera pan may be adjusted. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	MaxPanVelocity() (float32, bool)
}

type cameraState struct {
	sync.RWMutex
	// tilt represents camera tilt in degrees
	tilt *float32
	// pan represents camera pan in degrees
	pan *float32
	// centerTilt represents the camera tilt, in degrees, that is considered
	// "center"
	centerTilt *float32
	// centerPan represents the camera pan, in degrees, that is considered
	// "center"
	centerPan *float32
	// maxTiltVelocity represents the maximum velocity, in degrees per second, at
	// which camera tilt may be adjusted.
	maxTiltVelocity *float32
	// maxPanVelocity represents the maximum velocity, in degrees per second, at
	// which camera pan may be adjusted.
	maxPanVelocity *float32
}

func (c *cameraState) ID() uint8 {
	return 25
}

func (c *cameraState) Name() string {
	return "CameraState"
}

func (c *cameraState) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"Orientation",
			[]interface{}{
				int8(0), // tilt,
				int8(0), // pan,
			},
			c.orientation,
			log,
		),
		arcommands.NewD2CCommand(
			1,
			"defaultCameraOrientation",
			[]interface{}{
				int8(0), // tilt,
				int8(0), // pan,
			},
			c.defaultCameraOrientation,
			log,
		),
		arcommands.NewD2CCommand(
			2,
			"OrientationV2",
			[]interface{}{
				float32(0), // tilt,
				float32(0), // pan,
			},
			c.orientationV2,
			log,
		),
		arcommands.NewD2CCommand(
			3,
			"defaultCameraOrientationV2",
			[]interface{}{
				float32(0), // tilt,
				float32(0), // pan,
			},
			c.defaultCameraOrientationV2,
			log,
		),
		arcommands.NewD2CCommand(
			4,
			"VelocityRange",
			[]interface{}{
				float32(0), // max_tilt,
				float32(0), // max_pan,
			},
			c.velocityRange,
			log,
		),
	}
}

// orientation is deprecated in favor of orientationV2, but since we can still
// see this command being invoked, we'll implement the command to avoid a
// warning, but the implementation will remain a no-op unless / until such time
// that it becomes clear that older versions of the firmware might require us to
// support both commands.
func (c *cameraState) orientation(args []interface{}, log *log.Entry) error {
	// tilt := args[0].(int8)
	//   Tilt camera consign for the drone [-100;100]
	// pan := args[1].(int8)
	//   Pan camera consign for the drone [-100;100]
	log.Debug("camera orientation changed-- this is a no-op")
	return nil
}

// defaultCameraOrientation is deprecated in favor of
// defaultCameraOrientationV2, but since we can still see this command being
// invoked, we'll implement the command to avoid a warning, but the
// implementation will remain a no-op unless / until such time that it becomes
// clear that older versions of the firmware might require us to support both
// commands.
func (c *cameraState) defaultCameraOrientation(
	args []interface{},
	log *log.Entry,
) error {
	// tilt := args[0].(int8)
	//   Tilt value (in degree)
	// pan := args[1].(int8)
	//   Pan value (in degree)
	log.Debug("default camera orientation changed-- this is a no-op")
	return nil
}

// orientationV2 is invoked by the device when the camera orientation chnages.
func (c *cameraState) orientationV2(args []interface{}, log *log.Entry) error {
	c.Lock()
	defer c.Unlock()
	c.tilt = ptr.ToFloat32(args[0].(float32))
	c.pan = ptr.ToFloat32(args[1].(float32))
	log.WithField(
		"tilt", *c.tilt,
	).WithField(
		"pan", *c.pan,
	).Debug("camera orientation changed")
	return nil
}

// defaultCameraOrientationV2 is invoked by the device at connection time to
// indicate the tilt and pan values that reflect a centered camera. Use these
// values to recenter the camera if/when desired.
func (c *cameraState) defaultCameraOrientationV2(
	args []interface{},
	log *log.Entry,
) error {
	c.Lock()
	defer c.Unlock()
	c.centerTilt = ptr.ToFloat32(args[0].(float32))
	c.centerPan = ptr.ToFloat32(args[1].(float32))
	log.WithField(
		"tilt", *c.tilt,
	).WithField(
		"pan", *c.pan,
	).Debug("default camera orientation changed")
	return nil
}

// velocityRange is invoked by the device at connection time to communicate
// the maximum velocity, in degrees per second, at which the camera can be
// reoriented.
func (c *cameraState) velocityRange(args []interface{}, log *log.Entry) error {
	c.Lock()
	defer c.Unlock()
	c.maxTiltVelocity = ptr.ToFloat32(args[0].(float32))
	c.maxPanVelocity = ptr.ToFloat32(args[1].(float32))
	log.WithField(
		"max_tilt", *c.maxTiltVelocity,
	).WithField(
		"max_pan", *c.maxPanVelocity,
	).Debug("max camera reorientation velocity changed")
	return nil
}

func (c *cameraState) Tilt() (float32, bool) {
	if c.tilt == nil {
		return 0, false
	}
	return *c.tilt, true
}

func (c *cameraState) Pan() (float32, bool) {
	if c.pan == nil {
		return 0, false
	}
	return *c.pan, true
}

func (c *cameraState) CenterTilt() (float32, bool) {
	if c.centerTilt == nil {
		return 0, false
	}
	return *c.centerTilt, true
}

func (c *cameraState) CenterPan() (float32, bool) {
	if c.centerPan == nil {
		return 0, false
	}
	return *c.centerPan, true
}

func (c *cameraState) MaxTiltVelocity() (float32, bool) {
	if c.maxTiltVelocity == nil {
		return 0, false
	}
	return *c.maxTiltVelocity, true
}

func (c *cameraState) MaxPanVelocity() (float32, bool) {
	if c.maxPanVelocity == nil {
		return 0, false
	}
	return *c.maxPanVelocity, true
}
