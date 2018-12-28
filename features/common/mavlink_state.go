package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Mavlink flight plans states commands

// MavlinkState ...
// TODO: Document this
type MavlinkState interface{}

type mavlinkState struct{}

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
		arcommands.NewD2CCommand(
			2,
			"MissionItemExecuted",
			[]interface{}{
				uint32(0), // idx,
			},
			m.missionItemExecuted,
		),
	}
}

// TODO: Implement this
// Title: Playing state of a FlightPlan
// Description: Playing state of a FlightPlan.
// Support: 0901:2.0.29;090c;090e
// Triggered: by [StartFlightPlan](#0-11-0), [PauseFlightPlan](#0-11-1) or
//   [StopFlightPlan](#0-11-2).
// Result:
func (m *mavlinkState) mavlinkFilePlayingStateChanged(args []interface{}) error {
	// state := args[0].(int32)
	//   State of the mavlink
	//   0: playing: Mavlink file is playing
	//   1: stopped: Mavlink file is stopped (arg filepath and type are useless in
	//      this state)
	//   2: paused: Mavlink file is paused
	//   3: loaded: Mavlink file is loaded (it will be played at take-off)
	// filepath := args[1].(string)
	//   flight plan file path from the mavlink ftp root
	// type := args[2].(int32)
	//   type of the played mavlink file
	//   0: flightPlan: Mavlink file for FlightPlan
	//   1: mapMyHouse: Mavlink file for MapMyHouse
	log.Info("common.mavlinkFilePlayingStateChanged() called")
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

// TODO: Implement this
// Title: Mission item executed
// Description: Mission item has been executed.
// Support: 090c:4.2.0;090e:1.4.0
// Triggered: when a mission item has been executed during a flight plan.
// Result:
func (m *mavlinkState) missionItemExecuted(args []interface{}) error {
	// idx := args[0].(uint32)
	//   Index of the mission item. This is the place of the mission item in the
	//   list of the items of the mission. Begins at 0.
	log.Info("common.missionItemExecuted() called")
	return nil
}
