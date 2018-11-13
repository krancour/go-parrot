package ardrone3

import (
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/protocols/arnetwork"
)

// Feature ...
// TODO: Document this
type Feature interface {
	arcommands.D2CFeature
}

type feature struct {
	c2dChs map[uint8]chan<- arnetwork.Frame
}

// NewFeature ...
// TODO: Document this
func NewFeature(
	c2dChs map[uint8]chan<- arnetwork.Frame,
) Feature {
	return &feature{
		c2dChs: c2dChs,
	}
}

func (f *feature) ID() uint8 {
	return 1
}

func (f *feature) Name() string {
	return "ardrone3"
}

// TODO: Add stuff!
func (f *feature) D2CClasses() map[uint8]arcommands.D2CClass {
	return map[uint8]arcommands.D2CClass{

		4: arcommands.NewD2CClass(
			"PilotingState",
			map[uint16]arcommands.D2CCommand{
				4: arcommands.NewD2CCommand(
					"PositionChanged",
					[]interface{}{
						float64(0), // latitude
						float64(0), // longitude
						float64(0), // altitude
					},
					f.positionChanged,
				),
				5: arcommands.NewD2CCommand(
					"SpeedChanged",
					[]interface{}{
						float32(0), // speedX
						float32(0), // speedY
						float32(0), // speedZ
					},
					f.speedChanged,
				),
				6: arcommands.NewD2CCommand(
					"AttitudeChanged",
					[]interface{}{
						float32(0), // roll
						float32(0), // pitch
						float32(0), // yaw
					},
					f.attitudeChanged,
				),
				8: arcommands.NewD2CCommand(
					"AltitudeChanged",
					[]interface{}{
						float64(0), // altitude
					},
					f.altitudeChanged,
				),
				9: arcommands.NewD2CCommand(
					"GpsLocationChanged",
					[]interface{}{
						float64(0), // latitude
						float64(0), // longitude
						float64(0), // altitude
						int8(0),    // latitude_accuracy
						int8(0),    // longitude_accuracy
						int8(0),    // altitude_accuracy
					},
					f.gpsLocationChanged,
				),
			},
		),

		25: arcommands.NewD2CClass(
			"CameraState",
			map[uint16]arcommands.D2CCommand{
				0: arcommands.NewD2CCommand(
					"Orientation",
					[]interface{}{
						int8(0), // tilt
						int8(0), // pan
					},
					f.orientation,
				),
			},
		),

		31: arcommands.NewD2CClass(
			"GPSState",
			map[uint16]arcommands.D2CCommand{
				0: arcommands.NewD2CCommand(
					"NumberOfSatelliteChanged",
					[]interface{}{
						uint8(0), // numberOfSatellite
					},
					f.numberOfSatelliteChanged,
				),
			},
		),
	}
}
