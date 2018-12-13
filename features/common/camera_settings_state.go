package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Status of the camera settings

// CameraSettingsState ...
// TODO: Document this
type CameraSettingsState interface{}

type cameraSettingsState struct{}

func (c *cameraSettingsState) ID() uint8 {
	return 15
}

func (c *cameraSettingsState) Name() string {
	return "CameraSettingsState"
}

func (c *cameraSettingsState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"CameraSettingsChanged",
			[]interface{}{
				float32(0), // fov,
				float32(0), // panMax,
				float32(0), // panMin,
				float32(0), // tiltMax,
				float32(0), // tiltMin,
			},
			c.cameraSettingsChanged,
		),
	}
}

// TODO: Implement this
// Title: Camera info
// Description: Camera info.
// Support: 0901;090c;090e
// Triggered: at connection.
// Result:
func (c *cameraSettingsState) cameraSettingsChanged(args []interface{}) error {
	// fov := args[0].(float32)
	//   Value of the camera horizontal fov (in degree)
	// panMax := args[1].(float32)
	//   Value of max pan (right pan) (in degree)
	// panMin := args[2].(float32)
	//   Value of min pan (left pan) (in degree)
	// tiltMax := args[3].(float32)
	//   Value of max tilt (top tilt) (in degree)
	// tiltMin := args[4].(float32)
	//   Value of min tilt (bottom tilt) (in degree)
	log.Info("common.cameraSettingsChanged() called")
	return nil
}
