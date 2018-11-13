package ardrone3

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

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
			4,
			"PositionChanged",
			[]interface{}{
				float64(0), // latitude
				float64(0), // longitude
				float64(0), // altitude
			},
			p.positionChanged,
		),
		arcommands.NewD2CCommand(
			5,
			"SpeedChanged",
			[]interface{}{
				float32(0), // speedX
				float32(0), // speedY
				float32(0), // speedZ
			},
			p.speedChanged,
		),
		arcommands.NewD2CCommand(
			6,
			"AttitudeChanged",
			[]interface{}{
				float32(0), // roll
				float32(0), // pitch
				float32(0), // yaw
			},
			p.attitudeChanged,
		),
		arcommands.NewD2CCommand(
			8,
			"AltitudeChanged",
			[]interface{}{
				float64(0), // altitude
			},
			p.altitudeChanged,
		),
		arcommands.NewD2CCommand(
			9,
			"GpsLocationChanged",
			[]interface{}{
				float64(0), // latitude
				float64(0), // longitude
				float64(0), // altitude
				int8(0),    // latitude_accuracy
				int8(0),    // longitude_accuracy
				int8(0),    // altitude_accuracy
			},
			p.gpsLocationChanged,
		),
	}
}

// TODO: Implement this
func (p *pilotingState) positionChanged(args []interface{}) error {
	log.Debugf("the position changed: %v", args)
	return nil
}

// TODO: Implement this
func (p *pilotingState) speedChanged(args []interface{}) error {
	log.Debugf("the speed changed: %v", args)
	return nil
}

// TODO: Implement this
func (p *pilotingState) attitudeChanged(args []interface{}) error {
	log.Debugf("the attitude changed: %v", args)
	return nil
}

// TODO: Implement this
func (p *pilotingState) altitudeChanged(args []interface{}) error {
	log.Debugf("the altitude changed: %v", args)
	return nil
}

// TODO: Implement this
func (p *pilotingState) gpsLocationChanged(args []interface{}) error {
	log.Debugf("the gps location changed: %v", args)
	return nil
}
