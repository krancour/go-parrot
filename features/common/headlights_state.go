package common

// import (
// 	log "github.com/Sirupsen/logrus"
// 	"github.com/krancour/go-parrot/protocols/arcommands"
// )

// // Get information about the state of the Evo variants&#39; LEDs.

// // HeadlightsState ...
// // TODO: Document this
// type HeadlightsState interface{}

// type headlightsState struct{}

// func (h *headlightsState) ID() uint8 {
// 	return 23
// }

// func (h *headlightsState) Name() string {
// 	return "HeadlightsState"
// }

// func (h *headlightsState) D2CCommands() []arcommands.D2CCommand {
// 	return []arcommands.D2CCommand{
// 		arcommands.NewD2CCommand(
// 			0,
// 			"intensityChanged",
// 			[]interface{}{
// 				uint8(0), // left,
// 				uint8(0), // right,
// 			},
// 			h.intensityChanged,
// 		),
// 	}
// }

// // TODO: Implement this
// // Title: LEDs intensity
// // Description: Lighting LEDs intensity.
// // Support: 0905;0906;0907
// // Triggered: by [SetLedsIntensity](#0-22-0).
// // Result:
// func (h *headlightsState) intensityChanged(args []interface{}) error {
// 	// left := args[0].(uint8)
// 	//   The intensity value for the left LED (0 through 255).
// 	// right := args[1].(uint8)
// 	//   The intensity value for the right LED (0 through 255).
// 	log.Info("common.intensityChanged() called")
// 	return nil
// }
