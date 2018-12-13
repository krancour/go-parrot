package ardrone3

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Information about the connected accessories

// AccessoryState ...
// TODO: Document this
type AccessoryState interface{}

type accessoryState struct{}

func (a *accessoryState) ID() uint8 {
	return 33
}

func (a *accessoryState) Name() string {
	return "AccessoryState"
}

func (a *accessoryState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"ConnectedAccessories",
			[]interface{}{
				uint8(0),  // id,
				int32(0),  // accessory_type,
				string(0), // uid,
				string(0), // swVersion,
				uint8(0),  // list_flags,
			},
			a.connectedAccessories,
		),
		arcommands.NewD2CCommand(
			1,
			"Battery",
			[]interface{}{
				uint8(0), // id,
				uint8(0), // batteryLevel,
				uint8(0), // list_flags,
			},
			a.battery,
		),
	}
}

// TODO: Implement this
// Title: List of connected accessories
// Description: List of all connected accessories. This event presents the list
//   of all connected accessories. To actually use the component, use the
//   component dedicated feature.
// Support: 090e:1.5.0
// Triggered: at connection or when an accessory is connected.
// Result:
func (a *accessoryState) connectedAccessories(args []interface{}) error {
	// id := args[0].(uint8)
	//   Id of the accessory for the session.
	// accessory_type := args[1].(int32)
	//   Accessory type
	//   0: sequoia: Parrot Sequoia (multispectral camera for agriculture)
	//   1: flir: FLIR camera (thermal&#43;rgb camera)
	// uid := args[2].(string)
	//   Unique Id of the accessory. This id is unique by accessory_type.
	// swVersion := args[3].(string)
	//   Software Version of the accessory.
	// list_flags := args[4].(uint8)
	//   List entry attribute Bitfield. 0x01: First: indicate it&#39;s the first
	//   element of the list. 0x02: Last: indicate it&#39;s the last element of
	//   the list. 0x04: Empty: indicate the list is empty (implies First/Last).
	//   All other arguments should be ignored. 0x08: Remove: This value should be
	//   removed from the existing list.
	log.Info("ardrone3.connectedAccessories() called")
	return nil
}

// TODO: Implement this
// Title: Connected accessories battery
// Description: Connected accessories battery.
// Support: none
// Triggered:
// Result:
func (a *accessoryState) battery(args []interface{}) error {
	// id := args[0].(uint8)
	//   Id of the accessory for the session.
	// batteryLevel := args[1].(uint8)
	//   Battery level in percentage.
	// list_flags := args[2].(uint8)
	//   List entry attribute Bitfield. 0x01: First: indicate it&#39;s the first
	//   element of the list. 0x02: Last: indicate it&#39;s the last element of
	//   the list. 0x04: Empty: indicate the list is empty (implies First/Last).
	//   All other arguments should be ignored. 0x08: Remove: This value should be
	//   removed from the existing list.
	log.Info("ardrone3.battery() called")
	return nil
}
