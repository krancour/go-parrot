package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// State from drone

const (
	// FlyingStateLanded represents the state where the device is landed.
	FlyingStateLanded int32 = 0
	// FlyingStateTakingOff represents the state where the device is taking off.
	FlyingStateTakingOff int32 = 1
	// FlyingStateHovering represents the state where the device is hovering
	// (copters) or circling (fixed wing).
	FlyingStateHovering int32 = 2
	// FlyingStateFlying represents the state where the device is flying.
	FlyingStateFlying int32 = 3
	// FlyingStateLanding represents the state where the device is landing.
	FlyingStateLanding int32 = 4
	// FlyingStateEmergency represents the state where the device is experiencing
	// an emergency.
	FlyingStateEmergency int32 = 5
	// FlyingStateUserTakeOff represents the state where the device is waiting for
	// user action to effect take off.
	FlyingStateUserTakeOff int32 = 6
	// FlyingStateMotorRamping represents the motor ramping state.
	FlyingStateMotorRamping int32 = 7
	// FlyingStateEmergencyLanding represents the state where the device has
	// detected faulty sensor(s) and is effecting an emergency landing. When in
	// this state, the device will only respond to the PCMD command's ya argument.
	// All other flying commands are ignored.
	FlyingStateEmergencyLanding int32 = 8
)

// PilotingState ...
// TODO: Document this
type PilotingState interface {
	lock.ReadLockable
	// SpeedX returns the velocity relative to the north in m/s. When the drone
	// moves to the north, the value is > 0. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	SpeedX() (float32, bool)
	// SpeedY returns the velocity relative to the east in m/s. When the drone
	// moves to the east, the value is > 0. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	SpeedY() (float32, bool)
	// SpeedZ returns velocity relative to the ground in m/s. When the drone moves
	// down, the value is > 0. A boolean value is also returned, indicating
	// whether the first value was reported by the device (true) or a default
	// value (false). This permits callers to distinguish real zero values from
	// default zero values.
	SpeedZ() (float32, bool)
	// Altitude returns the altitude, relative to the take off point, in meters.
	// A boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	Altitude() (float64, bool)
	// Roll returns the roll in radians. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	Roll() (float32, bool)
	// Pitch returns the pitch in radians. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	Pitch() (float32, bool)
	// Yaw returns the yaw in radians. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	Yaw() (float32, bool)
	// Latitude returns the latitude, as determined by GPS, in degrees or 500.0 if
	// unavailable. A boolean value is also returned, indicating whether the first
	// value was reported by the device (true) or a default value (false). This
	// permits callers to distinguish real zero values from default zero values.
	Latitude() (float64, bool)
	// Longitude returns the longitude, as determined by GPS, in degrees or 500.0
	// if unavailable. A boolean value is also returned, indicating whether the
	// first value was reported by the device (true) or a default value (false).
	// This permits callers to distinguish real zero values from default zero
	// values.
	Longitude() (float64, bool)
	// GPSAltitude returns the altitude relative to sea level, as determined by
	// GPS, in meters. A boolean value is also returned, indicating whether the
	// first value was reported by the device (true) or a default value (false).
	// This permits callers to distinguish real zero values from default zero
	// values.
	GPSAltitude() (float64, bool)
	// LatitudeAccuracy returns GPS latitude location error (1 sigma/standard
	// deviation) in meters or -1 if unavailable. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	LatitudeAccuracy() (int8, bool)
	// LongitudeAccuracy returns GPS longitude location error (1 sigma/standard
	// deviation) in meters or -1 if unavailable. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	LongitudeAccuracy() (int8, bool)
	// GPSAltitudeAccuracy returns GPS altitude location error (1 sigma/standard
	// deviation) in meters or -1 if unavailable. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	GPSAltitudeAccuracy() (int8, bool)
	// FlyingState returns the current flying state of the drone-- whether it is
	// currently landed, taking off, landing, etc. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	FlyingState() (int32, bool)
}

