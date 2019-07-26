package common

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// FlightPlan state commands

// FlightPlanState ...
// TODO: Document this
type FlightPlanState interface {
	lock.ReadLockable
	GPSOK() (bool, bool)
	CalibrationOK() (bool, bool)
	MavlinkFileOK() (bool, bool)
	TakeOffOK() (bool, bool)
	WaypointBeyondGeofenceOK() (bool, bool)
	RunFlightPlanFileAvailable() (bool, bool)
}

type flightPlanState struct {
	sync.RWMutex
	gpsOK                      *bool
	calibrationOK              *bool
	mavlinkFileOK              *bool
	takeOffOK                  *bool
	waypointBeyondGeofenceOK   *bool
	runFlightPlanFileAvailable *bool
}

func newFlightPlanState() *flightPlanState {
	return &flightPlanState{}
}

func (f *flightPlanState) ClassID() uint8 {
	return 17
}

func (f *flightPlanState) ClassName() string {
	return "FlightPlanState"
}

func (f *flightPlanState) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"AvailabilityStateChanged",
			[]interface{}{
				uint8(0), // AvailabilityState,
			},
			f.availabilityStateChanged,
			log,
		),
		arcommands.NewD2CCommand(
			1,
			"ComponentStateListChanged",
			[]interface{}{
				int32(0), // component,
				uint8(0), // State,
			},
			f.componentStateListChanged,
			log,
		),
		arcommands.NewD2CCommand(
			2,
			"LockStateChanged",
			[]interface{}{
				uint8(0), // LockState,
			},
			f.lockStateChanged,
			log,
		),
	}
}

// availabilityStateChanged is invoked by the device to indicate whether
// running a flight plan file is available.
func (f *flightPlanState) availabilityStateChanged(
	args []interface{},
	log *log.Entry,
) error {
	f.Lock()
	defer f.Unlock()
	f.runFlightPlanFileAvailable = ptr.ToBool(args[0].(uint8) == 1)
	log.WithField(
		"runFlightPlanFileAvailable", *f.runFlightPlanFileAvailable,
	).Debug("availability of running a flight plan file changed")
	return nil
}

// componentStateListChanged is invoked by the device when the state of various
// device components is changed.
func (f *flightPlanState) componentStateListChanged(
	args []interface{},
	log *log.Entry,
) error {
	f.Lock()
	defer f.Unlock()
	componentID := args[0].(int32)
	state := ptr.ToBool(args[1].(uint8) == 1)
	log = log.WithField(
		"state", *state,
	)
	switch componentID {
	case 0:
		f.gpsOK = state
		log = log.WithField("component", "gps")
	case 1:
		f.calibrationOK = state
		log = log.WithField("component", "calibration")
	case 2:
		f.mavlinkFileOK = state
		log = log.WithField("component", "mavlinkFile")
	case 3:
		f.takeOffOK = state
		log = log.WithField("component", "takeOff")
	case 4:
		f.waypointBeyondGeofenceOK = state
		log = log.WithField("component", "waypointBeyondGeofence")
	default:
		log.WithField(
			"componentID", componentID,
		).Warn("component state changed with unknown component id")
		return nil
	}
	log.Debug("component state changed")
	return nil
}

// TODO: Implement this
// Title: FlightPlan lock state
// Description: FlightPlan lock state.\n Represents the fact that the controller
//   is able or not to stop or pause a playing FlightPlan
// Support: 0901:2.0.29;090c;090e
// Triggered: when the lock changes.
// Result:
func (f *flightPlanState) lockStateChanged(
	args []interface{},
	log *log.Entry,
) error {
	// LockState := args[0].(uint8)
	//   1 if FlightPlan is locked: can&#39;t pause or stop FlightPlan. 0 if
	//   FlightPlan is unlocked: pause or stop available.
	log.Warn("command not implemented")
	return nil
}

func (f *flightPlanState) GPSOK() (bool, bool) {
	if f.gpsOK == nil {
		return false, false
	}
	return *f.gpsOK, true
}

func (f *flightPlanState) CalibrationOK() (bool, bool) {
	if f.calibrationOK == nil {
		return false, false
	}
	return *f.calibrationOK, true
}

func (f *flightPlanState) MavlinkFileOK() (bool, bool) {
	if f.mavlinkFileOK == nil {
		return false, false
	}
	return *f.mavlinkFileOK, true
}

func (f *flightPlanState) TakeOffOK() (bool, bool) {
	if f.takeOffOK == nil {
		return false, false
	}
	return *f.takeOffOK, true
}

func (f *flightPlanState) WaypointBeyondGeofenceOK() (bool, bool) {
	if f.waypointBeyondGeofenceOK == nil {
		return false, false
	}
	return *f.waypointBeyondGeofenceOK, true
}

func (f *flightPlanState) RunFlightPlanFileAvailable() (bool, bool) {
	if f.runFlightPlanFileAvailable == nil {
		return false, false
	}
	return *f.runFlightPlanFileAvailable, true
}
