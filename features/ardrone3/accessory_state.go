package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Information about the connected accessories

// AccessoryState ...
// TODO: Document this
type AccessoryState interface {
	lock.ReadLockable
	// Accessories returns a map of Accessories indexed by ID.
	Accessories() map[uint8]Accessory
}

type accessoryState struct {
	sync.RWMutex
	accessories map[uint8]Accessory
}

func newAccessoryState() *accessoryState {
	return &accessoryState{}
}

func (a *accessoryState) ID() uint8 {
	return 33
}

func (a *accessoryState) Name() string {
	return "AccessoryState"
}

func (a *accessoryState) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
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
			log,
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
			log,
		),
	}
}

// connectedAccessories is invoked by the device to list all connected
// accessories.
func (a *accessoryState) connectedAccessories(
	args []interface{},
	log *log.Entry,
) error {
	a.Lock()
	defer a.Unlock()
	flags := args[4].(uint8)
	// 0x01: First: indicates it's the first element of the list.
	// 0x02: Last: indicates it's the last element of the list.
	// 0x04: Empty: indicates the list is empty. All other arguments should be
	//   ignored.
	// 0x08: Remove: This value should be removed from the existing list.
	if flags&4 == 4 {
		log.Debug("connected accessories updated with empty list")
		return nil
	}
	accessoryID := args[0].(uint8)
	if flags&8 == 8 {
		delete(a.accessories, accessoryID)
		log.WithField(
			"id", accessoryID,
		).Debug("connected accessory removed")
		return nil
	}
	// If we get to here, we should add or update the accessory.
	accessoryIface, ok := a.accessories[accessoryID]
	var acc *accessory
	if ok {
		acc = accessoryIface.(*accessory)
	} else {
		acc = &accessory{
			id: accessoryID,
		}
		a.accessories[accessoryID] = acc
	}
	acc.tipe = args[1].(int32)
	acc.uid = args[2].(string)
	acc.softwareVersion = args[3].(string)
	log.WithField(
		"id", acc.id,
	).WithField(
		"accessory_type", acc.tipe,
	).WithField(
		"uid", acc.uid,
	).WithField(
		"swVersion", acc.softwareVersion,
	).Debug("accessory added or updated")
	return nil
}

// battery is invoked by the device when the battery level of a connected
// accessory changes.
func (a *accessoryState) battery(args []interface{}, log *log.Entry) error {
	a.Lock()
	defer a.Unlock()
	flags := args[2].(uint8)
	// 0x01: First: indicates it's the first element of the list.
	// 0x02: Last: indicates it's the last element of the list.
	// 0x04: Empty: indicates the list is empty. All other arguments should be
	//   ignored.
	// 0x08: Remove: This value should be removed from the existing list.
	if flags&4 == 4 {
		log.Debug("connected accessories battery levels updated with empty list")
		return nil
	}
	accessoryID := args[0].(uint8)
	if flags&8 == 8 {
		delete(a.accessories, accessoryID)
		log.WithField(
			"id", accessoryID,
		).Debug("connected accessory battery level removed")
		return nil
	}
	// If we get to here, we should add or update the accessory.
	accessoryIface, ok := a.accessories[accessoryID]
	var acc *accessory
	if ok {
		acc = accessoryIface.(*accessory)
	} else {
		acc = &accessory{
			id: accessoryID,
		}
		a.accessories[accessoryID] = acc
	}
	acc.batteryPercent = args[1].(uint8)
	log.WithField(
		"id", acc.id,
	).WithField(
		"batteryLevel", acc.batteryPercent,
	).Debug("accessory battery level added or updated")
	return nil
}

func (a *accessoryState) Accessories() map[uint8]Accessory {
	return a.accessories
}
