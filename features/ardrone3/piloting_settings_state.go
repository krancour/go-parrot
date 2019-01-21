package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Piloting Settings state from product

// PilotingSettingsState ...
// TODO: Document this
type PilotingSettingsState interface {
	// RLock blocks until a read lock is obtained. This permits callers to procede
	// with querying any or all attributes of the piloting state without worry
	// that some attributes will be overwritten as others are read. i.e. It
	// permits the possibility of taking an atomic snapshop of piloting state.
	// Note that use of this function is not obligatory for applications that do
	// not require such guarantees. Callers MUST call RUnlock() or else piloting
	// state will never resume updating.
	RLock()
	// RUnlock releases a read lock on the GPS state. See RLock().
	RUnlock()
	// MotionDetectionEnabled returns a boolean indicating whether motion
	// detection is enabled. A boolean value is also returned, indicating whether
	// the first value was reported by the device (true) or a default value
	// (false). This permits callers to distinguish real zero values from default
	// zero values.
	MotionDetectionEnabled() (bool, bool)
}

type pilotingSettingsState struct {
	motionDetectionEnabled *bool
	lock                   sync.RWMutex
}

func (p *pilotingSettingsState) ID() uint8 {
	return 6
}

func (p *pilotingSettingsState) Name() string {
	return "PilotingSettingsState"
}

func (p *pilotingSettingsState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"MaxAltitudeChanged",
			[]interface{}{
				float32(0), // current,
				float32(0), // min,
				float32(0), // max,
			},
			p.maxAltitudeChanged,
		),
		arcommands.NewD2CCommand(
			1,
			"MaxTiltChanged",
			[]interface{}{
				float32(0), // current,
				float32(0), // min,
				float32(0), // max,
			},
			p.maxTiltChanged,
		),
		// arcommands.NewD2CCommand(
		// 	2,
		// 	"AbsolutControlChanged",
		// 	[]interface{}{
		// 		uint8(0), // on,
		// 	},
		// 	p.absolutControlChanged,
		// ),
		arcommands.NewD2CCommand(
			3,
			"MaxDistanceChanged",
			[]interface{}{
				float32(0), // current,
				float32(0), // min,
				float32(0), // max,
			},
			p.maxDistanceChanged,
		),
		arcommands.NewD2CCommand(
			4,
			"NoFlyOverMaxDistanceChanged",
			[]interface{}{
				uint8(0), // shouldNotFlyOver,
			},
			p.noFlyOverMaxDistanceChanged,
		),
		arcommands.NewD2CCommand(
			5,
			"AutonomousFlightMaxHorizontalSpeed",
			[]interface{}{
				float32(0), // value,
			},
			p.autonomousFlightMaxHorizontalSpeed,
		),
		arcommands.NewD2CCommand(
			6,
			"AutonomousFlightMaxVerticalSpeed",
			[]interface{}{
				float32(0), // value,
			},
			p.autonomousFlightMaxVerticalSpeed,
		),
		arcommands.NewD2CCommand(
			7,
			"AutonomousFlightMaxHorizontalAcceleration",
			[]interface{}{
				float32(0), // value,
			},
			p.autonomousFlightMaxHorizontalAcceleration,
		),
		arcommands.NewD2CCommand(
			8,
			"AutonomousFlightMaxVerticalAcceleration",
			[]interface{}{
				float32(0), // value,
			},
			p.autonomousFlightMaxVerticalAcceleration,
		),
		arcommands.NewD2CCommand(
			9,
			"AutonomousFlightMaxRotationSpeed",
			[]interface{}{
				float32(0), // value,
			},
			p.autonomousFlightMaxRotationSpeed,
		),
		arcommands.NewD2CCommand(
			10,
			"BankedTurnChanged",
			[]interface{}{
				uint8(0), // state,
			},
			p.bankedTurnChanged,
		),
		// arcommands.NewD2CCommand(
		// 	11,
		// 	"MinAltitudeChanged",
		// 	[]interface{}{
		// 		float32(0), // current,
		// 		float32(0), // min,
		// 		float32(0), // max,
		// 	},
		// 	p.minAltitudeChanged,
		// ),
		// arcommands.NewD2CCommand(
		// 	12,
		// 	"CirclingDirectionChanged",
		// 	[]interface{}{
		// 		int32(0), // value,
		// 	},
		// 	p.circlingDirectionChanged,
		// ),
		// arcommands.NewD2CCommand(
		// 	13,
		// 	"CirclingRadiusChanged",
		// 	[]interface{}{
		// 		uint16(0), // current,
		// 		uint16(0), // min,
		// 		uint16(0), // max,
		// 	},
		// 	p.circlingRadiusChanged,
		// ),
		// arcommands.NewD2CCommand(
		// 	14,
		// 	"CirclingAltitudeChanged",
		// 	[]interface{}{
		// 		uint16(0), // current,
		// 		uint16(0), // min,
		// 		uint16(0), // max,
		// 	},
		// 	p.circlingAltitudeChanged,
		// ),
		// arcommands.NewD2CCommand(
		// 	15,
		// 	"PitchModeChanged",
		// 	[]interface{}{
		// 		int32(0), // value,
		// 	},
		// 	p.pitchModeChanged,
		// ),
		arcommands.NewD2CCommand(
			16,
			"MotionDetection",
			[]interface{}{
				uint8(0), // enabled,
			},
			p.motionDetection,
		),
	}
}

