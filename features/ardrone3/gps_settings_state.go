package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// GPS settings state

const (
	// HomeTypeTakeoff represents a setting where the drone will return to the
	// takeoff position when instructed to return home.
	HomeTypeTakeoff int32 = 0
	// HomeTypePilot represents a setting where the drone will return to the
	// pilot's position when instructed to return home.
	HomeTypePilot int32 = 1
	// HomeTypeFollowee represents a setting where the drone will return to the
	// current (or last) "follow me" target when instructed to return home.
	HomeTypeFollowee int32 = 2
)

// GPSSettingsState ...
// TODO: Document this
type GPSSettingsState interface {
	lock.ReadLockable
	// HomeLatitude returns the home (the place the drone will return to) latitude
	// in degrees. A boolean value is also returned, indicating whether the first
	// value was reported by the device (true) or a default value (false). This
	// permits callers to distinguish real zero values from default zero values.
	HomeLatitude() (float64, bool)
	// HomeLatitude returns the home (the place the drone will return to)
	// longitude in degrees. A boolean value is also returned, indicating whether
	// the first value was reported by the device (true) or a default value
	// (false). This permits callers to distinguish real zero values from default
	// zero values.
	HomeLongitude() (float64, bool)
	// HomeLatitude returns the home (the place the drone will return to) altitude
	// in meters. A boolean value is also returned, indicating whether the first
	// value was reported by the device (true) or a default value (false). This
	// permits callers to distinguish real zero values from default zero values.
	HomeAltitude() (float64, bool)
	// HomeType returns the home type, which represents whether the drone should
	// return to the takeoff position, pilot position, or followee when instructed
	// to return home. A boolean value is also returned, indicating whether the
	// first value was reported by the device (true) or a default value (false).
	// This permits callers to distinguish real zero values from default zero
	// values.
	HomeType() (int32, bool)
	// IsGPSFixed returns a boolean indicating whether the device's GPS is
	// currently fixed. A second boolean value is also returned, indicating
	// whether the first value was reported by the device (true) or a default
	// value (false). This permits callers to distinguish real zero values from
	// default zero values.
	IsGPSFixed() (bool, bool)
	// ReturnHomeDelay returns the return home delay (the time after which return
	// to home is automatically triggered after a disconnection) in seconds. A
	// boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	ReturnHomeDelay() (uint16, bool)
}

type gpsSettingsState struct {
	sync.RWMutex
	// homeLatitude is the home latitude in degrees
	homeLatitude *float64
	// homeLongitude is the home longitude in degrees
	homeLongitude *float64
	// homeAltitude is the home altitude in meters
	homeAltitude *float64
	// homeType represents whether the drone should return to the takeoff
	// position, pilot position, or followee when instructed to return home
	homeType *int32
	// gpsFixed represents whether the device's GPS is currently fixed
	gpsFixed *bool
	// returnHomeDelay represents the return home delay (the time after which
	// return to home is automatically triggered after a disconnection) in
	// seconds
	returnHomeDelay *uint16
}

func newGPSSettingsState() *gpsSettingsState {
	return &gpsSettingsState{}
}

func (g *gpsSettingsState) ClassID() uint8 {
	return 24
}

func (g *gpsSettingsState) ClassName() string {
	return "GPSSettingsState"
}

func (g *gpsSettingsState) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
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
			log,
		),
		arcommands.NewD2CCommand(
			2,
			"GPSFixStateChanged",
			[]interface{}{
				uint8(0), // fixed,
			},
			g.gPSFixStateChanged,
			log,
		),
		arcommands.NewD2CCommand(
			4,
			"HomeTypeChanged",
			[]interface{}{
				int32(0), // type,
			},
			g.homeTypeChanged,
			log,
		),
		arcommands.NewD2CCommand(
			5,
			"ReturnHomeDelayChanged",
			[]interface{}{
				uint16(0), // delay,
			},
			g.returnHomeDelayChanged,
			log,
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

// homeChanged is invoked by the device when the home location (the place the
// drone will return to) is changed.
func (g *gpsSettingsState) homeChanged(
	args []interface{},
	log *log.Entry,
) error {
	g.Lock()
	defer g.Unlock()
	g.homeLatitude = ptr.ToFloat64(args[0].(float64))
	g.homeLongitude = ptr.ToFloat64(args[1].(float64))
	g.homeAltitude = ptr.ToFloat64(args[2].(float64))
	log.WithField(
		"latitude", *g.homeLatitude,
	).WithField(
		"longitude", *g.homeLongitude,
	).WithField(
		"altitude", *g.homeAltitude,
	).Debug("home location changed")
	return nil
}

// gPSFixStateChanged is invoked by the device when the GPS fix is changed.
func (g *gpsSettingsState) gPSFixStateChanged(
	args []interface{},
	log *log.Entry,
) error {
	g.Lock()
	defer g.Unlock()
	g.gpsFixed = ptr.ToBool(args[0].(uint8) == 1)
	log.WithField(
		"fixed", *g.gpsFixed,
	).Debug("GPS fixed state changed")
	return nil
}

// homeTypeChanged is invoked by the device when the home type is changed.
func (g *gpsSettingsState) homeTypeChanged(
	args []interface{},
	log *log.Entry,
) error {
	g.Lock()
	defer g.Unlock()
	g.homeType = ptr.ToInt32(args[0].(int32))
	log.WithField(
		"type", *g.homeType,
	).Debug("home type changed")
	return nil
}

// returnHomeDelayChanged is triggered by the device when the return home delay
// (the time after which return to home is automatically triggered after a
// disconnection) is changed.
func (g *gpsSettingsState) returnHomeDelayChanged(
	args []interface{},
	log *log.Entry,
) error {
	g.Lock()
	defer g.Unlock()
	g.returnHomeDelay = ptr.ToUint16(args[0].(uint16))
	log.WithField(
		"delay", *g.returnHomeDelay,
	).Debug("return home delay changed")
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
// 	log.Warn("command not implemented")
// 	return nil
// }

func (g *gpsSettingsState) HomeLatitude() (float64, bool) {
	if g.homeLatitude == nil {
		return 0, false
	}
	return *g.homeLatitude, true
}

func (g *gpsSettingsState) HomeLongitude() (float64, bool) {
	if g.homeLongitude == nil {
		return 0, false
	}
	return *g.homeLongitude, true
}

func (g *gpsSettingsState) HomeAltitude() (float64, bool) {
	if g.homeAltitude == nil {
		return 0, false
	}
	return *g.homeAltitude, true
}

func (g *gpsSettingsState) HomeType() (int32, bool) {
	if g.homeType == nil {
		return 0, false
	}
	return *g.homeType, true
}

func (g *gpsSettingsState) IsGPSFixed() (bool, bool) {
	if g.gpsFixed == nil {
		return false, false
	}
	return *g.gpsFixed, true
}

func (g *gpsSettingsState) ReturnHomeDelay() (uint16, bool) {
	if g.returnHomeDelay == nil {
		return 0, false
	}
	return *g.returnHomeDelay, true
}
