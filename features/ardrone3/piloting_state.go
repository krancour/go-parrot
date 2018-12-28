package ardrone3

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// State from drone

// PilotingState ...
// TODO: Document this
type PilotingState interface{}

type pilotingState struct{}

func (p *pilotingState) ID() uint8 {
	return 4
}

func (p *pilotingState) Name() string {
	return "PilotingState"
}

func (p *pilotingState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"FlatTrimChanged",
			[]interface{}{},
			p.flatTrimChanged,
		),
		arcommands.NewD2CCommand(
			1,
			"FlyingStateChanged",
			[]interface{}{
				int32(0), // state,
			},
			p.flyingStateChanged,
		),
		arcommands.NewD2CCommand(
			2,
			"AlertStateChanged",
			[]interface{}{
				int32(0), // state,
			},
			p.alertStateChanged,
		),
		arcommands.NewD2CCommand(
			3,
			"NavigateHomeStateChanged",
			[]interface{}{
				int32(0), // state,
				int32(0), // reason,
			},
			p.navigateHomeStateChanged,
		),
		arcommands.NewD2CCommand(
			4,
			"PositionChanged",
			[]interface{}{
				float64(0), // latitude,
				float64(0), // longitude,
				float64(0), // altitude,
			},
			p.positionChanged,
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
		),
		// arcommands.NewD2CCommand(
		// 	7,
		// 	"AutoTakeOffModeChanged",
		// 	[]interface{}{
		// 		uint8(0), // state,
		// 	},
		// 	p.autoTakeOffModeChanged,
		// ),
		arcommands.NewD2CCommand(
			8,
			"AltitudeChanged",
			[]interface{}{
				float64(0), // altitude,
			},
			p.altitudeChanged,
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
		),
		arcommands.NewD2CCommand(
			10,
			"LandingStateChanged",
			[]interface{}{
				int32(0), // state,
			},
			p.landingStateChanged,
		),
		arcommands.NewD2CCommand(
			11,
			"AirSpeedChanged",
			[]interface{}{
				float32(0), // airSpeed,
			},
			p.airSpeedChanged,
		),
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
		),
		arcommands.NewD2CCommand(
			13,
			"MotionState",
			[]interface{}{
				int32(0), // state,
			},
			p.motionState,
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
		),
		arcommands.NewD2CCommand(
			15,
			"ReturnHomeBatteryCapacity",
			[]interface{}{
				int32(0), // status,
			},
			p.returnHomeBatteryCapacity,
		),
	}
}

// TODO: Implement this
// Title: Flat trim changed
// Description: Drone acknowledges that flat trim was correctly processed.
// Support: 0901;090c;090e
// Triggered: by [FlatTrim](#1-0-0).
// Result:
func (p *pilotingState) flatTrimChanged(args []interface{}) error {
	log.Info("ardrone3.flatTrimChanged() called")
	return nil
}

// TODO: Implement this
// Title: Flying state
// Description: Flying state.
// Support: 0901;090c;090e
// Triggered: when the flying state changes.
// Result:
func (p *pilotingState) flyingStateChanged(args []interface{}) error {
	// state := args[0].(int32)
	//   Drone flying state
	//   0: landed: Landed state
	//   1: takingoff: Taking off state
	//   2: hovering: Hovering / Circling (for fixed wings) state
	//   3: flying: Flying state
	//   4: landing: Landing state
	//   5: emergency: Emergency state
	//   6: usertakeoff: User take off state. Waiting for user action to take off.
	//   7: motor_ramping: Motor ramping state.
	//   8: emergency_landing: Emergency landing state. Drone autopilot has
	//      detected defective sensor(s). Only Yaw argument in PCMD is taken into
	//      account. All others flying commands are ignored.
	log.Info("ardrone3.flyingStateChanged() called")
	return nil
}

