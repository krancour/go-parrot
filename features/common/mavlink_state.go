package common

import (
	"sync"

	"github.com/krancour/go-parrot/lock"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Mavlink flight plans states commands

const (
	MavlinkStatePlaying int32 = 0
	MavlinkStateStopped int32 = 1
	MavlinkStatePaused  int32 = 2
	MavlinkStateLoaded  int32 = 3

	MavlinkTypeFlightPlan int32 = 0
	MavlinkTypeMapMyHouse int32 = 1
)

// MavlinkState ...
// TODO: Document this
type MavlinkState interface {
	lock.ReadLockable
	// MavlinkState returns an int32 representing the current state of mavlink. A
	// boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	MavlinkState() (int32, bool)
	// MavlinkFilePath returns a string representing the ftp path for the mavlink
	// file. A boolean value is also returned, indicating whether the first value
	// was reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	MavlinkFilePath() (string, bool)
	// MavlinkType returns an int32 representing the type of mavlink file. A
	// boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	MavlinkType() (int32, bool)
}

type mavlinkState struct {
	sync.RWMutex
	mavlinkState    *int32
	mavlinkFilePath *string
	mavlinkType     *int32
}

func (m *mavlinkState) ID() uint8 {
	return 12
}

func (m *mavlinkState) Name() string {
	return "MavlinkState"
}

func (m *mavlinkState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"MavlinkFilePlayingStateChanged",
			[]interface{}{
				int32(0),  // state,
				string(0), // filepath,
				int32(0),  // type,
			},
			m.mavlinkFilePlayingStateChanged,
		),
		// arcommands.NewD2CCommand(
		// 	1,
		// 	"MavlinkPlayErrorStateChanged",
		// 	[]interface{}{
		// 		int32(0), // error,
		// 	},
		// 	m.mavlinkPlayErrorStateChanged,
		// ),
		// arcommands.NewD2CCommand(
		// 	2,
		// 	"MissionItemExecuted",
		// 	[]interface{}{
		// 		uint32(0), // idx,
		// 	},
		// 	m.missionItemExecuted,
		// ),
	}
}

// mavlinkFilePlayingStateChanged is invoked by the device when a flight plan
// is started, paused, or stopped.
func (m *mavlinkState) mavlinkFilePlayingStateChanged(args []interface{}) error {
	m.Lock()
	defer m.Unlock()
	m.mavlinkState = ptr.ToInt32(args[0].(int32))
	m.mavlinkFilePath = ptr.ToString(args[1].(string))
	m.mavlinkType = ptr.ToInt32(args[2].(int32))
	log.WithField(
		"mavlinkState", *m.mavlinkState,
	).WithField(
		"mavlinkFilePath", *m.mavlinkFilePath,
	).WithField(
		"mavlinkType", *m.mavlinkType,
	).Debug("mavlink file playing state changed")
	return nil
}

// // TODO: Implement this
// // Title: FlightPlan error
// // Description: FlightPlan error.
// // Support: 0901:2.0.29;090c;090e
// // Triggered: by [StartFlightPlan](#0-11-0) if an error occurs.
// // Result:
// // WARNING: Deprecated
// func (m *mavlinkState) mavlinkPlayErrorStateChanged(args []interface{}) error {
// 	// error := args[0].(int32)
// 	//   State of play error
// 	//   0: none: There is no error
// 	//   1: notInOutDoorMode: The drone is not in outdoor mode
// 	//   2: gpsNotFixed: The gps is not fixed
// 	//   3: notCalibrated: The magnetometer of the drone is not calibrated
// 	log.Info("common.mavlinkPlayErrorStateChanged() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Mission item executed
// // Description: Mission item has been executed.
// // Support: 090c:4.2.0;090e:1.4.0
// // Triggered: when a mission item has been executed during a flight plan.
// // Result:
// func (m *mavlinkState) missionItemExecuted(args []interface{}) error {
// 	// idx := args[0].(uint32)
// 	//   Index of the mission item. This is the place of the mission item in the
// 	//   list of the items of the mission. Begins at 0.
// 	log.Info("common.missionItemExecuted() called")
// 	return nil
// }

func (m *mavlinkState) MavlinkState() (int32, bool) {
	if m.mavlinkState == nil {
		return 0, false
	}
	return *m.mavlinkState, true
}

func (m *mavlinkState) MavlinkFilePath() (string, bool) {
	if m.mavlinkFilePath == nil {
		return "", false
	}
	return *m.mavlinkFilePath, true
}

func (m *mavlinkState) MavlinkType() (int32, bool) {
	if m.mavlinkType == nil {
		return 0, false
	}
	return *m.mavlinkType, true
}
