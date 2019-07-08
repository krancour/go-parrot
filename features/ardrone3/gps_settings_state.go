package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// GPS settings state

// GPSSettingsState ...
// TODO: Document this
type GPSSettingsState interface {
	lock.ReadLockable
}

type gpsSettingsState struct {
	lock sync.RWMutex
}

func (g *gpsSettingsState) ID() uint8 {
	return 24
}

func (g *gpsSettingsState) Name() string {
	return "GPSSettingsState"
}

func (g *gpsSettingsState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"HomeChanged",
			[]interface{}{
				float64(0), // latitude,
				float64(0), // longitude,
				float64(0), // altitude,
			},
			g.homeChanged,
		),
		// arcommands.NewD2CCommand(
		// 	1,
		// 	"ResetHomeChanged",
		// 	[]interface{}{
		// 		float64(0), // latitude,
		// 		float64(0), // longitude,
		// 		float64(0), // altitude,
		// 	},
		// 	g.resetHomeChanged,
		// ),
		arcommands.NewD2CCommand(
			2,
			"GPSFixStateChanged",
			[]interface{}{
				uint8(0), // fixed,
			},
			g.gPSFixStateChanged,
		),
		// arcommands.NewD2CCommand(
		// 	3,
		// 	"GPSUpdateStateChanged",
		// 	[]interface{}{
		// 		int32(0), // state,
		// 	},
		// 	g.gpsUpdateStateChanged,
		// ),
		arcommands.NewD2CCommand(
			4,
			"HomeTypeChanged",
			[]interface{}{
				int32(0), // type,
			},
			g.homeTypeChanged,
		),
		arcommands.NewD2CCommand(
			5,
			"ReturnHomeDelayChanged",
			[]interface{}{
				uint16(0), // delay,
			},
			g.returnHomeDelayChanged,
		),
		// arcommands.NewD2CCommand(
		// 	6,
		// 	"GeofenceCenterChanged",
		// 	[]interface{}{
		// 		float64(0), // latitude,
		// 		float64(0), // longitude,
		// 	},
		// 	g.geofenceCenterChanged,
		// ),
	}
}

// TODO: Implement this
// Title: Home location
// Description: Home location.
// Support: 0901;090c;090e
// Triggered: when [HomeType](#1-31-2) changes. Or by [SetHomeLocation](#1-23-2)
//   when [HomeType](#1-31-2) is Pilot. Or regularly after
//   [SetControllerGPS](#140-1) when [HomeType](#1-31-2) is FollowMeTarget. Or
//   at take off [HomeType](#1-31-2) is Takeoff. Or when the first fix occurs
//   and the [HomeType](#1-31-2) is FirstFix.
// Result:
func (g *gpsSettingsState) homeChanged(args []interface{}) error {
	// latitude := args[0].(float64)
	//   Home latitude in decimal degrees
	// longitude := args[1].(float64)
	//   Home longitude in decimal degrees
	// altitude := args[2].(float64)
	//   Home altitude in meters
	log.Info("ardrone3.homeChanged() called")
	return nil
}

// // TODO: Implement this
// // Title: Home location has been reset
// // Description: Home location has been reset.
// // Support: 0901;090c
// // Triggered: by [ResetHomeLocation](#1-23-1).
// // Result:
// // WARNING: Deprecated
// func (g *gpsSettingsState) resetHomeChanged(args []interface{}) error {
// 	// latitude := args[0].(float64)
// 	//   Home latitude in decimal degrees
// 	// longitude := args[1].(float64)
// 	//   Home longitude in decimal degrees
// 	// altitude := args[2].(float64)
// 	//   Home altitude in meters
// 	log.Info("ardrone3.resetHomeChanged() called")
// 	return nil
// }

// TODO: Implement this
// Title: Gps fix info
// Description: Gps fix info.
// Support: 0901;090c;090e
// Triggered: on change.
// Result:
func (g *gpsSettingsState) gPSFixStateChanged(args []interface{}) error {
	// fixed := args[0].(uint8)
	//   1 if gps on drone is fixed, 0 otherwise
	log.Info("ardrone3.gPSFixStateChanged() called")
	return nil
}

// // TODO: Implement this
// // Title: Gps update state
// // Description: Gps update state.
// // Support: 0901;090c;090e
// // Triggered: on change.
// // Result:
// // WARNING: Deprecated
// func (g *gpsSettingsState) gpsUpdateStateChanged(args []interface{}) error {
// 	// state := args[0].(int32)
// 	//   The state of the gps update
// 	//   0: updated: Drone GPS update succeed
// 	//   1: inProgress: Drone GPS update In progress
// 	//   2: failed: Drone GPS update failed
// 	log.Info("ardrone3.gpsUpdateStateChanged() called")
// 	return nil
// }

// TODO: Implement this
// Title: Preferred home type
// Description: User preference for the home type.\n See [HomeType](#1-31-2) to
//   get the drone actual home type.
// Support: 0901;090c;090e
// Triggered: by [SetPreferredHomeType](#1-23-3).
// Result:
func (g *gpsSettingsState) homeTypeChanged(args []interface{}) error {
	// type := args[0].(int32)
	//   The type of the home position
	//   0: TAKEOFF: The drone will try to return to the take off position
	//   1: PILOT: The drone will try to return to the pilot position
	//   2: FOLLOWEE: The drone will try to return to the target of the current
	//      (or last) follow me
	log.Info("ardrone3.homeTypeChanged() called")
	return nil
}

// TODO: Implement this
// Title: Return home delay
// Description: Return home trigger delay. This delay represents the time after
//   which the return home is automatically triggered after a disconnection.
// Support: 0901;090c;090e
// Triggered: by [SetReturnHomeDelay](#1-23-4).
// Result:
func (g *gpsSettingsState) returnHomeDelayChanged(args []interface{}) error {
	// delay := args[0].(uint16)
	//   Delay in second
	log.Info("ardrone3.returnHomeDelayChanged() called")
	return nil
}

// // TODO: Implement this
// // Title: Geofence center
// // Description: Geofence center location. This location represents the center of
// //   the geofence zone. This is updated at a maximum frequency of 1 Hz.
// // Support:
// // Triggered: when [HomeChanged](#1-24-0) and when [GpsLocationChanged](#1-4-9)
// //   before takeoff.
// // Result:
// func (g *gpsSettingsState) geofenceCenterChanged(args []interface{}) error {
// 	// latitude := args[0].(float64)
// 	//   GPS latitude in decimal degrees
// 	// longitude := args[1].(float64)
// 	//   GPS longitude in decimal degrees
// 	log.Info("ardrone3.geofenceCenterChanged() called")
// 	return nil
// }

func (g *gpsSettingsState) RLock() {
	g.lock.RLock()
}

func (g *gpsSettingsState) RUnlock() {
	g.lock.RUnlock()
}
