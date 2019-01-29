package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Speed Settings state from product

// SpeedSettingsState ...
// TODO: Document this
type SpeedSettingsState interface {
	// RLock blocks until a read lock is obtained. This permits callers to procede
	// with querying any or all attributes of the speed settings state without
	// worry that some attributes will be overwritten as others are read. i.e. It
	// permits the possibility of taking an atomic snapshop of speed settings
	// state. Note that use of this function is not obligatory for applications
	// that do not require such guarantees. Callers MUST call RUnlock() or else
	// speed settings state will never resume updating.
	RLock()
	// RUnlock releases a read lock on the GPS state. See RLock().
	RUnlock()
	// MaxVerticalSpeed returns the configured maximum vertical speed of the
	// device in meters per second. A boolean value is also returned, indicating
	// whether the first value was reported by the device (true) or a default
	// value (false). This permits callers to distinguish real zero values from
	// default zero values.
	MaxVerticalSpeed() (float32, bool)
	// MaxVerticalSpeedRangeMin returns the minimum speed, in meters per second,
	// that the device's maximum vertical speed can be configured to. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	MaxVerticalSpeedRangeMin() (float32, bool)
	// MaxVerticalSpeedRangeMax returns the maximum speed, in meters per second,
	// that the device's maximum vertical speed can be configured to. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	MaxVerticalSpeedRangeMax() (float32, bool)
	// MaxRotationSpeed returns the configured maximum rotational speed of the
	// device in degrees per second. A boolean value is also returned, indicating
	// whether the first value was reported by the device (true) or a default
	// value (false). This permits callers to distinguish real zero values from
	// default zero values.
	MaxRotationSpeed() (float32, bool)
	// MaxRotationSpeedRangeMin returns the minimum rotational speed, in degrees
	// per second, that the device's maximum rotational speed can be configured
	// to. A boolean value is also returned, indicating whether the first value
	// was reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	MaxRotationSpeedRangeMin() (float32, bool)
	// MaxRotationSpeedRangeMax returns the maximum rotational speed, in degrees
	// per second, that the device's maximum rotational speed can be configured
	// to. A boolean value is also returned, indicating whether the first value
	// was reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	MaxRotationSpeedRangeMax() (float32, bool)
	// MaxPitchRollRotationSpeed returns the configured maximum roational speed
	// when changing tilt (pitch or roll) in degrees per second. A boolean value
	// is also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	MaxPitchRollRotationSpeed() (float32, bool)
	// MaxPitchRollRotationSpeedRangeMin returns the minimum rotation speed, in
	// degrees per second, that the device's maximum pitch/roll rotation speed can
	// be configured to. A boolean value is also returned, indicating whether the
	// first value was reported by the device (true) or a default value (false).
	// This permits callers to distinguish real zero values from default zero
	// values.
	MaxPitchRollRotationSpeedRangeMin() (float32, bool)
	// MaxPitchRollRotationSpeedRangeMax returns the maximum rotation speed, in
	// degrees per second, that the device's maximum pitch/roll rotation speed can
	// be configured to. A boolean value is also returned, indicating whether the
	// first value was reported by the device (true) or a default value (false).
	// This permits callers to distinguish real zero values from default zero
	// values.
	MaxPitchRollRotationSpeedRangeMax() (float32, bool)
}