type pilotingState struct {
	sync.RWMutex
	// speedX is velocity relative to the north in m/s. When the drone moves to
	// the north, the value is > 0
	speedX *float32
	// speedY is velocity relative to the east in m/s. When the drone moves to the
	// east, the value is > 0
	speedY *float32
	// speedZ is velocity relative to the ground in m/s. When the drone moves
	// down, the value is > 0
	speedZ *float32
	// altitude, relative to the take off point, in meters
	altitude *float64
	// roll in radians
	roll *float32
	// pitch in radians
	pitch *float32
	// yaw in radians
	yaw *float32
	// latitude, as determined by GPS, in degrees (500.0 if unavailable)
	latitude *float64
	// longitude, as determined by GPS, in degrees (500.0 if unavailable)
	longitude *float64
	// gpsAltitude is altitude relative to sea level, as determined by GPS, in
	// meters
	gpsAltitude *float64
	// latitudeAccuracy represents latitude location error (1 sigma/standard
	// deviation) in meters (-1 if unavailable)
	latitudeAccuracy *int8
	// longitudeAccuracy represents longitude location error (1 sigma/standard
	// deviation) in meters (-1 if unavailable)
	longitudeAccuracy *int8
	// gpsAltitudeAccuracy represents altitude location error (1 sigma/standard
	// deviation) in meters (-1 if unavailable)
	gpsAltitudeAccuracy *int8
	// flyingState represents the current flying state of the device-- whether it
	// is landed, taking off, landed, etc.
	flyingState *int32
}

func newPilotingState() *pilotingState {
	return &pilotingState{}
}

func (p *pilotingState) ID() uint8 {
	return 4
}

func (p *pilotingState) Name() string {
	return "PilotingState"
}

func (p *pilotingState) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"FlatTrimChanged",
			[]interface{}{},
			p.flatTrimChanged,
			log,
		),
		arcommands.NewD2CCommand(
			1,
			"FlyingStateChanged",
			[]interface{}{
				int32(0), // state,
			},
			p.flyingStateChanged,
			log,
		),
		arcommands.NewD2CCommand(
			2,
			"AlertStateChanged",
			[]interface{}{
				int32(0), // state,
			},
			p.alertStateChanged,
			log,
		),
		arcommands.NewD2CCommand(
			3,
			"NavigateHomeStateChanged",
			[]interface{}{
				int32(0), // state,
				int32(0), // reason,
			},
			p.navigateHomeStateChanged,
			log,
		),
		// According to Parrot developer documentation, they seem poised to replace
		// PositionChanged with GpsLocationChanged. Since GpsLocationChanged is
		// functioning as expected, but we can still see this command being invoked,
		// we'll implement the command to avoid a warning, but the implementation
		// will remain a no-op unless / until such time that it becomes clear that
		// older versions of the firmware might require us to support both commands.
		arcommands.NewD2CCommand(
			4,
			"PositionChanged",
			[]interface{}{
				float64(0), // latitude,
				float64(0), // longitude,
				float64(0), // altitude,
			},
			p.positionChanged,
			log,
		),
		arcommands.NewD2CCommand(
			5,
			"SpeedChanged",
			[]interface{}{
				float32(0), // speedX,
				float32(0), // speedY,
				float32(0), // speedZ,
			},
			p.speedChanged,
			log,
		),
		arcommands.NewD2CCommand(
			6,
			"AttitudeChanged",
			[]interface{}{
				float32(0), // roll,
				float32(0), // pitch,
				float32(0), // yaw,
			},
			p.attitudeChanged,
			log,
		),
		arcommands.NewD2CCommand(
			8,
			"AltitudeChanged",
			[]interface{}{
				float64(0), // altitude,
			},
			p.altitudeChanged,
			log,
		),
		arcommands.NewD2CCommand(
			9,
			"GpsLocationChanged",
			[]interface{}{
				float64(0), // latitude,
				float64(0), // longitude,
				float64(0), // altitude,
				int8(0),    // latitude_accuracy,
				int8(0),    // longitude_accuracy,
				int8(0),    // altitude_accuracy,
			},
			p.gpsLocationChanged,
			log,
		),
		// arcommands.NewD2CCommand(
		// 	11,
		// 	"AirSpeedChanged",
		// 	[]interface{}{
		// 		float32(0), // airSpeed,
		// 	},
		// 	p.airSpeedChanged,
		// ),
		arcommands.NewD2CCommand(
			12,
			"moveToChanged",
			[]interface{}{
				float64(0), // latitude,
				float64(0), // longitude,
				float64(0), // altitude,
				int32(0),   // orientation_mode,
				float32(0), // heading,
				int32(0),   // status,
			},
			p.moveToChanged,
			log,
		),
		arcommands.NewD2CCommand(
			13,
			"MotionState",
			[]interface{}{
				int32(0), // state,
			},
			p.motionState,
			log,
		),
		arcommands.NewD2CCommand(
			14,
			"PilotedPOI",
			[]interface{}{
				float64(0), // latitude,
				float64(0), // longitude,
				float64(0), // altitude,
				int32(0),   // status,
			},
			p.pilotedPOI,
			log,
		),
		// arcommands.NewD2CCommand(
		// 	15,
		// 	"ReturnHomeBatteryCapacity",
		// 	[]interface{}{
		// 		int32(0), // status,
		// 	},
		// 	p.returnHomeBatteryCapacity,
		// ),
	}
}

