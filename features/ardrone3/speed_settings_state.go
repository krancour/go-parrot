package ardrone3

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Speed Settings state from product

// SpeedSettingsState ...
// TODO: Document this
type SpeedSettingsState interface{}

type speedSettingsState struct{}

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
	// current := args[0].(float32)
	//   Current max vertical speed in m/s
	// min := args[1].(float32)
	//   Range min of vertical speed
	// max := args[2].(float32)
	//   Range max of vertical speed
	log.Info("ardrone3.maxVerticalSpeedChanged() called")
	return nil
}

// TODO: Implement this
// Title: Max rotation speed
// Description: Max rotation speed.
// Support: 0901;090c
// Triggered: by [SetMaxRotationSpeed](#1-11-1).
// Result:
func (s *speedSettingsState) maxRotationSpeedChanged(args []interface{}) error {
	// current := args[0].(float32)
	//   Current max yaw rotation speed in degree/s
	// min := args[1].(float32)
	//   Range min of yaw rotation speed
	// max := args[2].(float32)
	//   Range max of yaw rotation speed
	log.Info("ardrone3.maxRotationSpeedChanged() called")
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

// TODO: Implement this
// Title: Outdoor mode
// Description: Outdoor mode.
// Support:
// Triggered:
// Result:
// WARNING: Deprecated
func (s *speedSettingsState) outdoorChanged(args []interface{}) error {
	// outdoor := args[0].(uint8)
	//   1 if outdoor flight, 0 if indoor flight
	log.Info("ardrone3.outdoorChanged() called")
	return nil
}

// TODO: Implement this
// Title: Max pitch/roll rotation speed
// Description: Max pitch/roll rotation speed.
// Support: 0901;090c
// Triggered: by [SetMaxPitchRollRotationSpeed](#1-11-4).
// Result:
func (s *speedSettingsState) maxPitchRollRotationSpeedChanged(
	args []interface{},
) error {
	// current := args[0].(float32)
	//   Current max pitch/roll rotation speed in degree/s
	// min := args[1].(float32)
	//   Range min of pitch/roll rotation speed
	// max := args[2].(float32)
	//   Range max of pitch/roll rotation speed
	log.Info("ardrone3.maxPitchRollRotationSpeedChanged() called")
	return nil
}
