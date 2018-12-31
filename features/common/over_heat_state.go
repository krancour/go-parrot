package common

// import (
// 	"github.com/krancour/go-parrot/protocols/arcommands"
// )

// // Overheat state from product

// // OverHeatState ...
// // TODO: Document this
// type OverHeatState interface{}

// type overHeatState struct{}

// func (o *overHeatState) ID() uint8 {
// 	return 7
// }

// func (o *overHeatState) Name() string {
// 	return "OverHeatState"
// }

// func (o *overHeatState) D2CCommands() []arcommands.D2CCommand {
// 	return []arcommands.D2CCommand{
// 		arcommands.NewD2CCommand(
// 			0,
// 			"OverHeatChanged",
// 			[]interface{}{},
// 			o.overHeatChanged,
// 		),
// 		arcommands.NewD2CCommand(
// 			1,
// 			"OverHeatRegulationChanged",
// 			[]interface{}{
// 				uint8(0), // regulationType,
// 			},
// 			o.overHeatRegulationChanged,
// 		),
// 	}
// }

// // TODO: Implement this
// // Title: Overheat
// // Description: Overheat temperature reached.
// // Support:
// // Triggered:
// // Result:
// // WARNING: Deprecated
// func (o *overHeatState) overHeatChanged(args []interface{}) error {
// 	log.Info("common.overHeatChanged() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Overheat regulation type
// // Description: Overheat regulation type.
// // Support:
// // Triggered:
// // Result:
// func (o *overHeatState) overHeatRegulationChanged(args []interface{}) error {
// 	// regulationType := args[0].(uint8)
// 	//   Type of overheat regulation : 0 for ventilation, 1 for switch off
// 	log.Info("common.overHeatRegulationChanged() called")
// 	return nil
// }