// TODO: Implement this
// Title: Flat trim changed
// Description: Drone acknowledges that flat trim was correctly processed.
// Support: 0901;090c;090e
// Triggered: by [FlatTrim](#1-0-0).
// Result:
func (p *pilotingState) flatTrimChanged(
	args []interface{},
	log *log.Entry,
) error {
	log.Warn("command not implemented")
	return nil
}

// flyingStateChanged is invoked by the device when the flying state changes
func (p *pilotingState) flyingStateChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.flyingState = ptr.ToInt32(args[0].(int32))
	log.WithField(
		"flyingState", *p.flyingState,
	).Debug("flying state changed")
	return nil
}

// TODO: Implement this
// Title: Alert state
// Description: Alert state.
// Support: 0901;090c;090e
// Triggered: when an alert happens on the drone.
// Result:
func (p *pilotingState) alertStateChanged(
	args []interface{},
	log *log.Entry,
) error {
	// state := args[0].(int32)
	//   Drone alert state
	//   0: none: No alert
	//   1: user: User emergency alert
	//   2: cut_out: Cut out alert
	//   3: critical_battery: Critical battery alert
	//   4: low_battery: Low battery alert
	//   5: too_much_angle: The angle of the drone is too high
	log.Warn("command not implemented")
	return nil
}

// TODO: Implement this
// Title: Return home state
// Description: Return home state.\n Availability is related to gps fix,
//   magnetometer calibration.
// Support: 0901;090c;090e
// Triggered: by [ReturnHome](#1-0-5) or when the state of the return home
//   changes.
// Result:
func (p *pilotingState) navigateHomeStateChanged(
	args []interface{},
	log *log.Entry,
) error {
	// state := args[0].(int32)
	//   State of navigate home
	//   0: available: Navigate home is available
	//   1: inProgress: Navigate home is in progress
	//   2: unavailable: Navigate home is not available
	//   3: pending: Navigate home has been received, but its process is pending
	// reason := args[1].(int32)
	//   Reason of the state
	//   0: userRequest: User requested a navigate home (available-&gt;inProgress)
	//   1: connectionLost: Connection between controller and product lost
	//      (available-&gt;inProgress)
	//   2: lowBattery: Low battery occurred (available-&gt;inProgress)
	//   3: finished: Navigate home is finished (inProgress-&gt;available)
	//   4: stopped: Navigate home has been stopped (inProgress-&gt;available)
	//   5: disabled: Navigate home disabled by product
	//      (inProgress-&gt;unavailable or available-&gt;unavailable)
	//   6: enabled: Navigate home enabled by product (unavailable-&gt;available)
	log.Warn("command not implemented")
	return nil
}

// positionChanged does not appear deprecated (yet), according to Parrot
// developer documentation, but according to the same documentation, they appear
// poised to replace PositionChanged with GpsLocationChanged. Since
// GpsLocationChanged is functioning as expected, but we can still see this
// command being invoked, we'll implement the command to avoid a warning, but
// the implementation will remain a no-op unless / until such time that it
// becomes clear that older versions of the firmware might require us to support
// both commands.
func (p *pilotingState) positionChanged(
	args []interface{},
	log *log.Entry,
) error {
	// latitude := args[0].(float64)
	//   Latitude position in decimal degrees (500.0 if not available)
	// longitude := args[1].(float64)
	//   Longitude position in decimal degrees (500.0 if not available)
	// altitude := args[2].(float64)
	//   Altitude in meters (from GPS)
	log.Debug("piloting state position changed-- this is a no-op")
	return nil
}

// speedChanged is invoked when the the device reports velocity at regular
// intervals.
func (p *pilotingState) speedChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.speedX = ptr.ToFloat32(args[0].(float32))
	p.speedY = ptr.ToFloat32(args[1].(float32))
	p.speedZ = ptr.ToFloat32(args[2].(float32))
	log.WithField(
		"speedX", *p.speedX,
	).WithField(
		"speedY", *p.speedY,
	).WithField(
		"speedZ", *p.speedZ,
	).Debug("piloting state speed updated")
	return nil
}

