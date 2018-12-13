package ardrone3

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// GPS related States

// GPSState ...
// TODO: Document this
type GPSState interface{}

type gpsState struct{}

func (g *gpsState) ID() uint8 {
	return 31
}

func (g *gpsState) Name() string {
	return "GPSState"
}

func (g *gpsState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"NumberOfSatelliteChanged",
			[]interface{}{
				uint8(0), // numberOfSatellite,
			},
			g.numberOfSatelliteChanged,
		),
		arcommands.NewD2CCommand(
			1,
			"HomeTypeAvailabilityChanged",
			[]interface{}{
				int32(0), // type,
				uint8(0), // available,
			},
			g.homeTypeAvailabilityChanged,
		),
		arcommands.NewD2CCommand(
			2,
			"HomeTypeChosenChanged",
			[]interface{}{
				int32(0), // type,
			},
			g.homeTypeChosenChanged,
		),
	}
}

// TODO: Implement this
// Title: Number of GPS satellites
// Description: Number of GPS satellites.
// Support: 0901;090c;090e
// Triggered: on change.
// Result:
func (g *gpsState) numberOfSatelliteChanged(args []interface{}) error {
	// numberOfSatellite := args[0].(uint8)
	//   The number of satellite
	log.Info("ardrone3.numberOfSatelliteChanged() called")
	return nil
}

// TODO: Implement this
// Title: Home type availability
// Description: Home type availability.
// Support: 0901;090c;090e
// Triggered: when the availability of, at least, one type changes.\n This might
//   be due to controller position availability, gps fix before take off or
//   other reason.
// Result:
func (g *gpsState) homeTypeAvailabilityChanged(args []interface{}) error {
	// type := args[0].(int32)
	//   The type of the return home
	//   0: TAKEOFF: The drone has enough information to return to the take off
	//      position
	//   1: PILOT: The drone has enough information to return to the pilot
	//      position
	//   2: FIRST_FIX: The drone has not enough information, it will return to the
	//      first GPS fix
	//   3: FOLLOWEE: The drone has enough information to return to the target of
	//      the current (or last) follow me
	// available := args[1].(uint8)
	//   1 if this type is available, 0 otherwise
	log.Info("ardrone3.homeTypeAvailabilityChanged() called")
	return nil
}

// TODO: Implement this
// Title: Home type
// Description: Home type.\n This choice is made by the drone, according to the
//   [PreferredHomeType](#1-24-4) and the [HomeTypeAvailability](#1-31-1). The
//   drone will choose the type matching with the user preference only if this
//   type is available. If not, it will chose a type in this order:\n FOLLOWEE ;
//   TAKEOFF ; PILOT ; FIRST_FIX
// Support: 0901;090c;090e
// Triggered: when the return home type chosen by the drone changes.\n This
//   might be produced by a user preference triggered by
//   [SetPreferedHomeType](#1-23-3) or by a change in the
//   [HomeTypesAvailabilityChanged](#1-31-1).
// Result:
func (g *gpsState) homeTypeChosenChanged(args []interface{}) error {
	// type := args[0].(int32)
	//   The type of the return home chosen
	//   0: TAKEOFF: The drone will return to the take off position
	//   1: PILOT: The drone will return to the pilot position In this case, the
	//      drone will use the position given by ARDrone3-SendControllerGPS
	//   2: FIRST_FIX: The drone has not enough information, it will return to the
	//      first GPS fix
	//   3: FOLLOWEE: The drone will return to the target of the current (or last)
	//      follow me In this case, the drone will use the position of the target
	//      of the followMe (given by ControllerInfo-GPS)
	log.Info("ardrone3.homeTypeChosenChanged() called")
	return nil
}