// TODO: Implement this
// Title: Alert state
// Description: Alert state.
// Support: 0901;090c;090e
// Triggered: when an alert happens on the drone.
// Result:
func (p *pilotingState) alertStateChanged(args []interface{}) error {
	// state := args[0].(int32)
	//   Drone alert state
	//   0: none: No alert
	//   1: user: User emergency alert
	//   2: cut_out: Cut out alert
	//   3: critical_battery: Critical battery alert
	//   4: low_battery: Low battery alert
	//   5: too_much_angle: The angle of the drone is too high
	log.Info("ardrone3.alertStateChanged() called")
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
func (p *pilotingState) navigateHomeStateChanged(args []interface{}) error {
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
	log.Info("ardrone3.navigateHomeStateChanged() called")
	return nil
}

// TODO: Implement this
// Title: Drone&#39;s position changed
// Description: Drone&#39;s position changed.
// Support: 0901;090c;090e
// Triggered: regularly.
// Result:
func (p *pilotingState) positionChanged(args []interface{}) error {
	// latitude := args[0].(float64)
	//   Latitude position in decimal degrees (500.0 if not available)
	// longitude := args[1].(float64)
	//   Longitude position in decimal degrees (500.0 if not available)
	// altitude := args[2].(float64)
	//   Altitude in meters (from GPS)
	log.Info("ardrone3.positionChanged() called")
	return nil
}

// TODO: Implement this
// Title: Drone&#39;s speed changed
// Description: Drone&#39;s speed changed.\n Expressed in the NED referential
//   (North-East-Down).
// Support: 0901;090c;090e
// Triggered: regularly.
// Result:
func (p *pilotingState) speedChanged(args []interface{}) error {
	// speedX := args[0].(float32)
	//   Speed relative to the North (when drone moves to the north, speed is
	//   &gt; 0) (in m/s)
	// speedY := args[1].(float32)
	//   Speed relative to the East (when drone moves to the east, speed is
	//   &gt; 0) (in m/s)
	// speedZ := args[2].(float32)
	//   Speed on the z axis (when drone moves down, speed is &gt; 0) (in m/s)
	log.Info("ardrone3.speedChanged() called")
	return nil
}

// TODO: Implement this
// Title: Drone&#39;s attitude changed
// Description: Drone&#39;s attitude changed.
// Support: 0901;090c;090e
// Triggered: regularly.
// Result:
func (p *pilotingState) attitudeChanged(args []interface{}) error {
	// roll := args[0].(float32)
	//   roll value (in radian)
	// pitch := args[1].(float32)
	//   Pitch value (in radian)
	// yaw := args[2].(float32)
	//   Yaw value (in radian)
	log.Info("ardrone3.attitudeChanged() called")
	return nil
}

// // TODO: Implement this
// // Title: Auto takeoff mode
// // Description: Auto takeoff mode
// // Support:
// // Triggered:
// // Result:
// // WARNING: Deprecated
// func (p *pilotingState) autoTakeOffModeChanged(args []interface{}) error {
// 	// state := args[0].(uint8)
// 	//   State of automatic take off mode (1 if enabled)
// 	log.Info("ardrone3.autoTakeOffModeChanged() called")
// 	return nil
// }

// TODO: Implement this
// Title: Drone&#39;s altitude changed
// Description: Drone&#39;s altitude changed.\n The altitude reported is the
//   altitude above the take off point.\n To get the altitude above sea level,
//   see [PositionChanged](#1-4-4).
// Support: 0901;090c;090e
// Triggered: regularly.
// Result:
func (p *pilotingState) altitudeChanged(args []interface{}) error {
	// altitude := args[0].(float64)
	//   Altitude in meters
	log.Info("ardrone3.altitudeChanged() called")
	return nil
}

// TODO: Implement this
// Title: Drone&#39;s location changed
// Description: Drone&#39;s location changed.\n This event is meant to replace
//   [PositionChanged](#1-4-4).
// Support: 0901:4.0.0;090c:4.0.0
// Triggered: regularly.
// Result:
func (p *pilotingState) gpsLocationChanged(args []interface{}) error {
	// latitude := args[0].(float64)
	//   Latitude location in decimal degrees (500.0 if not available)
	// longitude := args[1].(float64)
	//   Longitude location in decimal degrees (500.0 if not available)
	// altitude := args[2].(float64)
	//   Altitude location in meters.
	// latitude_accuracy := args[3].(int8)
	//   Latitude location error in meters (1 sigma/standard deviation) -1 if not
	//   available.
	// longitude_accuracy := args[4].(int8)
	//   Longitude location error in meters (1 sigma/standard deviation) -1 if not
	//   available.
	// altitude_accuracy := args[5].(int8)
	//   Altitude location error in meters (1 sigma/standard deviation) -1 if not
	//   available.
	log.Info("ardrone3.gpsLocationChanged() called")
	return nil
}

// TODO: Implement this
// Title: Landing state
// Description: Landing state.\n Only available for fixed wings (which have two
//   landing modes).
// Support: 090e
// Triggered: when the landing state changes.
// Result:
func (p *pilotingState) landingStateChanged(args []interface{}) error {
	// state := args[0].(int32)
	//   Drone landing state
	//   0: linear: Linear landing
	//   1: spiral: Spiral landing
	log.Info("ardrone3.landingStateChanged() called")
	return nil
}

// TODO: Implement this
// Title: Drone&#39;s air speed changed
// Description: Drone&#39;s air speed changed\n Expressed in the drone&#39;s
//   referential.
// Support: 090e:1.2.0
// Triggered: regularly.
// Result:
func (p *pilotingState) airSpeedChanged(args []interface{}) error {
	// airSpeed := args[0].(float32)
	//   Speed relative to air on x axis (speed is always &gt; 0) (in m/s)
	log.Info("ardrone3.airSpeedChanged() called")
	return nil
}

// TODO: Implement this
// Title: Move to changed
// Description: The drone moves or moved to a given location.
// Support: 090c:4.3.0
// Triggered: by [MoveTo](#1-0-10) or when the drone did reach the given
//   position.
// Result:
func (p *pilotingState) moveToChanged(args []interface{}) error {
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
	log.Info("ardrone3.moveToChanged() called")
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
func (p *pilotingState) motionState(args []interface{}) error {
	// state := args[0].(int32)
	//   Motion state
	//   0: steady: Drone is steady
	//   1: moving: Drone is moving
	log.Info("ardrone3.motionState() called")
	return nil
}

// TODO: Implement this
// Title: Piloted POI state
// Description: Piloted POI state.
// Support: 090c:4.3.0
// Triggered: by [StartPilotedPOI](#1-0-12) or [StopPilotedPOI](#1-0-13) or when
//   piloted POI becomes unavailable.
// Result:
func (p *pilotingState) pilotedPOI(args []interface{}) error {
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
	log.Info("ardrone3.pilotedPOI() called")
	return nil
}

// TODO: Implement this
// Title: Return home battery capacity
// Description: Battery capacity status to return home.
// Support: 090c:4.3.0
// Triggered: when the status of the battery capacity to do a return home
//   changes. This means that it is triggered either when the battery level
//   changes, when the distance to the home changes or when the position of the
//   home changes.
// Result:
func (p *pilotingState) returnHomeBatteryCapacity(args []interface{}) error {
	// status := args[0].(int32)
	//   Status of battery to return home
	//   0: OK: The battery is full enough to do a return home
	//   1: WARNING: The battery is about to be too discharged to do a return home
	//   2: CRITICAL: The battery level is too low to return to the home position
	//   3: UNKNOWN: Battery capacity to do a return home is unknown. This can be
	//      either because the home is unknown or the position of the drone is
	//      unknown, or the drone has not enough information to determine how long
	//      it takes to fly home.
	log.Info("ardrone3.returnHomeBatteryCapacity() called")
	return nil
}
