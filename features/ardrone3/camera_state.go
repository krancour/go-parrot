package ardrone3

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Camera state

// CameraState ...
// TODO: Document this
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
		// arcommands.NewD2CCommand(
		// 	0,
		// 	"Orientation",
		// 	[]interface{}{
		// 		int8(0), // tilt,
		// 		int8(0), // pan,
		// 	},
		// 	c.orientation,
		// ),
		// arcommands.NewD2CCommand(
		// 	1,
		// 	"defaultCameraOrientation",
		// 	[]interface{}{
		// 		int8(0), // tilt,
		// 		int8(0), // pan,
		// 	},
		// 	c.defaultCameraOrientation,
		// ),
		arcommands.NewD2CCommand(
			2,
			"OrientationV2",
			[]interface{}{
				float32(0), // tilt,
				float32(0), // pan,
			},
			c.orientationV2,
		),
		arcommands.NewD2CCommand(
			3,
			"defaultCameraOrientationV2",
			[]interface{}{
				float32(0), // tilt,
				float32(0), // pan,
			},
			c.defaultCameraOrientationV2,
		),
		arcommands.NewD2CCommand(
			4,
			"VelocityRange",
			[]interface{}{
				float32(0), // max_tilt,
				float32(0), // max_pan,
			},
			c.velocityRange,
		),
	}
}

// // TODO: Implement this
// // Title: Camera orientation
// // Description: Camera orientation.
// // Support: 0901;090c;090e
// // Triggered: by [SetCameraOrientation](#1-1-0).
// // Result:
// // WARNING: Deprecated
// func (c *cameraState) orientation(args []interface{}) error {
// 	// tilt := args[0].(int8)
// 	//   Tilt camera consign for the drone [-100;100]
// 	// pan := args[1].(int8)
// 	//   Pan camera consign for the drone [-100;100]
// 	log.Info("ardrone3.orientation() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Orientation of the camera center
// // Description: Orientation of the center of the camera.\n This is the value to
// //   send when you want to center the camera.
// // Support: 0901;090c;090e
// // Triggered: at connection.
// // Result:
// // WARNING: Deprecated
// func (c *cameraState) defaultCameraOrientation(args []interface{}) error {
// 	// tilt := args[0].(int8)
// 	//   Tilt value (in degree)
// 	// pan := args[1].(int8)
// 	//   Pan value (in degree)
// 	log.Info("ardrone3.defaultCameraOrientation() called")
// 	return nil
// }

// TODO: Implement this
// Title: Camera orientation
// Description: Camera orientation with float arguments.
// Support: 0901;090c;090e
// Triggered: by [SetCameraOrientationV2](#1-1-1)
// Result:
func (c *cameraState) orientationV2(args []interface{}) error {
	// tilt := args[0].(float32)
	//   Tilt camera consign for the drone [deg]
	// pan := args[1].(float32)
	//   Pan camera consign for the drone [deg]
	log.Info("ardrone3.orientationV2() called")
	return nil
}

// TODO: Implement this
// Title: Orientation of the camera center
// Description: Orientation of the center of the camera.\n This is the value to
//   send when you want to center the camera.
// Support: 0901;090c;090e
// Triggered: at connection.
// Result:
func (c *cameraState) defaultCameraOrientationV2(args []interface{}) error {
	// tilt := args[0].(float32)
	//   Tilt value [deg]
	// pan := args[1].(float32)
	//   Pan value [deg]
	log.Info("ardrone3.defaultCameraOrientationV2() called")
	return nil
}

// TODO: Implement this
// Title: Camera velocity range
// Description: Camera Orientation velocity limits.
// Support: 0901;090c;090e
// Triggered: at connection.
// Result:
func (c *cameraState) velocityRange(args []interface{}) error {
	// max_tilt := args[0].(float32)
	//   Absolute max tilt velocity [deg/s]
	// max_pan := args[1].(float32)
	//   Absolute max pan velocity [deg/s]
	log.Info("ardrone3.velocityRange() called")
	return nil
}
