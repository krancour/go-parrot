package ardrone3

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

type CameraState interface{}

type cameraState struct{}

func (c *cameraState) ID() uint8 {
	return 25
}

func (c *cameraState) Name() string {
	return "CameraState"
}

func (c *cameraState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"Orientation",
			[]interface{}{
				int8(0), // tilt
				int8(0), // pan
			},
			c.orientation,
		),
	}
}

// TODO: Implement this
func (c *cameraState) orientation(args []interface{}) error {
	log.Debugf("the camera orientation changed: %v", args)
	return nil
}
