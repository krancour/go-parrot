package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// FlightPlan state commands

// FlightPlanState ...
// TODO: Document this
type FlightPlanState interface{}

type flightPlanState struct{}

func (f *flightPlanState) ID() uint8 {
	return 17
}

func (f *flightPlanState) Name() string {
	return "FlightPlanState"
}

func (f *flightPlanState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"AvailabilityStateChanged",
			[]interface{}{
				uint8(0), // AvailabilityState,
			},
			f.availabilityStateChanged,
		),
		arcommands.NewD2CCommand(
			1,
			"ComponentStateListChanged",
			[]interface{}{
				int32(0), // component,
				uint8(0), // State,
			},
			f.componentStateListChanged,
		),
		arcommands.NewD2CCommand(
			2,
			"LockStateChanged",
			[]interface{}{
				uint8(0), // LockState,
			},
			f.lockStateChanged,
		),
	}
}

// TODO: Implement this
// Title: FlightPlan availability
// Description: FlightPlan availability.\n Availability is linked to GPS fix,
//   magnetometer calibration, sensor states...
// Support: 0901:2.0.29;090c;090e
// Triggered: on change.
// Result:
func (f *flightPlanState) availabilityStateChanged(args []interface{}) error {
	// AvailabilityState := args[0].(uint8)
	//   Running a flightPlan file is available (1 running a flightPlan file is
	//   available, otherwise 0)
	log.Info("common.availabilityStateChanged() called")
	return nil
}

// TODO: Implement this
// Title: FlightPlan components state list
// Description: FlightPlan components state list.
// Support: 0901:2.0.29;090c;090e
// Triggered: when the state of required components changes. \n GPS component is
//   triggered when the availability of the GPS of the drone changes. \n
//   Calibration component is triggered when the calibration state of the drone
//   sensors changes \n Mavlink_File component is triggered when the command
//   [StartFlightPlan](#0-11-0) is received. \n Takeoff component is triggered
//   when the drone needs to take-off to continue the FlightPlan. \n
//   WaypointsBeyondGeofence component is triggered when the command
//   [StartFlightPlan](#0-11-0) is received.
// Result:
func (f *flightPlanState) componentStateListChanged(args []interface{}) error {
	// component := args[0].(int32)
	//   Drone FlightPlan component id (unique)
	//   0: GPS: Drone GPS component. State is 0 when the drone needs a GPS fix.
	//   1: Calibration: Drone Calibration component. State is 0 when the sensors
	//      of the drone needs to be calibrated.
	//   2: Mavlink_File: Mavlink file component. State is 0 when the mavlink
	//      file is missing or contains error.
	//   3: TakeOff: Drone Take off component. State is 0 when the drone cannot
	//      take-off.
	//   4: WaypointsBeyondGeofence: Component for waypoints beyond the geofence.
	//      State is 0 when one or more waypoints are beyond the geofence.
	// State := args[1].(uint8)
	//   State of the FlightPlan component (1 FlightPlan component OK, otherwise
	//   0)
	log.Info("common.componentStateListChanged() called")
	return nil
}

// TODO: Implement this
// Title: FlightPlan lock state
// Description: FlightPlan lock state.\n Represents the fact that the controller
//   is able or not to stop or pause a playing FlightPlan
// Support: 0901:2.0.29;090c;090e
// Triggered: when the lock changes.
// Result:
func (f *flightPlanState) lockStateChanged(args []interface{}) error {
	// LockState := args[0].(uint8)
	//   1 if FlightPlan is locked: can&#39;t pause or stop FlightPlan. 0 if
	//   FlightPlan is unlocked: pause or stop available.
	log.Info("common.lockStateChanged() called")
	return nil
}