// TODO: Implement this
// Title: Max altitude
// Description: Max altitude.\n The drone will not fly higher than this altitude
//   (above take off point).
// Support: 0901;090c;090e
// Triggered: by [SetMaxAltitude](#1-2-0).
// Result:
func (p *pilotingSettingsState) maxAltitudeChanged(args []interface{}) error {
	// current := args[0].(float32)
	//   Current altitude max
	// min := args[1].(float32)
	//   Range min of altitude
	// max := args[2].(float32)
	//   Range max of altitude
	log.Info("ardrone3.maxAltitudeChanged() called")
	return nil
}

// TODO: Implement this
// Title: Max pitch/roll
// Description: Max pitch/roll.\n The drone will not fly higher than this
//   altitude (above take off point).
// Support: 0901;090c
// Triggered: by [SetMaxAltitude](#1-2-0).
// Result:
func (p *pilotingSettingsState) maxTiltChanged(args []interface{}) error {
	// current := args[0].(float32)
	//   Current max tilt
	// min := args[1].(float32)
	//   Range min of tilt
	// max := args[2].(float32)
	//   Range max of tilt
	log.Info("ardrone3.maxTiltChanged() called")
	return nil
}

// // TODO: Implement this
// // Title: Absolut control
// // Description: Absolut control.
// // Support:
// // Triggered:
// // Result:
// // WARNING: Deprecated
// func (p *pilotingSettingsState) absolutControlChanged(
// 	args []interface{},
// ) error {
// 	// on := args[0].(uint8)
// 	//   1 if enabled, 0 if disabled
// 	log.Info("ardrone3.absolutControlChanged() called")
// 	return nil
// }

// TODO: Implement this
// Title: Max distance
// Description: Max distance.
// Support: 0901;090c;090e
// Triggered: by [SetMaxDistance](#1-2-3).
// Result:
func (p *pilotingSettingsState) maxDistanceChanged(args []interface{}) error {
	// current := args[0].(float32)
	//   Current max distance in meter
	// min := args[1].(float32)
	//   Minimal possible max distance
	// max := args[2].(float32)
	//   Maximal possible max distance
	log.Info("ardrone3.maxDistanceChanged() called")
	return nil
}

// TODO: Implement this
// Title: Geofencing
// Description: Geofencing.\n If set, the drone won&#39;t fly over the
//   [MaxDistance](#1-6-3).
// Support: 0901;090c;090e
// Triggered: by [EnableGeofence](#1-2-4).
// Result:
func (p *pilotingSettingsState) noFlyOverMaxDistanceChanged(
	args []interface{},
) error {
	// shouldNotFlyOver := args[0].(uint8)
	//   1 if the drone won&#39;t fly further than max distance, 0 if no
	//   limitation on the drone will be done
	log.Info("ardrone3.noFlyOverMaxDistanceChanged() called")
	return nil
}