type speedSettingsState struct {
	// maxVerticalSpeed is the configured maximum vertical speed of the device in
	// meters per second.
	maxVerticalSpeed *float32
	// maxVerticalSpeedRangeMin is the minimum speed, in meters per second, that
	// the device's maximum vertical speed can be configured to.
	maxVerticalSpeedRangeMin *float32
	// maxVerticalSpeedRangeMax is the maximum speed, in meters per second, that
	// the device's maximum vertical speed can be configured to.
	maxVerticalSpeedRangeMax *float32
	// maxRotationSpeed is the configured maximum rotational speed of the device
	// in degrees per second.
	maxRotationSpeed *float32
	// maxRotationSpeedRangeMin is the minimum rotational speed, in degrees per
	// second, that the device's maximum rotational speed can be configured to.
	maxRotationSpeedRangeMin *float32
	// maxRotationSpeedRangeMax is the maximum rotational speed, in degrees per
	// second, that the device's maximum rotational speed can be configured to.
	maxRotationSpeedRangeMax *float32
	// maxPitchRollRotationSpeed is the configured maximum roational speed when
	// changing tilt (pitch or roll) in degrees per second.
	maxPitchRollRotationSpeed *float32
	// maxPitchRollRotationSpeedRangeMin is the minimum rotation speed, in degrees
	// per second, that the device's maximum pitch/roll rotation speed can be
	// configured to.
	maxPitchRollRotationSpeedRangeMin *float32
	// maxPitchRollRotationSpeedRangeMax is the maximum rotation speed, in degrees
	// per second, that the device's maximum pitch/roll rotation speed can be
	// configured to.
	maxPitchRollRotationSpeedRangeMax *float32
	lock                              sync.RWMutex
}

func (s *speedSettingsState) ID() uint8 {
	return 12
}

func (s *speedSettingsState) Name() string {
	return "SpeedSettingsState"
}

func (s *speedSettingsState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"MaxVerticalSpeedChanged",
			[]interface{}{
				float32(0), // current,
				float32(0), // min,
				float32(0), // max,
			},
			s.maxVerticalSpeedChanged,
		),
		arcommands.NewD2CCommand(
			1,
			"MaxRotationSpeedChanged",
			[]interface{}{
				float32(0), // current,
				float32(0), // min,
				float32(0), // max,
			},
			s.maxRotationSpeedChanged,
		),
		arcommands.NewD2CCommand(
			2,
			"HullProtectionChanged",
			[]interface{}{
				uint8(0), // present,
			},
			s.hullProtectionChanged,
		),
		arcommands.NewD2CCommand(
			3,
			"OutdoorChanged",
			[]interface{}{
				uint8(0), // outdoor,
			},
			s.outdoorChanged,
		),
		arcommands.NewD2CCommand(
			4,
			"MaxPitchRollRotationSpeedChanged",
			[]interface{}{
				float32(0), // current,
				float32(0), // min,
				float32(0), // max,
			},
			s.maxPitchRollRotationSpeedChanged,
		),
	}
}

// TODO: Implement this
// Title: Max vertical speed
// Description: Max vertical speed.
// Support: 0901;090c
// Triggered: by [SetMaxVerticalSpeed](#1-11-0).
// Result:
func (s *speedSettingsState) maxVerticalSpeedChanged(args []interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.maxVerticalSpeed = ptr.ToFloat32(args[0].(float32))
	s.maxVerticalSpeedRangeMin = ptr.ToFloat32(args[1].(float32))
	s.maxVerticalSpeedRangeMax = ptr.ToFloat32(args[2].(float32))
	log.WithField(
		"maxVerticalSpeed", *s.maxVerticalSpeed,
	).WithField(
		"maxVerticalSpeedRangeMin", *s.maxVerticalSpeedRangeMin,
	).WithField(
		"maxVerticalSpeedRangeMax", *s.maxVerticalSpeedRangeMax,
	).Debug("max vertical speed changed")
	return nil
}

// TODO: Implement this
// Title: Max rotation speed
// Description: Max rotation speed.
// Support: 0901;090c
// Triggered: by [SetMaxRotationSpeed](#1-11-1).
// Result:
func (s *speedSettingsState) maxRotationSpeedChanged(args []interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.maxRotationSpeed = ptr.ToFloat32(args[0].(float32))
	s.maxRotationSpeedRangeMin = ptr.ToFloat32(args[1].(float32))
	s.maxRotationSpeedRangeMax = ptr.ToFloat32(args[2].(float32))
	log.WithField(
		"maxRotationSpeed", *s.maxRotationSpeed,
	).WithField(
		"maxRotationSpeedRangeMin", *s.maxRotationSpeedRangeMin,
	).WithField(
		"maxRotationSpeedRangeMax", *s.maxRotationSpeedRangeMax,
	).Debug("max rotation speed changed")
	return nil
}

