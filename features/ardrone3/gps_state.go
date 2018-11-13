package ardrone3

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

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
				uint8(0), // numberOfSatellite
			},
			g.numberOfSatelliteChanged,
		),
	}
}

// TODO: Implement this
func (g *gpsState) numberOfSatelliteChanged(args []interface{}) error {
	log.Debugf("the number of satellites changed: %v", args)
	return nil
}