// TODO: Implement this
// Title: Autonomous flight max horizontal speed
// Description: Autonomous flight max horizontal speed.
// Support: 0901:3.3.0;090c:3.3.0
// Triggered: by [SetAutonomousFlightMaxHorizontalSpeed](#1-2-5).
// Result:
func (p *pilotingSettingsState) autonomousFlightMaxHorizontalSpeed(
	args []interface{},
) error {
	// value := args[0].(float32)
	//   maximum horizontal speed [m/s]
	log.Info("ardrone3.autonomousFlightMaxHorizontalSpeed() called")
	return nil
}

// TODO: Implement this
// Title: Autonomous flight max vertical speed
// Description: Autonomous flight max vertical speed.
// Support: 0901:3.3.0;090c:3.3.0
// Triggered: by [SetAutonomousFlightMaxVerticalSpeed](#1-2-6).
// Result:
func (p *pilotingSettingsState) autonomousFlightMaxVerticalSpeed(
	args []interface{},
) error {
	// value := args[0].(float32)
	//   maximum vertical speed [m/s]
	log.Info("ardrone3.autonomousFlightMaxVerticalSpeed() called")
	return nil
}

// TODO: Implement this
// Title: Autonomous flight max horizontal acceleration
// Description: Autonomous flight max horizontal acceleration.
// Support: 0901:3.3.0;090c:3.3.0
// Triggered: by [SetAutonomousFlightMaxHorizontalAcceleration](#1-2-7).
// Result:
func (p *pilotingSettingsState) autonomousFlightMaxHorizontalAcceleration(
	args []interface{},
) error {
	// value := args[0].(float32)
	//   maximum horizontal acceleration [m/s2]
	log.Info("ardrone3.autonomousFlightMaxHorizontalAcceleration() called")
	return nil
}

// TODO: Implement this
// Title: Autonomous flight max vertical acceleration
// Description: Autonomous flight max vertical acceleration.
// Support: 0901:3.3.0;090c:3.3.0
// Triggered: by [SetAutonomousFlightMaxVerticalAcceleration](#1-2-8).
// Result:
func (p *pilotingSettingsState) autonomousFlightMaxVerticalAcceleration(
	args []interface{},
) error {
	// value := args[0].(float32)
	//   maximum vertical acceleration [m/s2]
	log.Info("ardrone3.autonomousFlightMaxVerticalAcceleration() called")
	return nil
}

// TODO: Implement this
// Title: Autonomous flight max rotation speed
// Description: Autonomous flight max rotation speed.
// Support: 0901:3.3.0;090c:3.3.0
// Triggered: by [SetAutonomousFlightMaxRotationSpeed](#1-2-9).
// Result:
func (p *pilotingSettingsState) autonomousFlightMaxRotationSpeed(
	args []interface{},
) error {
	// value := args[0].(float32)
	//   maximum yaw rotation speed [deg/s]
	log.Info("ardrone3.autonomousFlightMaxRotationSpeed() called")
	return nil
}

// TODO: Implement this
// Title: Banked Turn mode
// Description: Banked Turn mode.\n If banked turn mode is enabled, the drone
//   will use yaw values from the piloting command to infer with roll and pitch
//   on the drone when its horizontal speed is not null.
// Support: 0901:3.2.0;090c:3.2.0
// Triggered: by [SetBankedTurnMode](#1-2-10).
// Result:
func (p *pilotingSettingsState) bankedTurnChanged(args []interface{}) error {
	// state := args[0].(uint8)
	//   1 if enabled, 0 if disabled
	log.Info("ardrone3.bankedTurnChanged() called")
	return nil
}

