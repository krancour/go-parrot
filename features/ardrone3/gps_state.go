package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// GPS related States

const (
	// ReturnToHomeTypeTakeoff represents a return home type wherein the drone
	// will return to the take off position
	ReturnToHomeTypeTakeoff int32 = 0
	// ReturnToHomeTypePilot represents a return home type wherein the drone will
	// return to the pilot position
	ReturnToHomeTypePilot int32 = 1
	// ReturnToHomeTypeFirstFix represents a return home type wherein the drone
	// will return to the first GPS fix
	ReturnToHomeTypeFirstFix int32 = 2
	// ReturnToHomeTypeFollowee represents a return home type wherein the drone
	// will return to the current (or last) "follow me" target
	ReturnToHomeTypeFollowee int32 = 3
)

// GPSState ...
// TODO: Document this
type GPSState interface {
	lock.ReadLockable
	// NumberOfSatellites returns the number of satellites used to determine GPS
	// coordinates and a boolean value indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	NumberOfSatellites() (uint8, bool)
	// IsReturnToTakeoffAvailable returns a boolean indicating whether the device
	// has enough information to return to the takeoff position. A second boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	IsReturnToTakeoffAvailable() (bool, bool)
	// IsReturnToTakeoffAvailable returns a boolean indicating whether the device
	// has enough information to return to the pilot's position. A second boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	IsReturnToPilotAvailable() (bool, bool)
	// IsReturnToTakeoffAvailable returns a boolean indicating whether the device
	// has enough information to return to the first GPS fix. A second boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	IsReturnToFirstFixAvailable() (bool, bool)
	// IsReturnToTakeoffAvailable returns a boolean indicating whether the device
	// has enough information to return to the current (or last) "follow me"
	// target. A second boolean value is also returned, indicating whether the
	// first value was reported by the device (true) or a default value (false).
	// This permits callers to distinguish real zero values from default zero
	// values.
	IsReturnToFolloweeAvailable() (bool, bool)
	// ReturnToHomeType returns the current return to home type. A boolean value
	// is also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	ReturnToHomeType() (int32, bool)
}

type gpsState struct {
	sync.RWMutex
	// numberOfSatellites is the number of satellites used to determine GPS
	// coordinates.
	numberOfSatellites        *uint8
	returnToTakeoffAvailable  *bool
	returnToPilotAvailable    *bool
	returnToFirstFixAvailable *bool
	returnToFolloweeAvailable *bool
	returnToHomeType          *int32
}

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
			"NumberOfSatellitesChanged",
			[]interface{}{
				uint8(0), // numberOfSatellites,
			},
			g.numberOfSatellitesChanged,
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

// numberOfSatellitesChanged is invoked when the the device reports that the
// number of satellites being used to determine GPS coordinates has changed.
func (g *gpsState) numberOfSatellitesChanged(args []interface{}) error {
	g.Lock()
	defer g.Unlock()
	g.numberOfSatellites = ptr.ToUint8(args[0].(uint8))
	log.WithField(
		"numberOfSatellites", *g.numberOfSatellites,
	).Debug("gps state number of satellites updated")
	return nil
}

// homeTypeAvailabilityChanged is invoked by the device when the availability
// of different return to home types is changed.
func (g *gpsState) homeTypeAvailabilityChanged(args []interface{}) error {
	g.Lock()
	defer g.Unlock()
	tipe := args[0].(int32)
	available := args[1].(uint8) == 1
	switch tipe {
	case ReturnToHomeTypeTakeoff:
		g.returnToTakeoffAvailable = ptr.ToBool(available)
	case ReturnToHomeTypePilot:
		g.returnToPilotAvailable = ptr.ToBool(available)
	case ReturnToHomeTypeFirstFix:
		g.returnToFirstFixAvailable = ptr.ToBool(available)
	case ReturnToHomeTypeFollowee:
		g.returnToFolloweeAvailable = ptr.ToBool(available)
	}
	log.WithField(
		"type", tipe,
	).WithField(
		"available", available,
	).Debug("home type availability changed")
	return nil
}

// homeTypeChosenChanged is invoked by the device when the return to home type
// is changed. This may be changed due to a change in user preferences or
// because of the availability of various types based on GPS information
// available. Note that if the user's preferred type is unavailable, the device
// will choose the first available type in this order: FOLLOWEE, TAKEOFF, PILOT,
// FIRST_FIX.
func (g *gpsState) homeTypeChosenChanged(args []interface{}) error {
	g.Lock()
	defer g.Unlock()
	g.returnToHomeType = ptr.ToInt32(args[0].(int32))
	log.WithField(
		"type", *g.returnToHomeType,
	).Debug("return home type changed")
	return nil
}

func (g *gpsState) NumberOfSatellites() (uint8, bool) {
	if g.numberOfSatellites == nil {
		return 0, false
	}
	return *g.numberOfSatellites, true
}

func (g *gpsState) IsReturnToTakeoffAvailable() (bool, bool) {
	if g.returnToTakeoffAvailable == nil {
		return false, false
	}
	return *g.returnToTakeoffAvailable, true
}

func (g *gpsState) IsReturnToPilotAvailable() (bool, bool) {
	if g.returnToPilotAvailable == nil {
		return false, false
	}
	return *g.returnToPilotAvailable, true
}

func (g *gpsState) IsReturnToFirstFixAvailable() (bool, bool) {
	if g.returnToFirstFixAvailable == nil {
		return false, false
	}
	return *g.returnToFirstFixAvailable, true
}

func (g *gpsState) IsReturnToFolloweeAvailable() (bool, bool) {
	if g.returnToFolloweeAvailable == nil {
		return false, false
	}
	return *g.returnToFolloweeAvailable, true
}

func (g *gpsState) ReturnToHomeType() (int32, bool) {
	if g.returnToHomeType == nil {
		return 0, false
	}
	return *g.returnToHomeType, true
}