// TODO: Implement this
// Title: Presence of hull protection
// Description: Presence of hull protection.
// Support: 0901;090c
// Triggered: by [SetHullProtectionPresence](#1-11-2).
// Result:
func (s *speedSettingsState) hullProtectionChanged(args []interface{}) error {
	// present := args[0].(uint8)
	//   1 if present, 0 if not present
	log.Info("ardrone3.hullProtectionChanged() called")
	return nil
}

// outdoorChanged is deprecated, but since we can still see this command being
// invoked, we'll implement the command to avoid a warning, but the
// implementation will remain a no-op unless / until such time that it becomes
// clear that older versions of the firmware might require us to support it.
func (s *speedSettingsState) outdoorChanged(args []interface{}) error {
	// outdoor := args[0].(uint8)
	//   1 if outdoor flight, 0 if indoor flight
	log.Debug("outdoor changed-- this is a no-op")
	return nil
}

// maxPitchRollRotationSpeedChanged is invoked by the device when the max
// pitch/roll rotation speed is changed.
func (s *speedSettingsState) maxPitchRollRotationSpeedChanged(
	args []interface{},
) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.maxPitchRollRotationSpeed = ptr.ToFloat32(args[0].(float32))
	s.maxPitchRollRotationSpeedRangeMin = ptr.ToFloat32(args[1].(float32))
	s.maxPitchRollRotationSpeedRangeMax = ptr.ToFloat32(args[2].(float32))
	log.WithField(
		"maxPitchRollRotationSpeed", *s.maxPitchRollRotationSpeed,
	).WithField(
		"maxPitchRollRotationSpeedRangeMin", *s.maxPitchRollRotationSpeedRangeMin,
	).WithField(
		"maxPitchRollRotationSpeedRangeMax", *s.maxPitchRollRotationSpeedRangeMax,
	).Debug("max pitch/roll rotation speed changed")
	return nil
}

func (s *speedSettingsState) RLock() {
	s.lock.RLock()
}

func (s *speedSettingsState) RUnlock() {
	s.lock.RUnlock()
}

func (s *speedSettingsState) MaxVerticalSpeed() (float32, bool) {
	if s.maxVerticalSpeed == nil {
		return 0, false
	}
	return *s.maxVerticalSpeed, true
}

func (s *speedSettingsState) MaxVerticalSpeedRangeMin() (float32, bool) {
	if s.maxVerticalSpeedRangeMin == nil {
		return 0, false
	}
	return *s.maxVerticalSpeedRangeMin, true
}

func (s *speedSettingsState) MaxVerticalSpeedRangeMax() (float32, bool) {
	if s.maxVerticalSpeedRangeMax == nil {
		return 0, false
	}
	return *s.maxVerticalSpeedRangeMax, true
}

func (s *speedSettingsState) MaxRotationSpeed() (float32, bool) {
	if s.maxRotationSpeed == nil {
		return 0, false
	}
	return *s.maxRotationSpeed, true
}

func (s *speedSettingsState) MaxRotationSpeedRangeMin() (float32, bool) {
	if s.maxRotationSpeedRangeMin == nil {
		return 0, false
	}
	return *s.maxRotationSpeedRangeMin, true
}

func (s *speedSettingsState) MaxRotationSpeedRangeMax() (float32, bool) {
	if s.maxRotationSpeedRangeMax == nil {
		return 0, false
	}
	return *s.maxRotationSpeedRangeMax, true
}

func (s *speedSettingsState) MaxPitchRollRotationSpeed() (float32, bool) {
	if s.maxPitchRollRotationSpeed == nil {
		return 0, false
	}
	return *s.maxPitchRollRotationSpeed, true
}

func (
	s *speedSettingsState,
) MaxPitchRollRotationSpeedRangeMin() (float32, bool) {
	if s.maxPitchRollRotationSpeedRangeMin == nil {
		return 0, false
	}
	return *s.maxPitchRollRotationSpeedRangeMin, true
}

func (
	s *speedSettingsState,
) MaxPitchRollRotationSpeedRangeMax() (float32, bool) {
	if s.maxPitchRollRotationSpeedRangeMax == nil {
		return 0, false
	}
	return *s.maxPitchRollRotationSpeedRangeMax, true
}