// // TODO: Implement this
// // Title: Min altitude
// // Description: Min altitude.\n Only sent by fixed wings.
// // Support: 090e
// // Triggered: by [SetMinAltitude](#1-2-11).
// // Result:
// func (p *pilotingSettingsState) minAltitudeChanged(args []interface{}) error {
// 	// current := args[0].(float32)
// 	//   Current altitude min
// 	// min := args[1].(float32)
// 	//   Range min of altitude min
// 	// max := args[2].(float32)
// 	//   Range max of altitude min
// 	log.Info("ardrone3.minAltitudeChanged() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Circling direction
// // Description: Circling direction.\n Only sent by fixed wings.
// // Support: 090e
// // Triggered: by [SetCirclingDirection](#1-2-12).
// // Result:
// func (p *pilotingSettingsState) circlingDirectionChanged(
// 	args []interface{},
// ) error {
// 	// value := args[0].(int32)
// 	//   The circling direction
// 	//   0: CW: Circling ClockWise
// 	//   1: CCW: Circling Counter ClockWise
// 	log.Info("ardrone3.circlingDirectionChanged() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Circling radius
// // Description: Circling radius.\n Only sent by fixed wings.
// // Support: none
// // Triggered: by [SetCirclingRadius](#1-2-13).
// // Result:
// // WARNING: Deprecated
// func (p *pilotingSettingsState) circlingRadiusChanged(
// 	args []interface{},
// ) error {
// 	// current := args[0].(uint16)
// 	//   The current circling radius in meter
// 	// min := args[1].(uint16)
// 	//   Range min of circling radius in meter
// 	// max := args[2].(uint16)
// 	//   Range max of circling radius in meter
// 	log.Info("ardrone3.circlingRadiusChanged() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Circling altitude
// // Description: Circling altitude.\n Bounds will be automatically adjusted
// //   according to the [MaxAltitude](#1-6-0).\n Only sent by fixed wings.
// // Support: 090e
// // Triggered: by [SetCirclingRadius](#1-2-14) or when bounds change due to
// //   [SetMaxAltitude](#1-2-0).
// // Result:
// func (p *pilotingSettingsState) circlingAltitudeChanged(
// 	args []interface{},
// ) error {
// 	// current := args[0].(uint16)
// 	//   The current circling altitude in meter
// 	// min := args[1].(uint16)
// 	//   Range min of circling altitude in meter
// 	// max := args[2].(uint16)
// 	//   Range max of circling altitude in meter
// 	log.Info("ardrone3.circlingAltitudeChanged() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Pitch mode
// // Description: Pitch mode.
// // Support: 090e
// // Triggered: by [SetPitchMode](#1-2-15).
// // Result:
// func (p *pilotingSettingsState) pitchModeChanged(args []interface{}) error {
// 	// value := args[0].(int32)
// 	//   The Pitch mode
// 	//   0: NORMAL: Positive pitch values will make the drone lower its nose.
// 	//      Negative pitch values will make the drone raise its nose.
// 	//   1: INVERTED: Pitch commands are inverted. Positive pitch values will
// 	//      make the drone raise its nose. Negative pitch values will make the
// 	//      drone lower its nose.
// 	log.Info("ardrone3.pitchModeChanged() called")
// 	return nil
// }

// motionDetection is invoked by the device to indicate whether motion
// detection is enabled.
func (p *pilotingSettingsState) motionDetection(args []interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.motionDetectionEnabled = ptr.ToBool(args[0].(uint8) == 1)
	log.WithField(
		"motionDetectionEnabled", *p.motionDetectionEnabled,
	).Debug("motion detection enabled changed")
	return nil
}

func (p *pilotingSettingsState) RLock() {
	p.lock.RLock()
}

func (p *pilotingSettingsState) RUnlock() {
	p.lock.RUnlock()
}

func (p *pilotingSettingsState) MotionDetectionEnabled() (bool, bool) {
	if p.motionDetectionEnabled == nil {
		return false, false
	}
	return *p.motionDetectionEnabled, true
}
