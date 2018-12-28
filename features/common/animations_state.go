package common

// import (
// 	log "github.com/Sirupsen/logrus"
// 	"github.com/krancour/go-parrot/protocols/arcommands"
// )

// // Animations-related notification/feedback commands.

// // AnimationsState ...
// // TODO: Document this
// type AnimationsState interface{}

// type animationsState struct{}

// func (a *animationsState) ID() uint8 {
// 	return 25
// }

// func (a *animationsState) Name() string {
// 	return "AnimationsState"
// }

// func (a *animationsState) D2CCommands() []arcommands.D2CCommand {
// 	return []arcommands.D2CCommand{
// 		arcommands.NewD2CCommand(
// 			0,
// 			"List",
// 			[]interface{}{
// 				int32(0), // anim,
// 				int32(0), // state,
// 				int32(0), // error,
// 			},
// 			a.list,
// 		),
// 	}
// }

// // TODO: Implement this
// // Title: Animation state list
// // Description: Paramaterless animations state list.
// // Support: 0902;0905;0906;0907;0909
// // Triggered: when the list of available animations changes and also when an
// //   animation state changes (can be triggered by [StartAnim](#0-24-0),
// //   [StopAnim](#0-24-1) or [StopAllAnims](#0-24-2).
// // Result:
// func (a *animationsState) list(args []interface{}) error {
// 	// anim := args[0].(int32)
// 	//   Animation type.
// 	//   0: HEADLIGHTS_FLASH: Flash headlights.
// 	//   1: HEADLIGHTS_BLINK: Blink headlights.
// 	//   2: HEADLIGHTS_OSCILLATION: Oscillating headlights.
// 	//   3: SPIN: Spin animation.
// 	//   4: TAP: Tap animation.
// 	//   5: SLOW_SHAKE: Slow shake animation.
// 	//   6: METRONOME: Metronome animation.
// 	//   7: ONDULATION: Standing dance animation.
// 	//   8: SPIN_JUMP: Spin jump animation.
// 	//   9: SPIN_TO_POSTURE: Spin that end in standing posture, or in jumper if it
// 	//      was standing animation.
// 	//   10: SPIRAL: Spiral animation.
// 	//   11: SLALOM: Slalom animation.
// 	//   12: BOOST: Boost animation.
// 	//   13: LOOPING: Make a looping. (Only for WingX)
// 	//   14: BARREL_ROLL_180_RIGHT: Make a barrel roll of 180 degree turning on
// 	//       right. (Only for WingX)
// 	//   15: BARREL_ROLL_180_LEFT: Make a barrel roll of 180 degree turning on
// 	//       left. (Only for WingX)
// 	//   16: BACKSWAP: Put the drone upside down. (Only for WingX)
// 	// state := args[1].(int32)
// 	//   State of the animation
// 	//   0: stopped: animation is stopped
// 	//   1: started: animation is started
// 	//   2: notAvailable: The animation is not available
// 	// error := args[2].(int32)
// 	//   Error to explain the state
// 	//   0: ok: No Error
// 	//   1: unknown: Unknown generic error
// 	log.Info("common.list() called")
// 	return nil
// }
