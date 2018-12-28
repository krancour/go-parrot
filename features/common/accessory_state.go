package common

// // Accessories-related commands.

// // AccessoryState ...
// // TODO: Document this
// type AccessoryState interface{}

// type accessoryState struct{}

// func (a *accessoryState) ID() uint8 {
// 	return 27
// }

// func (a *accessoryState) Name() string {
// 	return "AccessoryState"
// }

// func (a *accessoryState) D2CCommands() []arcommands.D2CCommand {
// 	return []arcommands.D2CCommand{
// 		arcommands.NewD2CCommand(
// 			0,
// 			"SupportedAccessoriesListChanged",
// 			[]interface{}{
// 				int32(0), // accessory,
// 			},
// 			a.supportedAccessoriesListChanged,
// 		),
// 		arcommands.NewD2CCommand(
// 			1,
// 			"AccessoryConfigChanged",
// 			[]interface{}{
// 				int32(0), // newAccessory,
// 				int32(0), // error,
// 			},
// 			a.accessoryConfigChanged,
// 		),
// 		arcommands.NewD2CCommand(
// 			2,
// 			"AccessoryConfigModificationEnabled",
// 			[]interface{}{
// 				uint8(0), // enabled,
// 			},
// 			a.accessoryConfigModificationEnabled,
// 		),
// 	}
// }

// // TODO: Implement this
// // Title: Supported accessories list
// // Description: Supported accessories list.
// // Support: 0902;0905;0906;0907;0909;090a
// // Triggered: at connection.
// // Result:
// func (a *accessoryState) supportedAccessoriesListChanged(
// 	args []interface{},
// ) error {
// 	// accessory := args[0].(int32)
// 	//   Accessory configurations supported by the product.
// 	//   0: NO_ACCESSORY: No accessory.
// 	//   1: STD_WHEELS: Standard wheels
// 	//   2: TRUCK_WHEELS: Truck wheels
// 	//   3: HULL: Hull
// 	//   4: HYDROFOIL: Hydrofoil
// 	log.Info("common.supportedAccessoriesListChanged() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Accessory config
// // Description: Accessory config.
// // Support: 0902;0905;0906;0907;0909;090a
// // Triggered: by [DeclareAccessory](#0-26-0).
// // Result:
// func (a *accessoryState) accessoryConfigChanged(args []interface{}) error {
// 	// newAccessory := args[0].(int32)
// 	//   Accessory configuration reported by firmware.
// 	//   0: UNCONFIGURED: No accessory configuration set. Controller needs to set
// 	//      one.
// 	//   1: NO_ACCESSORY: No accessory.
// 	//   2: STD_WHEELS: Standard wheels
// 	//   3: TRUCK_WHEELS: Truck wheels
// 	//   4: HULL: Hull
// 	//   5: HYDROFOIL: Hydrofoil
// 	//   6: IN_PROGRESS: Configuration in progress.
// 	// error := args[1].(int32)
// 	//   Error code.
// 	//   0: OK: No error. Accessory config change successful.
// 	//   1: UNKNOWN: Cannot change accessory configuration for some reason.
// 	//   2: FLYING: Cannot change accessory configuration while flying.
// 	log.Info("common.accessoryConfigChanged() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Accessory declaration availability
// // Description: Availability to declare or not an accessory.
// // Support: 0902;0905;0906;0907;0909;090a
// // Triggered: when the availability changes.
// // Result:
// func (a *accessoryState) accessoryConfigModificationEnabled(
// 	args []interface{},
// ) error {
// 	// enabled := args[0].(uint8)
// 	//   1 if the modification of the accessory Config is enabled, 0 otherwise
// 	log.Info("common.accessoryConfigModificationEnabled() called")
// 	return nil
// }