// attitudeChanged is invoked when the device reports attitude at regular
// intervals.
func (p *pilotingState) attitudeChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.roll = ptr.ToFloat32(args[0].(float32))
	p.pitch = ptr.ToFloat32(args[1].(float32))
	p.yaw = ptr.ToFloat32(args[2].(float32))
	log.WithField(
		"roll", *p.roll,
	).WithField(
		"pitch", *p.pitch,
	).WithField(
		"yaw", *p.yaw,
	).Debug("piloting state attitude updated")
	return nil
}

// altitudeChanged is invoked when the device reports attitude relative to the
// take off point at regular intervals.
func (p *pilotingState) altitudeChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.altitude = ptr.ToFloat64(args[0].(float64))
	log.WithField(
		"altitude", *p.altitude,
	).Debug("piloting state altitude updated")
	return nil
}

// gpsLocationChanged is invoked when the device reports gps coordinates at
// regular intervals.
func (p *pilotingState) gpsLocationChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.latitude = ptr.ToFloat64(args[0].(float64))
	p.longitude = ptr.ToFloat64(args[1].(float64))
	p.gpsAltitude = ptr.ToFloat64(args[2].(float64))
	p.latitudeAccuracy = ptr.ToInt8(args[3].(int8))
	p.longitudeAccuracy = ptr.ToInt8(args[4].(int8))
	p.gpsAltitudeAccuracy = ptr.ToInt8(args[5].(int8))
	log.WithField(
		"latitude", *p.latitude,
	).WithField(
		"longitude", *p.longitude,
	).WithField(
		"altitude", *p.gpsAltitude,
	).WithField(
		"latitudeAccuracy", *p.latitudeAccuracy,
	).WithField(
		"longitudeAccuracy", *p.longitudeAccuracy,
	).WithField(
		"altitudeAccuracy", *p.gpsAltitudeAccuracy,
	).Debug("piloting state gps coordinates updated")
	return nil
}

// // TODO: Implement this
// // Title: Drone&#39;s air speed changed
// // Description: Drone&#39;s air speed changed\n Expressed in the drone&#39;s
// //   referential.
// // Support: 090e:1.2.0
// // Triggered: regularly.
// // Result:
// func (p *pilotingState) airSpeedChanged(args []interface{}) error {
// 	// airSpeed := args[0].(float32)
// 	//   Speed relative to air on x axis (speed is always &gt; 0) (in m/s)
// 	log.Warn("command not implemented")
// 	return nil
// }

// TODO: Implement this
// Title: Move to changed
// Description: The drone moves or moved to a given location.
// Support: 090c:4.3.0
// Triggered: by [MoveTo](#1-0-10) or when the drone did reach the given
//   position.
// Result:
func (p *pilotingState) moveToChanged(
	args []interface{},
	log *log.Entry,
) error {
	// latitude := args[0].(float64)
	//   Latitude of the location (in degrees) to reach
	// longitude := args[1].(float64)
	//   Longitude of the location (in degrees) to reach
	// altitude := args[2].(float64)
	//   Altitude above sea level (in m) to reach
	// orientation_mode := args[3].(int32)
	//   Orientation mode of the move to
	//   0: NONE: The drone won&#39;t change its orientation
	//   1: TO_TARGET: The drone will make a rotation to look in direction of the
	//      given location
	//   2: HEADING_START: The drone will orientate itself to the given heading
	//      before moving to the location
	//   3: HEADING_DURING: The drone will orientate itself to the given heading
	//      while moving to the location
	// heading := args[4].(float32)
	//   Heading (relative to the North in degrees). This value is only used if
	//   the orientation mode is HEADING_START or HEADING_DURING
	// status := args[5].(int32)
	//   Status of the move to
	//   0: RUNNING: The drone is actually flying to the given position
	//   1: DONE: The drone has reached the target
	//   2: CANCELED: The move to has been canceled, either by a CancelMoveTo
	//      command or when a disconnection appears.
	//   3: ERROR: The move to has not been finished or started because of an
	//      error.
	log.Warn("command not implemented")
	return nil
}

// TODO: Implement this
// Title: Motion state
// Description: Motion state.\n If [MotionDetection](#1-6-16) is disabled,
//   motion is steady.\n This information is only valid when the drone is not
//   flying.
// Support: 090c:4.3.0
// Triggered: when the [FlyingState](#1-4-1) is landed and the
//   [MotionDetection](#1-6-16) is enabled and the motion state changes.\n This
//   event is triggered at a filtered rate.
// Result:
func (p *pilotingState) motionState(args []interface{}, log *log.Entry) error {
	// state := args[0].(int32)
	//   Motion state
	//   0: steady: Drone is steady
	//   1: moving: Drone is moving
	log.Warn("command not implemented")
	return nil
}

// TODO: Implement this
// Title: Piloted POI state
// Description: Piloted POI state.
// Support: 090c:4.3.0
// Triggered: by [StartPilotedPOI](#1-0-12) or [StopPilotedPOI](#1-0-13) or when
//   piloted POI becomes unavailable.
// Result:
func (p *pilotingState) pilotedPOI(args []interface{}, log *log.Entry) error {
	// latitude := args[0].(float64)
	//   Latitude of the location (in degrees) to look at. This information is
	//   only valid when the state is pending or running.
	// longitude := args[1].(float64)
	//   Longitude of the location (in degrees) to look at. This information is
	//   only valid when the state is pending or running.
	// altitude := args[2].(float64)
	//   Altitude above sea level (in m) to look at. This information is only
	//   valid when the state is pending or running.
	// status := args[3].(int32)
	//   Status of the move to
	//   0: UNAVAILABLE: The piloted POI is not available
	//   1: AVAILABLE: The piloted POI is available
	//   2: PENDING: Piloted POI has been requested. Waiting to be in state that
	//      allow the piloted POI to start
	//   3: RUNNING: Piloted POI is running
	log.Warn("command not implemented")
	return nil
}

// // TODO: Implement this
// // Title: Return home battery capacity
// // Description: Battery capacity status to return home.
// // Support: 090c:4.3.0
// // Triggered: when the status of the battery capacity to do a return home
// //   changes. This means that it is triggered either when the battery level
// //   changes, when the distance to the home changes or when the position of the
// //   home changes.
// // Result:
// func (p *pilotingState) returnHomeBatteryCapacity(args []interface{}) error {
// 	// status := args[0].(int32)
// 	//   Status of battery to return home
// 	//   0: OK: The battery is full enough to do a return home
// 	//   1: WARNING: The battery is about to be too discharged to do a return home
// 	//   2: CRITICAL: The battery level is too low to return to the home position
// 	//   3: UNKNOWN: Battery capacity to do a return home is unknown. This can be
// 	//      either because the home is unknown or the position of the drone is
// 	//      unknown, or the drone has not enough information to determine how long
// 	//      it takes to fly home.
// 	log.Warn("command not implemented")
// 	return nil
// }

func (p *pilotingState) SpeedX() (float32, bool) {
	if p.speedX == nil {
		return 0, false
	}
	return *p.speedX, true
}

func (p *pilotingState) SpeedY() (float32, bool) {
	if p.speedY == nil {
		return 0, false
	}
	return *p.speedY, true
}

func (p *pilotingState) SpeedZ() (float32, bool) {
	if p.speedZ == nil {
		return 0, false
	}
	return *p.speedZ, true
}

func (p *pilotingState) Altitude() (float64, bool) {
	if p.altitude == nil {
		return 0, false
	}
	return *p.altitude, true
}

func (p *pilotingState) Roll() (float32, bool) {
	if p.roll == nil {
		return 0, false
	}
	return *p.roll, true
}

func (p *pilotingState) Pitch() (float32, bool) {
	if p.pitch == nil {
		return 0, false
	}
	return *p.pitch, true
}

func (p *pilotingState) Yaw() (float32, bool) {
	if p.yaw == nil {
		return 0, false
	}
	return *p.yaw, true
}

func (p *pilotingState) Latitude() (float64, bool) {
	if p.latitude == nil {
		return 0, false
	}
	return *p.latitude, true
}

func (p *pilotingState) Longitude() (float64, bool) {
	if p.longitude == nil {
		return 0, false
	}
	return *p.longitude, true
}

func (p *pilotingState) GPSAltitude() (float64, bool) {
	if p.gpsAltitude == nil {
		return 0, false
	}
	return *p.gpsAltitude, true
}

func (p *pilotingState) LatitudeAccuracy() (int8, bool) {
	if p.latitudeAccuracy == nil {
		return 0, false
	}
	return *p.latitudeAccuracy, true
}

func (p *pilotingState) LongitudeAccuracy() (int8, bool) {
	if p.longitudeAccuracy == nil {
		return 0, false
	}
	return *p.longitudeAccuracy, true
}

func (p *pilotingState) GPSAltitudeAccuracy() (int8, bool) {
	if p.gpsAltitudeAccuracy == nil {
		return 0, false
	}
	return *p.gpsAltitudeAccuracy, true
}

func (p *pilotingState) FlyingState() (int32, bool) {
	if p.flyingState == nil {
		return 0, false
	}
	return *p.flyingState, true
}
